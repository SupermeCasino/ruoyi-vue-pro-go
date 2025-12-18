package file

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileClient 文件客户端接口
type FileClient interface {
	Upload(content []byte, path string) (string, error)
	Delete(path string) error
	GetContent(path string) ([]byte, error)
	GetURL(path string) string
	GetPresignedURL(path string) (string, error)
}

// ClientConfig 客户端配置通用结构 (用于解析 JSON)
type ClientConfig struct {
	Domain   string `json:"domain"`
	BasePath string `json:"basePath"` // Local 使用
	// S3 相关字段可按需添加
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Bucket    string `json:"bucket"`
}

// LocalFileClient 本地文件客户端
type LocalFileClient struct {
	Config ClientConfig
}

func NewLocalFileClient(config json.RawMessage) (*LocalFileClient, error) {
	var cfg ClientConfig
	if err := json.Unmarshal(config, &cfg); err != nil {
		return nil, err
	}
	return &LocalFileClient{Config: cfg}, nil
}

func (c *LocalFileClient) Upload(content []byte, path string) (string, error) {
	fullPath := filepath.Join(c.Config.BasePath, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(fullPath, content, 0644); err != nil {
		return "", err
	}
	// 返回完整 URL
	return c.GetURL(path), nil
}

func (c *LocalFileClient) Delete(path string) error {
	fullPath := filepath.Join(c.Config.BasePath, path)
	return os.Remove(fullPath)
}

func (c *LocalFileClient) GetContent(path string) ([]byte, error) {
	fullPath := filepath.Join(c.Config.BasePath, path)
	return ioutil.ReadFile(fullPath)
}

func (c *LocalFileClient) GetURL(path string) string {
	return c.Config.Domain + "/" + path
}

func (c *LocalFileClient) GetPresignedURL(path string) (string, error) {
	// Local 模式下不支持真正的预签名上传，返回上传接口地址
	// 前端需特殊处理：如果是 Local，直接调用 /upload
	return c.Config.Domain + "/admin-api/infra/file/upload", nil
}

// FileClientFactory 简单工厂
func NewFileClient(storage int32, config json.RawMessage) (FileClient, error) {
	switch storage {
	case 10: // Local
		return NewLocalFileClient(config)
	case 20: // S3
		return NewS3FileClient(config)
	default:
		return nil, errors.New("unknown storage type")
	}
}
