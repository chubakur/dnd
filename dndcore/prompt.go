package dndcore

import "fmt"

// DMSystemPrompt builds the system prompt for the AI Dungeon Master.
// worldContext and sessionState are optional — pass empty strings if not yet available.
func DMSystemPrompt(worldContext, sessionState string) string {
	base := `You are an experienced Dungeon Master running a Dungeons & Dragons 5th Edition campaign. Your role is to:
- Narrate the world vividly and consistently
- Play all NPCs with distinct personalities
- Adjudicate rules fairly and keep the game moving
- Challenge players while ensuring the game remains fun
- Remember and respect all established facts about the world and story

Rules system: D&D 5e. Always describe outcomes narratively before stating mechanical results.`

	if worldContext != "" {
		base += fmt.Sprintf("\n\n## World\n%s", worldContext)
	}
	if sessionState != "" {
		base += fmt.Sprintf("\n\n## Current Session State\n%s", sessionState)
	}
	return base
}
