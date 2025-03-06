import { NextApiRequest, NextApiResponse } from 'next';
import { getSession } from 'next-auth/react';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  try {
    const session = await getSession({ req });
    return res.status(200).json({
      loggedIn: !!session,
      user: session?.user || null
    });
  } catch (error) {
    console.error('Auth status check failed:', error);
    return res.status(500).json({ error: 'Internal server error' });
  }
}