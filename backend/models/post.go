package models

import (
	"time"

	"gorm.io/gorm"
)

// Post 帖子主体结构
type Post struct {
	gorm.Model
	Title   string      `json:"title" gorm:"size:100;not null"`    // 帖子标题
	Content string      `json:"content" gorm:"type:text;not null"` // 帖子内容
	Author  string      `json:"author" gorm:"size:50;not null"`    // 作者
	Images  []PostImage `json:"images" gorm:"foreignKey:PostID"`   // 帖子图片
}

// PostImage 帖子图片结构
type PostImage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PostID    uint      `json:"post_id" gorm:"index;not null"`         // 关联的帖子ID
	URL       string    `json:"url" gorm:"type:varchar(255);not null"` // 图片URL
	CreatedAt time.Time `json:"created_at"`                            // 创建时间
}

// TableName 指定帖子图片表名
func (PostImage) TableName() string {
	return "post_images"
}

// BeforeDelete 删除帖子时同时删除相关的图片记录
func (p *Post) BeforeDelete(tx *gorm.DB) error {
	return tx.Delete(&PostImage{}, "post_id = ?", p.ID).Error
}

// CreatePostWithImages 创建帖子并添加图片
func CreatePostWithImages(post *Post, imageURLs []string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 创建帖子
		if err := tx.Create(post).Error; err != nil {
			return err
		}

		// 如果有图片，创建图片记录
		if len(imageURLs) > 0 {
			images := make([]PostImage, len(imageURLs))
			for i, url := range imageURLs {
				images[i] = PostImage{
					PostID: post.ID,
					URL:    url,
				}
			}
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetPostWithImages 获取帖子及其图片
func GetPostWithImages(postID uint) (*Post, error) {
	var post Post
	err := DB.Preload("Images").First(&post, postID).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// UpdatePostImages 更新帖子的图片
func UpdatePostImages(postID uint, newImageURLs []string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 删除原有的图片记录
		if err := tx.Delete(&PostImage{}, "post_id = ?", postID).Error; err != nil {
			return err
		}

		// 添加新的图片记录
		if len(newImageURLs) > 0 {
			images := make([]PostImage, len(newImageURLs))
			for i, url := range newImageURLs {
				images[i] = PostImage{
					PostID: postID,
					URL:    url,
				}
			}
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
