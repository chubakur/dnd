package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/chubakur/dnd/async"
	"github.com/chubakur/dnd/transport"
)

type tasksProduceRequest struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func TasksProduceHandler(ctx context.Context, req *tasksProduceRequest) (*Response, error) {
	if req.Type == "" {
		return errorMsg(fmt.Errorf("Empty type"))
	}
	parsed, e := async.ParseAsync(req.Type, req.Data)
	if e != nil {
		slog.ErrorContext(ctx, "parse async task failed", "type", req.Type, "err", e)
		return errorMsg(e)
	}
	t, c, e := transport.InitTransport(ctx)
	if e != nil {
		slog.ErrorContext(ctx, "transport init failed", "err", e)
		return errorMsg(e)
	}
	defer c()
	jdata, e := json.Marshal(parsed)
	if e != nil {
		return errorMsg(e)
	}
	e = transport.ProduceMsg(t, "jobs", string(jdata))
	if e != nil {
		slog.ErrorContext(ctx, "produce message failed", "type", req.Type, "err", e)
		return errorMsg(e)
	}
	slog.InfoContext(ctx, "task queued", "type", req.Type)
	return &Response{
		StatusCode: 200,
		Body:       "OK",
	}, nil
}
