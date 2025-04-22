package myutils

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func ExtractMeetingParticipants(meeting map[string]interface{}) []string {
	users := make(map[string]struct{})
	contents, ok := meeting["contents"].([]interface{})
	if !ok {
		return nil
	}

	for _, item := range contents {
		entry, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		user, ok := entry["user"].(string)
		if ok {
			users[user] = struct{}{}
		}
	}

	// Convert map keys to a slice
	result := make([]string, 0, len(users))
	for user := range users {
		result = append(result, user)
	}
	return result
}

func ExtractBeginAndEndTime(meeting map[string]interface{}) (string, string) {
	contents, ok := meeting["contents"].([]interface{})
	if !ok || len(contents) == 0 {
		return "", ""
	}

	// 获取第一个项目
	firstItem, ok1 := contents[0].(map[string]interface{})
	// 获取最后一个项目
	lastItem, ok2 := contents[len(contents)-1].(map[string]interface{})
	if !ok1 || !ok2 {
		return "", ""
	}

	// 提取 time_from 和 time_to
	timeFrom, ok3 := firstItem["time_from"].(string)
	timeTo, ok4 := lastItem["time_to"].(string)
	if !ok3 || !ok4 {
		return "", ""
	}

	return timeFrom, timeTo
}

func ExtractALLtext(meeting map[string]interface{}) string {
	contents, ok := meeting["contents"].([]interface{})
	if !ok {
		return ""
	}
	var allText string
	for _, contentItem := range contents {
		// 将每个内容项转换为 map[string]interface{} 类型
		itemMap, ok := contentItem.(map[string]interface{})
		if !ok {
			fmt.Println("转换内容项类型失败")
			continue
		}
		user, ok := itemMap["user"].(string)
		if !ok {
			fmt.Println("获取 user 失败")
			continue
		}
		// 获取 content 字段
		content, ok := itemMap["content"].(map[string]interface{})
		if !ok {
			fmt.Println("获取 content 失败")
			continue
		}
		text, ok := content["text"].(string)
		if !ok {
			fmt.Println("获取 text 失败")
			continue
		}
		allText += user + "说：" + text
	}
	return allText
}

func GetProjectRoot() string {
	// 获取当前这个文件的绝对路径（例如 myutils/helpers.go）
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("无法获取项目根路径")
	}

	// filename = .../meetingagent/myutils/helpers.go
	// 返回 .../meetingagent
	projectRoot := filepath.Join(filepath.Dir(filename), "../")
	return filepath.Clean(projectRoot)
}
