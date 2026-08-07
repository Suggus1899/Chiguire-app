package fiscaldevice

import (
	"log"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

var errNotFound = errors.New("not found")

type FiscalDevice struct {
	ID         string    `json:"id"`
	BranchID   string    `json:"branch_id"`
	Name       string    `json:"name"`
	DeviceType string    `json:"device_type"`
	Model      string    `json:"model"`
	Serial     string    `json:"serial"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

const fdCols = `id, COALESCE(branch_id::text,''), name, device_type, COALESCE(model,''), COALESCE(serial,''), is_active, created_at, updated_at`

type DocumentSequence struct {
	ID             string    `json:"id"`
	FiscalDeviceID string    `json:"fiscal_device_id"`
	DocType        string    `json:"doc_type"`
	Prefix         string    `json:"prefix"`
	Suffix         string    `json:"suffix"`
	LastSeq        int       `json:"last_seq"`
	CreatedAt      time.Time `json:"created_at"`
}

const dsCols = `id, COALESCE(fiscal_device_id::text,''), doc_type, prefix, suffix, last_seq, created_at`

type Contingency struct {
	ID        string    `json:"id"`
	BranchID  string    `json:"branch_id"`
	DocType   string    `json:"doc_type"`
	Prefix    string    `json:"prefix"`
	StartSeq  int       `json:"start_seq"`
	EndSeq    int       `json:"end_seq"`
	CurrentSeq int      `json:"current_seq"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

const cbCols = `id, COALESCE(branch_id::text,''), doc_type, prefix, start_seq, end_seq, current_seq, is_active, created_at`

func HandleCreate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			BranchID   string `json:"branch_id"`
			Name       string `json:"name"`
			DeviceType string `json:"device_type"`
			Model      string `json:"model"`
			Serial     string `json:"serial"`
			IsActive   *bool  `json:"is_active"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.DeviceType == "" {
			http.Error(w, "name and device_type are required", http.StatusBadRequest)
			return
		}
		isActive := true
		if body.IsActive != nil {
			isActive = *body.IsActive
		}

		var fd FiscalDevice
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO fiscal_devices (tenant_id, branch_id, name, device_type, model, serial, is_active)
				 VALUES ($1,NULLIF($2,'')::uuid,$3,$4,$5,$6,$7)
				 RETURNING `+fdCols,
				mw.TenantIDFrom(r.Context()), body.BranchID, body.Name, body.DeviceType, body.Model, body.Serial, isActive,
			).Scan(&fd.ID, &fd.BranchID, &fd.Name, &fd.DeviceType, &fd.Model, &fd.Serial, &fd.IsActive, &fd.CreatedAt, &fd.UpdatedAt)
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fd)
	}
}

func HandleList(pool *pgxpool.Pool) http.HandlerFunc {
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
		var fds []FiscalDevice
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+fdCols+` FROM fiscal_devices ORDER BY created_at DESC LIMIT $1 OFFSET $2`, perPage, offset)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var fd FiscalDevice
				if err := rows.Scan(&fd.ID, &fd.BranchID, &fd.Name, &fd.DeviceType, &fd.Model, &fd.Serial, &fd.IsActive, &fd.CreatedAt, &fd.UpdatedAt); err != nil {
					return err
				}
				fds = append(fds, fd)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if fds == nil {
			fds = []FiscalDevice{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fds)
	}
}

func HandleGet(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var fd FiscalDevice
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`SELECT `+fdCols+` FROM fiscal_devices WHERE id=$1`, id,
			).Scan(&fd.ID, &fd.BranchID, &fd.Name, &fd.DeviceType, &fd.Model, &fd.Serial, &fd.IsActive, &fd.CreatedAt, &fd.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fd)
	}
}

func HandleUpdate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			BranchID   string `json:"branch_id"`
			Name       string `json:"name"`
			DeviceType string `json:"device_type"`
			Model      string `json:"model"`
			Serial     string `json:"serial"`
			IsActive   *bool  `json:"is_active"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		var fd FiscalDevice
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if body.IsActive != nil {
				return tx.QueryRow(r.Context(),
					`UPDATE fiscal_devices SET branch_id=NULLIF($1,'')::uuid, name=$2, device_type=$3, model=$4, serial=$5, is_active=$6, updated_at=NOW()
					 WHERE id=$7
					 RETURNING `+fdCols,
					body.BranchID, body.Name, body.DeviceType, body.Model, body.Serial, *body.IsActive, id,
				).Scan(&fd.ID, &fd.BranchID, &fd.Name, &fd.DeviceType, &fd.Model, &fd.Serial, &fd.IsActive, &fd.CreatedAt, &fd.UpdatedAt)
			}
			return tx.QueryRow(r.Context(),
				`UPDATE fiscal_devices SET branch_id=NULLIF($1,'')::uuid, name=$2, device_type=$3, model=$4, serial=$5, updated_at=NOW()
				 WHERE id=$6
				 RETURNING `+fdCols,
				body.BranchID, body.Name, body.DeviceType, body.Model, body.Serial, id,
			).Scan(&fd.ID, &fd.BranchID, &fd.Name, &fd.DeviceType, &fd.Model, &fd.Serial, &fd.IsActive, &fd.CreatedAt, &fd.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fd)
	}
}

func HandleDelete(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			ct, err := tx.Exec(r.Context(), `DELETE FROM fiscal_devices WHERE id=$1`, id)
			if err != nil {
				return err
			}
			if ct.RowsAffected() == 0 {
				return errNotFound
			}
			return nil
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleCreateSequence(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			FiscalDeviceID string `json:"fiscal_device_id"`
			DocType        string `json:"doc_type"`
			Prefix         string `json:"prefix"`
			Suffix         string `json:"suffix"`
			LastSeq        int    `json:"last_seq"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.DocType == "" {
			http.Error(w, "doc_type is required", http.StatusBadRequest)
			return
		}

		var ds DocumentSequence
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO document_sequences (tenant_id, fiscal_device_id, doc_type, prefix, suffix, last_seq)
				 VALUES ($1,NULLIF($2,'')::uuid,$3,$4,$5,$6)
				 RETURNING `+dsCols,
				mw.TenantIDFrom(r.Context()), body.FiscalDeviceID, body.DocType, body.Prefix, body.Suffix, body.LastSeq,
			).Scan(&ds.ID, &ds.FiscalDeviceID, &ds.DocType, &ds.Prefix, &ds.Suffix, &ds.LastSeq, &ds.CreatedAt)
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ds)
	}
}

func HandleListSequences(pool *pgxpool.Pool) http.HandlerFunc {
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
		var dss []DocumentSequence
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+dsCols+` FROM document_sequences ORDER BY created_at DESC LIMIT $1 OFFSET $2`, perPage, offset)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var ds DocumentSequence
				if err := rows.Scan(&ds.ID, &ds.FiscalDeviceID, &ds.DocType, &ds.Prefix, &ds.Suffix, &ds.LastSeq, &ds.CreatedAt); err != nil {
					return err
				}
				dss = append(dss, ds)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if dss == nil {
			dss = []DocumentSequence{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dss)
	}
}

func HandleCreateContingency(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			BranchID  string `json:"branch_id"`
			DocType   string `json:"doc_type"`
			Prefix    string `json:"prefix"`
			StartSeq  int    `json:"start_seq"`
			EndSeq    int    `json:"end_seq"`
			IsActive  *bool  `json:"is_active"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.DocType == "" || body.Prefix == "" {
			http.Error(w, "doc_type and prefix are required", http.StatusBadRequest)
			return
		}
		isActive := true
		if body.IsActive != nil {
			isActive = *body.IsActive
		}

		var cb Contingency
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO contingency_books (tenant_id, branch_id, doc_type, prefix, start_seq, end_seq, current_seq, is_active)
				 VALUES ($1,NULLIF($2,'')::uuid,$3,$4,$5,$6,$5,$7)
				 RETURNING `+cbCols,
				mw.TenantIDFrom(r.Context()), body.BranchID, body.DocType, body.Prefix, body.StartSeq, body.EndSeq, isActive,
			).Scan(&cb.ID, &cb.BranchID, &cb.DocType, &cb.Prefix, &cb.StartSeq, &cb.EndSeq, &cb.CurrentSeq, &cb.IsActive, &cb.CreatedAt)
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(cb)
	}
}

func HandleListContingency(pool *pgxpool.Pool) http.HandlerFunc {
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
		var cbs []Contingency
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+cbCols+` FROM contingency_books ORDER BY created_at DESC LIMIT $1 OFFSET $2`, perPage, offset)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var cb Contingency
				if err := rows.Scan(&cb.ID, &cb.BranchID, &cb.DocType, &cb.Prefix, &cb.StartSeq, &cb.EndSeq, &cb.CurrentSeq, &cb.IsActive, &cb.CreatedAt); err != nil {
					return err
				}
				cbs = append(cbs, cb)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if cbs == nil {
			cbs = []Contingency{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cbs)
	}
}
