package delivery

import (
	"log"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

type Route struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id,omitempty"`
	DriverID  string    `json:"driver_id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Stop struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id,omitempty"`
	RouteID      string     `json:"route_id"`
	CustomerID   string     `json:"customer_id"`
	InvoiceID    *string    `json:"invoice_id,omitempty"`
	StopOrder    int        `json:"stop_order"`
	Status       string     `json:"status"`
	SignatureData *string   `json:"signature_data,omitempty"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
}

func HandleCreateRoute(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			DriverID string `json:"driver_id"`
			Date     string `json:"date"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.DriverID == "" {
			http.Error(w, "driver_id required", http.StatusBadRequest)
			return
		}
		var date time.Time
		var err error
		if body.Date != "" {
			date, err = time.Parse("2006-01-02", body.Date)
			if err != nil {
				http.Error(w, "invalid date", http.StatusBadRequest)
				return
			}
		} else {
			date = time.Now()
		}
		var rt Route
		err = mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO delivery_routes (tenant_id, driver_id, date, status)
				 VALUES ($1,$2,$3,'pending')
				 RETURNING id, tenant_id, driver_id, date, status, created_at`,
				mw.TenantIDFrom(r.Context()), body.DriverID, date,
			).Scan(&rt.ID, &rt.TenantID, &rt.DriverID, &rt.Date, &rt.Status, &rt.CreatedAt)
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(rt)
	}
}

func HandleListRoutes(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		if perPage < 1 {
			perPage = 50
		}
		if perPage > 200 {
			perPage = 200
		}
		offset := (page - 1) * perPage
		var routes []Route
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT id, tenant_id, driver_id, date, status, created_at
				 FROM delivery_routes ORDER BY date DESC LIMIT $1 OFFSET $2`, perPage, offset)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var rt Route
				if err := rows.Scan(&rt.ID, &rt.TenantID, &rt.DriverID, &rt.Date, &rt.Status, &rt.CreatedAt); err != nil {
					return err
				}
				routes = append(routes, rt)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if routes == nil {
			routes = []Route{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes)
	}
}

func HandleGetRoute(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var rt Route
		var stops []Stop
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`SELECT id, tenant_id, driver_id, date, status, created_at
				 FROM delivery_routes WHERE id=$1`, id,
			).Scan(&rt.ID, &rt.TenantID, &rt.DriverID, &rt.Date, &rt.Status, &rt.CreatedAt); err != nil {
				return err
			}
			rows, err := tx.Query(r.Context(),
				`SELECT id, tenant_id, route_id, customer_id, invoice_id, stop_order, status, signature_data, delivered_at, notes
				 FROM delivery_stops WHERE route_id=$1 ORDER BY stop_order`, id)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var s Stop
				if err := rows.Scan(&s.ID, &s.TenantID, &s.RouteID, &s.CustomerID, &s.InvoiceID,
					&s.StopOrder, &s.Status, &s.SignatureData, &s.DeliveredAt, &s.Notes); err != nil {
					return err
				}
				stops = append(stops, s)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if stops == nil {
			stops = []Stop{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"route": rt, "stops": stops})
	}
}

func HandleAddStop(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			CustomerID string  `json:"customer_id"`
			InvoiceID  *string `json:"invoice_id"`
			StopOrder  int     `json:"stop_order"`
			Notes      *string `json:"notes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.CustomerID == "" {
			http.Error(w, "customer_id required", http.StatusBadRequest)
			return
		}
		var s Stop
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO delivery_stops (tenant_id, route_id, customer_id, invoice_id, stop_order, status, notes)
				 VALUES ($1,$2,$3,$4,$5,'pending',$6)
				 RETURNING id, tenant_id, route_id, customer_id, invoice_id, stop_order, status, signature_data, delivered_at, notes`,
				mw.TenantIDFrom(r.Context()), id, body.CustomerID, body.InvoiceID, body.StopOrder, body.Notes,
			).Scan(&s.ID, &s.TenantID, &s.RouteID, &s.CustomerID, &s.InvoiceID,
				&s.StopOrder, &s.Status, &s.SignatureData, &s.DeliveredAt, &s.Notes)
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(s)
	}
}

func HandleUpdateStopStatus(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stopID := chi.URLParam(r, "stopId")
		var body struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		switch body.Status {
		case "pending", "delivered", "rejected", "skipped":
		default:
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if body.Status == "delivered" {
				_, err := tx.Exec(r.Context(),
					`UPDATE delivery_stops SET status=$2, delivered_at=NOW() WHERE id=$1`,
					stopID, body.Status)
				return err
			}
			_, err := tx.Exec(r.Context(),
				`UPDATE delivery_stops SET status=$2 WHERE id=$1`, stopID, body.Status)
			return err
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleRecordSignature(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stopID := chi.URLParam(r, "stopId")
		var body struct {
			SignatureData string `json:"signature_data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.SignatureData == "" {
			http.Error(w, "signature_data required", http.StatusBadRequest)
			return
		}
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			_, err := tx.Exec(r.Context(),
				`UPDATE delivery_stops SET signature_data=$2, status='delivered', delivered_at=NOW() WHERE id=$1`,
				stopID, body.SignatureData)
			return err
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
