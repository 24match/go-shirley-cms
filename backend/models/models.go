package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Role     string `gorm:"default:'admin'" json:"role"`
}

type Image struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Filename        string         `gorm:"not null" json:"filename"`
	FilePath        string         `gorm:"not null" json:"file_path"`
	FileSize        int64          `json:"file_size"`
	Description     string         `json:"description"`
	LongDescription string         `json:"long_description"`
	Category        string         `json:"category"`
	SortOrder       int            `gorm:"default:0" json:"sort_order"`
}

type PageConfig struct {
	gorm.Model
	PageName   string `gorm:"uniqueIndex;not null" json:"pageName"`
	ConfigData string `gorm:"type:text" json:"configData"`
}

type ModuleConfig struct {
	gorm.Model
	ModuleName  string `gorm:"uniqueIndex;not null" json:"moduleName"`
	Enabled     bool   `gorm:"default:true" json:"enabled"`
	ZhTitle     string `json:"zhTitle"`
	EnTitle     string `json:"enTitle"`
	ZhSubtitle  string `json:"zhSubtitle"`
	EnSubtitle  string `json:"enSubtitle"`
	ZhContent   string `gorm:"type:text" json:"zhContent"`
	EnContent   string `gorm:"type:text" json:"enContent"`
	Title       string `json:"title"`    // 向后兼容
	Subtitle    string `json:"subtitle"` // 向后兼容
	Content     string `gorm:"type:text" json:"content"` // 向后兼容
	ImagePath   string `json:"imagePath"`
	SortOrder   int    `gorm:"default:0" json:"sortOrder"`
	ExtraData   string `gorm:"type:text" json:"extraData"`
	ZhDescription string `json:"zhDescription"`
	EnDescription string `json:"enDescription"`
	Description string `json:"description"` // 向后兼容
}

type ContentItem struct {
	gorm.Model
	Section     string `gorm:"not null" json:"section"`
	ZhTitle     string `json:"zhTitle"`
	EnTitle     string `json:"enTitle"`
	ZhDescription string `json:"zhDescription"`
	EnDescription string `json:"enDescription"`
	Title       string `json:"title"`       // 向后兼容
	Description string `json:"description"` // 向后兼容
	Icon        string `json:"icon"`
	ImagePath   string `json:"image_path"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
}

type LanguageText struct {
	gorm.Model
	Key         string `gorm:"uniqueIndex;not null" json:"key"`
	Module      string `gorm:"not null" json:"module"`
	EnText      string `gorm:"type:text" json:"enText"`
	ZhText      string `gorm:"type:text" json:"zhText"`
	Description string `json:"description"`
	Version     int    `gorm:"default:1" json:"version"`
}

type LanguageTextVersion struct {
	gorm.Model
	LanguageTextID uint      `gorm:"not null" json:"languageTextId"`
	Key            string    `gorm:"not null" json:"key"`
	Module         string    `gorm:"not null" json:"module"`
	EnText         string    `gorm:"type:text" json:"enText"`
	ZhText         string    `gorm:"type:text" json:"zhText"`
	Description    string    `json:"description"`
	Version        int       `gorm:"not null" json:"version"`
	UpdatedAt      time.Time `json:"updatedAt"`
}