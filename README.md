# ollama-mcp

An MCP server that exposes Ollama language models as tools

## Prerequisites

- Go 1.23+
- A running [Ollama](https://ollama.com/) instance with a pulled model
  (e.g. `ollama pull llama3`)

## Build

```sh
go build -o ollama-mcp
```

Or using the Makefile:

```sh
make
```

## Configuration

Create a `config.yaml` (or copy and edit the included one):

```yaml
host: localhost
port: 11434
model: llama3
max_tokens: 1024
```

The configuration is optional. If the file is missing or some fields are not
specified, default values are used. You can also override the configuration
using environment variables:

- `OLLAMA_HOST`
- `OLLAMA_PORT`
- `OLLAMA_MODEL`
- `OLLAMA_MAX_TOKENS`

| Field        | Description                              | Default     |
|--------------|------------------------------------------|-------------|
| `host`       | Ollama server hostname / base URL        | `localhost` |
| `port`       | Ollama server port                       | `11434`     |
| `model`      | Ollama model to use                      | `llama3`    |
| `max_tokens` | Maximum tokens to generate (num_predict) | `1024`      |

## Usage

Run the server:

```sh
./ollama-mcp
```

Use a custom config path:

```sh
./ollama-mcp -config /path/to/config.yaml
```

The server communicates over stdio using the MCP protocol (JSON-RPC).
Connect it to any MCP host by configuring the host to launch this
binary as a stdio transport.

### Claude Desktop example

Add to your Claude Desktop MCP config:

```json
{
  "mcpServers": {
    "ollama": {
      "command": "/path/to/ollama-mcp",
      "args": ["-config", "/path/to/config.yaml"]
    }
  }
}
```

## Tools

### `generate`

Generate text using the configured Ollama model.

**Input:**

- `prompt` (string, required) — the prompt to send
- `system` (string, optional) — system message override

**Output:** The generated text.

### `chat`

Multi-turn conversation with the configured Ollama model.

**Input:**

- `messages` (array of objects, required) — the list of messages in the conversation
  - `role` (string, required) — the role of the message (`system`, `user`, `assistant`)
  - `content` (string, required) — the content of the message

**Output:** The generated response.

### `list_models`

List available Ollama models.

**Input:** None.

**Output:** Comma-separated list of models.
