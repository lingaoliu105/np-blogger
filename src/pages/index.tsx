import { useEffect, useState } from 'react';
import { Box, Container, Typography, Button, CircularProgress } from '@mui/material';
import axios from 'axios';

export default function Home() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // 检查用户登录状态
    axios.get('/api/auth/status')
      .then(response => {
        setIsLoggedIn(response.data.loggedIn);
      })
      .catch(error => {
        console.error('Failed to check auth status:', error);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  const handleLogin = () => {
    window.location.href = '/api/auth/github';
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
      <Box sx={{ my: 4, textAlign: 'center' }}>
        <Typography variant="h2" component="h1" gutterBottom>
          NP Blogger
        </Typography>
        <Typography variant="h5" component="h2" gutterBottom>
          将你的技术博客自动同步到多个平台
        </Typography>
        
        {!isLoggedIn ? (
          <Button
            variant="contained"
            color="primary"
            size="large"
            onClick={handleLogin}
            sx={{ mt: 2 }}
          >
            使用 GitHub 登录
          </Button>
        ) : (
          <Button
            variant="contained"
            color="primary"
            size="large"
            href="/dashboard"
            sx={{ mt: 2 }}
          >
            进入控制台
          </Button>
        )}
      </Box>
    </Container>
  );
}