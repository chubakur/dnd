package main

import (
	"context"
	"log"
	"os"

	ye "github.com/ydb-platform/ydb-go-sdk-auth-environ"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

type DeferFunc func()

type transport struct {
	ydbClient      *ydb.Driver
	deepSeekclient *deepSeekClient
	ctx            context.Context
}

func InitTransport(ctx context.Context) (*transport, DeferFunc, error) {
	connStr := os.Getenv("YDB_CONNECTION_STRING")
	if connStr == "" {
		log.Fatal("Set YDB_CONNECTION_STRING")
		panic("Set YDB_CONNECTION_STRING")
	}
	db, err := ydb.Open(ctx,
		connStr,
		ye.WithEnvironCredentials(),
	)
	if err != nil {
		return nil, func() {}, err
	}
	tp := &transport{
		ydbClient: db,
		ctx:       ctx,
	}
	return tp, func() {
		db.Close(ctx)
	}, nil
}
