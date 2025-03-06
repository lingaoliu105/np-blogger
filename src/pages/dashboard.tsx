import { useEffect, useState } from 'react';
import { Box, Container, Typography, List, ListItem, ListItemText, ListItemSecondaryAction, IconButton, CircularProgress, Button } from '@mui/material';
import { Sync as SyncIcon, Settings as SettingsIcon } from '@mui/icons-material';
import axios from 'axios';

interface Repository {
  id: number;
  name: string;
  full_name: string;
  description: string;
  synced: boolean;
}

export default function Dashboard() {
  const [repositories, setRepositories] = useState<Repository[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // 获取用户的GitHub仓库列表
    axios.get('/api/github/repositories')
      .then(response => {
        setRepositories(response.data);
      })
      .catch(error => {
        console.error('Failed to fetch repositories:', error);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  const handleSync = async (repoId: number) => {
    try {
      await axios.post(`/api/sync/repository/${repoId}`);
      // 更新仓库状态
      setRepositories(repos =>
        repos.map(repo =>
          repo.id === repoId ? { ...repo, synced: true } : repo
        )
      );
    } catch (error) {
      console.error('Failed to sync repository:', error);
    }
  };

  const handleSettings = (repoId: number) => {
    window.location.href = `/settings/${repoId}`;
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Container maxWidth="md">
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          我的仓库
        </Typography>

        <List>
          {repositories.map(repo => (
            <ListItem key={repo.id} divider>
              <ListItemText
                primary={repo.name}
                secondary={repo.description || '暂无描述'}
              />
              <ListItemSecondaryAction>
                <IconButton
                  edge="end"
                  aria-label="sync"
                  onClick={() => handleSync(repo.id)}
                  sx={{ mr: 1 }}
                >
                  <SyncIcon color={repo.synced ? 'success' : 'action'} />
                </IconButton>
                <IconButton
                  edge="end"
                  aria-label="settings"
                  onClick={() => handleSettings(repo.id)}
                >
                  <SettingsIcon />
                </IconButton>
              </ListItemSecondaryAction>
            </ListItem>
          ))}
        </List>

        {repositories.length === 0 && (
          <Box sx={{ textAlign: 'center', mt: 4 }}>
            <Typography variant="body1" color="text.secondary" gutterBottom>
              暂无可同步的仓库
            </Typography>
            <Button
              variant="contained"
              color="primary"
              href="https://github.com/new"
              target="_blank"
              sx={{ mt: 2 }}
            >
              创建新仓库
            </Button>
          </Box>
        )}
      </Box>
    </Container>
  );
}