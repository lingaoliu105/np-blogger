package config

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	GitHub   GitHubConfig   `mapstructure:"github"`
	RAG      RAGConfig      `mapstructure:"rag"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

// GitHubConfig GitHub相关配置
type GitHubConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

// RAGConfig RAG模块配置
type RAGConfig struct {
	MilvusHost string `mapstructure:"milvus_host"`
	MilvusPort int    `mapstructure:"milvus_port"`
	GeminiKey  string `mapstructure:"gemini_key"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host: "localhost",
			Port: 5432,
			User: "postgres",
			DBName: "np_blogger",
		},
		GitHub: GitHubConfig{
			RedirectURL: "http://localhost:8080/auth/callback",
		},
		RAG: RAGConfig{
			MilvusHost: "localhost",
			MilvusPort: 19530,
		},
	}
}