package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/viniciuslsdias/support-portal/config"
)

type key int

const (
	DatabasePoolContextKeyID     key = iota + 100
	appName                          = "support-portal"
	ErrorCodeDuplicatePrimaryKey     = "23505"
)

var (
	dbPool *pgxpool.Pool
)

func ConnectPool(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	connection := fmt.Sprintf(
		"host='%s' port='%s' user='%s' dbname='%s' sslmode='%s' password='%s' application_name='%s'",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUsername, cfg.PostgresDatabase, cfg.PostgresSSLMode, cfg.PostgresPassword, appName,
	)

	poolCfg, err := pgxpool.ParseConfig(connection)
	if err != nil {
		log.Fatalf("pgx_pool:config error. err=%v", err)
	}
	poolCfg.ConnConfig.PreferSimpleProtocol = true // is true to pgbouncer or pg-pool
	poolCfg.ConnConfig.LogLevel, _ = pgx.LogLevelFromString(cfg.LogLevel)
	poolCfg.MaxConns = int32(cfg.PostgresMaxPoolSize)
	poolCfg.MinConns = int32(cfg.PostgresMinPoolSize)

	if !cfg.PostgresLogEnabled {
		poolCfg.ConnConfig.LogLevel = pgx.LogLevelError
	}

	dbPool, err = pgxpool.ConnectConfig(ctx, poolCfg)
	if err != nil {
		log.Fatalf("pgx_pool:connection error. err=%v", err)
	}

	log.Println("pgx_pool:Database connection opened")

	return dbPool, nil
}

func ClosePool() {
	if dbPool != nil {
		dbPool.Close()
	}
	log.Println("pgx_pool:Database connection closed.")
}

func GetPool() *pgxpool.Pool {
	howOpenConnections("GetPool called.")
	return dbPool
}

func Begin(ctx context.Context) (context.Context, error) {
	howOpenConnections("Begin called.")

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		log.Printf("pgx_pool:begin_tx: err=%v", err)
		return ctx, err
	}

	return context.WithValue(ctx, DatabasePoolContextKeyID, tx), nil
}

func CurrentTX(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(DatabasePoolContextKeyID).(pgx.Tx)
	if !ok {
		return nil, errors.New("pgx_pool:Could not get database transaction from context")
	}
	return tx, nil
}

func Rollback(ctx context.Context) error {
	tx, err := CurrentTX(ctx)
	if err != nil {
		return err
	}

	if err := tx.Rollback(ctx); err != nil {
		return errors.Wrap(err, "pgx_pool:Error on rollback")
	}

	return nil
}

func Commit(ctx context.Context) error {
	tx, err := CurrentTX(ctx)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "pgx_pool:Error on commit")
	}

	return nil
}

func IsDuplicateKeyError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if ok && pgErr.Code == ErrorCodeDuplicatePrimaryKey {
		return true
	}
	return false
}

func howOpenConnections(message string) {
	stats := dbPool.Stat()

	cfg := config.GetConfig()

	if cfg.LogLevel == "debug" {

		log.Println(message,
			"max_size", stats.MaxConns(),
			"total", stats.TotalConns(),
			"acquired", stats.AcquiredConns(),
			"idle", stats.IdleConns(),
			"constructing", stats.ConstructingConns(),
			"acquire_sum", stats.AcquireCount(),
			"acquire_duration", stats.AcquireDuration().String(),
			"acquire_canceled", stats.CanceledAcquireCount(),
			"acquire_empty", stats.EmptyAcquireCount(),
		)

	}
}
