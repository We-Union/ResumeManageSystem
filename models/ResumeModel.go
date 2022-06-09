package models

import (
	"ResumeMamageSystem/dao"
	"time"
)

type ResumeModel struct {
	ID        int    `json:"id,omitempty" gorm:"primary_key"`
	Name      string `json:"name,omitempty"  validate:"required"`
	Target    string `json:"target"`
	File      string `json:"file"`
	OwnerID   int
	Owner     UserModel `json:"owner,omitempty" gorm:"Foreignkey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"autoCreateTime,omitempty"`
}

func CreateResume(resume *ResumeModel) (err error) {
	err = dao.DB.Create(&resume).Error
	return err
}

func GetResumeByID(groupID int) (resume *ResumeModel, err error) {
	resume = new(ResumeModel)
	err = dao.DB.Where("id = ?", groupID).First(&resume).Error
	if err != nil {
		return nil, err
	}
	err = dao.DB.Model(&resume).Select("id,name").Association("Owner").Find(&resume.Owner)
	if err != nil {
		return nil, err
	}
	return resume, nil
}

func GetResumesByOwnerID(ownerID int, start int, end int) (resumes *[]ResumeModel, err error) {
	resumes = new([]ResumeModel)
	if end-start > 100 {
		end = start + 100
	}
	err = dao.DB.Where("owner_id = ?", ownerID).Order("created_at desc").Offset(start).Limit(end - start).Find(&resumes).Error
	if err != nil {
		return nil, err
	}
	return resumes, nil
}

func GetResumeNumByOwnerID(ownerID int) (num int64, err error) {
	err = dao.DB.Model(&ResumeModel{}).Where("owner_id = ?", ownerID).Count(&num).Error
	if err != nil {
		return 0, err
	}
	return num, nil
}
func DeleteResumeByID(groupID int) (err error) {
	err = dao.DB.Where("id=?", groupID).Delete(&ResumeModel{}).Error
	return
}

func UpdateResume(resume *ResumeModel) (err error) {
	err = dao.DB.Save(resume).Error
	return err
}

func ValidateResume(resume *ResumeModel) (message string) {
	if resume.Name == "" {
		return "简历名称不能为空"
	}
	if resume.Target == "" {
		return "投递意向不能为空"
	}
	return ""
}
