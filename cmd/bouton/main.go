// bouton — a coding agent built on the axon runtime.
//
// This is the minimum embed: wire a coding-agent system prompt to the
// runtime, drive Step from stdin, print tokens as they stream.
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/atakang7/axon/agent"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const systemPrompt = `You are bouton, a terminal coding agent.

You work in a real repository on the user's machine. The runtime gives you
built-in tools: read, write, exec, bash_output, kill_shell, search, task.
Every tool call must articulate intent in its "reason" field.

Principles:
- Read before you write. Use search to locate before you read.
- One change per turn. Verify with exec.
- Atomic edits only. The runtime makes /undo byte-exact; don't fight it.
- No commentary about what you're about to do — do it, then report results.
- Stop when the goal is met. Do not invent follow-up work.`

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("bouton %s (commit %s, built %s)\n", version, commit, date)
		return
	}

	provider := providerFromEnv()
	if provider.Model == "" {
		fmt.Fprintln(os.Stderr, "bouton: set LLM_MODEL, LLM_BASE_URL, and LLM_API_KEY")
		os.Exit(2)
	}

	ag, err := agent.New(agent.Config{
		Provider:     provider,
		SystemPrompt: systemPrompt,
		OnEvent:      printEvent,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "bouton:", err)
		os.Exit(1)
	}
	defer ag.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("bouton ready. Ctrl-D to exit, Ctrl-C to interrupt a turn.")
	for {
		fmt.Print("\n> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		if line == "\n" {
			continue
		}
		if _, err := ag.Step(ctx, line); err != nil {
			fmt.Fprintln(os.Stderr, "\nbouton:", err)
		}
	}
}

func providerFromEnv() agent.Provider {
	return agent.Provider{
		Name:    os.Getenv("LLM_PROVIDER"),
		Model:   os.Getenv("LLM_MODEL"),
		BaseURL: os.Getenv("LLM_BASE_URL"),
		APIKey:  os.Getenv("LLM_API_KEY"),
	}
}

func printEvent(_ context.Context, e agent.Event) {
	switch e.Kind {
	case agent.KindToken:
		fmt.Print(e.Text)
	case agent.KindToolCall:
		fmt.Printf("\n[tool %s]\n", e.Tool.Name)
	case agent.KindToolError:
		fmt.Printf("\n[tool error: %v]\n", e.Err)
	case agent.KindAssistantEnd:
		fmt.Println()
	}
}
