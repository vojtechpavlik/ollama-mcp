# ollama-mcp

An MCP server that exposes Ollama language models as tools over stdio.

## Prerequisites

- Go 1.25+
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

| Field        | Description                              | Default     |
|--------------|------------------------------------------|-------------|
| `host`       | Ollama server hostname                   | `localhost` |
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

## Tool

### `generate`

Generate text using the configured Ollama model.

**Input:**

- `prompt` (string, required) — the prompt to send
- `system` (string, optional) — system message override

**Output:** The generated text.
