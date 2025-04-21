package mytools

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"log"
	"meetingagent/models"
	"meetingagent/store"
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

	meetingID := "meeting_20231010123456"

	store.Meetings = append(store.Meetings, models.Meeting{
		ID: meetingID,
		Content: map[string]any{
			"title":        "会议标题",
			"description":  "会议描述",
			"participants": []string{"Andy", "Tom", "Lily"},
			"start_time":   "2023-10-10 12:00:00",
			"end_time":     "2023-10-10 13:00:00",
			"content":      "会议内容",
		},
	})

	store.Summaries = make(map[string]store.Summary)
	store.Summaries[meetingID] = store.Summary{
		Content: "123",
		Todos: []string{
			"Andy的任务是负责语音识别部分，跟进同声传译团队技术方案，做详细调研报告，需要Tom帮忙收集不同类型会议录音数据做测试。",
			"Tom的任务是负责会议总结功能，调研现有会议总结方案并做对比分析，帮忙收集会议录音数据。",
			"Lily的任务是负责任务管理部分，与产品团队刘工详细讨论任务管理系统对接方案，建群方便后续沟通。",
		},
	}

	template := prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个会议助手，需要按照给出的内容生成待办事项，标题为主要事件，使用会议名称作为所属清单名，如果有负责人则将assignee设置为负责人"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("会议ID:{meetingID}, 会议名称：{list}, 待办事项：{todo}"),
	)

	for _, todo := range store.Summaries[meetingID].Todos {
		messages, _ := template.Format(context.Background(), map[string]any{
			"meetingID": store.Meetings[0].ID,
			"list":      store.Meetings[0].Content["title"],
			"todo":      todo,
		})

		result, _ := agent.Generate(ctx, messages)
		fmt.Printf("%+v\n", result)
	}

}
