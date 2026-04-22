package main

import (
	"fmt"
	"time"

	"github.com/chubakur/dnd/transport"
	"github.com/google/uuid"
)

type TgBinding struct {
	PlayerId uuid.UUID `sql:"player_id"`
	TgId     int64     `sql:"tg_id"`
	BindTime time.Time `sql:"bind_time"`
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
