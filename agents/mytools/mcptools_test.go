package mytools

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"log"
	"testing"

	"meetingagent/agents/myllm"
)

func TestGetMCPTools(t *testing.T) {
	tools := GetMCPTools()
	if len(tools) == 0 {
		t.Error("Expected to get MCP tools, but got none")
	}

	for _, tool := range tools {
		info, err := tool.Info(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Name:", info.Name)
		fmt.Println("Desc:", info.Desc)
		fmt.Println()
	}
}

func TestAddTodo(t *testing.T) {
	ctx := context.Background()

	tools := GetMCPTools()
	llm := myllm.CreateArkChatModel(ctx)

	agent, _ := react.NewAgent(ctx, &react.AgentConfig{
		Model:       llm,
		ToolsConfig: compose.ToolsNodeConfig{Tools: tools},
	})

	agent.Generate(ctx, []*schema.Message{schema.UserMessage("添加一个写代码的todo")})
	message, _ := agent.Generate(ctx, []*schema.Message{schema.UserMessage("显示所有的todo")})

	fmt.Printf("%+v\n", message)
}
