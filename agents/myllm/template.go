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
	"fmt"
	"log"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func createTemplate(funcType string) prompt.ChatTemplate {
	systemMessage := ""

	if funcType == "title" {
		systemMessage = "你是一个{role}，你需要用{style}的语言总结会议主题。请注意，直接输出会议主题，前面不需要前缀（如会议主题：），会议内容可能会包含一些技术术语和行业术语，你需要根据上下文进行理解和总结。"
	} else if funcType == "description" {
		systemMessage = "你是一个{role}，你需要用{style}的语言总结核心议程。请注意，直接输出核心议程，前面不需要前缀（如核心议程：），会议内容可能会包含一些技术术语和行业术语，你需要根据上下文进行理解和总结。"
	} else if funcType == "summary" {
		systemMessage = "你是一个{role}，你需要用{style}的语言提炼会议摘要。至少包括以下内容1. 会议主题 2. 会议参与者3. 会议时间4. 会议内容5. 关键任务提取(这部分必须以以下格式总结，小红的任务是扫地，把其中的人称代词转换成具体的负责人，且需要把相同人们的做的事组合在一起，比如小红的任务是扫地、拖地)。请注意，直接输出会议摘要，前面不需要前缀（如会议摘要：），会议内容可能会包含一些技术术语和行业术语，你需要根据上下文进行理解和总结。"
	} else if funcType == "chat" {
		systemMessage = "你是一个{role}，你需要用{style}的语言。请注意，你需要结合会议的全文，回答用户问题，会议内容可能会包含一些技术术语和行业术语，你需要根据上下文进行理解和总结。"
	}

	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage(systemMessage),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("{question}"),
	)
}

func CreateMessagesFromTemplate(funcType string, allText string, ask string, HistoryChat []*schema.Message) []*schema.Message {
	template := createTemplate(funcType)
	//创建一个字符串，如果ask不为空，则为alltext+ask，若为空则为alltext
	question := allText
	if ask != "" {
		question += ("下面是我的问题：" + ask)
	}
	fmt.Printf("HistoryChat2: %s", HistoryChat)
	for _, v := range HistoryChat {
		fmt.Printf("HistoryChat3: %s", v.Content)
	}
	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"role":  "会议总结助手",
		"style": "专业",
		//"question": "我的代码一直报错，感觉好沮丧，该怎么办？",
		"question": question,
		// 对话历史（这个例子里模拟两轮对话历史）
		//"chat_history": []*schema.Message{
		//	schema.UserMessage("bob在踢球吗"),
		//	schema.AssistantMessage("bob很用力的踢球", nil),
		//	schema.UserMessage("rich在跑步吗"),
		//	schema.AssistantMessage("rich在飞快的跑步", nil),
		//},
		"chat_history": HistoryChat,
	})

	if err != nil {
		log.Fatalf("format template failed: %v\n", err)
	}

	return messages
}

// 输出结果
//func main() {
//	messages := createMessagesFromTemplate()
//	fmt.Printf("formatted message: %v", messages)
//}

// formatted message: [system: 你是一个程序员鼓励师。你需要用积极、温暖且专业的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。 user: 你好 assistant: 嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？ user: 我觉得自己写的代码太烂了 assistant: 每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。 user: 问题: 我的代码一直报错，感觉好沮丧，该怎么办？]
