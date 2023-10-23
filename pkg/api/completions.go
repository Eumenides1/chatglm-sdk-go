package api

import (
	"fmt"
	"time"
)

const (
	CHATGLM_6B_SSE   = "chatglm_6b_sse"   // ChatGLM-6B 测试模型
	CHATGLM_LITE     = "chatglm_lite"     // 轻量版模型，适用对推理速度和成本敏感的场景
	CHATGLM_LITE_32K = "chatglm_lite_32k" // 标准版模型，适用兼顾效果和成本的场景
	CHATGLM_STD      = "chatglm_std"      // 适用于对知识量、推理能力、创造力要求较高的场景
	CHATGLM_PRO      = "chatglm_pro"      // 适用于对知识量、推理能力、创造力要求较高的场景
)

type ChatCompletionRequest struct {
	Model       string   `json:"model"`       // 模型
	RequestId   string   `json:"request_id"`  // 请求ID
	Temperature float64  `json:"temperature"` // 温度【随机性】
	TopP        float64  `json:"top_p"`       // 多样性控制
	Prompt      []Prompt `json:"prompt"`      // 输入给模型的会话信息,用户输入的内容；role=user,挟带历史的内容；role=assistant
	Incremental bool     `json:"incremental"` // 智普AI sse 固定参数 incremental = true 【增量返回】
	SSEFormat   string   `json:"sseformat"`   // 用于兼容解决sse增量模式okhttpsse截取data:后面空格问题, [data: hello]。只在增量模式下使用sseFormat。
}

type Prompt struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewChatCompletionRequest(model string, temperature, topP float64, prompt []Prompt, sseFormat string) *ChatCompletionRequest {
	if model == "" {
		model = CHATGLM_6B_SSE
	}
	if temperature == 0 {
		temperature = 0.9
	}
	if topP == 0 {
		topP = 0.7
	}
	if sseFormat == "" {
		sseFormat = "data"
	}
	return &ChatCompletionRequest{
		Model:       model,
		RequestId:   fmt.Sprintf("%s-%d", "Jaguarliu", time.Now().Unix()),
		Temperature: temperature,
		TopP:        topP,
		Prompt:      prompt,
		Incremental: true,
		SSEFormat:   sseFormat,
	}
}

type Completions interface {
	Completions(model string)
}
