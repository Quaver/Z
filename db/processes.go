package db

import (
	"crypto/md5"
	"encoding/hex"
)

type Process struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

// GetMD5 Returns an MD5 of Process.Name
func (p *Process) GetMD5() string {
	hash := md5.Sum([]byte(p.Name))
	return hex.EncodeToString(hash[:])
}

// FetchProcesses Fetches all the tracked processes from the database
func FetchProcesses() ([]*Process, error) {
	processes := make([]*Process, 0)

	err := SQL.Select(&processes, "SELECT * FROM processes")

	if err != nil {
		return nil, err
	}

	return processes, nil
}
