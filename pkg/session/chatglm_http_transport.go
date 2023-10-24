package session

import (
	"fmt"
	"net/http"

	"chatglm-sdk-go/pkg/utils"
)

const (
	SSE_CONTENT_TYPE   = "text/event-stream"
	DEFAULT_USER_AGENT = "Mozilla/4.0 (compatible; MSIE 5.0; Windows NT; DigExt)"
	APPLICATION_JSON   = "application/json"
	JSON_CONTENT_TYPE  = "application/json; charset=utf-8"

	ConnectTimeout = 450
	WriteTimeout   = 450
	ReadTimeout    = 450
)

type ChatGLMHTTPTransport struct {
	transport http.RoundTripper
	conf      *CredentialsConfiguration
}

func (c *ChatGLMHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 1. 获取原始 Request
	_ = req.Clone(req.Context())
	token := utils.NewGlmToken(&utils.APICredentials{
		ApiKey:    c.conf.ApiKey,
		ApiSecret: c.conf.ApiSecret,
	})
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", token))
	req.Header.Set("Content-Type", APPLICATION_JSON)
	req.Header.Set("User-Agent", DEFAULT_USER_AGENT)
	req.Header.Set("Accept", SSE_CONTENT_TYPE)
	// 3. 返回执行结果
	resp, err := c.transport.RoundTrip(req)
	if err != nil {
		fmt.Println("RoundTrip Error:", err)
		return nil, err
	}
	return resp, nil
}
