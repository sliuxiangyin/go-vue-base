package models

// TableInfo 表结构信息
type TableInfo struct {
	Name    string `json:"name" gorm:"column:TABLE_NAME"`
	Comment string `json:"comment" gorm:"column:TABLE_COMMENT"`
}
type TableDDL struct {
	Name string `json:"name"`
	DDL  string `json:"ddl"`
}
