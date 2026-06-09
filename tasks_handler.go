package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/chubakur/dnd/async"
	"github.com/chubakur/dnd/executions"
	"github.com/chubakur/dnd/transport"
	"github.com/google/uuid"
)

type tasksProduceRequest struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// llmCall params, sent by the workflow's llmAsync step.
type llmCallParams struct {
	PlayerId uuid.UUID `json:"player_id"`
	ChatId   uuid.UUID `json:"chat_id"`
	Message  string    `json:"message"`
}

// getResult params, sent by the workflow's polling step.
type getResultParams struct {
	ExecutionId uuid.UUID `json:"execution_id"`
}

func TasksProduceHandler(ctx context.Context, req *tasksProduceRequest) (*Response, error) {
	if req.Type == "" {
		return errorMsg(fmt.Errorf("Empty type"))
	}
	t, c, e := transport.InitTransport(ctx)
	if e != nil {
		slog.ErrorContext(ctx, "transport init failed", "err", e)
		return errorMsg(e)
	}
	defer c()

	switch req.Type {
	case "llmCall":
		return handleLlmCall(ctx, t, req.Data)
	case "getResult":
		return handleGetResult(ctx, t, req.Data)
	}
	return errorMsg(fmt.Errorf("Invalid type: %s", req.Type))
}

// handleLlmCall registers a pending execution, enqueues the job, and returns
// the execution_id so the workflow can poll for the result.
func handleLlmCall(ctx context.Context, t *transport.Transport, data string) (*Response, error) {
	var p llmCallParams
	if err := json.Unmarshal([]byte(data), &p); err != nil {
		slog.ErrorContext(ctx, "llmCall: bad params", "err", err)
		return errorMsg(err)
	}

	execId, err := executions.Create(t, p.PlayerId, p.ChatId)
	if err != nil {
		slog.ErrorContext(ctx, "llmCall: create execution failed", "err", err)
		return errorMsg(err)
	}

	job := async.AsyncTaskChatLlmStruct{
		AsyncTaskCommon: async.AsyncTaskCommon{Type: "llmCall"},
		ExecutionId:     execId,
		PlayerId:        p.PlayerId,
		ChatId:          p.ChatId,
	}
	jdata, err := json.Marshal(job)
	if err != nil {
		return errorMsg(err)
	}
	if err := transport.ProduceMsg(t, "jobs", string(jdata)); err != nil {
		slog.ErrorContext(ctx, "llmCall: produce failed", "execution_id", execId, "err", err)
		return errorMsg(err)
	}

	slog.InfoContext(ctx, "llmCall queued", "execution_id", execId, "player_id", p.PlayerId, "chat_id", p.ChatId)
	body, _ := json.Marshal(map[string]string{"execution_id": execId.String()})
	return &Response{StatusCode: 200, Body: string(body)}, nil
}

// handleGetResult returns the current status (and response/error) of an execution.
func handleGetResult(ctx context.Context, t *transport.Transport, data string) (*Response, error) {
	var p getResultParams
	if err := json.Unmarshal([]byte(data), &p); err != nil {
		slog.ErrorContext(ctx, "getResult: bad params", "err", err)
		return errorMsg(err)
	}

	exec, err := executions.Get(t, p.ExecutionId)
	if err != nil {
		slog.ErrorContext(ctx, "getResult: lookup failed", "execution_id", p.ExecutionId, "err", err)
		return errorMsg(err)
	}
	if exec == nil {
		return errorMsg(fmt.Errorf("execution %s not found", p.ExecutionId))
	}

	body, err := json.Marshal(map[string]string{
		"status":   exec.Status,
		"response": exec.Response,
		"error":    exec.Error,
	})
	if err != nil {
		return errorMsg(err)
	}
	return &Response{StatusCode: 200, Body: string(body)}, nil
}
