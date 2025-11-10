package repo

import (
	"databaseAi/internal/app/tools/models"
	"databaseAi/internal/infra/database"
	"databaseAi/internal/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type DataBaseRepo struct {
	db *database.DB
}

func NewDataBaseRepo(db *database.DB) *DataBaseRepo {
	return &DataBaseRepo{db: db}
}

// GetAllTables 获取数据库所有表名和表注释
func (r *DataBaseRepo) GetAllTables() ([]models.TableInfo, error) {
	var tables []models.TableInfo
	var gdb *gorm.DB = r.db.GetDB()

	// 获取当前数据库名
	var dbName string
	if err := gdb.Raw("SELECT DATABASE()").Scan(&dbName).Error; err != nil {
		return nil, err
	}

	// 查询 information_schema.tables
	sql := `
	SELECT 
		TABLE_NAME, 
		IFNULL(TABLE_COMMENT, '') AS TABLE_COMMENT
	FROM information_schema.tables 
	WHERE table_schema = '%s'
	ORDER BY TABLE_NAME
`
	query := fmt.Sprintf(sql, dbName)
	if err := gdb.Raw(query).Scan(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

// GetTablesDDL 根据多个表名获取各自的 DDL 语句
func (r *DataBaseRepo) GetTablesDDL(tableNames []string) ([]models.TableDDL, error) {
	var gdb *gorm.DB = r.db.GetDB()
	if gdb == nil {
		return nil, errors.New("刷新页面")
	}
	var results []models.TableDDL

	for _, tableName := range tableNames {
		var ddl string
		// SHOW CREATE TABLE 返回两列：Table 和 Create Table
		row := gdb.Raw(fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)).Row()
		var name string
		if err := row.Scan(&name, &ddl); err != nil {
			// 出错则跳过该表，但不中断整个流程
			continue
		}
		results = append(results, models.TableDDL{
			Name: tableName,
			DDL:  utils.CompressDDL(ddl),
		})
	}

	return results, nil
}
