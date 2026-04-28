package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/chubakur/dnd/transport"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type TgBinding struct {
	PlayerId uuid.UUID `sql:"player_id" json:"player_id"`
	TgId     int64     `sql:"tg_id" json:"tg_id"`
	BindTime time.Time `sql:"bind_time" json:"bind_time"`
}

type TgBindingReq struct {
	TgId int64 `json:"tg_id"`
}

func TgBindHandler(ctx context.Context, req *TgBindingReq) (*Response, error) {
	t, c, e := transport.InitTransport(ctx)
	if e != nil {
		return errorMsg(e)
	}
	defer c()
	if req.TgId == 0 {
		errorMsg(fmt.Errorf("Empty tg_id"))
	}
	bind, e := getBindingByTgId(t, req.TgId)
	if e != nil {
		return errorMsg(e)
	}
	jbind, e := json.Marshal(bind)
	if e != nil {
		return errorMsg(e)
	}
	return &Response{
		StatusCode: 200,
		Body:       string(jbind),
	}, nil
}

func getBindingByTgId(t *transport.Transport, tgId int64) (*TgBinding, error) {
	result, err := t.YdbClient.Query().QueryResultSet(
		t.Ctx,
		fmt.Sprintf("SELECT player_id, tg_id, bind_time FROM tg_bindings WHERE tg_id = %d LIMIT 1", tgId),
		query.WithIdempotent())
	if err != nil {
		return nil, err
	}
	defer result.Close(t.Ctx)

	row, err := result.NextRow(t.Ctx)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, err
	}

	var tgb TgBinding
	err = row.Scan(&tgb.PlayerId, &tgb.TgId, &tgb.BindTime)
	if err != nil {
		return nil, err
	}

	return &tgb, nil
}
