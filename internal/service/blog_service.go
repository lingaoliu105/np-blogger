package service

import (
	"fmt"
	"np-blogger/internal/ai"
	"np-blogger/internal/rag"
	"np-blogger/internal/repository"
	"strings"
)

// BlogService 博客服务
type BlogService struct {
	gemini       *ai.GeminiService
	rag          *rag.RAGService
	githubClient *repository.GitHubClient
}

// NewBlogService 创建新的博客服务实例
func NewBlogService(
	gemini *ai.GeminiService,
	rag *rag.RAGService,
	githubClient *repository.GitHubClient,
) *BlogService {
	return &BlogService{
		gemini:       gemini,
		rag:          rag,
		githubClient: githubClient,
	}
}

// GenerateAndPublishBlog 生成并发布博客
func (s *BlogService) GenerateAndPublishBlog(userID uint, commitMessage string) error {
	// 1. 从提交信息中提取主题
	topic := s.extractTopic(commitMessage)

	// 2. 生成文本嵌入
	embedding, err := s.gemini.GenerateEmbedding(topic)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// 3. 搜索相关内容
	references, err := s.rag.SearchSimilar("blog_contents", embedding, 5)
	if err != nil {
		return fmt.Errorf("failed to search similar content: %w", err)
	}

	// 4. 生成博客内容
	content, err := s.gemini.GenerateBlogPost(topic, references)
	if err != nil {
		return fmt.Errorf("failed to generate blog post: %w", err)
	}

	// 5. 发布到GitHub
	fileName := s.generateFileName(topic)
	err = s.githubClient.CreateFile(
		"owner", // TODO: 从配置获取
		"repo",  // TODO: 从配置获取
		fileName,
		content,
		fmt.Sprintf("Add blog post: %s", topic),
	)
	if err != nil {
		return fmt.Errorf("failed to publish blog post: %w", err)
	}

	// 6. 存储文本向量用于后续检索
	err = s.rag.StoreEmbedding("blog_contents", content, embedding)
	if err != nil {
		return fmt.Errorf("failed to store embedding: %w", err)
	}

	return nil
}

// extractTopic 从提交信息中提取主题
func (s *BlogService) extractTopic(commitMessage string) string {
	// TODO: 实现更复杂的主题提取逻辑
	return commitMessage
}

// generateFileName 生成博客文件名
func (s *BlogService) generateFileName(topic string) string {
	// 将主题转换为URL友好的格式
	slug := strings.ToLower(topic)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	return fmt.Sprintf("posts/%s.md", slug)
}