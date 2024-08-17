package services

import (
	"database/sql"
	"fmt"

	_ "github.com/godror/godror"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/passwordless-auth/config"
)

func Connection(dbname string, opts *config.DBConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"%s/%s@//%s:%d/%s",
		opts.EserviceUsername,
		opts.EservicePassword,
		opts.EserviceHost,
		opts.EservicePort,
		opts.EserviceServiceName,
	)
	if dbname == "ebill" {
		connectionString = fmt.Sprintf(
			"%s/%s@//%s:%d/%s",
			opts.EbillUsername,
			opts.EbillPassword,
			opts.EbillHost,
			opts.EbillPort,
			opts.EbillServiceName,
		)
	}
	conn, err := sql.Open("godror", connectionString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
