package model

import "time"

type TBArticleInfo struct {
	ID         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	Title      string     `gorm:"column:title;type:varchar(100)" description:"文章标题"`
	Topic      string     `gorm:"column:topic;type:varchar(100)" description:"文章主题"`
	Author     string     `gorm:"column:author;type:varchar(100)" description:"作者"`
	IsOriginal int8       `gorm:"column:is_original; default 0" description:"是否原创 0-否 1-是"`
}

func (t *TBArticleInfo) TableName() string {
	return "tb_article_info"
}

func (t *TBArticleInfo) Create() error {
	return DB.Create(t).Error
}

func GetArticleInfo(u *TBArticleInfo) (*TBArticleInfo, error) {
	var articleInfo TBArticleInfo
	err := DB.Where(u).First(&articleInfo).Error
	return &articleInfo, err
}

func DeleteArticleInfo(id uint) error {
	articleInfo := TBArticleInfo{
		ID: id,
	}
	return DB.Delete(&articleInfo).Error
}

func (t *TBArticleInfo) Update() error {
	return DB.Save(t).Error
}
