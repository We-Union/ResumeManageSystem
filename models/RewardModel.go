package models

import (
	"ResumeMamageSystem/dao"
	"time"
)

type RewardModel struct {
	ID        int    `json:"id,omitempty" gorm:"primary_key"`
	Name      string `json:"name,omitempty"  validate:"required"`
	Rank      string `json:"rank"`
	Host      string `json:"host"`
	Date      string `json:"date"`
	File      string `json:"file"`
	OwnerID   int
	Owner     UserModel `json:"owner,omitempty" gorm:"Foreignkey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"autoCreateTime,omitempty"`
}

func CreateReward(reward *RewardModel) (err error) {
	err = dao.DB.Create(&reward).Error
	return err
}

func GetRewardByID(id int) (reward *RewardModel, err error) {
	reward = new(RewardModel)
	err = dao.DB.Where("id = ?", id).First(&reward).Error
	if err != nil {
		return nil, err
	}
	err = dao.DB.Model(&reward).Select("id,name").Association("Owner").Find(&reward.Owner)
	if err != nil {
		return nil, err
	}
	return reward, nil
}

func GetRewardsByOwnerID(ownerID int, start int, end int) (rewards *[]RewardModel, err error) {
	rewards = new([]RewardModel)
	if end-start > 100 {
		end = start + 100
	}
	err = dao.DB.Where("owner_id = ?", ownerID).Order("created_at desc").Offset(start).Limit(end - start).Find(&rewards).Error
	if err != nil {
		return nil, err
	}
	return rewards, nil
}

func GetRewardNumByOwnerID(ownerID int) (num int64, err error) {
	err = dao.DB.Model(&RewardModel{}).Where("owner_id = ?", ownerID).Count(&num).Error
	if err != nil {
		return 0, err
	}
	return num, nil
}
func DeleteRewardByID(id int) (err error) {
	err = dao.DB.Where("id=?", id).Delete(&RewardModel{}).Error
	return
}

func UpdateReward(reward *RewardModel) (err error) {
	err = dao.DB.Save(reward).Error
	return err
}

func ValidateReward(reward *RewardModel) (message string) {
	if reward.Date == "" {
		return "获奖日期不能为空"
	}
	if reward.Name == "" {
		return "名称不能为空"
	}
	return ""
}
