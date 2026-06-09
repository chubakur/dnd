package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/chubakur/dnd/async"
	"github.com/chubakur/dnd/dndcore"
	"github.com/chubakur/dnd/executions"
	"github.com/chubakur/dnd/llmcore"
	"github.com/chubakur/dnd/mcp"
	"github.com/chubakur/dnd/transport"
)

const agentExecutorLimit = 10

type Response struct {
	StatusCode int `json:"statusCode"`
	Body       any `json:"body"`
}

// ydbTopicEvent is the payload a YDB Topic / Data Streams trigger delivers.
// The actual message bytes live (base64-encoded) under one of a few paths
// depending on the trigger type, so we probe several.
type ydbTopicEvent struct {
	Messages []ydbTopicMessage `json:"messages"`
}

type ydbTopicMessage struct {
	Details struct {
		Data    string `json:"data"`
		Message struct {
			Data string `json:"data"`
			Body string `json:"body"`
		} `json:"message"`
	} `json:"details"`
}

// payload returns the decoded message bytes, probing the known trigger paths.
func (m ydbTopicMessage) payload() string {
	for _, candidate := range []string{m.Details.Data, m.Details.Message.Data, m.Details.Message.Body} {
		if candidate == "" {
			continue
		}
		// Trigger payloads are usually base64-encoded; fall back to raw.
		if decoded, err := base64.StdEncoding.DecodeString(candidate); err == nil {
			return string(decoded)
		}
		return candidate
	}
	return ""
}

// QueueHandler consumes LLM jobs from the YDB topic trigger, runs the agent,
// and writes the result back to the executions table for the workflow to poll.
func QueueHandler(ctx context.Context, raw json.RawMessage) (*Response, error) {
	slog.InfoContext(ctx, "queue handler invoked", "raw", string(raw))

	var event ydbTopicEvent
	if err := json.Unmarshal(raw, &event); err != nil {
		slog.ErrorContext(ctx, "failed to parse trigger event", "err", err)
		return errorMsg(err)
	}
	if len(event.Messages) == 0 {
		slog.WarnContext(ctx, "queue event had no messages")
		return &Response{StatusCode: 200, Body: "no messages"}, nil
	}

	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return errorMsg(fmt.Errorf("DEEPSEEK_API_KEY not set"))
	}

	t, closeT, err := transport.InitTransport(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "transport init failed", "err", err)
		return errorMsg(err)
	}
	defer closeT()

	processed := 0
	for _, m := range event.Messages {
		payload := m.payload()
		if payload == "" {
			slog.WarnContext(ctx, "empty message payload, skipping")
			continue
		}
		if err := processJob(ctx, t, apiKey, []byte(payload)); err != nil {
			// processJob already records the error on the execution; log and continue.
			slog.ErrorContext(ctx, "job processing failed", "err", err)
		}
		processed++
	}

	return &Response{StatusCode: 200, Body: fmt.Sprintf("processed %d", processed)}, nil
}

// processJob runs one LLM job and persists its result/error on the execution.
func processJob(ctx context.Context, t *transport.Transport, apiKey string, payload []byte) error {
	var job async.AsyncTaskChatLlmStruct
	if err := json.Unmarshal(payload, &job); err != nil {
		slog.ErrorContext(ctx, "bad job payload", "err", err)
		return err
	}
	slog.InfoContext(ctx, "processing job", "execution_id", job.ExecutionId, "chat_id", job.ChatId)

	chain, err := job.Handle(t)
	if err != nil {
		_ = executions.SetError(t, job.ExecutionId, err.Error())
		return err
	}

	client := llmcore.NewDeepSeekClient(apiKey, mcp.MCPGetTools())
	pc := &dndcore.GameContext{PlayerId: job.PlayerId, ChatId: job.ChatId}
	_, resp, err := client.AgentExecutor(t, pc, chain, agentExecutorLimit)
	if err != nil {
		_ = executions.SetError(t, job.ExecutionId, err.Error())
		return err
	}

	choice := resp.GetFirstChoice()
	if choice == nil {
		err := fmt.Errorf("empty response from LLM")
		_ = executions.SetError(t, job.ExecutionId, err.Error())
		return err
	}

	if err := executions.SetDone(t, job.ExecutionId, choice.Message.Content); err != nil {
		slog.ErrorContext(ctx, "failed to store result", "execution_id", job.ExecutionId, "err", err)
		return err
	}
	slog.InfoContext(ctx, "job done", "execution_id", job.ExecutionId)
	return nil
}

func errorMsg(e error) (*Response, error) {
	return &Response{
		StatusCode: 500,
		Body:       e.Error(),
	}, e
}

// main is a local debug entrypoint; the serverless runtime uses QueueHandler.
func main() {
	if os.Getenv("DEEPSEEK_API_KEY") == "" {
		panic("Set DEEPSEEK_API_KEY")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, closeT, err := transport.InitTransport(ctx)
	if err != nil {
		panic(err)
	}
	defer closeT()
	fmt.Println("transport initialized")
}
