package utils

import (
	"regexp"
	"strings"
)

// CompressDDL 去除换行和多余空格，把 DDL 整理成一行
func CompressDDL(ddl string) string {
	// 去掉换行符、制表符
	ddl = strings.ReplaceAll(ddl, "\n", " ")
	ddl = strings.ReplaceAll(ddl, "\r", " ")
	ddl = strings.ReplaceAll(ddl, "\t", " ")

	// 合并多个空格为一个
	re := regexp.MustCompile(`\s+`)
	ddl = re.ReplaceAllString(ddl, " ")

	// 去除前后空格
	ddl = strings.TrimSpace(ddl)

	// 修复一些常见格式问题： ( 前后多余空格
	ddl = strings.ReplaceAll(ddl, "( ", "(")
	ddl = strings.ReplaceAll(ddl, " )", ")")
	ddl = strings.ReplaceAll(ddl, " ,", ",")

	return ddl
}

// SplitDDLToStrings 把多个 DDL 按 maxSize 分批拼接成字符串
// 每个分片长度不超过 maxSize，且保持每条 DDL 完整
func SplitDDLToStrings(ddls []string, maxSize int) []string {
	var result []string
	var batch []string
	currentSize := 0

	for _, ddl := range ddls {
		ddlSize := len(ddl)

		// 如果单条 DDL 本身超过 maxSize，直接作为单独一批
		if ddlSize > maxSize {
			if len(batch) > 0 {
				result = append(result, strings.Join(batch, ";\n"))
				batch = nil
				currentSize = 0
			}
			result = append(result, ddl)
			continue
		}

		// 累积超过限制则换批
		if currentSize+ddlSize > maxSize && len(batch) > 0 {
			result = append(result, strings.Join(batch, ";\n"))
			batch = []string{ddl}
			currentSize = ddlSize
		} else {
			batch = append(batch, ddl)
			currentSize += ddlSize
		}
	}

	if len(batch) > 0 {
		result = append(result, strings.Join(batch, ";\n"))
	}

	return result
}
