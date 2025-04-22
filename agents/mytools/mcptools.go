package mytools

import (
	einomcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"log"
	"os"

	"context"
)

func initMCP(ctx context.Context) *client.Client {
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "todo-client",
		Version: "1.0.0",
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("无法获取工作目录:", err)
	}
	path := wd + "/agents/mytools/mcp-todo-list/dist/index.js"
	// stdio client
	cli, _ := client.NewStdioMCPClient("node", nil, path)
	_, _ = cli.Initialize(ctx, initRequest)

	return cli
}

func GetMCPTools() []tool.BaseTool {
	ctx := context.Background()

	cli := initMCP(ctx)
	tools, _ := einomcp.GetTools(ctx, &einomcp.Config{Cli: cli})

	return tools
}
