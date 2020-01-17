package model

import "time"

type TBMoney struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Uid       string     `gorm:"column:uid;type:varchar(100);unique_index" description:"uid"`
	Coin      int32      `gorm:"column:coin; default 0" description:"铜币"`
}

func (t *TBMoney) TableName() string {
	return "tb_money"
}

func (t *TBMoney) Create() error {
	return DB.Create(t).Error
}

func GetMoney(u *TBMoney) (*TBMoney, error) {
	var money TBMoney
	err := DB.Where(u).First(&money).Error
	return &money, err
}

func DeleteMoney(id uint) error {
	money := TBMoney{
		ID: id,
	}
	return DB.Delete(&money).Error
}

func (t *TBMoney) Update() error {
	return DB.Save(t).Error
}
