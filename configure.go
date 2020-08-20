package main

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

//Configure takes an input and decides which db to use
func Configure(hostname string, port string, username string, password string) (DBAccess, error) {
	var dba DBAccess
	if hostname == "fake" {
		dba = NewMockSQLAccessor(hostname, geoRecord)
	} else {
		configRead := mysql.Config{
			User:   username,
			Passwd: password,
			Net:    "tcp",
			Addr:   fmt.Sprintf("%s:%d", hostname, port),
		}
		dbRead, err := sql.Open("mysql", configRead.FormatDSN())
		if err != nil {
			return nil, errors.Wrap(err, "failed to open db global read client")
		}
		if err := dbRead.Ping(); err != nil {
			return nil, errors.Wrap(err, "failed to ping db global read server")
		}
		dba = NewSQLAccessor(dbRead, geoRecord)
	}
	return dba, nil
}
