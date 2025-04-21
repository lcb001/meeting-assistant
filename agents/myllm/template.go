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
	//contents := input["contents"].([]interface{})
	////lastItem := contents[len(contents)-1].(map[string]interface{})
	//var allText string
	//for _, contentItem := range contents {
	//	// 将每个内容项转换为 map[string]interface{} 类型
	//	itemMap, ok := contentItem.(map[string]interface{})
	//	if !ok {
	//		fmt.Println("转换内容项类型失败")
	//		continue
	//	}
	//	user, ok := itemMap["user"].(string)
	//	if !ok {
	//		fmt.Println("获取 user 失败")
	//		continue
	//	}
	//	// 获取 content 字段
	//	content, ok := itemMap["content"].(map[string]interface{})
	//	if !ok {
	//		fmt.Println("获取 content 失败")
	//		continue
	//	}
	//	text, ok := content["text"].(string)
	//	if !ok {
	//		fmt.Println("获取 text 失败")
	//		continue
	//	}
	//	allText += user + "说：" + text
	//}
	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"role":  "会议总结助手",
		"style": "专业",
		//"question": "我的代码一直报错，感觉好沮丧，该怎么办？",
		"question": allText,
		// 对话历史（这个例子里模拟两轮对话历史）
		//"chat_history": []*schema.Message{
		//	schema.UserMessage("你好"),
		//	schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
		//	schema.UserMessage("我觉得自己写的代码太烂了"),
		//	schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
		//},
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
