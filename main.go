package main

import (
	"ResumeMamageSystem/dao"
	"ResumeMamageSystem/models"
	"ResumeMamageSystem/routers"
	"ResumeMamageSystem/setting"
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	// 加载配置文件
	if err := setting.Init("conf/config.ini"); err != nil {
		fmt.Printf("load config from file failed, err:%v\n", err)
		return
	}

	err := dao.InitMySQL(setting.Conf.MySQLConfig)
	if err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}

	err = dao.DB.AutoMigrate(&models.UserModel{}, &models.RewardModel{}, &models.ResumeModel{})
	if err != nil {
		fmt.Printf("migrate failed, err:%v\n", err)
		return
	}

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	err = os.MkdirAll(filepath.Join(dir, "resumes"), os.ModePerm)
	if err != nil {
		fmt.Printf("创建简历文件夹失败")
		return
	}
	err = os.MkdirAll(filepath.Join(dir, "rewards"), os.ModePerm)
	if err != nil {
		fmt.Printf("创建文件夹失败")
		return
	}
	// 注册路由
	r := routers.SetupRouter()

	if err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port)); err != nil {
		fmt.Printf("server startup failed, err:%v\n", err)
	}
}
