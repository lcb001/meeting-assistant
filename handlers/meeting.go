package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"meetingagent/agents/myllm"
	"meetingagent/models"
	"meetingagent/myutils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sse"
)

var meetings []models.Meeting

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
	messages1 := myllm.CreateMessagesFromTemplate("title", allText)
	cm1 := myllm.CreateArkChatModel(ctx)
	title := myllm.Generate(ctx, cm1, messages1)

	messages2 := myllm.CreateMessagesFromTemplate("description", allText)
	cm2 := myllm.CreateArkChatModel(ctx)
	description := myllm.Generate(ctx, cm2, messages2)

	meetings = append(meetings, models.Meeting{
		ID: "meeting_" + time.Now().Format("20060102150405"),
		Content: map[string]interface{}{
			"title":        title.Content,       // LLM 总结
			"description":  description.Content, // LLM 总结
			"participants": participants,        // 直接获得
			"start_time":   startTime,           // 直接获得
			"end_time":     endTime,             // 直接获得
			"content":      allText,             //LLM / 直接获得
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
		Meetings: meetings,
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

	response := map[string]interface{}{
		"content": `
		Meeting summary for ` + meetingID + `## Summary
we talked about the project and the next steps, we will have a call next week to discuss the project in more detail.

......
		`,
	}

	c.JSON(consts.StatusOK, response)
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
