package utils

import (
	"os"
	"path/filepath"
)

// ProjectRoot 自动定位项目根目录（包含 go.mod 的目录）
func ProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("未找到 go.mod，无法定位项目根目录")
		}
		dir = parent
	}
}
