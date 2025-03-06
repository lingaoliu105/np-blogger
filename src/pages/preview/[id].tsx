import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { Box, Container, Typography, Paper, CircularProgress, Button, TextField } from '@mui/material';
import { Edit as EditIcon, Save as SaveIcon } from '@mui/icons-material';
import axios from 'axios';

interface BlogPost {
  id: string;
  title: string;
  content: string;
  createdAt: string;
}

export default function BlogPreview() {
  const router = useRouter();
  const { id } = router.query;

  const [post, setPost] = useState<BlogPost | null>(null);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const [editedContent, setEditedContent] = useState('');

  useEffect(() => {
    if (id) {
      fetchBlogPost();
    }
  }, [id]);

  const fetchBlogPost = async () => {
    try {
      const response = await axios.get(`/api/posts/${id}`);
      setPost(response.data);
      setEditedContent(response.data.content);
    } catch (error) {
      console.error('Failed to fetch blog post:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = () => {
    setEditing(true);
  };

  const handleSave = async () => {
    try {
      await axios.put(`/api/posts/${id}`, {
        content: editedContent
      });
      setPost(prev => prev ? { ...prev, content: editedContent } : null);
      setEditing(false);
    } catch (error) {
      console.error('Failed to save blog post:', error);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh">
        <CircularProgress />
      </Box>
    );
  }

  if (!post) {
    return (
      <Container maxWidth="md">
        <Box sx={{ my: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom>
            博客未找到
          </Typography>
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="md">
      <Box sx={{ my: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1">
            {post.title}
          </Typography>
          <Button
            variant="contained"
            color={editing ? 'success' : 'primary'}
            startIcon={editing ? <SaveIcon /> : <EditIcon />}
            onClick={editing ? handleSave : handleEdit}
          >
            {editing ? '保存' : '编辑'}
          </Button>
        </Box>

        <Paper sx={{ p: 3 }}>
          {editing ? (
            <TextField
              fullWidth
              multiline
              minRows={10}
              value={editedContent}
              onChange={(e) => setEditedContent(e.target.value)}
              variant="outlined"
            />
          ) : (
            <Typography variant="body1" component="div">
              {post.content}
            </Typography>
          )}
        </Paper>

        <Box sx={{ mt: 2, textAlign: 'right' }}>
          <Typography variant="body2" color="text.secondary">
            创建时间：{new Date(post.createdAt).toLocaleString()}
          </Typography>
        </Box>
      </Box>
    </Container>
  );
}