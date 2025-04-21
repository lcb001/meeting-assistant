package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"meetingagent/agents/myllm"
	"meetingagent/models"
	"meetingagent/myutils"
	"meetingagent/store"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sse"
)

// CreateMeeting handles the creation of a new meeting
func CreateMeeting(ctx context.Context, c *app.RequestContext) {
	var reqBody map[string]interface{}
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
		return
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
		return
	}

	fmt.Printf("create meeting: %s\n", string(jsonBody))

	// 使用llm处理reqBody，转换成目标Json格式

	participants := myutils.ExtractMeetingParticipants(reqBody)
	startTime, endTime := myutils.ExtractBeginAndEndTime(reqBody)
	allText := myutils.ExtractALLtext(reqBody)

	model := myllm.CreateArkChatModel(ctx)

	titleMessages := myllm.CreateMessagesFromTemplate("title", allText)
	title := myllm.Generate(ctx, model, titleMessages)

	descMessages := myllm.CreateMessagesFromTemplate("description", allText)
	description := myllm.Generate(ctx, model, descMessages)

	store.Meetings = append(store.Meetings, models.Meeting{
		ID: "meeting_" + time.Now().Format("20060102150405"),
		Content: map[string]interface{}{
			"title":        title.Content,       // LLM 总结
			"description":  description.Content, // LLM 总结
			"participants": participants,        // 直接获得
			"start_time":   startTime,           // 直接获得
			"end_time":     endTime,             // 直接获得
			"content":      allText,             // LLM / 直接获得
		},
	})

	// TODO: Implement actual meeting creation logic

	response := models.PostMeetingResponse{
		ID: "meeting_" + time.Now().Format("20060102150405"),
	}

	c.JSON(consts.StatusOK, response)
}

// ListMeetings handles listing all meetings
func ListMeetings(ctx context.Context, c *app.RequestContext) {
	// TODO: Implement actual meeting retrieval logic
	response := models.GetMeetingsResponse{
		Meetings: store.Meetings,
	}
	//response := models.GetMeetingsResponse{
	//	Meetings: []models.Meeting{
	//		{
	//			ID: "meeting_123",
	//			Content: map[string]interface{}{
	//				"title":        "Sample Meeting", // LLM 总结
	//				"description":  "This is a sample meeting", // LLM 总结
	//				"participants": []string{"John Doe", "Jane Smith"}, // 直接获得
	//				"start_time":   "2025-04-20 08:00:00", // 直接获得
	//				"end_time":     "2025-04-20 09:00:00", // 直接获得
	//				"content":      "This is the content of the meeting", //LLM / 直接获得
	//			},
	//		},
	//	},
	//}
	c.JSON(consts.StatusOK, response)
}

// GetMeetingSummary handles retrieving a meeting summary
func GetMeetingSummary(ctx context.Context, c *app.RequestContext) {
	meetingID := c.Query("meeting_id")
	if meetingID == "" {
		c.JSON(consts.StatusBadRequest, utils.H{"error": "meeting_id is required"})
		return
	}
	fmt.Printf("meetingID: %s\n", meetingID)

	// TODO: Implement actual summary retrieval logic

	// 调用llm获取总结，写入response
	// 总结包括：
	// 1. 会议主题
	// 2. 会议参与者
	// 3. 会议时间
	// 4. 会议内容
	// 5. 关键任务提取
	// 6. 关键任务管理器

	summary := "Meeting summary for ` + meetingID + `## Summary\nwe talked about the project and the next steps, we will have a call next week to discuss the project in more detail.\n\n......"
	todos := []string{
		"Andy的任务是负责语音识别部分，跟进同声传译团队技术方案，做详细调研报告，需要Tom帮忙收集不同类型会议录音数据做测试。",
		"Tom的任务是负责会议总结功能，调研现有会议总结方案并做对比分析，帮忙收集会议录音数据。",
		"Lily的任务是负责任务管理部分，与产品团队刘工详细讨论任务管理系统对接方案，建群方便后续沟通。",
	}

	store.Summarys[meetingID] = store.Summary{
		summary,
		todos,
	}

	response := map[string]interface{}{
		"content": summary,
	}

	c.JSON(consts.StatusOK, response)
}

func GetMeetingTodo(ctx context.Context, c *app.RequestContext) {
	meetingID := c.Query("meeting_id")
	if meetingID == "" {
		c.JSON(consts.StatusBadRequest, utils.H{"error": "meeting_id is required"})
		return
	}
	fmt.Printf("meetingID: %s\n", meetingID)

	// 连接 SQLite 查询TODOS
}

// HandleChat handles the SSE chat session
func HandleChat(ctx context.Context, c *app.RequestContext) {
	meetingID := c.Query("meeting_id")
	sessionID := c.Query("session_id")
	message := c.Query("message")

	if meetingID == "" || sessionID == "" {
		c.JSON(consts.StatusBadRequest, utils.H{"error": "meeting_id and session_id are required"})
		return
	}

	if message == "" {
		c.JSON(consts.StatusBadRequest, utils.H{"error": "message is required"})
		return
	}

	fmt.Printf("meetingID: %s, sessionID: %s, message: %s\n", meetingID, sessionID, message)

	// Set SSE headers
	c.Response.Header.Set("Content-Type", "text/event-stream")
	c.Response.Header.Set("Cache-Control", "no-cache")
	c.Response.Header.Set("Connection", "keep-alive")

	// Create SSE stream
	stream := sse.NewStream(c)

	// TODO: Implement actual chat logic
	// This is a simple example that sends a message every second
	ticker := time.NewTicker(time.Millisecond * 100)
	stopChan := make(chan struct{})
	go func() {
		time.AfterFunc(time.Second, func() {
			ticker.Stop()
			close(stopChan)
		})
	}()

	msg := fmt.Sprintf("Fake sample chat message: %s\n", time.Now().Format(time.RFC3339))

	for {
		select {
		case <-ticker.C:
			res := models.ChatMessage{
				Data: msg,
			}

			data, err := json.Marshal(res)
			if err != nil {
				return
			}

			event := &sse.Event{
				Data: data,
			}

			if err := stream.Publish(event); err != nil {
				return
			}
		case <-stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}
