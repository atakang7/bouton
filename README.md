# bouton

**A terminal coding agent built on the [axon](https://github.com/atakang7/axon) runtime.**

The runtime ships the loop (streaming, tool dispatch, append-only session, pruning). bouton ships the coding-agent personality: system prompt, defaults, and a thin CLI wrapper around `agent.Step`.

```
github.com/atakang7/axon/agent  ← runtime (signals)
github.com/atakang7/bouton      ← coding agent (terminal of the axon)
```

## Install

Pre-built binaries for linux/darwin × amd64/arm64 are attached to each [GitHub Release](https://github.com/atakang7/bouton/releases).

Or build from source:

```sh
go install github.com/atakang7/bouton/cmd/bouton@latest
```

## Run

Point it at any OpenAI-compatible endpoint:

```sh
LLM_MODEL=gpt-4o \
LLM_BASE_URL=https://api.openai.com \
LLM_API_KEY=sk-... \
bouton
```

Talk to it. It reads, writes, and executes in the current directory.

```
> add a /version flag to cmd/foo
[tool search]
[tool read]
[tool write] cmd/foo/main.go
[tool exec] go build ./...
done.
```

## What's built in

bouton inherits the runtime's tool set from axon: `read`, `write`, `exec`, `search`, `task`, `bash_output`, `kill_shell`. Every tool call requires a `reason` field. Edits are atomic. Ctrl-C cancels the in-flight turn — both the model stream and any running subprocess.

There is no sandbox. Use Ctrl-C and `git` as your guardrails.

## Configuration

bouton reads three environment variables:

| Variable        | Purpose                                          |
| --------------- | ------------------------------------------------ |
| `LLM_MODEL`     | Model name (`gpt-4o`, `claude-3-opus`, …)        |
| `LLM_BASE_URL`  | OpenAI-compatible API base URL                   |
| `LLM_API_KEY`   | API key                                          |
| `LLM_PROVIDER`  | Optional. Name hint for the provider.            |

Bring your own provider — anything OpenAI-compatible works (OpenAI, Anthropic via proxy, Ollama, LM Studio, vLLM).

## Why a separate repo

axon is the runtime — a library that drives the agent loop and can host any agent personality. bouton is one such personality: opinionated for coding work in a terminal. Keeping them separate means:

- axon stays minimal and importable for other agents (reviewer, researcher, …).
- bouton can ship binaries, slash commands, and a coding-flavored system prompt without bloating the library.

## License

MIT. See [LICENSE](LICENSE).
