package main

import (
	"context"
	"fmt"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	StatusCode int `json:"statusCode"`
	Body       any `json:"body"`
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

func main() {
	// playerId := "8a852931-c090-42c9-b7d4-e9b69721174f"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	connectors, close, err := InitTransport(ctx)
	if err != nil {
		panic(err)
	}
	defer close()
	test, err := getBindingByTgId(connectors, 1234)
	fmt.Println(test)
	fmt.Println(err)
	// playerUuid, err := uuid.Parse(playerId)
	// if err != nil {
	// 	panic(err)
	// }
	// activeSessions, err := GetActivePlayerSessions(connectors, playerUuid)
	// fmt.Println(activeSessions)
	// fmt.Println(err)
	// panic("end")

	// mc := newMessageChain()
	// mc.addUserMessage("Привет, дружище, подскажи, какие сеттинги для игры ты знаешь?")
	// mc.addSystemMessage(fmt.Sprintf("Ты gamemaster, проводящий игры. Данный пользователь имеет uuid: %s", playerId))
	// mc.addUserMessage("Привет, дружище, подскажи, какие у меня есть активные сессии?")
	// res, err := connectors.deepSeekclient.Query(mc)
	// fmt.Println(err)
	// fmt.Println(res)
	// for _, choice := range res.Choices {
	// 	mc.addMessage(choice.Message)
	// 	for _, toolCall := range choice.Message.ToolCalls {
	// 		mcp_result := MCPCall(connectors, toolCall)
	// 		fmt.Println(mcp_result)
	// 		mc.addToolMessage(mcp_result)
	// 	}
	// }
	// res, err = connectors.deepSeekclient.Query(mc)
	// fmt.Println(err)
	// fmt.Println(res)
	// fmt.Println("END")

	// req := &Request{
	// 	Action: "get_worlds",
	// }
	// resp, err := WebhookHandler(ctx, req)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(resp)
}
