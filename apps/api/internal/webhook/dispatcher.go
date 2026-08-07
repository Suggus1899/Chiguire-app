package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const maxAttempts = 5

func signPayload(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func deliver(url, secret, eventType string, payload []byte) (string, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Event-Type", eventType)
	if secret != "" {
		req.Header.Set("X-Signature-256", "sha256="+signPayload(secret, payload))
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return resp.Status, nil
	}
	return string(body), nil
}

func DispatchEvent(pool *pgxpool.Pool, tenantID, eventType string, payload any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, "SET LOCAL app.tenant_id = $1", tenantID); err != nil {
		return err
	}
	rows, err := tx.Query(ctx,
		`SELECT id, url, secret FROM outgoing_webhooks WHERE is_active=true AND ($1 = ANY(string_to_array(events, ',')) OR events='*')`,
		eventType)
	if err != nil {
		return err
	}
	type wh struct {
		id, url, secret string
	}
	var matches []wh
	for rows.Next() {
		var w wh
		if err := rows.Scan(&w.id, &w.url, &w.secret); err != nil {
			rows.Close()
			return err
		}
		matches = append(matches, w)
	}
	rows.Close()
	for _, m := range matches {
		if _, err := tx.Exec(ctx,
			`INSERT INTO webhook_deliveries (tenant_id, webhook_id, event_type, payload, status, attempts)
			 VALUES ($1,$2,$3,$4,'pending',0)`,
			tenantID, m.id, eventType, payloadBytes); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func StartWebhookWorker(ctx context.Context, pool *pgxpool.Pool) {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				processPending(pool)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func processPending(pool *pgxpool.Pool) {
	ctx := context.Background()
	rows, err := pool.Query(ctx,
		`SELECT wd.id, wd.tenant_id, wd.webhook_id, wd.event_type, wd.payload, wd.attempts, ow.url, ow.secret
		 FROM webhook_deliveries wd
		 JOIN outgoing_webhooks ow ON ow.id = wd.webhook_id
		 WHERE wd.status='pending' AND (wd.next_retry_at IS NULL OR wd.next_retry_at <= NOW())
		 LIMIT 20`)
	if err != nil {
		return
	}
	type pending struct {
		id, tenantID, webhookID, eventType string
		payload                           []byte
		attempts                          int
		url, secret                       string
	}
	var items []pending
	for rows.Next() {
		var p pending
		if err := rows.Scan(&p.id, &p.tenantID, &p.webhookID, &p.eventType, &p.payload,
			&p.attempts, &p.url, &p.secret); err != nil {
			rows.Close()
			return
		}
		items = append(items, p)
	}
	rows.Close()

	for _, p := range items {
		resp, err := deliver(p.url, p.secret, p.eventType, p.payload)
		status := "delivered"
		if err != nil {
			status = "failed"
		}
		nextAttempts := p.attempts + 1
		tx, txErr := pool.Begin(ctx)
		if txErr != nil {
			continue
		}
		var nextRetry *time.Time
		if status == "failed" && nextAttempts < maxAttempts {
			backoff := time.Duration(1<<uint(nextAttempts)) * time.Minute
			nr := time.Now().Add(backoff)
			nextRetry = &nr
			status = "pending"
		}
		var deliveredAt *time.Time
		if status == "delivered" {
			now := time.Now()
			deliveredAt = &now
		}
		respPtr := &resp
		if err != nil {
			respPtr = nil
		}
		_, _ = tx.Exec(ctx,
			`UPDATE webhook_deliveries SET status=$2, attempts=$3, last_response=$4, next_retry_at=$5, delivered_at=$6 WHERE id=$1`,
			p.id, status, nextAttempts, respPtr, nextRetry, deliveredAt)
		if err := tx.Commit(ctx); err != nil {
			tx.Rollback(ctx)
		}
	}
}
