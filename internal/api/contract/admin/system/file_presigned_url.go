package system

// FilePresignedUrlResp 文件预签名地址 Response
type FilePresignedUrlResp struct {
	ConfigID  int64  `json:"configId"`
	UploadURL string `json:"uploadUrl"`
	URL       string `json:"url"`
	Path      string `json:"path"`
}
