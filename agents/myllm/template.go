/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package myllm

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func createTemplate(funcType string) prompt.ChatTemplate {
	systemMessage := ""

	if funcType == "title" {
		systemMessage = "你是一个{role}，你需要用{style}的语言总结会议主题，不超过20字。请注意，直接输出会议主题，前面不需要前缀（如会议主题：）会议内容可能会包含一些技术术语和行业术语，你需要根据上下文进行理解和总结。"
	} else if funcType == "description" {
		systemMessage = "你是一个{role}，你需要用{style}的语言总结核心议程。请注意，直接输出核心议程，前面不需要前缀（如核心议程：）会议内容可能会包含一些技术术语和行业术语，你需要根据上下文进行理解和总结。"
	} else if funcType == "summary" {
		systemMessage = "你是一个{role}，你需要用{style}的语言提炼会议内容。请注意，会议内容可能会包含一些技术术语和行业术语，你需要根据上下文进行理解和总结。"
	} else if funcType == "chat" {

	}

	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage(systemMessage),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("会议内容: {question}"),
	)
}

func CreateMessagesFromTemplate(funcType string, allText string) []*schema.Message {
	template := createTemplate(funcType)
	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"role":     "会议总结助手",
		"style":    "专业",
		"question": allText,
	})
	if err != nil {
		log.Fatalf("format template failed: %v\n", err)
	}
	return messages
}
