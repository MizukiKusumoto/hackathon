package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/oklog/ulid"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Message struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`
	ChannelID  string    `json:"channel_id"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

var db *sql.DB

func init() {
	// DB接続のための準備
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PWD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	_db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}
	if err := _db.Ping(); err != nil {
		log.Fatalf("fail: _db.Ping, %v\n", err)
	}
	db = _db
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		log.Printf("fail: json.NewDecoder, %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if message.Content == "" {
		log.Println("fail: content is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t := time.Now()
	message.ID = generateID()
	message.CreatedAt = t
	message.ModifiedAt = t

	_, err = db.Exec("INSERT INTO message (id, content, channel_id, created_at, modified_at) VALUES (?, ?, ?, ?, ?)",
		message.ID, message.Content, message.ChannelID, message.CreatedAt, message.ModifiedAt)
	if err != nil {
		log.Printf("fail: db.Exec, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	channelID := r.URL.Query().Get("channel_id")
	if channelID == "" {
		log.Println("fail: channel_id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT id, content, channel_id, created_at, modified_at FROM message WHERE channel_id = ?", channelID)
	if err != nil {
		log.Printf("fail: db.Query, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	messages := make([]Message, 0)
	for rows.Next() {
		var message Message
		var createdAt, modifiedAt string
		if err := rows.Scan(&message.ID, &message.Content, &message.ChannelID, &createdAt, &modifiedAt); err != nil {
			log.Printf("fail: rows.Scan, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		message.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		message.ModifiedAt, _ = time.Parse("2006-01-02 15:04:05", modifiedAt)

		messages = append(messages, message)
	}

	bytes, err := json.Marshal(messages)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func editMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	switch r.Method {
	case http.MethodOptions:
		w.Header()

	case http.MethodPut:
		var message Message
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			log.Printf("fail: json.NewDecoder, %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if message.ID == "" {
			log.Println("fail: id is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if message.Content == "" {
			log.Println("fail: content is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		message.ModifiedAt = time.Now()

		_, err = db.Exec("UPDATE message SET content = ?, modified_at = ? WHERE id = ?",
			message.Content, message.ModifiedAt, message.ID)
		if err != nil {
			log.Printf("fail: db.Exec, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		log.Printf("fail: HTTP Method is %s\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func deleteMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("fail: id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM message WHERE id = ?", id)
	if err != nil {
		log.Printf("fail: db.Exec, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func generateID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy).String()
	return id
}

func main() {
	http.HandleFunc("/message", postMessage)
	http.HandleFunc("/messages", getMessages)
	http.HandleFunc("/message/edit", editMessage)
	http.HandleFunc("/message/delete", deleteMessage)
	closeDBWithSysCall()
	log.Println("Listening...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func closeDBWithSysCall() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("received syscall, %v", s)

		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("success: db.Close()")
		os.Exit(0)
	}()
}
