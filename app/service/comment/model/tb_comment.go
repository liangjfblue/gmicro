package model

import "time"

type TBComment struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	ArticleId uint       `gorm:"column:article_id;not null" description:"文章id"`
	Comment   string     `gorm:"column:comment;type:varchar(200)" description:"评论内容"`
	FromId    string     `gorm:"column:from_id;type:varchar(200)" description:"评论者id"`
	ToId      string     `gorm:"column:to_id;type:varchar(200)" description:"被评论者id"`
}

func (t *TBComment) TableName() string {
	return "tb_comment"
}

func (t *TBComment) Create() error {
	return DB.Create(t).Error
}

func GetComment(u *TBComment) (*TBComment, error) {
	var comment TBComment
	err := DB.Where(u).First(&comment).Error
	return &comment, err
}

func DeleteComment(id uint) error {
	comment := TBComment{
		ID: id,
	}
	return DB.Delete(&comment).Error
}

func (t *TBComment) Update() error {
	return DB.Save(t).Error
}

func ListComments(articleId uint, offset, limit int32) ([]*TBComment, uint64, error) {
	var (
		err      error
		comments = make([]*TBComment, 0)
		count    uint64
	)

	offset, limit = CheckPageSize(offset, limit)

	err = DB.Model(&TBComment{}).Where("article_id = ?", articleId).Count(&count).Error
	err = DB.Where("article_id = ?", articleId).Offset(offset).Limit(limit).Order("id desc").Find(&comments).Error
	return comments, count, err
}
