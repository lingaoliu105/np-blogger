package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"np-blogger/internal/config"
	"np-blogger/internal/database"
	"np-blogger/internal/model"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	oauthConfig *oauth2.Config
)

// InitializeGitHub 初始化GitHub OAuth配置
func InitializeGitHub(cfg *config.GitHubConfig) {
	oauthConfig = &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"user",
			"repo",
		},
		Endpoint: github.Endpoint,
	}
}

// GetAuthURL 获取GitHub授权URL
func GetAuthURL() string {
	return oauthConfig.AuthCodeURL("state")
}

// HandleCallback 处理GitHub OAuth回调
func HandleCallback(code string) (*model.User, error) {
	// 使用授权码获取token
	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// 获取用户信息
	client := oauthConfig.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	// 解析用户信息
	var githubUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// 查找或创建用户
	db := database.GetDB()
	var user model.User
	result := db.Where("github_id = ?", githubUser.ID).First(&user)
	if result.Error != nil {
		// 创建新用户
		user = model.User{
			GitHubID:    githubUser.ID,
			Username:    githubUser.Login,
			Email:       githubUser.Email,
			AvatarURL:   githubUser.AvatarURL,
			AccessToken: token.AccessToken,
		}
		if err := db.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// 更新现有用户
		user.Username = githubUser.Login
		user.Email = githubUser.Email
		user.AvatarURL = githubUser.AvatarURL
		user.AccessToken = token.AccessToken
		if err := db.Save(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	return &user, nil
}