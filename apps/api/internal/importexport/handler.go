package importexport

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

type customerRow struct {
	Name        string `json:"name"`
	TaxID       string `json:"tax_id"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Currency    string `json:"currency"`
}

type productRow struct {
	Name        string  `json:"name"`
	SKU         string  `json:"sku"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Cost        float64 `json:"cost"`
	CategoryID  string  `json:"category_id"`
}

func HandleImportCustomers(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		var rows []customerRow
		if strings.HasPrefix(ct, "application/json") {
			if err := json.NewDecoder(r.Body).Decode(&rows); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
		} else {
			reader := csv.NewReader(r.Body)
			reader.FieldsPerRecord = -1
			header, err := reader.Read()
			if err != nil {
				http.Error(w, "invalid csv", http.StatusBadRequest)
				return
			}
			idx := map[string]int{}
			for i, h := range header {
				idx[strings.ToLower(strings.TrimSpace(h))] = i
			}
			for {
				rec, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					http.Error(w, "csv parse error", http.StatusBadRequest)
					return
				}
				var c customerRow
				if i, ok := idx["name"]; ok && i < len(rec) {
					c.Name = rec[i]
				}
				if i, ok := idx["tax_id"]; ok && i < len(rec) {
					c.TaxID = rec[i]
				}
				if i, ok := idx["email"]; ok && i < len(rec) {
					c.Email = rec[i]
				}
				if i, ok := idx["phone"]; ok && i < len(rec) {
					c.Phone = rec[i]
				}
				if i, ok := idx["address"]; ok && i < len(rec) {
					c.Address = rec[i]
				}
				if i, ok := idx["currency"]; ok && i < len(rec) {
					c.Currency = rec[i]
				}
				rows = append(rows, c)
			}
		}

		var inserted int
		tid := mw.TenantIDFrom(r.Context())
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			for _, c := range rows {
				if c.Name == "" {
					continue
				}
				if _, err := tx.Exec(r.Context(),
					`INSERT INTO customers (tenant_id, name, tax_id, email, phone, address, currency)
					 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
					tid, c.Name, c.TaxID, c.Email, c.Phone, c.Address, c.Currency,
				); err != nil {
					return err
				}
				inserted++
			}
			return nil
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"imported": inserted})
	}
}

func HandleImportProducts(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		var rows []productRow
		if strings.HasPrefix(ct, "application/json") {
			if err := json.NewDecoder(r.Body).Decode(&rows); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
		} else {
			reader := csv.NewReader(r.Body)
			reader.FieldsPerRecord = -1
			header, err := reader.Read()
			if err != nil {
				http.Error(w, "invalid csv", http.StatusBadRequest)
				return
			}
			idx := map[string]int{}
			for i, h := range header {
				idx[strings.ToLower(strings.TrimSpace(h))] = i
			}
			for {
				rec, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					http.Error(w, "csv parse error", http.StatusBadRequest)
					return
				}
				var p productRow
				if i, ok := idx["name"]; ok && i < len(rec) {
					p.Name = rec[i]
				}
				if i, ok := idx["sku"]; ok && i < len(rec) {
					p.SKU = rec[i]
				}
				if i, ok := idx["description"]; ok && i < len(rec) {
					p.Description = rec[i]
				}
				if i, ok := idx["price"]; ok && i < len(rec) {
					p.Price, _ = strconv.ParseFloat(rec[i], 64)
				}
				if i, ok := idx["cost"]; ok && i < len(rec) {
					p.Cost, _ = strconv.ParseFloat(rec[i], 64)
				}
				if i, ok := idx["category_id"]; ok && i < len(rec) {
					p.CategoryID = rec[i]
				}
				rows = append(rows, p)
			}
		}

		var inserted int
		tid := mw.TenantIDFrom(r.Context())
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			for _, p := range rows {
				if p.Name == "" {
					continue
				}
				if _, err := tx.Exec(r.Context(),
					`INSERT INTO products (tenant_id, name, sku, description, price, cost, category_id)
					 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
					tid, p.Name, p.SKU, p.Description, p.Price, p.Cost, p.CategoryID,
				); err != nil {
					return err
				}
				inserted++
			}
			return nil
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"imported": inserted})
	}
}

func HandleExportInvoices(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=invoices.csv")
		cw := csv.NewWriter(w)
		cw.Write([]string{"id", "number", "customer_id", "status", "currency", "subtotal", "tax_total", "total", "created_at"})
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT id, number, customer_id, status, currency, subtotal, tax_total, total, created_at
				 FROM invoices ORDER BY created_at DESC LIMIT 1000`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var id, number, customerID, status, currency, createdAt string
				var subtotal, taxTotal, total float64
				if err := rows.Scan(&id, &number, &customerID, &status, &currency, &subtotal, &taxTotal, &total, &createdAt); err != nil {
					return err
				}
				cw.Write([]string{
					id, number, customerID, status, currency,
					strconv.FormatFloat(subtotal, 'f', 2, 64),
					strconv.FormatFloat(taxTotal, 'f', 2, 64),
					strconv.FormatFloat(total, 'f', 2, 64),
					createdAt,
				})
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		cw.Flush()
	}
}

func HandleExportCustomers(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=customers.csv")
		cw := csv.NewWriter(w)
		cw.Write([]string{"id", "name", "tax_id", "email", "phone", "address", "currency", "created_at"})
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT id, name, tax_id, email, phone, address, currency, created_at
				 FROM customers ORDER BY created_at DESC LIMIT 1000`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var id, name, taxID, email, phone, address, currency, createdAt string
				if err := rows.Scan(&id, &name, &taxID, &email, &phone, &address, &currency, &createdAt); err != nil {
					return err
				}
				cw.Write([]string{id, name, taxID, email, phone, address, currency, createdAt})
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		cw.Flush()
	}
}
