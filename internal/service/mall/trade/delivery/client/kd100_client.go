package client

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/wxlbd/ruoyi-mall-go/pkg/config"
)

type Kd100ExpressClient struct {
	conf config.Kd100Config
}

func NewKd100ExpressClient(conf config.Kd100Config) *Kd100ExpressClient {
	return &Kd100ExpressClient{conf: conf}
}

func (c *Kd100ExpressClient) GetExpressTrackList(req *ExpressTrackQueryReqDTO) ([]ExpressTrackRespDTO, error) {
	// 1. 准备请求参数
	param := map[string]string{
		"com":      req.ExpressCode,
		"num":      req.LogisticsNo,
		"phone":    req.Phone,
		"from":     "",
		"to":       "",
		"resultv2": "1", // 开启行政区域解析
		"show":     "0", // 返回 json
		"order":    "desc",
	}
	paramJson, _ := json.Marshal(param)

	// 2. 签名
	// 签名规则：MD5(param + key + customer) 转大写
	signStr := string(paramJson) + c.conf.Key + c.conf.Customer
	hasher := md5.New()
	hasher.Write([]byte(signStr))
	sign := strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))

	// 3. 发送请求
	formData := url.Values{}
	formData.Set("customer", c.conf.Customer)
	formData.Set("sign", sign)
	formData.Set("param", string(paramJson))

	resp, err := http.PostForm("https://poll.kuaidi100.com/poll/query.do", formData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 4. 解析响应
	var respData struct {
		Message string                `json:"message"`
		State   string                `json:"state"`
		Status  string                `json:"status"`
		Data    []ExpressTrackRespDTO `json:"data"`
	}

	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, err
	}

	if respData.Status != "200" {
		return nil, fmt.Errorf("查询失败: %s", respData.Message)
	}

	return respData.Data, nil
}
