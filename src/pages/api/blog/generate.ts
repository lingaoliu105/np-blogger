import { NextApiRequest, NextApiResponse } from 'next';
import { RAGService } from '../../../internal/rag/rag';
import { GeminiService } from '../../../internal/ai/gemini';
import { config } from '../../../internal/config/config';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  const user = req.session.user;
  if (!user) {
    return res.status(401).json({ error: 'Unauthorized' });
  }

  const { topic, content } = req.body;
  if (!topic || !content) {
    return res.status(400).json({ error: 'Topic and content are required' });
  }

  try {
    // 初始化服务
    const ragService = await RAGService.NewRAGService(config.RAG);
    const geminiService = await GeminiService.NewGeminiService(config.RAG.GeminiKey);

    // 生成文本嵌入
    const embedding = await geminiService.GenerateEmbedding(content);

    // 存储到向量数据库
    await ragService.StoreEmbedding('blog_posts', content, embedding);

    // 搜索相似内容
    const similarTexts = await ragService.SearchSimilar('blog_posts', embedding, 5);

    // 生成优化后的博客文章
    const blogPost = await geminiService.GenerateBlogPost(topic, similarTexts);

    // 关闭服务连接
    await ragService.Close();
    geminiService.Close();

    res.status(200).json({
      topic,
      content: blogPost
    });
  } catch (error) {
    console.error('Failed to generate blog post:', error);
    res.status(500).json({ error: 'Failed to generate blog post' });
  }
}