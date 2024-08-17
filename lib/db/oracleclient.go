package db

import (
	"database/sql"
	"fmt"

	"github.com/go-kit/log"
	_ "github.com/sijms/go-ora/v2"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/config"
)

type OracleClient interface {
	ConnectOracle() *sql.DB
}

type oracleClient struct {
	logger log.Logger
	config config.Oracle
}

func NewOracleClient(logger log.Logger, config config.Oracle) OracleClient {
	return oracleClient{
		logger: logger,
		config: config,
	}
}

func (o oracleClient) ConnectOracle() *sql.DB {
	connectionString := fmt.Sprintf("oracle://%s:%s@%s:%s/%s", o.config.Username, o.config.Password, o.config.HostName, o.config.Port, o.config.ServiceName)

	db, err := sql.Open("oracle", connectionString)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}
	return db
}
