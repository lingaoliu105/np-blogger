package rag

import (
	"context"
	"fmt"
	"np-blogger/internal/config"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

// RAGService RAG服务封装
type RAGService struct {
	milvusClient client.Client
	geminiKey    string
	ctx          context.Context
}

// NewRAGService 创建新的RAG服务实例
func NewRAGService(cfg *config.RAGConfig) (*RAGService, error) {
	ctx := context.Background()

	// 连接Milvus
	milvusClient, err := client.NewGrpcClient(
		ctx,
		fmt.Sprintf("%s:%d", cfg.MilvusHost, cfg.MilvusPort),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to milvus: %w", err)
	}

	return &RAGService{
		milvusClient: milvusClient,
		geminiKey:    cfg.GeminiKey,
		ctx:          ctx,
	}, nil
}

// Close 关闭连接
func (s *RAGService) Close() error {
	return s.milvusClient.Close()
}

// StoreEmbedding 存储文本向量
func (s *RAGService) StoreEmbedding(collectionName string, text string, embedding []float32) error {
	// 确保集合存在
	exists, err := s.milvusClient.HasCollection(s.ctx, collectionName)
	if err != nil {
		return fmt.Errorf("failed to check collection: %w", err)
	}

	if !exists {
		// 创建集合
		schema := &entity.Schema{
			CollectionName: collectionName,
			Description:   "Text embeddings collection",
			Fields: []*entity.Field{
				{
					Name:       "id",
					DataType:   entity.FieldTypeInt64,
					PrimaryKey: true,
					AutoID:     true,
				},
				{
					Name:     "text",
					DataType: entity.FieldTypeString,
				},
				{
					Name:     "embedding",
					DataType: entity.FieldTypeFloatVector,
					TypeParams: map[string]string{
						"dim": "768", // 根据实际embedding维度设置
					},
				},
			},
		}

		err = s.milvusClient.CreateCollection(s.ctx, schema, 1)
		if err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
	}

	// 插入数据
	columns := []*entity.Column{
		entity.NewColumnString("text", []string{text}),
		entity.NewColumnFloatVector("embedding", 768, [][]float32{embedding}),
	}

	_, err = s.milvusClient.Insert(s.ctx, collectionName, "", columns...)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	return nil
}

// SearchSimilar 搜索相似文本
func (s *RAGService) SearchSimilar(collectionName string, queryEmbedding []float32, topK int) ([]string, error) {
	// 加载集合
	err := s.milvusClient.LoadCollection(s.ctx, collectionName, false)
	if err != nil {
		return nil, fmt.Errorf("failed to load collection: %w", err)
	}
	defer s.milvusClient.ReleaseCollection(s.ctx, collectionName)

	// 执行搜索
	searchParams := entity.NewIndexFlatSearchParam()
	vectors := []entity.Vector{entity.FloatVector(queryEmbedding)}

	results, err := s.milvusClient.Search(
		s.ctx,
		collectionName,
		"",
		[]string{"text"},
		vectors,
		"embedding",
		entity.L2,
		topK,
		searchParams,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	// 提取结果
	var texts []string
	for _, result := range results {
		for i := 0; i < topK && i < len(result.IDs); i++ {
			text, ok := result.Fields["text"].(*entity.ColumnString)
			if !ok {
				continue
			}
			texts = append(texts, text.Data()[i])
		}
	}

	return texts, nil
}