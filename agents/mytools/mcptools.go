package mytools

import (
	einomcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"

	"context"
)

func initMCP(ctx context.Context) *client.Client {
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "todo-client",
		Version: "1.0.0",
	}

	// stdio client
	cli, _ := client.NewStdioMCPClient("node", nil, "/Users/leyan/Code/todo-list-mcp/dist/index.js")
	_, _ = cli.Initialize(ctx, initRequest)

	return cli
}

func GetMCPTools() []tool.BaseTool {
	ctx := context.Background()

	cli := initMCP(ctx)
	tools, _ := einomcp.GetTools(ctx, &einomcp.Config{Cli: cli})

	return tools
}
