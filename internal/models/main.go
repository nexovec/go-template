package models

import (
	"configuration"

	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
)

var totalPgxConnections uint32
var Pool *pgxpool.Pool

func Initialize() {
	dsn := lo.Must(configuration.GetAppDeploymentConfiguration()).DbDSNStringPgx
	slog.Info("Initializing database")
	initCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// configures the connection pool
	cfg := lo.Must(configuration.GetAppConfiguration())
	if !strings.Contains(dsn, "pool_max_conns") {
		size_str := fmt.Sprintf(" pool_max_conns=%d", cfg.DbConnectionPoolSize)
		slog.Info("Setting pool_max_conns", "value", cfg.DbConnectionPoolSize)
		dsn += size_str
	}
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}
	config.ConnConfig.StatementCacheCapacity = lo.Must(configuration.GetAppConfiguration()).DbStatementCacheCapacity

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		atomic.AddUint32(&totalPgxConnections, 1)
		slog.Info("Total connections", "total", atomic.LoadUint32(&totalPgxConnections))
		return nil
	}
	config.BeforeClose = func(conn *pgx.Conn) {
		atomic.AddUint32(&totalPgxConnections, ^uint32(1)) // this decrements the counter by 1 by using integer overflow
		slog.Info("Total connections", "total", atomic.LoadUint32(&totalPgxConnections))
	}
	pool, err := pgxpool.NewWithConfig(initCtx, config)
	if err != nil {
		panic(err)
	}
	Pool = pool
	// pings connection if not in production
	if os.Getenv("DEPLOYMENT") == configuration.EnumDeploymentDev || os.Getenv("DEPLOYMENT") == configuration.EnumDeploymentDebug {
		conns := pool.AcquireAllIdle(initCtx)
		for _, conn := range conns {
			err := conn.Ping(initCtx)
			if err != nil {
				panic(err)
			}
			conn.Release()
		}
		// creates mock accounts
		pwd := "admin"
		_, err = New(pool).InsertUser(initCtx, "admin", "admin@admin.admin", pwd)
		if err != nil && err.Error() != configuration.ErrUserExists && err.Error() != configuration.ErrNoRows {
			panic(err)
		}

		pwd = "johndoe"
		_, err = New(pool).InsertUser(initCtx, "johndoe", "john@doe.com", pwd)
		if err != nil && err.Error() != configuration.ErrUserExists && err.Error() != configuration.ErrNoRows {
			panic(err.Error())
		}
	}
}
