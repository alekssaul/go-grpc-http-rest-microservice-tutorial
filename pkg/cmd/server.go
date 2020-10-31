package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/alekssaul/go-grpc-http-rest-microservice-tutorial/pkg/protocol/grpc"
	v1 "github.com/alekssaul/go-grpc-http-rest-microservice-tutorial/pkg/service/v1"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort string

	// DB Datastore parameters section
	// DatastoreDBHost is host of database
	DatastoreDBHost string
	// DatastoreDBUser is username to connect to database
	DatastoreDBUser string
	// DatastoreDBPassword password to connect to database
	DatastoreDBPassword string
	// DatastoreDBSchema is schema of database
	DatastoreDBSchema string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	cfg.GRPCPort = os.Getenv("PORT")
	cfg.DatastoreDBHost = os.Getenv("MYSQL_HOST")
	cfg.DatastoreDBUser = os.Getenv("MYSQL_USER")
	cfg.DatastoreDBPassword = os.Getenv("MYSQL_PASSWORD")
	cfg.DatastoreDBSchema = os.Getenv("MYSQL_DB")

	if len(cfg.GRPCPort) == 0 {
		fmt.Printf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
		cfg.GRPCPort = "8080"
	}

	// add MySQL driver specific parameter to parse date/time
	// Drop it for another database
	param := "parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBSchema,
		param)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1API := v1.NewToDoServiceServer(db)

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
