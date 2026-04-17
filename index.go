package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

func errorMsg(e error) (*Response, error) {
	return &Response{
		StatusCode: 500,
		Body:       e.Error(),
	}, nil
}

func QueueHandler(ctx context.Context, r queueRequest) (any, error) {

	fmt.Println(r)

	return r.Messages[0].Details.Message.Body, nil
}

func WebhookHandler(ctx context.Context, r *Request) (*Response, error) {
	connectors, close, err := InitTransport(ctx)
	if err != nil {
		return errorMsg(err)
	}
	defer close()
	mc := newMessageChain()
	mc.addUserMessage(r.Message)
	res, err := connectors.deepSeekclient.Query(mc)
	if err != nil {
		return &Response{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}
	return &Response{
		StatusCode: 200,
		Body:       res.Choices[0].Message.Content,
	}, nil
}

func worldsHandler(ctx context.Context, r *Request, connections *transport) (*Response, error) {
	descs, err := GetWorldDescriptions(ctx, connections)
	if err != nil {
		return errorMsg(err)
	}
	b, err := json.Marshal(descs)
	if err != nil {
		return errorMsg(err)
	}
	return &Response{
		StatusCode: 200,
		Body:       string(b),
	}, nil
}

func main() {
	playerId := "8a852931-c090-42c9-b7d4-e9b69721174f"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	connectors, close, err := InitTransport(ctx)
	if err != nil {
		panic(err)
	}
	defer close()
	// playerUuid, err := uuid.Parse(playerId)
	// if err != nil {
	// 	panic(err)
	// }
	// activeSessions, err := GetActivePlayerSessions(connectors, playerUuid)
	// fmt.Println(activeSessions)
	// fmt.Println(err)
	// panic("end")

	mc := newMessageChain()
	// mc.addUserMessage("Привет, дружище, подскажи, какие сеттинги для игры ты знаешь?")
	mc.addSystemMessage(fmt.Sprintf("Ты gamemaster, проводящий игры. Данный пользователь имеет uuid: %s", playerId))
	mc.addUserMessage("Привет, дружище, подскажи, какие у меня есть активные сессии?")
	res, err := connectors.deepSeekclient.Query.Query(mc)
	fmt.Println(err)
	fmt.Println(res)
	for _, choice := range res.Choices {
		mc.addMessage(choice.Message)
		for _, toolCall := range choice.Message.ToolCalls {
			mcp_result := MCPCall(connectors, toolCall)
			fmt.Println(mcp_result)
			mc.addToolMessage(mcp_result)
		}
	}
	res, err = connectors.deepSeekclient.Query(mc)
	fmt.Println(err)
	fmt.Println(res)
	fmt.Println("END")

	// req := &Request{
	// 	Action: "get_worlds",
	// }
	// resp, err := WebhookHandler(ctx, req)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(resp)
}
