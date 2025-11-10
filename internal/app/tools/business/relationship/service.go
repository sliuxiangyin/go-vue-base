package relationship

import (
	"databaseAi/internal/app/tools/business/dbaccount"
	"databaseAi/internal/app/tools/infra/response/api"
	"databaseAi/internal/app/tools/models"
	"databaseAi/internal/app/tools/repo"
	"databaseAi/internal/app/tools/repo/aliyun_bailian_relationship"
	"databaseAi/internal/utils"
	"encoding/json"
	"errors"
)

// Service 关系服务
type Service struct {
	repo                          *RelationshipRepo
	dbAccountRepo                 *dbaccount.DBAccountRepo
	manageDatabaseRepo            *repo.ManageDatabaseRepo
	aliYunBaiLianRelationshipRepo *aliyun_bailian_relationship.AliYunBaiLianRelationshipRepo
}

// NewService 创建新的关系服务实例
func NewService(repo *RelationshipRepo, manageDatabaseRepo *repo.ManageDatabaseRepo,
	dbAccountRepo *dbaccount.DBAccountRepo,
	aliYunBaiLianRelationshipRepo *aliyun_bailian_relationship.AliYunBaiLianRelationshipRepo,
) *Service {
	return &Service{
		repo:                          repo,
		manageDatabaseRepo:            manageDatabaseRepo,
		dbAccountRepo:                 dbAccountRepo,
		aliYunBaiLianRelationshipRepo: aliYunBaiLianRelationshipRepo,
	}
}

// Create 创建关系记录
func (s *Service) Create(relationship *models.Relationship) (*models.Relationship, error) {
	// 验证必要字段
	_, err := s.dbAccountRepo.GetByID(relationship.AccountID)
	if err != nil {
		return nil, err
	}
	if relationship.Title == "" {
		return nil, errors.New("标题不能为空")
	}
	if relationship.AccountID <= 0 {
		return nil, errors.New("账号ID必须大于0")
	}

	relationship, err = s.repo.Create(relationship)
	return relationship, err
}

// GetByID 根据ID获取关系记录
func (s *Service) GetByID(id uint) (*models.Relationship, error) {
	return s.repo.GetByID(id)
}

// GetByAccountID 根据账号ID获取关系记录列表
func (s *Service) GetByAccountID(accountID uint) ([]models.Relationship, error) {
	return s.repo.GetByAccountID(accountID)
}

// GetAll 获取所有关系记录
func (s *Service) GetAll() ([]models.Relationship, error) {
	return s.repo.GetAll()
}

// Update 更新关系记录
func (s *Service) Update(relationship *models.Relationship) error {
	// 验证必要字段
	if relationship.Title == "" {
		return errors.New("标题不能为空")
	}
	if relationship.AccountID <= 0 {
		return errors.New("账号ID必须大于0")
	}

	// 检查关系记录是否存在
	existing, err := s.repo.GetByID(relationship.ID)
	if err != nil {
		return errors.New("关系记录不存在")
	}

	// 保留创建时间
	relationship.CreatedTime = existing.CreatedTime

	return s.repo.Update(relationship)
}

// Delete 删除关系记录
func (s *Service) Delete(id uint) error {
	// 检查关系记录是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("关系记录不存在")
	}

	return s.repo.Delete(id)
}

// Analyze 分析关系
func (s *Service) Analyze(id uint, msgChan chan<- api.Response) error {
	relationship, err := s.GetByID(id)
	if err != nil {
		msgChan <- api.Response{
			Code: 400,
			Msg:  err.Error(),
		}
		return nil
	}

	var tables []models.TableInfo
	err = json.Unmarshal([]byte(relationship.Tables), &tables)
	if err != nil {
		msgChan <- api.Response{
			Code: 400,
			Msg:  err.Error(),
		}
		return nil
	}
	var names = make([]string, 0)
	for _, v := range tables {
		names = append(names, v.Name)
	}

	ddl, err := s.manageDatabaseRepo.GetDatabaseRepo(relationship.AccountID).GetTablesDDL(names)
	if err != nil {
		msgChan <- api.Response{
			Code: 400,
			Msg:  err.Error(),
		}
		return nil
	}
	//读取ddl
	ddlstr, err := json.Marshal(ddl)
	if err != nil {
		msgChan <- api.Response{
			Code: 400,
			Msg:  err.Error(),
		}
		return nil
	}
	//更新ddl
	relationship.Ddls = string(ddlstr)

	var ddlArr = make([]string, 0)
	for _, v := range ddl {
		ddlArr = append(ddlArr, v.DDL)
	}

	s.aliYunBaiLianRelationshipRepo.ClearAllMessages()
	s.aliYunBaiLianRelationshipRepo.AddBaseMessage(utils.SplitDDLToStrings(ddlArr, 5000))
	s.aliYunBaiLianRelationshipRepo.AddMessage([]string{aliyun_bailian_relationship.UserSummaryPrompt()})

	chatSteam, err := s.aliYunBaiLianRelationshipRepo.Chat()
	if err != nil {
		msgChan <- api.Response{
			Code: 400,
			Msg:  err.Error(),
		}
		return nil
	}
	allText := ""
	for text := range chatSteam {
		allText = allText + text
		msgChan <- api.Response{
			Code: 200,
			Data: text,
		}
	}
	relationship.Result = allText
	err = s.Update(relationship)
	if err != nil {
		msgChan <- api.Response{
			Code: 400,
			Msg:  err.Error(),
		}
		return nil
	}
	return nil
}
