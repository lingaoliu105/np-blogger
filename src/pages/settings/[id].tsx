import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { Box, Container, Typography, TextField, FormControlLabel, Switch, Button, CircularProgress, Alert } from '@mui/material';
import axios from 'axios';

interface RepositorySettings {
  id: number;
  name: string;
  syncEnabled: boolean;
  syncInterval: number;
  targetPlatforms: {
    juejin: boolean;
    csdn: boolean;
    zhihu: boolean;
  };
}

export default function RepositorySettings() {
  const router = useRouter();
  const { id } = router.query;

  const [settings, setSettings] = useState<RepositorySettings | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    if (!id) return;

    // 获取仓库设置
    axios.get(`/api/settings/repository/${id}`)
      .then(response => {
        setSettings(response.data);
      })
      .catch(error => {
        console.error('Failed to fetch repository settings:', error);
        setError('获取仓库设置失败');
      })
      .finally(() => {
        setLoading(false);
      });
  }, [id]);

  const handleSave = async () => {
    if (!settings) return;

    setSaving(true);
    setError('');
    setSuccess(false);

    try {
      await axios.put(`/api/settings/repository/${id}`, settings);
      setSuccess(true);
    } catch (error) {
      console.error('Failed to save repository settings:', error);
      setError('保存设置失败');
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh">
        <CircularProgress />
      </Box>
    );
  }

  if (!settings) {
    return (
      <Container maxWidth="md">
        <Box sx={{ my: 4 }}>
          <Alert severity="error">找不到仓库设置</Alert>
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="md">
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {settings.name} - 同步设置
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {success && (
          <Alert severity="success" sx={{ mb: 2 }}>
            设置已保存
          </Alert>
        )}

        <Box component="form" sx={{ mt: 3 }}>
          <FormControlLabel
            control={
              <Switch
                checked={settings.syncEnabled}
                onChange={(e) => setSettings(prev => prev ? {
                  ...prev,
                  syncEnabled: e.target.checked
                } : null)}
              />
            }
            label="启用自动同步"
          />

          <TextField
            fullWidth
            label="同步间隔（分钟）"
            type="number"
            value={settings.syncInterval}
            onChange={(e) => setSettings(prev => prev ? {
              ...prev,
              syncInterval: parseInt(e.target.value)
            } : null)}
            sx={{ mt: 2 }}
          />

          <Typography variant="h6" sx={{ mt: 3, mb: 2 }}>
            目标平台
          </Typography>

          <FormControlLabel
            control={
              <Switch
                checked={settings.targetPlatforms.juejin}
                onChange={(e) => setSettings(prev => prev ? {
                  ...prev,
                  targetPlatforms: {
                    ...prev.targetPlatforms,
                    juejin: e.target.checked
                  }
                } : null)}
              />
            }
            label="掘金"
          />

          <FormControlLabel
            control={
              <Switch
                checked={settings.targetPlatforms.csdn}
                onChange={(e) => setSettings(prev => prev ? {
                  ...prev,
                  targetPlatforms: {
                    ...prev.targetPlatforms,
                    csdn: e.target.checked
                  }
                } : null)}
              />
            }
            label="CSDN"
          />

          <FormControlLabel
            control={
              <Switch
                checked={settings.targetPlatforms.zhihu}
                onChange={(e) => setSettings(prev => prev ? {
                  ...prev,
                  targetPlatforms: {
                    ...prev.targetPlatforms,
                    zhihu: e.target.checked
                  }
                } : null)}
              />
            }
            label="知乎"
          />

          <Box sx={{ mt: 4 }}>
            <Button
              variant="contained"
              color="primary"
              onClick={handleSave}
              disabled={saving}
            >
              {saving ? '保存中...' : '保存设置'}
            </Button>
          </Box>
        </Box>
      </Box>
    </Container>
  );
}