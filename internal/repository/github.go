package repository

import (
	"context"
	"fmt"
	"np-blogger/internal/database"
	"np-blogger/internal/model"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// GitHubClient GitHub客户端封装
type GitHubClient struct {
	client *github.Client
	ctx    context.Context
}

// NewGitHubClient 创建新的GitHub客户端
func NewGitHubClient(accessToken string) *GitHubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHubClient{
		client: client,
		ctx:    ctx,
	}
}

// ListRepositories 获取用户的仓库列表
func (c *GitHubClient) ListRepositories(userID uint) ([]model.Repository, error) {
	// 获取用户信息
	db := database.GetDB()
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 获取GitHub仓库列表
	opts := &github.RepositoryListOptions{
		Sort:      "updated",
		Direction: "desc",
	}
	repos, _, err := c.client.Repositories.List(c.ctx, "", opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	// 转换为内部模型
	var repositories []model.Repository
	for _, repo := range repos {
		repositories = append(repositories, model.Repository{
			UserID:      userID,
			GitHubID:    repo.GetID(),
			Name:        repo.GetName(),
			FullName:    repo.GetFullName(),
			Description: repo.GetDescription(),
		})
	}

	return repositories, nil
}

// GetRepository 获取仓库详情
func (c *GitHubClient) GetRepository(owner, name string) (*github.Repository, error) {
	repo, _, err := c.client.Repositories.Get(c.ctx, owner, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}
	return repo, nil
}

// CreateFile 在仓库中创建文件
func (c *GitHubClient) CreateFile(owner, repo, path, content, message string) error {
	// 创建或更新文件
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: []byte(content),
	}

	_, _, err := c.client.Repositories.CreateFile(c.ctx, owner, repo, path, opts)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	return nil
}

// UpdateFile 更新仓库中的文件
func (c *GitHubClient) UpdateFile(owner, repo, path, sha, content, message string) error {
	// 更新文件
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: []byte(content),
		SHA:     github.String(sha),
	}

	_, _, err := c.client.Repositories.UpdateFile(c.ctx, owner, repo, path, opts)
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}

	return nil
}

// GetFileContent 获取文件内容
func (c *GitHubClient) GetFileContent(owner, repo, path string) (string, string, error) {
	content, _, _, err := c.client.Repositories.GetContents(c.ctx, owner, repo, path, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to get file content: %w", err)
	}

	fileContent, err := content.GetContent()
	if err != nil {
		return "", "", fmt.Errorf("failed to decode content: %w", err)
	}

	return fileContent, content.GetSHA(), nil
}