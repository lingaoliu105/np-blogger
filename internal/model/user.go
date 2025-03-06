package model

import "time"

// User 用户模型
type User struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	GitHubID    int64  `gorm:"uniqueIndex"`
	Username    string `gorm:"size:255"`
	Email       string `gorm:"size:255"`
	AvatarURL   string `gorm:"size:255"`
	AccessToken string `gorm:"size:255"`

	Repositories []Repository `gorm:"foreignKey:UserID"`
	Blogs        []Blog       `gorm:"foreignKey:UserID"`
}

// Repository 仓库模型
type Repository struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID      uint   `gorm:"index"`
	GitHubID    int64  `gorm:"uniqueIndex"`
	Name        string `gorm:"size:255"`
	FullName    string `gorm:"size:255"`
	Description string
	Branch      string `gorm:"size:255;default:'main'"`
	BlogPath    string `gorm:"size:255;default:'content/posts'"`

	Blogs []Blog `gorm:"foreignKey:RepositoryID"`
}

// Blog 博客文章模型
type Blog struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID       uint   `gorm:"index"`
	RepositoryID uint   `gorm:"index"`
	Title        string `gorm:"size:255"`
	Slug         string `gorm:"size:255;uniqueIndex"`
	Content      string `gorm:"type:text"`
	CommitSHA    string `gorm:"size:40"`
	Status       string `gorm:"size:20;default:'draft'"` // draft, published
	PublishedAt  *time.Time
}