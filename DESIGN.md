# Design

## Goal

Allow an MCP host to run specialized LLMs via Ollama in addition to
its main LLM. Each server instance connects to a single Ollama model.

## Architecture

The server is a Go program that communicates over stdio using the MCP
protocol. It acts as a bridge: the MCP host calls one of the provided
tools, and the server forwards the request to an Ollama instance,
collects the response, and returns it as text content.

```text
MCP Host <--stdio/JSON-RPC--> ollama-mcp <--HTTP--> Ollama
```

## Configuration

A YAML config file controls the Ollama connection and generation parameters:

| Field        | Type   | Description                          |
|--------------|--------|--------------------------------------|
| `host`       | string | Ollama server hostname / URL         |
| `port`       | int    | Ollama server port                   |
| `model`      | string | Model name to use for generation     |
| `max_tokens` | int    | Max tokens to generate (num_predict) |

The config file path is specified via the `-config` flag (default: `config.yaml`).
Sensible defaults are used if the file is missing or values are unspecified.
Environment variables (`OLLAMA_HOST`, `OLLAMA_PORT`, `OLLAMA_MODEL`,
`OLLAMA_MAX_TOKENS`) override the configuration.

## Tools

### `generate`

**Description:** Generate text using the configured Ollama model.

**Input:**

- `prompt` (string, required) — the prompt to send to the model
- `system` (string, optional) — system message override

**Output:** The complete generated text, returned as MCP text content.

**Error handling:** Ollama errors are returned as tool-level errors
(via `SetError`) so the MCP host can see and handle them.

### `chat`

**Description:** Multi-turn conversation with the configured Ollama model.

**Input:**

- `messages` (array of objects, required) — list of messages (role, content)

**Output:** The generated response.

### `list_models`

**Description:** List available Ollama models.

**Input:** None.

**Output:** A list of model names.

## Dependencies

- `github.com/modelcontextprotocol/go-sdk` — MCP Go SDK
- `github.com/ollama/ollama` — Ollama Go client
- `gopkg.in/yaml.v3` — YAML parsing

## File Structure

```text
ollama-mcp/
├── main.go         # Entry point: Ollama client, tool handlers, server
├── config.go       # Config struct and loading (YAML, environment, defaults)
├── config_test.go  # Unit tests for configuration
├── config.yaml     # Default configuration
├── Makefile        # Build, test, and linting
├── ollama-mcp.1    # Man page
├── go.mod
└── go.sum
```
