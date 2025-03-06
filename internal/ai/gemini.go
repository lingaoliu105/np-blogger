package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiService Gemini API服务封装
type GeminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
	ctx    context.Context
}

// NewGeminiService 创建新的Gemini服务实例
func NewGeminiService(apiKey string) (*GeminiService, error) {
	ctx := context.Background()

	// 初始化Gemini客户端
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	// 获取生成模型
	model := client.GenerativeModel("gemini-pro")

	return &GeminiService{
		client: client,
		model:  model,
		ctx:    ctx,
	}, nil
}

// Close 关闭连接
func (s *GeminiService) Close() {
	s.client.Close()
}

// GenerateEmbedding 生成文本嵌入向量
func (s *GeminiService) GenerateEmbedding(text string) ([]float32, error) {
	// 使用模型生成文本嵌入
	embeddingModel := s.client.EmbeddingModel("embedding-001")
	result, err := embeddingModel.EmbedContent(s.ctx, genai.Text(text))
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	return result.Embedding, nil
}

// GenerateBlogPost 生成博客文章
func (s *GeminiService) GenerateBlogPost(topic string, references []string) (string, error) {
	// 构建提示信息
	prompt := fmt.Sprintf(
		"请根据以下主题和参考内容生成一篇技术博客文章：\n\n主题：%s\n\n参考内容：\n%s",
		topic,
		strings.Join(references, "\n"),
	)

	// 生成文章内容
	response, err := s.model.GenerateContent(s.ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate blog post: %w", err)
	}

	var content strings.Builder
	for _, part := range response.Candidates[0].Content.Parts {
		content.WriteString(fmt.Sprint(part))
	}

	return content.String(), nil
}