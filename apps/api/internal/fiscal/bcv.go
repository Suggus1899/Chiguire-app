package fiscal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FetchBCVRate scrapes the official BCV USD->VES rate from bcv.org.ve.
func FetchBCVRate(ctx context.Context) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.bcv.org.ve/", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ChiguireERP/1.0)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("bcv: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("bcv: read body: %w", err)
	}

	return parseBCVRate(string(body))
}

// parseBCVRate extracts the USD/VES rate from BCV HTML.
// Looks for id="dolar" then the first <strong> value.
// Venezuelan format: "36.500,00" → 36500.0
func parseBCVRate(html string) (float64, error) {
	idx := strings.Index(html, `id="dolar"`)
	if idx < 0 {
		return 0, errors.New("bcv: dolar div not found")
	}
	sub := html[idx:]

	start := strings.Index(sub, "<strong>")
	if start < 0 {
		return 0, errors.New("bcv: <strong> tag not found")
	}
	sub = sub[start+len("<strong>"):]

	end := strings.Index(sub, "</strong>")
	if end < 0 {
		return 0, errors.New("bcv: </strong> closing tag not found")
	}
	raw := strings.TrimSpace(sub[:end])

	// Venezuelan number format: periods are thousand separators, comma is decimal
	raw = strings.ReplaceAll(raw, ".", "")
	raw = strings.ReplaceAll(raw, ",", ".")

	rate, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("bcv: parse %q: %w", raw, err)
	}
	if rate <= 0 {
		return 0, fmt.Errorf("bcv: invalid rate %.4f", rate)
	}
	return rate, nil
}

// StartBCVCron ticks at the given interval, fetches the BCV rate, and inserts
// a row into exchange_rates. Blocks until parent ctx is cancelled.
func StartBCVCron(parent context.Context, pool *pgxpool.Pool, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(parent, 30*time.Second)
			rate, err := FetchBCVRate(ctx)
			if err != nil {
				log.Printf("bcv cron: fetch error: %v", err)
				cancel()
				continue
			}
			_, err = pool.Exec(ctx,
				`INSERT INTO exchange_rates (currency, rate_to_ves, source, effective_at)
				 VALUES ('USD', $1, 'bcv', NOW())`,
				rate,
			)
			if err != nil {
				log.Printf("bcv cron: insert error: %v", err)
			}
			cancel()
		case <-parent.Done():
			return
		}
	}
}
