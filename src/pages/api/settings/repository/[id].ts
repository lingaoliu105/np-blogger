import { NextApiRequest, NextApiResponse } from 'next';
import { getSession } from 'next-auth/react';
import prisma from '@/lib/prisma';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const session = await getSession({ req });
  if (!session) {
    return res.status(401).json({ error: 'Unauthorized' });
  }

  const { id } = req.query;
  const repoId = parseInt(id as string);

  switch (req.method) {
    case 'GET':
      try {
        const settings = await prisma.repositorySettings.findUnique({
          where: { id: repoId }
        });

        if (!settings) {
          // 如果设置不存在，返回默认设置
          return res.status(200).json({
            id: repoId,
            name: '',
            syncEnabled: false,
            autoGenerate: false,
            platformSettings: {
              juejin: false,
              csdn: false,
              zhihu: false
            },
            ragSettings: {
              enabled: false,
              maxReferences: 5
            }
          });
        }

        return res.status(200).json(settings);
      } catch (error) {
        console.error('Failed to fetch repository settings:', error);
        return res.status(500).json({ error: 'Internal server error' });
      }

    case 'PUT':
      try {
        const settings = await prisma.repositorySettings.upsert({
          where: { id: repoId },
          update: req.body,
          create: {
            id: repoId,
            ...req.body
          }
        });

        return res.status(200).json(settings);
      } catch (error) {
        console.error('Failed to update repository settings:', error);
        return res.status(500).json({ error: 'Internal server error' });
      }

    default:
      return res.status(405).json({ error: 'Method not allowed' });
  }
}