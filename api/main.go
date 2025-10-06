package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/stygian91/ow2-map-poll/db"
	_ "github.com/mattn/go-sqlite3"
)

const MAX_POLLS_PER_DAY = 20

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there")
}

func getClientHash(r *http.Request) (hash string, err error) {
	host, _, split_err := net.SplitHostPort(r.RemoteAddr)

	if split_err != nil {
		err = split_err
		return
	}

	hash = fmt.Sprintf("%x", sha256.Sum256([]byte(host)))
	return
}

func returnError(w http.ResponseWriter, err error) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "err: %s\n", err)
		fmt.Fprintf(w, `{"status": "error", "message": "Internal error."}`)
}

// TODO: check if user should be throttled
func me(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	client_hash, err := getClientHash(r)
	if err != nil {
		returnError(w, fmt.Errorf("getClientHash: %w", err))
		return
	}

	user, user_err := db.GetOrCreateUser(client_hash)
	if user_err != nil {
		returnError(w, fmt.Errorf("GetOrCreateUser: %w", user_err))
		return
	}

	poll_count, cnt_err := db.GetPollCountForToday(&user)
	if cnt_err != nil {
		returnError(w, fmt.Errorf("GetPollCountForToday: %w", cnt_err))
		return
	}

	if poll_count >= MAX_POLLS_PER_DAY {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, `{"status": "error", "message": "Too many requests"}`)
		return
	}

	poll, poll_err := db.GetOrCreateUserPoll(&user)
	if poll_err != nil {
		returnError(w, fmt.Errorf("GetOrCreateUserPoll: %w", poll_err))
		return
	}

	poll_json, poll_json_err := json.Marshal(poll)
	if poll_json_err != nil {
		returnError(w, poll_json_err)
		return
	}

	fmt.Fprintf(w, `{"poll": %s}`, string(poll_json))
}

func vote(w http.ResponseWriter, r *http.Request) {
	client_hash, err := getClientHash(r)
	if err != nil {
		returnError(w, fmt.Errorf("getClientHash: %w", err))
		return
	}

	user, user_err := db.GetOrCreateUser(client_hash)
	if user_err != nil {
		returnError(w, fmt.Errorf("GetOrCreateUser: %w", user_err))
		return
	}

	poll_count, cnt_err := db.GetPollCountForToday(&user)
	if cnt_err != nil {
		returnError(w, fmt.Errorf("GetPollCountForToday: %w", cnt_err))
		return
	}

	if poll_count >= MAX_POLLS_PER_DAY {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, `{"status": "error", "message": "Too many requests"}`)
		return
	}

	// TODO:
	// 1. Get current user poll and validate with req.poll_id
	// 2. Validate that the req.vote is valid for the current poll
	// 3. Update the poll vote in DB

	fmt.Fprintf(w, `{"success": true}`)
}

func main() {
	db_err := db.Open()
	if db_err != nil {
		fmt.Fprintf(os.Stderr, "db_err: %s\n", db_err)
		os.Exit(1)
	}

	defer db.Close()

	fmt.Println("Listening on 127.0.0.1:8080")
	http.HandleFunc("/api/me", me)
	http.HandleFunc("/api/vote", vote)
	http.HandleFunc("/", root)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
