package compile

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sync"
	"time"
)

// 启动前端开发服务器
func StartWebDevServer(root string, webName string) string {
	getwd, _ := os.Getwd()
	devDir := path.Join(getwd, root, webName)
	var webHost = ""
	// 检查 web 目录是否存在
	if _, err := os.Stat(devDir); os.IsNotExist(err) {
		log.Println("web directory not found, skipping web dev server start")
		return ""
	}
	// 检查 package.json 是否存在
	if _, err := os.Stat(path.Join(devDir, "package.json")); os.IsNotExist(err) {
		log.Println("web package.json not found, skipping web dev server start")
		return ""
	}

	log.Println("Starting web development server...")
	cmd := exec.Command("pnpm", "run", "dev")
	cmd.Dir = devDir

	var buf bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &buf)
	cmd.Stdout = mw
	cmd.Stderr = mw
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start web dev server: %v", err)
	} else {
		log.Println("web development server started successfully")
	}
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		defer wait.Done()
		for {
			urlPattern := regexp.MustCompile(`http://localhost:\d+`)
			cleaned := stripANSI(buf.String())
			match := urlPattern.FindString(cleaned)
			if match != "" {
				webHost = match
				break
			} else {
			}
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		// 等待命令完成
		if err := cmd.Wait(); err != nil {
			fmt.Println("命令退出:", err)
		}
	}()

	wait.Wait()
	return webHost

}
func stripANSI(str string) string {
	ansi := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansi.ReplaceAllString(str, "")
}
