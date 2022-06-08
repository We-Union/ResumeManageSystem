package models

import (
	"ResumeMamageSystem/dao"
	"ResumeMamageSystem/utils"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"time"
)

type UserModel struct {
	ID         int       `json:"id,omitempty" gorm:"primaryKey"`
	Username   string    `json:"username,omitempty"  gorm:"type:varchar(30);uniqueIndex"`
	Name       string    `json:"name,omitempty"`
	Password   string    `json:"-"  validate:"-"`
	Email      string    `json:"email,omitempty"`
	Status     int       `json:"status,omitempty"`
	Salt       string    `json:"-"`
	Validation string    `json:"-"`
	ResetCnt   int       `json:"-"`
	CreatedAt  time.Time `json:"created_at,omitempty" gorm:"autoCreateTime,omitempty"`
	LastLogin  time.Time `json:"last_login,omitempty"`
}

func CreateUser(user *UserModel) (err error) {
	err = dao.DB.Create(&user).Error
	user.Validation = ""
	user.ResetCnt = 0
	if err != nil {
		return err
	}
	return err
}
func GetUserByUserName(username string) (user *UserModel, err error) {
	user = new(UserModel)
	if err = dao.DB.Where("username = ?", username).First(user).Error; err != nil {
		return nil, err
	}
	return
}

func GetUserById(id int) (user *UserModel, err error) {
	user = new(UserModel)
	if err = dao.DB.Where("id=?", id).First(user).Error; err != nil {
		return nil, err
	}
	return
}
func GetEncryptPassword(passwd string, user *UserModel) string {
	if user.Salt == "" {
		user.Salt = utils.GetToken(32)
	}
	m := sha256.New()
	m.Write([]byte(passwd + user.Salt))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}
func SetUserPasswd(user *UserModel, passwd string) {
	user.Password = GetEncryptPassword(passwd, user)
}
func CheckUserPasswd(user *UserModel, passwd string) bool {
	return GetEncryptPassword(passwd, user) == user.Password
}
func UpdateUser(user *UserModel) (err error) {
	err = dao.DB.Save(user).Error
	return err
}

func ChangePassword(user *UserModel, passwd string) (err error) {
	user.Password = GetEncryptPassword(passwd, user)
	err = dao.DB.Save(user).Error
	return err
}
func ResetPassword(user *UserModel, passwd string) (err error) {
	user.Salt = ""
	user.Password = GetEncryptPassword(passwd, user)
	user.Validation = ""
	user.ResetCnt = 0
	err = dao.DB.Save(user).Error
	return err
}
func ResetUsername(user *UserModel, username string) (err error) {
	user.Username = username
	user.Validation = ""
	user.ResetCnt = 0
	err = dao.DB.Save(user).Error
	return err
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func ValidateUser(user *UserModel) (message string) {
	match, _ := regexp.MatchString("^[a-zA-Z0-9]{5,15}$", user.Username)
	if !match {
		return "学号格式错误"
	}

	if len(user.Password) != 64 {
		return "密码格式错误"
	}
	match, _ = regexp.MatchString("^[\u4e00-\u9fa5·]{2,18}$", user.Name)
	if !match {
		return "姓名格式错误"
	}
	match, _ = regexp.MatchString("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$", user.Email)
	if !match {
		return "邮箱格式错误"
	}
	return ""
}
