import { useState } from 'react';
import { Box, Container, Typography, Paper, TextField, Button, Snackbar, Alert } from '@mui/material';
import { Save as SaveIcon, Preview as PreviewIcon } from '@mui/icons-material';
import dynamic from 'next/dynamic';
import axios from 'axios';

// 动态导入Markdown编辑器，避免SSR问题
const MDEditor = dynamic(() => import('@uiw/react-md-editor'), { ssr: false });

export default function BlogEditor() {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  const handleSave = async () => {
    try {
      await axios.post('/api/posts', {
        title,
        content
      });
      setSnackbar({
        open: true,
        message: '博客保存成功',
        severity: 'success'
      });
    } catch (error) {
      setSnackbar({
        open: true,
        message: '保存失败，请重试',
        severity: 'error'
      });
    }
  };

  const handlePreview = () => {
    // 在新窗口中预览博客
    const previewWindow = window.open('', '_blank');
    if (previewWindow) {
      previewWindow.document.write(`
        <html>
          <head>
            <title>${title || '预览'}</title>
            <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/github-markdown-css/github-markdown.min.css">
            <style>
              .markdown-body {
                box-sizing: border-box;
                min-width: 200px;
                max-width: 980px;
                margin: 0 auto;
                padding: 45px;
              }
            </style>
          </head>
          <body class="markdown-body">
            <h1>${title}</h1>
            ${content}
          </body>
        </html>
      `);
    }
  };

  return (
    <Container maxWidth="lg">
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          编写博客
        </Typography>

        <Paper sx={{ p: 3, mb: 3 }}>
          <TextField
            fullWidth
            label="博客标题"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            sx={{ mb: 3 }}
          />

          <Box sx={{ mb: 3 }}>
            <MDEditor
              value={content}
              onChange={(value) => setContent(value || '')}
              height={400}
            />
          </Box>

          <Box sx={{ display: 'flex', gap: 2 }}>
            <Button
              variant="contained"
              color="primary"
              startIcon={<SaveIcon />}
              onClick={handleSave}
            >
              保存
            </Button>
            <Button
              variant="outlined"
              startIcon={<PreviewIcon />}
              onClick={handlePreview}
            >
              预览
            </Button>
          </Box>
        </Paper>
      </Box>

      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={() => setSnackbar({ ...snackbar, open: false })}
      >
        <Alert
          onClose={() => setSnackbar({ ...snackbar, open: false })}
          severity={snackbar.severity === 'success' ? 'success' : 'error'}
        >
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Container>
  );
}