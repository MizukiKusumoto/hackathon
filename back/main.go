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

type MessageForGet struct {
	Id         string    `json:"id"`
	Content    string    `json:"content"`
	ChannelId  string    `json:"channel_id"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}
type MessageForPost struct {
	Content   string `json:"content"`
	ChannelId string `json:"channel_id"`
}
type MessageForPut struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

type Channel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB

func init() {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PWD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	_db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase))
	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}
	if err := _db.Ping(); err != nil {
		log.Fatalf("fail: _db.Ping, %v\n", err)
	}
	db = _db
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	switch r.Method {
	case http.MethodOptions:
		w.Header()

	case http.MethodPost:
		var message MessageForPost
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

		tx, err := db.Begin()
		if err != nil {
			log.Printf("fail: db.Begin, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t := time.Now()
		createdAt := t
		modifiedAt := t
		entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
		id := ulid.MustNew(ulid.Timestamp(t), entropy).String()

		_, err = tx.Exec("INSERT INTO message (id, content, channel_id, created_at, modified_at) VALUES (?, ?, ?, ?, ?)",
			id, message.Content, message.ChannelId, createdAt, modifiedAt)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				log.Printf("fail: tx.Rollback, %v\n", err2)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Printf("fail: tx.Exec, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := tx.Commit(); err != nil {
			log.Printf("fail: tx.Commit, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	case http.MethodGet:
		channelID := r.URL.Query().Get("channel_id")
		if channelID == "" {
			log.Println("fail: channel_id is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rows, err := db.Query("SELECT * FROM message WHERE channel_id = ?", channelID)
		if err != nil {
			log.Printf("fail: db.Query, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		messages := make([]MessageForGet, 0)
		for rows.Next() {
			var message MessageForGet
			if err := rows.Scan(&message.Id, &message.Content, &message.ChannelId, &message.CreatedAt, &message.ModifiedAt); err != nil {
				log.Printf("fail: rows.Scan, %v\n", err)
				if err := rows.Close(); err != nil { // 500を返して終了するが、その前にrowsのClose処理が必要
					log.Printf("fail: rows.Close(), %v\n", err)
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			messages = append(messages, message)
		}

		bytes, err := json.Marshal(messages)
		if err != nil {
			log.Printf("fail: json.Marshal, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := rows.Close(); err != nil {
			log.Printf("fail: rows.Close(), %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(bytes)

	case http.MethodPut:
		var message MessageForPut
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			log.Printf("fail: json.NewDecoder, %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if message.Id == "" {
			log.Println("fail: id is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if message.Content == "" {
			log.Println("fail: content is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ModifiedAt := time.Now()

		_, err = db.Exec("UPDATE message SET content = ?, modified_at = ? WHERE id = ?",
			message.Content, ModifiedAt, message.Id)
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

func channelHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	switch r.Method {
	case http.MethodOptions:
		w.Header()

	case http.MethodGet:
		rows, err := db.Query("SELECT * FROM channel ")
		if err != nil {
			log.Printf("fail: db.Query, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		channels := make([]Channel, 0)
		for rows.Next() {
			var channel Channel
			if err := rows.Scan(&channel.Id, &channel.Name); err != nil {
				log.Printf("fail: rows.Scan, %v\n", err)
				if err := rows.Close(); err != nil {
					log.Printf("fail: rows.Close(), %v\n", err)
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			channels = append(channels, channel)
		}

		bytes, err := json.Marshal(channels)
		if err != nil {
			log.Printf("fail: json.Marshal, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := rows.Close(); err != nil {
			log.Printf("fail: rows.Close(), %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(bytes)

	default:
		log.Printf("fail: HTTP Method is %s\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	http.HandleFunc("/message", messageHandler)
	http.HandleFunc("/channel", channelHandler)
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
