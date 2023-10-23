package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

var cachedTokens *cache.Cache
var sugarLogger *zap.SugaredLogger

// ExpireMillis Token过期时间 这里的单位是分钟
const ExpireMillis = 30

type GlmToken struct {
	Token string
}

// APICredentials Api加密信息
type APICredentials struct {
	ApiKey         string // 登录创建 ApiKey <a href="https://open.bigmodel.cn/usercenter/apikeys">apikeys</a>
	ApiSecret      string // apiKey的后半部分 828902ec516c45307619708d3e780ae1.w5eKiLvhnLP8MtIf 取 w5eKiLvhnLP8MtIf 使用
	ExpireDuration int    // 过期时间
}

func init() {
	cachedTokens = cache.New(5*time.Minute, 10*time.Minute)
	logger, _ := zap.NewProduction()
	sugarLogger = logger.Sugar()
}

func NewGlmToken(credentials *APICredentials) *GlmToken {
	token, err := generateToken(credentials.ApiKey, []byte(credentials.ApiSecret), credentials.ExpireDuration)
	if err != nil {
		sugarLogger.Errorf(" generate token failed: %s", err)
		return nil
	}
	return &GlmToken{
		Token: token,
	}
}

func generateToken(apiKey string, apiSecret []byte, expireDuration int) (string, error) {
	var expireTime time.Duration
	if expireDuration == 0 {
		// 默认缓存半小时
		expireTime = ExpireMillis * time.Minute
	} else {
		expireTime = time.Duration(expireDuration) * time.Minute
	}

	// 缓存token
	if token := findTokenInCache(apiKey); token != "" {
		return token, nil
	}
	// 创建Token
	header := map[string]interface{}{
		"alg":       "HS256",
		"sign_type": "SIGN",
	}
	claims := jwt.MapClaims{
		"api_key":   apiKey,
		"exp":       time.Now().Add(expireTime).UnixNano() / int64(time.Millisecond),
		"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header = header
	signedToken, err := token.SignedString(apiSecret)
	if err != nil {
		sugarLogger.Errorf("signing token failed: %s", err)
		return "", err
	}
	err = cacheToken(apiKey, signedToken)
	if err != nil {
		sugarLogger.Errorf("cache token failed: %s", err)
		return "", err
	}
	return signedToken, nil
}

// cacheToken 函数接收两个参数 key 和 token，用于存储和检索令牌信息。
func cacheToken(key string, token string) error {
	// 使用 go-cache 的 Add 方法将令牌添加到缓存中，设置过期时间为 5 分钟。
	if _, exists := cachedTokens.Get(key); !exists {
		err := cachedTokens.Add(key, token, time.Duration(5)*time.Minute)
		if err != nil {
			sugarLogger.Errorf("cache token failed: %s", err)
			return err
		} else {
			sugarLogger.Infof("Added new token to cache for %s with value %s and expiration in 5 minutes.\n", key, token)
		}
	} else {
		sugarLogger.Errorf(" Token already exists in cache for %s, not adding new token.\n", key)
		return errors.New("token already exists in cache")
	}
	return nil
}

func findTokenInCache(key string) string {
	// 检查缓存是否包含给定的键，如果存在则获取其值。
	value, found := cachedTokens.Get(key)
	if found {
		sugarLogger.Infof("Value of token for %s is %v\n", key, value)
		return value.(string)
	} else {
		sugarLogger.Infof("No token found in cache for %s\n", key)
		return ""
	}
}
