package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	// postgresql driver
	_ "github.com/lib/pq"

	"github.com/sauravgsh16/api-grpc/pkg/protocol/grpc"
	"github.com/sauravgsh16/api-grpc/pkg/protocol/rest"

	v1 "github.com/sauravgsh16/api-grpc/pkg/service/v1"
)

// Config stores Server configuration
type Config struct {
	GRPCPort string
	HTTPPort string
	DBHost   string
	DBPort   int64
	DBUser   string
	DBPwd    string
	DBName   string
}

// RunServer runs gRPC server and opens DB
func RunServer() error {
	ctx := context.Background()

	// Get all configurations
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "gRPC-port", "", "port which gRPC server binds to")
	flag.StringVar(&cfg.HTTPPort, "http-port", "", "port which HTTP server binds to")
	flag.StringVar(&cfg.DBHost, "db-host", "", "Database Host")
	flag.Int64Var(&cfg.DBPort, "db-port", 5432, "Database port")
	flag.StringVar(&cfg.DBUser, "db-user", "", "Database User")
	flag.StringVar(&cfg.DBPwd, "db-password", "", "Database Password")
	flag.StringVar(&cfg.DBName, "db-name", "", "Database name")
	flag.Parse()

	/*
		cfg := Config{
			GRPCPort: "9090",
			HTTPPort: "",
			DBHost:   "localhost",
			DBPort:   5432,
			DBUser:   "postgres",
			DBPwd:    "postgres",
			DBName:   "tasktodo",
		}
	*/

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: %s", cfg.GRPCPort)
	}

	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP server: %s", cfg.HTTPPort)
	}

	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPwd,
		cfg.DBName,
	)
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1API := v1.NewToDoServiceServer(db)

	// Run REST server
	go func() {
		rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
