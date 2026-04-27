package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chubakur/dnd/transport"
	"github.com/google/uuid"
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
	row, err := t.YdbClient.Query().QueryRow(t.Ctx, fmt.Sprintf("SELECT player_id, tg_id, bind_time FROM tg_bindings WHERE tg_id = %d", tgId))
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}
	var tgb TgBinding
	err = row.ScanStruct(&tgb)
	if err != nil {
		return nil, err
	}
	return &tgb, nil
}
