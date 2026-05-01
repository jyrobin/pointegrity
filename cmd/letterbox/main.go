// letterbox is a tiny self-hosted newsletter-subscriber service for the
// Pointegrity workshop. It does one thing: captures email addresses
// posted from a plain HTML form and stores them in a SQLite file.
//
// What it deliberately does NOT do (yet):
//   - Double opt-in / confirmation emails. Add when abuse warrants it.
//   - Send the actual newsletter. The list-dump endpoint is enough for
//     now — pipe it into your mail client of choice.
//   - Track opens, clicks, anything behavioral. The whole point of the
//     Pointegrity stance is "no surveillance" — that applies internally
//     too.
//   - Auth beyond a shared-secret admin key for the list-dump endpoint.
//
// Endpoints:
//
//   POST /subscribe         (form-encoded: email, name?, source?)
//                           Stores row, 303-redirects to /journal/subscribed/.
//   GET  /unsubscribe?t=... 303 to /journal/unsubscribed/ after deleting.
//   GET  /list?key=...      Admin CSV dump (compares LETTERBOX_ADMIN_KEY).
//   GET  /healthz           "ok"
//
// Config via env vars:
//
//   LETTERBOX_ADDR       :3737 (default)
//   LETTERBOX_DB         ./data/letterbox.db
//   LETTERBOX_REDIRECT   https://www.pointegrity.com (no trailing slash)
//   LETTERBOX_ADMIN_KEY  shared-secret for the /list dump
//
// Storage: a single SQLite file. Easy to back up (cp), easy to migrate
// elsewhere (the binary is portable; the file is portable).
package main

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	addr := envOr("LETTERBOX_ADDR", ":3737")
	dbPath := envOr("LETTERBOX_DB", "./data/letterbox.db")
	redirectBase := strings.TrimRight(envOr("LETTERBOX_REDIRECT", "https://www.pointegrity.com"), "/")
	adminKey := os.Getenv("LETTERBOX_ADMIN_KEY")

	if adminKey == "" {
		log.Fatal("LETTERBOX_ADMIN_KEY must be set (shared secret for /list)")
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		log.Fatalf("open db %s: %v", dbPath, err)
	}
	defer db.Close()

	if err := initSchema(db); err != nil {
		log.Fatalf("schema: %v", err)
	}

	srv := &server{db: db, redirectBase: redirectBase, adminKey: adminKey}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /subscribe", srv.handleSubscribe)
	mux.HandleFunc("GET /unsubscribe", srv.handleUnsubscribe)
	mux.HandleFunc("GET /list", srv.handleList)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok\n"))
	})

	log.Printf("letterbox: listening on %s (db=%s, redirect=%s)", addr, dbPath, redirectBase)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// =============================================================================
// Storage
// =============================================================================

func initSchema(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS subscribers (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    email      TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL DEFAULT '',
    source     TEXT NOT NULL DEFAULT '',
    token      TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_subscribers_email ON subscribers(email);
CREATE INDEX IF NOT EXISTS idx_subscribers_token ON subscribers(token);
`)
	return err
}

// =============================================================================
// Handlers
// =============================================================================

type server struct {
	db           *sql.DB
	redirectBase string
	adminKey     string
}

// handleSubscribe accepts a form POST and stores the email address.
//
// Idempotency: if the email already exists we still redirect to the
// "subscribed" page — no error to the user, no second row inserted.
// This avoids leaking "this email is/isn't on our list" to anyone
// probing the endpoint.
func (s *server) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.redirect(w, r, "subscribe-error/")
		return
	}

	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	name := strings.TrimSpace(r.FormValue("name"))
	source := strings.TrimSpace(r.FormValue("source"))
	if source == "" {
		source = "journal"
	}

	if _, err := mail.ParseAddress(email); err != nil {
		log.Printf("[letterbox] invalid email rejected: %q", email)
		s.redirect(w, r, "subscribe-error/")
		return
	}

	token, err := newToken()
	if err != nil {
		log.Printf("[letterbox] token gen: %v", err)
		s.redirect(w, r, "subscribe-error/")
		return
	}

	// INSERT OR IGNORE: existing email returns rowsAffected=0, no error.
	res, err := s.db.ExecContext(r.Context(),
		`INSERT OR IGNORE INTO subscribers(email, name, source, token) VALUES (?,?,?,?)`,
		email, name, source, token)
	if err != nil {
		log.Printf("[letterbox] insert: %v", err)
		s.redirect(w, r, "subscribe-error/")
		return
	}
	if rows, _ := res.RowsAffected(); rows == 1 {
		log.Printf("[letterbox] subscribed: %s (source=%s)", email, source)
	} else {
		log.Printf("[letterbox] already subscribed: %s", email)
	}

	s.redirect(w, r, "subscribed/")
}

// handleUnsubscribe deletes a row by token. Tokens are unguessable
// random 32-byte values, so a leaked token is the only way to
// unsubscribe someone — which is exactly what every newsletter footer
// link does.
func (s *server) handleUnsubscribe(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("t")
	if token == "" {
		s.redirect(w, r, "subscribe-error/")
		return
	}

	res, err := s.db.ExecContext(r.Context(), `DELETE FROM subscribers WHERE token = ?`, token)
	if err != nil {
		log.Printf("[letterbox] unsubscribe delete: %v", err)
		s.redirect(w, r, "subscribe-error/")
		return
	}
	if rows, _ := res.RowsAffected(); rows == 1 {
		log.Printf("[letterbox] unsubscribed: token=%s...", token[:8])
	} else {
		log.Printf("[letterbox] unsubscribe: token not found")
	}

	s.redirect(w, r, "unsubscribed/")
}

// handleList streams a CSV dump of subscribers. Guarded by a
// constant-time compare of LETTERBOX_ADMIN_KEY — not auth-grade
// (no rate limit, no audit log) but enough to keep casual probing out.
// Run with `curl https://.../list?key=$LETTERBOX_ADMIN_KEY > subs.csv`.
func (s *server) handleList(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if subtle.ConstantTimeCompare([]byte(key), []byte(s.adminKey)) != 1 {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	rows, err := s.db.QueryContext(r.Context(),
		`SELECT email, name, source, token, created_at FROM subscribers ORDER BY created_at`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename="subscribers-%s.csv"`, time.Now().Format("2006-01-02")))

	cw := csv.NewWriter(w)
	cw.Write([]string{"email", "name", "source", "token", "created_at"})
	for rows.Next() {
		var email, name, source, token, createdAt string
		if err := rows.Scan(&email, &name, &source, &token, &createdAt); err != nil {
			log.Printf("[letterbox] list scan: %v", err)
			return
		}
		cw.Write([]string{email, name, source, token, createdAt})
	}
	cw.Flush()
}

// redirect 303s back to a static journal status page. All outcomes
// (subscribed, unsubscribed, error) bounce to a real HTML page on
// pointegrity.com so the user never sees a bare API response.
func (s *server) redirect(w http.ResponseWriter, r *http.Request, path string) {
	dest := s.redirectBase + "/journal/" + path
	http.Redirect(w, r, dest, http.StatusSeeOther)
}

// =============================================================================
// Helpers
// =============================================================================

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func newToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("rand: " + err.Error())
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}
