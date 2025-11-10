package compile

import (
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"strings"
)

func SetupDev(app *fiber.App) {

	// 使用构建标签来区分开发和生产环境
	// 在开发模式下启动前端开发服务器
	webHost := StartWebDevServer()

	// 在开发模式下代理前端请求到 Vite 开发服务器
	app.Use("/", func(c *fiber.Ctx) error {
		// 如果是 API 请求，继续处理
		if strings.HasPrefix(c.Path(), "/api/") {
			return c.Next()
		}

		// 创建一个请求体读取器
		bodyReader := strings.NewReader(string(c.Request().Body()))

		// 其他请求代理到前端开发服务器
		proxyReq, err := http.NewRequest(c.Method(), webHost+c.Path(), bodyReader)
		if err != nil {
			return c.Next() // 如果创建请求失败，继续到下一个处理器
		}

		// 复制请求头
		c.Request().Header.VisitAll(func(key, value []byte) {
			proxyReq.Header.Set(string(key), string(value))
		})

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			return c.Next() // 如果请求失败，继续到下一个处理器
		}
		defer resp.Body.Close()

		// 复制响应状态码
		c.Response().SetStatusCode(resp.StatusCode)

		// 复制响应头
		for key, values := range resp.Header {
			for _, value := range values {
				c.Response().Header.Add(key, value)
			}
		}

		// 读取并返回响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return c.Send(body)
	})
}
