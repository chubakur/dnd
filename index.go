package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type Request struct {
	Action string `json:"action"`
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
	if r.Action == "default" {
		return defaultHandler(ctx, r, nil)
	}
	connectors, close, err := InitTransport(ctx)
	if err != nil {
		return errorMsg(err)
	}
	defer close()
	if r.Action == "get_worlds" {
		return worldsHandler(ctx, r, connectors)
	}
	return &Response{
		StatusCode: 200,
		Body:       "TestQ: invalid",
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

func defaultHandler(_ context.Context, r *Request, _ *transport) (*Response, error) {
	return &Response{
		StatusCode: 200,
		Body:       fmt.Sprintf("TestQ: [%s]", r.Action),
	}, nil
}

func main() {
	playerId := "8a852931-c090-42c9-b7d4-e9b69721174f"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
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
	tools := MCPGetTools()
	client := NewDeepSeekClient(apiKey, tools)
	mc := newMessageChain()
	// mc.addUserMessage("Привет, дружище, подскажи, какие сеттинги для игры ты знаешь?")
	mc.addSystemMessage(fmt.Sprintf("Ты gamemaster, проводящий игры. Данный пользователь имеет uuid: %s", playerId))
	mc.addUserMessage("Привет, дружище, подскажи, какие у меня есть активные сессии?")
	res, err := client.Query(mc)
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
	res, err = client.Query(mc)
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
