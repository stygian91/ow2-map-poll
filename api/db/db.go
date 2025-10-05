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

var (
	instance           *sql.DB
	get_user_stmt      *sql.Stmt
	create_user_stmt   *sql.Stmt
	get_user_poll_stmt *sql.Stmt
	get_rand_maps_stmt *sql.Stmt
	create_poll_stmt   *sql.Stmt
)

func Open() error {
	_db, err := sql.Open("sqlite3", "../db/db.sqlite")
	if err != nil {
		return err
	}
	instance = _db

	get_user_stmt, err = instance.Prepare("select id, hash, created_at from users where hash = ? limit 1")
	if err != nil {
		return err
	}

	create_user_stmt, err = instance.Prepare("insert into users(hash, created_at) values (?, ?)")
	if err != nil {
		return err
	}

	get_user_poll_stmt, err = instance.Prepare("select id, created_at, user_id, map1_id, map2_id, map3_id, vote from polls where user_id = ? and (vote is null or vote = 0) order by id desc limit 1")
	if err != nil {
		return err
	}

	get_rand_maps_stmt, err = instance.Prepare("select id, name from maps order by random() limit 3")
	if err != nil {
		return err
	}

	create_poll_stmt, err = instance.Prepare("insert into polls(user_id, map1_id, map2_id, map3_id, created_at, vote) values (?, ?, ?, ?, ?, 0)")
	if err != nil {
		return err
	}

	return err
}

func Close() {
	get_user_stmt.Close()
	create_user_stmt.Close()
	instance.Close()
}

func createUser(hash string) error {
	_, err := create_user_stmt.Exec(hash, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

func getUser(hash string) (User, error) {
	user := User{}
	get_res := get_user_stmt.QueryRow(hash)
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

func getPoll(user *User) (Poll, error) {
	poll := Poll{}
	get_res := get_user_poll_stmt.QueryRow(user.Id)
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

	rows, err := get_rand_maps_stmt.Query()
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
	_, err := create_poll_stmt.Exec(user.Id, maps[0].Id, maps[1].Id, maps[2].Id, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

func GetOrCreateUserPoll(user *User) (Poll, error) {
	poll, get_err := getPoll(user)
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

		return getPoll(user)
	}

	return poll, nil
}
