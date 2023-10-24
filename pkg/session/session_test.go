package session

import (
	"chatglm-sdk-go/pkg/api"
	"fmt"
	"testing"
)

func TestCompletions(t *testing.T) {

	prompt := []api.Prompt{}
	prompt = append(prompt, api.Prompt{
		Role:    "user",
		Content: "写一段冒泡排序",
	})

	request := api.NewChatCompletionRequest(api.CHATGLM_6B_SSE, 0, 0, prompt, "")

	chat := NewChatGlmSession("32c4d2cf67f74327dec011bd2c81e093.bLzGFNLyB3eSzetp")

	completions, err := chat.Completions("32c4d2cf67f74327dec011bd2c81e093.bLzGFNLyB3eSzetp", request)
	if err != nil {
		return
	}
	fmt.Sprintf("result: %v", completions)

}
