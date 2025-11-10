package dbaccount

import (
	"databaseAi/internal/app/tools/models"
)

// DBAccountReq 数据库账号模型
type DBAccountReq struct {
	ID       uint   ` json:"id"`
	Name     string ` json:"name" form:"name" binding:"required"`
	Account  string `json:"account" form:"account" binding:"required"`
	Dbname   string `json:"dbname" form:"dbname" binding:"required"`
	Host     string ` json:"host" form:"host" binding:"required"`
	Password string ` json:"password" form:"password" binding:"required"`
	Port     int    ` json:"port" form:"port" binding:"required"`
}

func (receiver *DBAccountReq) ToModel() *models.DBAccount {
	return &models.DBAccount{
		ID:       receiver.ID,
		Name:     receiver.Name,
		Dbname:   receiver.Dbname,
		Account:  receiver.Account,
		Host:     receiver.Host,
		Password: receiver.Password,
		Port:     receiver.Port,
	}
}
