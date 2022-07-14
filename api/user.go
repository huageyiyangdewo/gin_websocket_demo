package api

import (
	"chat/service"
	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func UserRegister(c *gin.Context)  {
	var user service.UserRegisterService
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, "register error")
		logging.Infof("register error, %s \n", err)
	} else {
		res := user.Register()
		c.JSON(200, res)
	}
}
