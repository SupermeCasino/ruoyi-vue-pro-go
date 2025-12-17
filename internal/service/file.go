package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/file"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/samber/lo"
)

type FileService struct {
	q                 *query.Query
	fileConfigService *FileConfigService
}

func NewFileService(q *query.Query, fileConfigService *FileConfigService) *FileService {
	return &FileService{
		q:                 q,
		fileConfigService: fileConfigService,
	}
}

// CreateFile 上传/创建文件
func (s *FileService) CreateFile(ctx context.Context, name string, path string, content []byte) (string, error) {
	// 1. 获取 Master 配置
	config, err := s.fileConfigService.GetMasterFileConfig(ctx)
	if err != nil {
		return "", errors.New("请先配置主文件存储")
	}

	// 2. 初始化客户端
	client, err := file.NewFileClient(config.Storage, config.Config)
	if err != nil {
		return "", fmt.Errorf("初始化文件客户端失败: %v", err)
	}

	// 3. 处理路径 (如果未提供 path，则使用 name)
	// 这里可以添加日期分目录的逻辑，如 2023/10/01/uuid.jpg
	if path == "" {
		path = fmt.Sprintf("%s/%s", time.Now().Format("2006/01/02"), name)
	}

	// 4. 上传
	url, err := client.Upload(content, path)
	if err != nil {
		return "", err
	}

	// 5. 保存记录
	fileRecord := &model.InfraFile{
		ConfigId: config.ID,
		Name:     name,
		Path:     path,
		Url:      url,
		Type:     "", // 可以通过 http.DetectContentType(content) 获取
		Size:     len(content),
	}
	err = s.q.InfraFile.WithContext(ctx).Create(fileRecord)
	if err != nil {
		return "", err
	}

	return url, nil
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(ctx context.Context, id int64) error {
	f := s.q.InfraFile
	fileRecord, err := f.WithContext(ctx).Where(f.ID.Eq(id)).First()
	if err != nil {
		return errors.New("文件不存在")
	}

	// 获取配置
	config, err := s.fileConfigService.GetFileConfig(ctx, fileRecord.ConfigId)
	if err != nil {
		// 如果配置都不存在了，只删除数据库记录
		f.WithContext(ctx).Where(f.ID.Eq(id)).Delete()
		return nil
	}

	// 初始化客户端并删除物理文件
	if config.Config != nil {
		client, err := file.NewFileClient(config.Storage, *config.Config)
		if err == nil {
			_ = client.Delete(fileRecord.Path)
		}
	}

	_, err = f.WithContext(ctx).Where(f.ID.Eq(id)).Delete()
	return err
}

// GetFileContent 获取文件内容
func (s *FileService) GetFileContent(ctx context.Context, configId int64, path string) ([]byte, error) {
	config, err := s.fileConfigService.GetFileConfig(ctx, configId)
	if err != nil {
		return nil, errors.New("配置不存在")
	}
	if config.Config == nil {
		return nil, errors.New("配置内容为空")
	}

	client, err := file.NewFileClient(config.Storage, *config.Config)
	if err != nil {
		return nil, err
	}
	return client.GetContent(path)
}

// GetFilePage 获得文件分页
func (s *FileService) GetFilePage(ctx context.Context, req *req.FilePageReq) (*pagination.PageResult[*resp.FileRespVO], error) {
	f := s.q.InfraFile
	qb := f.WithContext(ctx)

	if req.Path != "" {
		qb = qb.Where(f.Path.Like("%" + req.Path + "%"))
	}
	if req.Type != "" {
		qb = qb.Where(f.Type.Like("%" + req.Type + "%"))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	list, err := qb.Order(f.ID.Desc()).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*resp.FileRespVO]{
		List:  lo.Map(list, func(item *model.InfraFile, _ int) *resp.FileRespVO { return s.convertResp(item) }),
		Total: total,
	}, nil
}

func (s *FileService) convertResp(item *model.InfraFile) *resp.FileRespVO {
	return &resp.FileRespVO{
		ID:         item.ID,
		ConfigId:   item.ConfigId,
		Name:       item.Name,
		Path:       item.Path,
		Url:        item.Url,
		Type:       item.Type,
		Size:       item.Size,
		CreateTime: item.CreatedAt,
	}
}

func (s *FileService) GetFilePresignedUrl(ctx context.Context, path string) (*resp.FilePresignedUrlResp, error) {
	config, err := s.fileConfigService.GetMasterFileConfig(ctx)
	if err != nil {
		return nil, errors.New("请先配置主文件存储")
	}

	client, err := file.NewFileClient(config.Storage, config.Config)
	if err != nil {
		return nil, err
	}

	presignedUrl, err := client.GetPresignedURL(path)
	if err != nil {
		return nil, err
	}

	return &resp.FilePresignedUrlResp{
		ConfigID:  config.ID,
		UploadURL: presignedUrl,
		URL:       client.GetURL(path),
		Path:      path,
	}, nil
}

func (s *FileService) CreateFileCallback(ctx context.Context, req *req.FileCreateReq) (int64, error) {
	// 验证配置是否存在
	_, err := s.fileConfigService.GetFileConfig(ctx, req.ConfigID)
	if err != nil {
		return 0, errors.New("配置不存在")
	}

	fileRecord := &model.InfraFile{
		ConfigId: req.ConfigID,
		Name:     req.Name,
		Path:     req.Path,
		Url:      req.URL,
		Type:     req.Type,
		Size:     req.Size,
	}
	err = s.q.InfraFile.WithContext(ctx).Create(fileRecord)
	if err != nil {
		return 0, err
	}
	return fileRecord.ID, nil
}
