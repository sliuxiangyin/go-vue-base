package relationship

import (
	"databaseAi/internal/app/tools/models"
)

// RelationshipReq 关系模型请求结构体
type RelationshipReq struct {
	ID        uint   `json:"id"`
	AccountID uint   `json:"account_id" form:"account_id" binding:"required"`
	Title     string `json:"title" form:"title" binding:"required"`
	Tables    string `json:"tables" form:"tables"`
	Ddls      string `json:"ddls" form:"ddls"`
	Result    string `json:"result" form:"result"`
}

func (receiver *RelationshipReq) ToModel() *models.Relationship {
	return &models.Relationship{
		ID:        receiver.ID,
		AccountID: receiver.AccountID,
		Title:     receiver.Title,
		Tables:    receiver.Tables,
		Ddls:      receiver.Ddls,
		Result:    receiver.Result,
	}
}
