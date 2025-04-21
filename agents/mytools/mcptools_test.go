package mytools

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/prompt"
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

	todos := []string{
		"Andy的任务是负责语音识别部分，跟进同声传译团队技术方案，做详细调研报告，需要Tom帮忙收集不同类型会议录音数据做测试。",
		"Tom的任务是负责会议总结功能，调研现有会议总结方案并做对比分析，帮忙收集会议录音数据。",
		"Lily的任务是负责任务管理部分，与产品团队刘工详细讨论任务管理系统对接方案，建群方便后续沟通。",
	}

	//agent.Generate(ctx, []*schema.Message{schema.UserMessage("添加一个写代码的todo，分配给John，下周三截止")})
	//message, _ := agent.Generate(ctx, []*schema.Message{schema.UserMessage("显示所有的todo")})

	template := prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个会议助手，需要按照给出的内容生成待办事项，标题为主要事件，使用会议名称作为所属清单名，如果有负责人则将assignee设置为负责人"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("会议名称：{list}，待办事项：{todo}"),
	)

	for _, todo := range todos {
		messages, _ := template.Format(context.Background(), map[string]any{
			"list": "同声传译团队技术对接会",
			"todo": todo,
		})

		result, _ := agent.Generate(ctx, messages)
		fmt.Printf("%+v\n", result)
	}

}
