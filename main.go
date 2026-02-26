package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
)

type GenerateArgs struct {
	Prompt string `json:"prompt" jsonschema:"The prompt to send to the model"`
	System string `json:"system,omitempty" jsonschema:"Optional system message override"`
}

type Message struct {
	Role    string `json:"role" jsonschema:"The role of the message (system, user, assistant)"`
	Content string `json:"content" jsonschema:"The content of the message"`
}

type ChatArgs struct {
	Messages []Message `json:"messages" jsonschema:"The list of messages in the conversation"`
}

type ListModelsArgs struct{}

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	rawURL := cfg.Host
	if !strings.Contains(rawURL, "://") {
		rawURL = fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
	}
	base, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("Invalid Ollama URL: %v", err)
	}

	// Create a custom HTTP client. We set a large timeout here to allow for
	// long-running generation, while still preventing indefinite hangs.
	// Individual requests are also governed by the context provided by the
	// MCP host.
	httpClient := &http.Client{
		Timeout: 5 * time.Minute,
	}
	ollamaClient := api.NewClient(base, httpClient)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "ollama-mcp",
		Version: "0.1.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate",
		Description: "Generate text using the configured Ollama model",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args GenerateArgs) (*mcp.CallToolResult, any, error) {
		genReq := &api.GenerateRequest{
			Model:  cfg.Model,
			Prompt: args.Prompt,
			Options: map[string]any{
				"num_predict": cfg.MaxTokens,
			},
		}
		if args.System != "" {
			genReq.System = args.System
		}

		var sb strings.Builder
		err := ollamaClient.Generate(ctx, genReq, func(resp api.GenerateResponse) error {
			sb.WriteString(resp.Response)
			return nil
		})
		if err != nil {
			result := &mcp.CallToolResult{}
			result.SetError(fmt.Errorf("ollama generate failed: %w", err))
			return result, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: sb.String()},
			},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "chat",
		Description: "Multi-turn conversation with the configured Ollama model",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args ChatArgs) (*mcp.CallToolResult, any, error) {
		ollamaMessages := make([]api.Message, len(args.Messages))
		for i, m := range args.Messages {
			ollamaMessages[i] = api.Message{
				Role:    m.Role,
				Content: m.Content,
			}
		}
		chatReq := &api.ChatRequest{
			Model:    cfg.Model,
			Messages: ollamaMessages,
			Options: map[string]any{
				"num_predict": cfg.MaxTokens,
			},
		}

		var sb strings.Builder
		err := ollamaClient.Chat(ctx, chatReq, func(resp api.ChatResponse) error {
			sb.WriteString(resp.Message.Content)
			return nil
		})
		if err != nil {
			result := &mcp.CallToolResult{}
			result.SetError(fmt.Errorf("ollama chat failed: %w", err))
			return result, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: sb.String()},
			},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_models",
		Description: "List available Ollama models",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args ListModelsArgs) (*mcp.CallToolResult, any, error) {
		resp, err := ollamaClient.List(ctx)
		if err != nil {
			result := &mcp.CallToolResult{}
			result.SetError(fmt.Errorf("ollama list failed: %w", err))
			return result, nil, nil
		}

		var names []string
		for _, m := range resp.Models {
			names = append(names, m.Name)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Available models: %s", strings.Join(names, ", "))},
			},
		}, nil, nil
	})

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
