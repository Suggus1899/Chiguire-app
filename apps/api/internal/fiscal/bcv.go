package fiscal

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FetchBCVRate scrapes the official BCV USD->VES rate.
// TODO: implement real scraping of bcv.org.ve in production.
func FetchBCVRate(ctx context.Context) (float64, error) {
	// placeholder: real implementation will HTTP GET bcv.org.ve and parse the rate
	return 0, nil
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
			if rate <= 0 {
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
