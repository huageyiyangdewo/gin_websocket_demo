package service

import (
	"chat/model"
	"chat/serializer"

	logging "github.com/sirupsen/logrus"

)

type UserRegisterService struct {
	Name string `json:"name"`
	Pwd string `json:"pwd"`
}

func (u *UserRegisterService) Register() serializer.Response {
	var user model.User
	var count int
	model.DB.Model(&user).Where("name=?", u.Name).Count(&count)
	if count == 1 {
		return serializer.Response{
			Code: 200,
			Msg: "user exist",
		}
	}

	user = model.User{
		Name: u.Name,
	}

	if err := user.SetPwd(u.Pwd); err != nil {
		return serializer.Response{
			Code: 400,
			Msg: "pwd bcrypt failed",
		}
	}

	if err := model.DB.Create(&user).Error; err != nil {
		logging.Errorf("create user failed, err: %s \n", err  )
		return serializer.Response{
			Code: 400,
			Msg: "create user failed",
		}
	}

	return serializer.Response{
		Code: 200,
		Msg: "create user success",
	}
}
