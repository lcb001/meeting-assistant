package myllm

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"log"
	"time"
)

func CreateArkChatModel(ctx context.Context) model.ChatModel {
	timeout := 30 * time.Second
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		// TODO: 需要从配置文件中读取
		APIKey:  "e4900371-625c-413c-83d2-c22f1f6efc9c",
		Region:  "cn-beijing",
		Model:   "doubao-pro-32k-241215",
		Timeout: &timeout,
	})
	if err != nil {
		log.Fatalf("create ark chat model failed: %v", err)
	}
	return chatModel
}
