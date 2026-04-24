package main

import (
	"fmt"
	"io"

	"github.com/chubakur/dnd/transport"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicoptions"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topictypes"
)

func AddConsumer(t *transport.Transport) error {
	return t.YdbClient.Topic().Alter(t.Ctx, "jobs", topicoptions.AlterWithAddConsumers(topictypes.Consumer{
		Name: "devcons1",
	}))
}

func ConsumeMsg(t *transport.Transport) error {
	reader, err := t.YdbClient.Topic().StartReader("devcons1", topicoptions.ReadTopic("jobs"))
	if err != nil {
		return err
	}
	defer reader.Close(t.Ctx) // Закрываем reader после завершения

	msg, err := reader.ReadMessage(t.Ctx)
	if err != nil {
		return err
	}
	content, err := io.ReadAll(msg)
	if err != nil {
		return err
	}
	fmt.Println(msg.WrittenAt)
	fmt.Println(msg.Metadata)
	fmt.Println(string(content))

	if err = reader.Commit(t.Ctx, msg); err != nil {
		return err
	}
	return nil
}
