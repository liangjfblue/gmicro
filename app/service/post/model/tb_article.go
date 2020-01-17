package model

import "time"

/**
后面上MongoDB
*/
type TBArticle struct {
	ID            uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
	Uid           string     `gorm:"column:uid;type:varchar(100)" description:"uid"`
	Content       string     `gorm:"column:content;type:text" description:"文章内容"`
	ArticleInfoId uint       `gorm:"column:article_info_id;not null" description:"文章信息id"`
}

func (t *TBArticle) TableName() string {
	return "tb_article"
}

func (t *TBArticle) Create() error {
	return DB.Create(t).Error
}

func GetArticle(u *TBArticle) (*TBArticle, error) {
	var article TBArticle
	err := DB.Where(u).First(&article).Error
	return &article, err
}

func DeleteArticle(id uint) error {
	article := TBArticle{
		ID: id,
	}
	return DB.Delete(&article).Error
}

func (t *TBArticle) Update() error {
	return DB.Save(t).Error
}
