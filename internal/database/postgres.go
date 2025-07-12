package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

const (
	strConnFormat = "postgres://%s:%s@%s:%s/%s"
)

var connectionInstance *pgx.Conn
var configInstance PostgresConf

type PostgresConf struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  string
}

func OpenPostgres(conf PostgresConf) error {
	configInstance = conf
	strConn := fmt.Sprintf(
		strConnFormat,
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DBName)

	var err error
	connectionInstance, err = pgx.Connect(context.Background(), strConn)
	return err
}

func ClosePostgres() {
	if connectionInstance != nil && !connectionInstance.IsClosed() {
		connectionInstance.Close(context.Background())
	}
}

func BeginTx(ctx context.Context) (pgx.Tx, error) {
	if connectionInstance.IsClosed() {
		err := OpenPostgres(configInstance)
		if err != nil {
			var tx pgx.Tx
			return tx, err
		}
	}

	return connectionInstance.Begin(ctx)
}

func GetConnection() *pgx.Conn {
	return connectionInstance
}
