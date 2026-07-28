package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Note struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

var pool *pgxpool.Pool

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://demo:demo@localhost:5432/demo?sslmode=disable"
	}

	ctx := context.Background()

	var err error
	pool, err = connectWithRetry(ctx, dsn, 10, 2*time.Second)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}
	defer pool.Close()
	log.Println("подключение к postgres установлено")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/notes", notesHandler)
	mux.HandleFunc("/api/notes/", noteHandler)
	mux.HandleFunc("/api/health", healthHandler)

	addr := ":8080"
	log.Printf("backend слушает %s", addr)
	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func connectWithRetry(ctx context.Context, dsn string, attempts int, delay time.Duration) (*pgxpool.Pool, error) {
	var lastErr error
	for i := 0; i < attempts; i++ {
		p, err := pgxpool.New(ctx, dsn)
		if err == nil {
			if pingErr := p.Ping(ctx); pingErr == nil {
				return p, nil
			} else {
				lastErr = pingErr
				p.Close()
			}
		} else {
			lastErr = err
		}
		log.Printf("БД ещё не готова (попытка %d/%d): %v", i+1, attempts, lastErr)
		time.Sleep(delay)
	}
	return nil, lastErr
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listNotes(w, r)
	case http.MethodPost:
		createNote(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func noteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	deleteNote(w, r, id)
}

func listNotes(w http.ResponseWriter, r *http.Request) {
	rows, err := pool.Query(r.Context(), "SELECT id, text, created_at FROM notes ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	notes := []Note{}
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Text, &n.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notes = append(notes, n)
	}

	writeJSON(w, http.StatusOK, notes)
}

func createNote(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "неверное тело запроса", http.StatusBadRequest)
		return
	}
	input.Text = strings.TrimSpace(input.Text)
	if input.Text == "" {
		http.Error(w, "text не может быть пустым", http.StatusBadRequest)
		return
	}

	var n Note
	err := pool.QueryRow(
		r.Context(),
		"INSERT INTO notes (text) VALUES ($1) RETURNING id, text, created_at",
		input.Text,
	).Scan(&n.ID, &n.Text, &n.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, n)
}

func deleteNote(w http.ResponseWriter, r *http.Request, id int) {
	_, err := pool.Exec(r.Context(), "DELETE FROM notes WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
