package transport

import (
	"fmt"
	"io"
	"strings"

	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicoptions"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topictypes"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicwriter"
)

func AddConsumer(t *Transport) error {
	return t.YdbClient.Topic().Alter(t.Ctx, "jobs", topicoptions.AlterWithAddConsumers(topictypes.Consumer{
		Name: "devcons1",
	}))
}

func ConsumeMsg(t *Transport) error {
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

func ProduceMsg(t *Transport, topic, msg string) error {
	writer, err := t.YdbClient.Topic().StartWriter(topic)
	if err != nil {
		return err
	}
	err = writer.Write(t.Ctx, topicwriter.Message{Data: strings.NewReader(msg)})
	if err != nil {
		return err
	}
	err = writer.Flush(t.Ctx)
	if err != nil {
		return err
	}
	return nil
}
