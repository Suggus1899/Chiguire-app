package product

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

// ---------- Products ----------

type Product struct {
	ID          string    `json:"id"`
	Sku         string    `json:"sku"`
	Barcode     string    `json:"barcode"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CategoryID  string    `json:"category_id"`
	UnitID      string    `json:"unit_id"`
	TaxCategory string    `json:"tax_category"`
	IsActive    bool      `json:"is_active"`
	MinStock    string    `json:"min_stock"`
	AvgCostUSD  string    `json:"avg_cost_usd"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const productCols = `id, sku, COALESCE(barcode,''), name, COALESCE(description,''), COALESCE(category_id::text,''), COALESCE(unit_id::text,''), tax_category, is_active, min_stock::text, avg_cost_usd::text, created_at, updated_at`

func HandleCreateProduct(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Sku         string `json:"sku"`
			Barcode     string `json:"barcode"`
			Name        string `json:"name"`
			Description string `json:"description"`
			CategoryID  string `json:"category_id"`
			UnitID      string `json:"unit_id"`
			TaxCategory string `json:"tax_category"`
			IsActive    *bool  `json:"is_active"`
			MinStock    string `json:"min_stock"`
			AvgCostUSD  string `json:"avg_cost_usd"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Sku == "" || body.Name == "" {
			http.Error(w, "sku and name are required", http.StatusBadRequest)
			return
		}
		if body.TaxCategory == "" {
			body.TaxCategory = "general"
		}
		if body.MinStock == "" {
			body.MinStock = "0"
		}
		if body.AvgCostUSD == "" {
			body.AvgCostUSD = "0"
		}
		active := true
		if body.IsActive != nil {
			active = *body.IsActive
		}

		var p Product
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO products (tenant_id, sku, barcode, name, description, category_id, unit_id, tax_category, is_active, min_stock, avg_cost_usd)
				 VALUES ($1,$2,$3,$4,$5,NULLIF($6,'')::uuid,NULLIF($7,'')::uuid,$8,$9,$10,$11)
				 RETURNING `+productCols,
				mw.TenantIDFrom(r.Context()), body.Sku, body.Barcode, body.Name, body.Description, body.CategoryID, body.UnitID, body.TaxCategory, active, body.MinStock, body.AvgCostUSD,
			).Scan(&p.ID, &p.Sku, &p.Barcode, &p.Name, &p.Description, &p.CategoryID, &p.UnitID, &p.TaxCategory, &p.IsActive, &p.MinStock, &p.AvgCostUSD, &p.CreatedAt, &p.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "sku already exists", http.StatusConflict)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}
}

func HandleListProducts(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var products []Product
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+productCols+` FROM products ORDER BY created_at DESC`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var p Product
				if err := rows.Scan(&p.ID, &p.Sku, &p.Barcode, &p.Name, &p.Description, &p.CategoryID, &p.UnitID, &p.TaxCategory, &p.IsActive, &p.MinStock, &p.AvgCostUSD, &p.CreatedAt, &p.UpdatedAt); err != nil {
					return err
				}
				products = append(products, p)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if products == nil {
			products = []Product{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func HandleGetProduct(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var p Product
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`SELECT `+productCols+` FROM products WHERE id=$1`, id,
			).Scan(&p.ID, &p.Sku, &p.Barcode, &p.Name, &p.Description, &p.CategoryID, &p.UnitID, &p.TaxCategory, &p.IsActive, &p.MinStock, &p.AvgCostUSD, &p.CreatedAt, &p.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

func HandleUpdateProduct(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			Sku         string `json:"sku"`
			Barcode     string `json:"barcode"`
			Name        string `json:"name"`
			Description string `json:"description"`
			CategoryID  string `json:"category_id"`
			UnitID      string `json:"unit_id"`
			TaxCategory string `json:"tax_category"`
			IsActive    *bool  `json:"is_active"`
			MinStock    string `json:"min_stock"`
			AvgCostUSD  string `json:"avg_cost_usd"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.MinStock == "" {
			body.MinStock = "0"
		}
		if body.AvgCostUSD == "" {
			body.AvgCostUSD = "0"
		}
		active := true
		if body.IsActive != nil {
			active = *body.IsActive
		}

		var p Product
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`UPDATE products SET sku=$1, barcode=$2, name=$3, description=$4, category_id=NULLIF($5,'')::uuid, unit_id=NULLIF($6,'')::uuid, tax_category=$7, is_active=$8, min_stock=$9, avg_cost_usd=$10, updated_at=NOW()
				 WHERE id=$11
				 RETURNING `+productCols,
				body.Sku, body.Barcode, body.Name, body.Description, body.CategoryID, body.UnitID, body.TaxCategory, active, body.MinStock, body.AvgCostUSD, id,
			).Scan(&p.ID, &p.Sku, &p.Barcode, &p.Name, &p.Description, &p.CategoryID, &p.UnitID, &p.TaxCategory, &p.IsActive, &p.MinStock, &p.AvgCostUSD, &p.CreatedAt, &p.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "not found or conflict", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

func HandleDeleteProduct(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			ct, err := tx.Exec(r.Context(), `DELETE FROM products WHERE id=$1`, id)
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

// ---------- Categories ----------

type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ParentID  string    `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

const categoryCols = `id, name, COALESCE(parent_id::text,''), created_at`

func HandleCreateCategory(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name     string `json:"name"`
			ParentID string `json:"parent_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		var c Category
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO product_categories (tenant_id, name, parent_id)
				 VALUES ($1,$2,NULLIF($3,'')::uuid)
				 RETURNING `+categoryCols,
				mw.TenantIDFrom(r.Context()), body.Name, body.ParentID,
			).Scan(&c.ID, &c.Name, &c.ParentID, &c.CreatedAt)
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(c)
	}
}

func HandleListCategories(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cats []Category
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+categoryCols+` FROM product_categories ORDER BY name`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var c Category
				if err := rows.Scan(&c.ID, &c.Name, &c.ParentID, &c.CreatedAt); err != nil {
					return err
				}
				cats = append(cats, c)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if cats == nil {
			cats = []Category{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cats)
	}
}

func HandleDeleteCategory(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			ct, err := tx.Exec(r.Context(), `DELETE FROM product_categories WHERE id=$1`, id)
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

// ---------- Units of measure ----------

type Unit struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Abbreviation string    `json:"abbreviation"`
	CreatedAt    time.Time `json:"created_at"`
}

func HandleCreateUnit(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name         string `json:"name"`
			Abbreviation string `json:"abbreviation"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.Abbreviation == "" {
			http.Error(w, "name and abbreviation are required", http.StatusBadRequest)
			return
		}

		var u Unit
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO units_of_measure (tenant_id, name, abbreviation)
				 VALUES ($1,$2,$3)
				 RETURNING id, name, abbreviation, created_at`,
				mw.TenantIDFrom(r.Context()), body.Name, body.Abbreviation,
			).Scan(&u.ID, &u.Name, &u.Abbreviation, &u.CreatedAt)
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	}
}

func HandleListUnits(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var units []Unit
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT id, name, abbreviation, created_at FROM units_of_measure ORDER BY name`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var u Unit
				if err := rows.Scan(&u.ID, &u.Name, &u.Abbreviation, &u.CreatedAt); err != nil {
					return err
				}
				units = append(units, u)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if units == nil {
			units = []Unit{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(units)
	}
}

func HandleDeleteUnit(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			ct, err := tx.Exec(r.Context(), `DELETE FROM units_of_measure WHERE id=$1`, id)
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

// ---------- Product prices ----------

type Price struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Currency  string    `json:"currency"`
	Amount    string    `json:"amount"`
	ValidFrom time.Time `json:"valid_from"`
	CreatedAt time.Time `json:"created_at"`
}

func HandleCreatePrice(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			ProductID string `json:"product_id"`
			Currency  string `json:"currency"`
			Amount    string `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.ProductID == "" || body.Currency == "" || body.Amount == "" {
			http.Error(w, "product_id, currency and amount are required", http.StatusBadRequest)
			return
		}

		var p Price
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO product_prices (tenant_id, product_id, currency, amount)
				 VALUES ($1,$2,$3,$4)
				 RETURNING id, product_id::text, currency, amount::text, valid_from, created_at`,
				mw.TenantIDFrom(r.Context()), body.ProductID, body.Currency, body.Amount,
			).Scan(&p.ID, &p.ProductID, &p.Currency, &p.Amount, &p.ValidFrom, &p.CreatedAt)
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}
}

func HandleListPrices(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID := chi.URLParam(r, "id")
		var prices []Price
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT id, product_id::text, currency, amount::text, valid_from, created_at
				 FROM product_prices WHERE product_id=$1 ORDER BY valid_from DESC`, productID)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var p Price
				if err := rows.Scan(&p.ID, &p.ProductID, &p.Currency, &p.Amount, &p.ValidFrom, &p.CreatedAt); err != nil {
					return err
				}
				prices = append(prices, p)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if prices == nil {
			prices = []Price{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prices)
	}
}
