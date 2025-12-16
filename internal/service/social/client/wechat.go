package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
)

type WeChatClient struct {
	Client *model.SocialClient
}

func NewWeChatClient(client *model.SocialClient) *WeChatClient {
	return &WeChatClient{Client: client}
}

type WeChatMiniSessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type WeChatMPAccessTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

type WeChatMPUserInfoResp struct {
	OpenID   string `json:"openid"`
	Nickname string `json:"nickname"`
	Sex      int    `json:"sex"`
	Province string `json:"province"`
	City     string `json:"city"`
	Country  string `json:"country"`
	HeadImg  string `json:"headimgurl"`
	UnionID  string `json:"unionid"`
}

func (c *WeChatClient) GetAuthUser(ctx context.Context, code string, state string) (*AuthUser, error) {
	if c.Client.SocialType == 31 {
		// 微信小程序
		return c.getMiniAuthUser(code)
	}
	// 微信公众号
	return c.getMPAuthUser(code)
}

// 小程序登录
func (c *WeChatClient) getMiniAuthUser(code string) (*AuthUser, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		c.Client.ClientId, c.Client.ClientSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var session WeChatMiniSessionResp
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, err
	}

	if session.ErrCode != 0 {
		return nil, core.NewBizError(1002004003, fmt.Sprintf("微信登录失败: %s", session.ErrMsg))
	}

	return &AuthUser{
		Openid:       session.OpenID,
		Token:        session.SessionKey, // 小程序使用 session_key 作为 token 凭证
		RawTokenInfo: toJson(session),
		// 小程序登录不直接返回用户信息，需要前端 getUserProfile 配合，这里先返回空
		Nickname:    "",
		Avatar:      "",
		RawUserInfo: "{}",
	}, nil
}

// 公众号登录
func (c *WeChatClient) getMPAuthUser(code string) (*AuthUser, error) {
	// 1. 获取 Access Token
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		c.Client.ClientId, c.Client.ClientSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResp WeChatMPAccessTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	if tokenResp.ErrCode != 0 {
		return nil, core.NewBizError(1002004003, fmt.Sprintf("微信登录失败: %s", tokenResp.ErrMsg))
	}

	// 2. 获取用户信息
	userUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		tokenResp.AccessToken, tokenResp.OpenID)

	userResp, err := http.Get(userUrl)
	if err != nil {
		return nil, err
	}
	defer userResp.Body.Close()

	var userInfo WeChatMPUserInfoResp
	if err := json.NewDecoder(userResp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &AuthUser{
		Openid:       tokenResp.OpenID,
		Token:        tokenResp.AccessToken,
		RawTokenInfo: toJson(tokenResp),
		Nickname:     userInfo.Nickname,
		Avatar:       userInfo.HeadImg,
		RawUserInfo:  toJson(userInfo),
	}, nil
}

func toJson(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func (c *WeChatClient) GetAuthUrl(state string, redirectUri string) string {
	if c.Client.SocialType == 31 {
		return "" // 小程序不支持跳转登录
	}
	// 公众号
	// scope: snsapi_userinfo (需关注?) or snsapi_base (静默).
	// RuoYi explicitly uses snsapi_userinfo usually for full info.
	return fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect",
		c.Client.ClientId, redirectUri, state)
}
