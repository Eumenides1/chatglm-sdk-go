package session

import (
	"chatglm-sdk-go/pkg/api"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	Api_Host = "https://open.bigmodel.cn/api/paas/"
)

type ChatGLMSession struct {
	Host   string
	Client *http.Client
}

type CredentialsConfiguration struct {
	ApiKey    string
	ApiSecret string
}

func NewChatGlmSession(apiSecretKey string) *ChatGLMSession {
	credentials := strings.Split(apiSecretKey, ".")
	// api 格式不合法，这里会panic
	if len(credentials) != 2 {
		panic("api secret key format error")
	}
	transport := &ChatGLMHTTPTransport{
		transport: http.DefaultTransport,
		conf: &CredentialsConfiguration{
			ApiKey:    credentials[0],
			ApiSecret: credentials[1],
		},
	}
	return &ChatGLMSession{
		Host: Api_Host,
		Client: &http.Client{
			// Timeout: time.Duration(5) * time.Second,
			Timeout:   10 * time.Second,
			Transport: transport,
		},
	}
}

func (c *ChatGLMSession) Completions(apiSecretKey string, request *api.ChatCompletionRequest) (string, error) {

	session := NewChatGlmSession(apiSecretKey)

	switch request.Model {
	case api.CHATGLM_6B_SSE, api.CHATGLM_LITE_32K, api.CHATGLM_LITE, api.CHATGLM_PRO, api.CHATGLM_STD:
	default:
		return "", api.ErrModelNotSupported
	}
	url := session.Host + "v3/completions/" + request.Model + "/sse-invoke"

	req, err := http.NewRequest("POST", url, strings.NewReader(request.String()))
	if err != nil {
		return "", err
	}

	resp, err := session.Client.Do(req)
	if err != nil {
		return "", err
	}
	// 读取 SSE 数据流
	for {
		body, _ := io.ReadAll(resp.Body)
		// 处理 SSE 事件
		fmt.Printf("Received event: %v\n", body)
	}
	defer resp.Body.Close()
	return "", nil
}
