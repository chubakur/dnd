package main

import (
	"context"
	"fmt"
	"log"
	"os"

	ye "github.com/ydb-platform/ydb-go-sdk-auth-environ"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

type Request struct {
	Action string `json:"action"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

func Handler(ctx context.Context, r *Request) (*Response, error) {
	return &Response{
		StatusCode: 200,
		Body:       fmt.Sprintf("TestQ: [%s]", r.Action),
	}, nil
}

func main() {
	// apiKey := os.Getenv("DEEPSEEK_API_KEY")
	// client := NewDeepSeekClient(apiKey)
	// fmt.Println(client.Query("Привет, дружище, я тут пытаюсь с тобой коннектнуться."))
	connStr := os.Getenv("YDB_CONNECTION_STRING")
	if connStr == "" {
		log.Fatal("Set YDB_CONNECTION_STRING")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := ydb.Open(ctx,
		connStr,
		ye.WithEnvironCredentials(),
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close(ctx)
	}()

	res, err := db.Query().Query(ctx, "SELECT state_id FROM `players_states` WHERE user_id = 123")

	if err != nil {
		panic(err)
	}

	fmt.Println(res)

}
