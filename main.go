package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// MCPサーバ本体を作成
	s := server.NewMCPServer(
		"mcp-trial",
		"0.3.0",
		server.WithRecovery(),
		server.WithToolCapabilities(false),
	)

	// --- hello ツール ---
	helloTool := mcp.NewTool(
		"hello",
		mcp.WithDescription("テスト用: 任意の名前に対して挨拶します"),
		mcp.WithString(
			"name",
			mcp.Description("挨拶する相手の名前（省略可）"),
		),
	)
	s.AddTool(helloTool, helloHandler)

	// --- ping ツール ---
	pingTool := mcp.NewTool(
		"ping",
		mcp.WithDescription("疎通確認用のツール。常に 'pong' を返します"),
	)
	s.AddTool(pingTool, pingHandler)

	// --- now ツール ---
	nowTool := mcp.NewTool(
		"now",
		mcp.WithDescription("サーバの現在時刻（RFC3339）を返します"),
	)
	s.AddTool(nowTool, nowHandler)

	// 🔴 ここから下を Streamable HTTP に変更
	//
	// OpenAI Remote MCP は Streamable HTTP を強くサポートしているので、
	// /mcp パスで stateless なエンドポイントを立てる。
	httpServer := server.NewStreamableHTTPServer(
		s,
		server.WithEndpointPath("/mcp"),
		server.WithStateLess(true), // セッション周りで 400 を避けるため stateless 推奨
	)

	log.Printf("Streamable HTTP MCP server listening on :8080/mcp")
	if err := httpServer.Start(":8080"); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}

// hello ツール
func helloHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := req.GetString("name", "")
	if name == "" {
		name = "world"
	}
	msg := fmt.Sprintf("hello from MCP, %s!", name)
	return mcp.NewToolResultText(msg), nil
}

// ping ツール
func pingHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("pong"), nil
}

// now ツール
func nowHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	now := time.Now().Format(time.RFC3339)
	return mcp.NewToolResultText(now), nil
}
