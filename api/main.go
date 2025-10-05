package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/itsstyg/ow2-map-poll/db"
	_ "github.com/mattn/go-sqlite3"
)

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

// TODO: 
func me(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	returnError := func(err error) {
		w.WriteHeader(500)
		fmt.Fprintf(os.Stderr, "err: %s\n", err)
		fmt.Fprintf(w, `{"status": "error", "message": "Internal error."}`)
	}

	client_hash, err := getClientHash(r)
	if err != nil {
		err = fmt.Errorf("getClientHash: %v", err)
		returnError(err)
		return
	}

	user, user_err := db.GetOrCreateUser(client_hash)
	if user_err != nil {
		user_err = fmt.Errorf("GetOrCreateUser: %v", user_err)
		returnError(user_err)
		return
	}

	poll, poll_err := db.GetOrCreateUserPoll(&user)
	if poll_err != nil {
		poll_err = fmt.Errorf("GetOrCreateUserPoll: %v", poll_err)
		returnError(poll_err)
		return
	}

	user_json, user_json_err := json.Marshal(user)
	if user_json_err != nil {
		returnError(user_json_err)
		return
	}

	poll_json, poll_json_err := json.Marshal(poll)
	if poll_json_err != nil {
		returnError(poll_json_err)
		return
	}

	fmt.Fprintf(w, `{"user": %s, "poll": %s}`, string(user_json), string(poll_json))
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
	http.HandleFunc("/", root)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
