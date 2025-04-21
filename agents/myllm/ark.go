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
		APIKey:  "3a7e16b2-b1fb-4b56-81db-c390b91d6839",
		Region:  "cn-beijing",
		Model:   "doubao-pro-32k-241215",
		Timeout: &timeout,
	})
	if err != nil {
		log.Fatalf("create ark chat model failed: %v", err)
	}
	return chatModel
}
