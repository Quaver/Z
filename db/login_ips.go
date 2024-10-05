package db

import (
	"database/sql"
	"time"
)

type LoginIp struct {
	UserId    int    `db:"user_id"`
	Ip        string `db:"ip"`
	Timestamp int64  `db:"timestamp"`
}

// InsertLoginIpAddress Logs the ip address of a user in the database
func InsertLoginIpAddress(userId int, ip string) error {
	var result LoginIp

	err := SQL.Get(&result, "SELECT * FROM login_ips WHERE user_id = ? AND ip = ? LIMIT 1", userId, ip)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	_, err = SQL.Exec("INSERT INTO login_ips (user_id, ip, timestamp) VALUES (?, ?, ?)",
		userId, ip, time.Now().UnixMilli())

	if err != nil {
		return err
	}

	return nil
}
