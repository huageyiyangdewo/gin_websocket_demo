package model

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name string `json:"name"`
	Pwd string `json:"pwd"`
}

const (
	PwdCost = 12 // 密码加密难度
)

func (u *User) SetPwd(pwd string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), PwdCost)
	if err != nil {
		return err
	}

	u.Pwd = string(bytes)
	return nil
}

func (u *User) CheckPwd(pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Pwd), []byte(pwd))
	return err == nil
}


