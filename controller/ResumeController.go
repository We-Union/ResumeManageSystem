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

func CreateResume(c *gin.Context) {

	uidInt := utils.GetUidInt(c)
	if uidInt == -1 {
		return
	}

	var resume models.ResumeModel

	err := c.BindJSON(&resume)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}
	validate := models.ValidateResume(&resume)
	if validate != "" {
		c.JSON(http.StatusOK, gin.H{"code": 2002, "msg": validate})
		return
	}

	resume.OwnerID = uidInt
	err = models.CreateResume(&resume)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{"id": resume.ID},
		})
		return
	}
}
func UploadResume(c *gin.Context) {
	uidInt := utils.GetUidInt(c)
	if uidInt == -1 {
		return
	}

	resumeId := c.Query("id")
	if resumeId == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	idInt, _ := strconv.Atoi(resumeId)
	resume, err := models.GetResumeByID(idInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的简历不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}

	if resume.OwnerID != uidInt {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有权限执行操作"})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}

	if !utils.IsExist(filepath.ToSlash(filepath.Join(dir, "resumes", strconv.Itoa(uidInt)))) {
		err = os.MkdirAll(filepath.ToSlash(filepath.Join(dir, "resumes", strconv.Itoa(uidInt))), os.ModePerm)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": "创建上传文件夹失败"})
			return
		}
	}
	resume.File = path.Join("resumes", strconv.Itoa(uidInt), strconv.Itoa(resume.ID)+path.Ext(file.Filename))
	fullPath := filepath.ToSlash(filepath.Join(dir, resume.File))
	err = c.SaveUploadedFile(file, fullPath)
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}
func GetResume(c *gin.Context) {
	uidInt := utils.GetUidInt(c)
	if uidInt == -1 {
		return
	}
	resumeId := c.Query("id")
	if resumeId == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	idInt, _ := strconv.Atoi(resumeId)
	resume, err := models.GetResumeByID(idInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的简历不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}

	if resume.OwnerID != uidInt {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有权限查看"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": resume})
	return
}

func DeleteResume(c *gin.Context) {
	uidInt := utils.GetUidInt(c)
	if uidInt == -1 {
		return
	}

	id := c.Query("id")

	if id == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	resumeId, _ := strconv.Atoi(id)
	resume, err := models.GetResumeByID(resumeId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的简历不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if resume.OwnerID != uidInt {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有权限执行该操作"})
		return
	}
	err = models.DeleteResumeByID(resume.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}

func GetMyResume(c *gin.Context) {

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
	resumes, err := models.GetResumesByOwnerID(ownerIDint, startInt, endInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	num, err := models.GetResumeNumByOwnerID(ownerIDint)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"totalNum": num, "items": resumes}})
	return
}

func UpdateResume(c *gin.Context) {
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
	resume, err := models.GetResumeByID(idInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的小组不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if resume.OwnerID != uidInt {
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
		resume.Name = name.(string)
	}

	target, exist := jsonData["target"]
	if exist {
		resume.Target = target.(string)
	}

	err = models.UpdateResume(resume)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}

func DownloadResume(c *gin.Context) {
	uidInt := utils.GetUidInt(c)
	if uidInt == -1 {
		return
	}

	resumeId := c.Query("id")
	if resumeId == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	idInt, _ := strconv.Atoi(resumeId)
	resume, err := models.GetResumeByID(idInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的简历不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}

	if resume.OwnerID != uidInt {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有权限查看"})
		return
	}
	if resume.File == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "该简历没有上传文件"})
		return
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	fullPath := path.Join(dir, resume.File)
	if !utils.IsExist(fullPath) {
		c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "文件不存在，这可能是系统错误，请联系管理员"})
		return
	}
	_, fileName := filepath.Split(fullPath)
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
	c.File(fullPath)
	return
}
