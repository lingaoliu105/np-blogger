import { NextApiRequest, NextApiResponse } from 'next';
import axios from 'axios';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  const { code } = req.query;
  const clientId = process.env.GITHUB_CLIENT_ID;
  const clientSecret = process.env.GITHUB_CLIENT_SECRET;

  if (!code || !clientId || !clientSecret) {
    return res.status(400).json({ error: 'Missing required parameters' });
  }

  try {
    // 使用授权码获取访问令牌
    const tokenResponse = await axios.post(
      'https://github.com/login/oauth/access_token',
      {
        client_id: clientId,
        client_secret: clientSecret,
        code: code,
      },
      {
        headers: { Accept: 'application/json' },
      }
    );

    const { access_token } = tokenResponse.data;

    // 获取用户信息
    const userResponse = await axios.get('https://api.github.com/user', {
      headers: {
        Authorization: `Bearer ${access_token}`,
      },
    });

    // 创建用户会话
    req.session.user = {
      id: userResponse.data.id,
      login: userResponse.data.login,
      accessToken: access_token,
    };
    await req.session.save();

    // 重定向到仪表板
    res.redirect('/dashboard');
  } catch (error) {
    console.error('GitHub authentication error:', error);
    res.status(500).json({ error: 'Authentication failed' });
  }
}