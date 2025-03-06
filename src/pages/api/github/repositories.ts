import { NextApiRequest, NextApiResponse } from 'next';
import { getSession } from 'next-auth/react';
import { Octokit } from '@octokit/rest';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  try {
    const session = await getSession({ req });
    if (!session) {
      return res.status(401).json({ error: 'Unauthorized' });
    }

    const accessToken = session.accessToken as string;
    const octokit = new Octokit({ auth: accessToken });

    // 获取用户的GitHub仓库列表
    const { data: repos } = await octokit.repos.listForAuthenticatedUser({
      visibility: 'all',
      sort: 'updated',
      per_page: 100
    });

    // 转换为前端需要的格式
    const repositories = repos.map(repo => ({
      id: repo.id,
      name: repo.name,
      full_name: repo.full_name,
      description: repo.description,
      synced: false // 默认未同步状态
    }));

    return res.status(200).json(repositories);
  } catch (error) {
    console.error('Failed to fetch repositories:', error);
    return res.status(500).json({ error: 'Internal server error' });
  }
}