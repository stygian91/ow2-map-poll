package db

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id        int
	Hash      string
	CreatedAt string
}

type Poll struct {
	Id        int
	CreatedAt string
	UserId    int
	Map1Id    int
	Map2Id    int
	Map3Id    int
	Vote      int
}

type Map struct {
	Id   int
	Name string
}

const DATE_FORMAT = "2006-01-02 15:04:05"

var (
	instance          *sql.DB
	cached_statements map[string]*sql.Stmt
)

var statement_queries map[string]string = map[string]string{
	"get_user":         "select id, hash, created_at from users where hash = ? limit 1",
	"create_user":      "insert into users(hash, created_at) values (?, ?)",
	"get_user_poll":    "select id, created_at, user_id, map1_id, map2_id, map3_id, vote from polls where user_id = ? and (vote is null or vote = 0) order by id desc limit 1",
	"get_rand_maps":    "select id, name from maps order by random() limit 3",
	"create_poll":      "insert into polls(user_id, map1_id, map2_id, map3_id, created_at, vote) values (?, ?, ?, ?, ?, 0)",
	"get_poll_count":   "select count(*) as cnt from polls where user_id = ? and created_at >= ? and vote is not null and vote != 0",
	"update_poll_vote": "update polls set vote = ? where id = ?",
}

func Open() error {
	_db, err := sql.Open("sqlite3", "../db/db.sqlite")
	if err != nil {
		return err
	}
	instance = _db

	cached_statements = map[string]*sql.Stmt{}

	for k, v := range statement_queries {
		stmt, stmt_err := instance.Prepare(v)
		if stmt_err != nil {
			return stmt_err
		}

		cached_statements[k] = stmt
	}

	return err
}

func Close() {
	for k := range cached_statements {
		cached_statements[k].Close()
	}
	instance.Close()
}

func createUser(hash string) error {
	_, err := cached_statements["create_user"].Exec(hash, time.Now().Format(DATE_FORMAT))
	return err
}

func getUser(hash string) (User, error) {
	user := User{}
	get_res := cached_statements["get_user"].QueryRow(hash)
	scan_err := get_res.Scan(&user.Id, &user.Hash, &user.CreatedAt)
	return user, scan_err
}

func GetOrCreateUser(hash string) (User, error) {
	user, get_err := getUser(hash)
	if get_err != nil && !errors.Is(get_err, sql.ErrNoRows) {
		return user, get_err
	}

	if errors.Is(get_err, sql.ErrNoRows) {
		create_err := createUser(hash)
		if create_err != nil {
			return user, create_err
		}

		return getUser(hash)
	}

	return user, nil
}

func GetPoll(user *User) (Poll, error) {
	poll := Poll{}
	get_res := cached_statements["get_user_poll"].QueryRow(user.Id)
	scan_err := get_res.Scan(
		&poll.Id,
		&poll.CreatedAt,
		&poll.UserId,
		&poll.Map1Id,
		&poll.Map2Id,
		&poll.Map3Id,
		&poll.Vote,
	)
	return poll, scan_err
}

func get3RandomMaps() ([3]Map, error) {
	maps := [3]Map{}

	rows, err := cached_statements["get_rand_maps"].Query()
	if err != nil {
		return maps, err
	}

	i := 0
	for rows.Next() {
		if err := rows.Scan(&maps[i].Id, &maps[i].Name); err != nil {
			return maps, err
		}
		i += 1
	}

	return maps, nil
}

func createPoll(user *User, maps [3]Map) error {
	_, err := cached_statements["create_poll"].Exec(user.Id, maps[0].Id, maps[1].Id, maps[2].Id, time.Now().Format(DATE_FORMAT))
	return err
}

func GetOrCreateUserPoll(user *User) (Poll, error) {
	poll, get_err := GetPoll(user)
	if get_err != nil && !errors.Is(get_err, sql.ErrNoRows) {
		return poll, get_err
	}

	if errors.Is(get_err, sql.ErrNoRows) {
		maps, get_maps_err := get3RandomMaps()
		if get_maps_err != nil {
			return poll, get_maps_err
		}

		if create_err := createPoll(user, maps); create_err != nil {
			return poll, create_err
		}

		return GetPoll(user)
	}

	return poll, nil
}

func GetPollCountForToday(user *User) (cnt int, err error) {
	yesterday := time.Now().Add(-time.Hour * 24).Format(DATE_FORMAT)
	res := cached_statements["get_poll_count"].QueryRow(user.Id, yesterday)
	err = res.Scan(&cnt)
	return
}

func UpdatePollVote(poll *Poll, vote int) error {
	_, err := cached_statements["update_poll_vote"].Exec(vote, poll.Id)
	return err
}
