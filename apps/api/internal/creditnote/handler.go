package creditnote

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

var errNotFound = errors.New("not found")

type CreditNote struct {
	ID        string     `json:"id"`
	InvoiceID string     `json:"invoice_id"`
	Number    string     `json:"number"`
	Type      string     `json:"type"`
	Reason    string     `json:"reason"`
	Subtotal  string     `json:"subtotal"`
	TaxTotal  string     `json:"tax_total"`
	Total     string     `json:"total"`
	Status    string     `json:"status"`
	IssuedAt  *time.Time `json:"issued_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Items     []Item     `json:"items,omitempty"`
}

type Item struct {
	ID          string `json:"id"`
	ProductID   string `json:"product_id"`
	Description string `json:"description"`
	Qty         string `json:"qty"`
	UnitPrice   string `json:"unit_price"`
	TaxPct      string `json:"tax_pct"`
	LineTotal   string `json:"line_total"`
}

const cnCols = `id, COALESCE(invoice_id::text,''), COALESCE(number,''), type, COALESCE(reason,''), subtotal::text, tax_total::text, total::text, status, issued_at, created_at, updated_at`
const itemCols = `id, COALESCE(product_id::text,''), description, qty::text, unit_price::text, tax_pct::text, line_total::text`

type itemInput struct {
	ProductID   string `json:"product_id"`
	Description string `json:"description"`
	Qty         string `json:"qty"`
	UnitPrice   string `json:"unit_price"`
	TaxPct      string `json:"tax_pct"`
}

func HandleCreate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			InvoiceID string      `json:"invoice_id"`
			Type      string      `json:"type"`
			Reason    string      `json:"reason"`
			Items     []itemInput `json:"items"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Type == "" {
			body.Type = "credit"
		}
		if body.Type != "credit" && body.Type != "debit" {
			http.Error(w, "type must be credit or debit", http.StatusBadRequest)
			return
		}
		if len(body.Items) == 0 {
			http.Error(w, "items required", http.StatusBadRequest)
			return
		}

		var cn CreditNote
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`INSERT INTO credit_notes (tenant_id, invoice_id, type, reason, status)
				 VALUES ($1,NULLIF($2,'')::uuid,$3,$4,'draft')
				 RETURNING `+cnCols,
				mw.TenantIDFrom(r.Context()), body.InvoiceID, body.Type, body.Reason,
			).Scan(&cn.ID, &cn.InvoiceID, &cn.Number, &cn.Type, &cn.Reason, &cn.Subtotal, &cn.TaxTotal, &cn.Total, &cn.Status, &cn.IssuedAt, &cn.CreatedAt, &cn.UpdatedAt); err != nil {
				return err
			}
			for _, it := range body.Items {
				if it.Qty == "" {
					it.Qty = "1"
				}
				if it.UnitPrice == "" {
					it.UnitPrice = "0"
				}
				if it.TaxPct == "" {
					it.TaxPct = "0"
				}
				var item Item
				if err := tx.QueryRow(r.Context(),
					`INSERT INTO credit_note_items (credit_note_id, product_id, description, qty, unit_price, tax_pct, line_total)
					 VALUES ($1,NULLIF($2,'')::uuid,$3,$4,$5,$6,$4*$5)
					 RETURNING `+itemCols,
					cn.ID, it.ProductID, it.Description, it.Qty, it.UnitPrice, it.TaxPct,
				).Scan(&item.ID, &item.ProductID, &item.Description, &item.Qty, &item.UnitPrice, &item.TaxPct, &item.LineTotal); err != nil {
					return err
				}
				cn.Items = append(cn.Items, item)
			}
			_, err := tx.Exec(r.Context(),
				`UPDATE credit_notes SET
				   subtotal=(SELECT COALESCE(SUM(qty*unit_price),0) FROM credit_note_items WHERE credit_note_id=$1),
				   tax_total=(SELECT COALESCE(SUM(qty*unit_price*tax_pct/100),0) FROM credit_note_items WHERE credit_note_id=$1),
				   total=(SELECT COALESCE(SUM(qty*unit_price + qty*unit_price*tax_pct/100),0) FROM credit_note_items WHERE credit_note_id=$1),
				   updated_at=NOW()
				 WHERE id=$1`, cn.ID)
			return err
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		err = mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`SELECT `+cnCols+` FROM credit_notes WHERE id=$1`, cn.ID,
			).Scan(&cn.ID, &cn.InvoiceID, &cn.Number, &cn.Type, &cn.Reason, &cn.Subtotal, &cn.TaxTotal, &cn.Total, &cn.Status, &cn.IssuedAt, &cn.CreatedAt, &cn.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(cn)
	}
}

func HandleList(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cnType := r.URL.Query().Get("type")
		status := r.URL.Query().Get("status")
		var cns []CreditNote
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+cnCols+` FROM credit_notes WHERE ($1='' OR type=$1) AND ($2='' OR status=$2) ORDER BY created_at DESC`, cnType, status)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var cn CreditNote
				if err := rows.Scan(&cn.ID, &cn.InvoiceID, &cn.Number, &cn.Type, &cn.Reason, &cn.Subtotal, &cn.TaxTotal, &cn.Total, &cn.Status, &cn.IssuedAt, &cn.CreatedAt, &cn.UpdatedAt); err != nil {
					return err
				}
				cns = append(cns, cn)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if cns == nil {
			cns = []CreditNote{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cns)
	}
}

func HandleGet(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var cn CreditNote
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`SELECT `+cnCols+` FROM credit_notes WHERE id=$1`, id,
			).Scan(&cn.ID, &cn.InvoiceID, &cn.Number, &cn.Type, &cn.Reason, &cn.Subtotal, &cn.TaxTotal, &cn.Total, &cn.Status, &cn.IssuedAt, &cn.CreatedAt, &cn.UpdatedAt); err != nil {
				return err
			}
			rows, err := tx.Query(r.Context(),
				`SELECT `+itemCols+` FROM credit_note_items WHERE credit_note_id=$1 ORDER BY created_at`, id)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var it Item
				if err := rows.Scan(&it.ID, &it.ProductID, &it.Description, &it.Qty, &it.UnitPrice, &it.TaxPct, &it.LineTotal); err != nil {
					return err
				}
				cn.Items = append(cn.Items, it)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if cn.Items == nil {
			cn.Items = []Item{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cn)
	}
}

func HandleVoid(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			ct, err := tx.Exec(r.Context(),
				`UPDATE credit_notes SET status='void', updated_at=NOW()
				 WHERE id=$1 AND status IN ('draft','issued')`, id)
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
