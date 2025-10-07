package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stygian91/ow2-map-poll/db"
)

const MAX_POLLS_PER_DAY = 20

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there")
}

func clientIPSimple(r *http.Request) (string, error) {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		for _, part := range strings.Split(xff, ",") {
			ip := strings.TrimSpace(part)
			if ip != "" {
				return stripPort(ip)
			}
		}
	}

	return stripPort(r.RemoteAddr)
}

func stripPort(ip string) (string, error) {
	if ! strings.Contains(ip, ":") {
		return ip, nil
	}

	host, _, split_err := net.SplitHostPort(ip)
	return host, split_err
}

func getClientHash(r *http.Request) (hash string, err error) {
	// host, _, split_err := net.SplitHostPort(r.RemoteAddr)
	host, split_err := clientIPSimple(r)

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
	r.Body = http.MaxBytesReader(w, r.Body, 1024*100)
	body, read_err := io.ReadAll(r.Body)
	if read_err != nil {
		returnError(w, fmt.Errorf("read body err: %w", read_err))
		return
	}

	query, parse_query_err := url.ParseQuery(string(body))
	if parse_query_err != nil {
		returnError(w, fmt.Errorf("parse query err: %w", parse_query_err))
		return
	}

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

	poll, poll_err := db.GetPoll(&user)
	if poll_err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"status": "error", "message": "Poll not found."}`)
		return
	}

	req_poll_id := getQueryInt(query, "poll_id")
	if req_poll_id != poll.Id {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "error", "message": "Invalid poll id."}`)
		return
	}

	req_vote_id := getQueryInt(query, "vote")
	if req_vote_id != poll.Map1Id && req_vote_id != poll.Map2Id && req_vote_id != poll.Map3Id {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "error", "message": "Invalid vote id."}`)
		return
	}

	if update_err := db.UpdatePollVote(&poll, req_vote_id); update_err != nil {
		returnError(w, fmt.Errorf("UpdatePollVote: %w", update_err))
		return
	}

	fmt.Fprintf(w, `{"success": true}`)
}

func getQueryInt(query url.Values, key string) int {
	map_val, ok := query[key]
	if !ok {
		return 0
	}

	if len(map_val) == 0 {
		return 0
	}

	int_val, strconv_err := strconv.Atoi(map_val[0])
	if strconv_err != nil {
		return 0
	}

	return int_val
}

func results(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	top_maps, top_maps_err := db.GetTopMaps()
	if top_maps_err != nil {
		returnError(w, fmt.Errorf("GetTopMaps: %w", top_maps_err))
		return
	}

	top_modes, top_modes_err := db.GetTopModes()
	if top_modes_err != nil {
		returnError(w, fmt.Errorf("GetTopModes: %w", top_modes_err))
		return
	}

	top_maps_json, top_maps_json_err := json.Marshal(top_maps)
	if top_maps_json_err != nil {
		returnError(w, fmt.Errorf("marshal: %w", top_maps_json_err))
		return
	}

	top_modes_json, top_modes_json_err := json.Marshal(top_modes)
	if top_modes_json_err != nil {
		returnError(w, fmt.Errorf("marshal: %w", top_modes_json_err))
		return
	}

	fmt.Fprintf(w, `{"top_maps": %s, "top_modes": %s}`, string(top_maps_json), string(top_modes_json))
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
	http.HandleFunc("/api/results", results)
	http.HandleFunc("/", root)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
