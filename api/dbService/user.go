package dbService

import (
	"SecKill/data"
	"SecKill/model"
)

func GetUser(userName string) (model.User, error) {
	user := model.User{}
	operation := data.Db.Where("username = ?", userName).First(&user)
	return user, operation.Error
}