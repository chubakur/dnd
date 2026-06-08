package main

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/chubakur/dnd/async"
	"github.com/chubakur/dnd/dndcore"
	"github.com/chubakur/dnd/llmcore"
	"github.com/google/uuid"
)

func TestLocalHealth(t *testing.T) {
	resp, err := HealthHandler(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if resp.Body != "Ok" {
		t.Errorf("expected 'Ok', got %v", resp.Body)
	}
}

func TestProdHealth(t *testing.T) {
	healthUrl := os.Getenv("HEALTH_CHECK_PROD_URL")
	if healthUrl == "" {
		t.Skip("HEALTH_CHECK_PROD_URL not set")
	}
	resp, err := http.Get(healthUrl)
	if err != nil {
		t.Fatal(err)
	}
	bRes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(bRes) != "Ok" {
		t.Errorf("Integration health-check error. Ok != %s", string(bRes))
	}
}

// --- dndcore.DMSystemPrompt ---

func TestDMSystemPrompt_Base(t *testing.T) {
	prompt := dndcore.DMSystemPrompt("", "")
	if !strings.Contains(prompt, "Dungeon Master") {
		t.Error("base prompt should mention Dungeon Master")
	}
	if !strings.Contains(prompt, "D&D 5e") {
		t.Error("base prompt should mention D&D 5e")
	}
}

func TestDMSystemPrompt_WithWorld(t *testing.T) {
	prompt := dndcore.DMSystemPrompt("A dark fantasy world", "")
	if !strings.Contains(prompt, "A dark fantasy world") {
		t.Error("world context should appear in prompt")
	}
	if !strings.Contains(prompt, "## World") {
		t.Error("World section header should be present")
	}
}

func TestDMSystemPrompt_WithSession(t *testing.T) {
	prompt := dndcore.DMSystemPrompt("", "Player is in a tavern")
	if !strings.Contains(prompt, "Player is in a tavern") {
		t.Error("session state should appear in prompt")
	}
	if !strings.Contains(prompt, "## Current Session State") {
		t.Error("Session State section header should be present")
	}
}

func TestDMSystemPrompt_Both(t *testing.T) {
	prompt := dndcore.DMSystemPrompt("World lore here", "Session state here")
	if !strings.Contains(prompt, "World lore here") || !strings.Contains(prompt, "Session state here") {
		t.Error("both world and session context should appear in prompt")
	}
}

// --- llmcore.MessageChain ---

func TestMessageChain_AddAndGrow(t *testing.T) {
	mc := llmcore.NewMessageChain()
	mc.AddSystemMessage("You are a DM")
	mc.AddUserMessage("Hello")
	msgs := mc.Messages()
	if len(msgs) != 2 {
		t.Errorf("expected 2 messages, got %d", len(msgs))
	}
	if msgs[0].Role != "system" {
		t.Errorf("first message should be system, got %s", msgs[0].Role)
	}
	if msgs[1].Role != "user" {
		t.Errorf("second message should be user, got %s", msgs[1].Role)
	}
}

// --- async.ParseAsync ---

func TestParseAsync_LlmCall(t *testing.T) {
	playerId := uuid.New()
	chatId := uuid.New()
	jdata := `{"type":"llmCall","player_id":"` + playerId.String() + `","chat_id":"` + chatId.String() + `"}`
	task, err := async.ParseAsync("llmCall", jdata)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task == nil {
		t.Fatal("expected non-nil task")
	}
}

func TestParseAsync_UnknownType(t *testing.T) {
	_, err := async.ParseAsync("unknown_type", "{}")
	if err == nil {
		t.Error("expected error for unknown task type")
	}
}

// --- ParseJob ---

func TestParseJob_GenerateWorld(t *testing.T) {
	playerId := uuid.New()
	jdata := `{"name":"generate_world","player_id":"` + playerId.String() + `","prompt":"A dark forest"}`
	job, err := ParseJob([]byte(jdata))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.GetName() != "generate_world" {
		t.Errorf("expected generate_world, got %s", job.GetName())
	}
}

func TestParseJob_Unknown(t *testing.T) {
	_, err := ParseJob([]byte(`{"name":"nonexistent"}`))
	if err == nil {
		t.Error("expected error for unknown job name")
	}
}
