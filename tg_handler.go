package main

import "context"

type TgRequest struct {
	Message string `json:"message"`
	TgId    string `json:"tg_id"`
}

func WebhookHandler(ctx context.Context, r *TgRequest) (*Response, error) {
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
