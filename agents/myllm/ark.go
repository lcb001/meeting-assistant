package myllm

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/spf13/viper"
	"log"
	"time"
)

func CreateArkChatModel(ctx context.Context) model.ChatModel {
	viper.SetConfigFile("config.yaml")
	viper.ReadInConfig()
	apiKey := viper.GetString("ark.api_key")
	region := viper.GetString("ark.region")
	model := viper.GetString("ark.model")

	timeout := 30 * time.Second
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		// TODO: 需要从配置文件中读取
		APIKey:  apiKey,
		Region:  region,
		Model:   model,
		Timeout: &timeout,
	})
	if err != nil {
		log.Fatalf("create ark chat model failed: %v", err)
	}
	return chatModel
}
