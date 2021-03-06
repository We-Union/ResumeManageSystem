package controller

import (
	"ResumeMamageSystem/models"
	"ResumeMamageSystem/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func CreateReward(c *gin.Context) {

	uidInt := utils.GetUidInt(c)
	if uidInt == -1 {
		return
	}

	var reward models.RewardModel

	err := c.BindJSON(&reward)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}
	validate := models.ValidateReward(&reward)
	if validate != "" {
		c.JSON(http.StatusOK, gin.H{"code": 2002, "msg": validate})
		return
	}

	reward.OwnerID = uidInt
	err = models.CreateReward(&reward)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{"id": reward.ID},
		})
		return
	}
}
func UploadReward(c *gin.Context) {
	uidInt := utils.GetUidInt(c)
	if uidInt == -1 {
		return
	}

	rewardId := c.Query("id")
	if rewardId == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	idInt, _ := strconv.Atoi(rewardId)
	reward, err := models.GetRewardByID(idInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的奖项不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}

	if reward.OwnerID != uidInt {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有权限执行操作"})
		return
	}
	file, err := c.FormFile("file")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if !utils.IsExist(filepath.ToSlash(filepath.Join(dir, "rewards", strconv.Itoa(uidInt)))) {
		err = os.MkdirAll(filepath.ToSlash(filepath.Join(dir, "rewards", strconv.Itoa(uidInt))), os.ModePerm)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": "创建上传文件夹失败"})
			return
		}
	}
	reward.File = path.Join("rewards", strconv.Itoa(uidInt), strconv.Itoa(reward.ID)+path.Ext(file.Filename))
	err = models.UpdateReward(reward)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	full_path := filepath.ToSlash(filepath.Join(dir, reward.File))
	err = c.SaveUploadedFile(file, full_path)
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}
func GetReward(c *gin.Context) {
	uidInt := utils.GetUidInt(c)
	if uidInt == -1 {
		return
	}

	rewardId := c.Query("id")
	if rewardId == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	idInt, _ := strconv.Atoi(rewardId)
	reward, err := models.GetRewardByID(idInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的奖项不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}

	if reward.OwnerID != uidInt {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有权限查看"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": reward})
	return
}

func DeleteReward(c *gin.Context) {
	uidInt := utils.GetUidInt(c)
	if uidInt == -1 {
		return
	}
	id := c.Query("id")

	if id == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	rewardId, _ := strconv.Atoi(id)
	reward, err := models.GetRewardByID(rewardId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的奖项不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if reward.OwnerID != uidInt {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有权限执行该操作"})
		return
	}
	err = models.DeleteRewardByID(reward.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}

func GetMyReward(c *gin.Context) {
	uidInt := utils.GetUidInt(c)
	if uidInt == -1 {
		return
	}

	start := c.Query("start")
	end := c.Query("end")
	if start == "" || end == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}
	ownerIDint := uidInt
	startInt, _ := strconv.Atoi(start)
	endInt, _ := strconv.Atoi(end)
	fmt.Println(startInt, endInt)
	rewards, err := models.GetRewardsByOwnerID(ownerIDint, startInt, endInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	num, err := models.GetRewardNumByOwnerID(ownerIDint)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"totalNum": num, "items": rewards}})
	return
}

func UpdateReward(c *gin.Context) {
	uidInt := utils.GetUidInt(c)
	if uidInt == -1 {
		return
	}

	id := c.Query("id")

	if id == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	idInt, _ := strconv.Atoi(id)
	reward, err := models.GetRewardByID(idInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的小组不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if reward.OwnerID != uidInt {
		c.JSON(http.StatusOK, gin.H{"ode": 4003, "msg": "您没有权限执行该操作"})
		return
	}

	jsonData := make(map[string]interface{})
	err = c.BindJSON(&jsonData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	name, exist := jsonData["name"]
	if exist {
		reward.Name = name.(string)
	}
	rank, exist := jsonData["rank"]
	if exist {
		reward.Rank = rank.(string)
	}
	host, exist := jsonData["host"]
	if exist {
		reward.Host = host.(string)
	}
	date, exist := jsonData["date"]
	if exist {
		reward.Date = date.(string)
	}

	err = models.UpdateReward(reward)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}

func DownloadReward(c *gin.Context) {
	uidInt := utils.GetUidInt(c)
	if uidInt == -1 {
		return
	}
	rewardId := c.Query("id")
	if rewardId == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	idInt, _ := strconv.Atoi(rewardId)
	reward, err := models.GetRewardByID(idInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的奖项不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}

	if reward.OwnerID != uidInt {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有权限查看"})
		return
	}

	if reward.File == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "该奖项没有上传文件"})
		return
	}
	fmt.Println("a")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	fullPath := path.Join(dir, reward.File)
	fmt.Println(reward.File)
	fmt.Println(fullPath)
	if !utils.IsExist(fullPath) {
		c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "文件不存在，这可能是系统错误，请联系管理员"})
		return
	}
	_, fileName := filepath.Split(fullPath)
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
	c.File(fullPath)
	return
}
