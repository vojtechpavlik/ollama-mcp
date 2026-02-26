package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
)

type GenerateArgs struct {
	Prompt string `json:"prompt" jsonschema:"description=The prompt to send to the model,required"`
	System string `json:"system,omitempty" jsonschema:"description=Optional system message override"`
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	base, err := url.Parse(fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Fatalf("Invalid Ollama URL: %v", err)
	}
	ollamaClient := api.NewClient(base, http.DefaultClient)

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

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
