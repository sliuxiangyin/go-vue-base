package event

import (
	"bufio"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// EventType 事件类型
type EventType string

const (
	EventTypeSuccess EventType = "success"
	EventTypeError   EventType = "error"
	EventTypeWarning EventType = "warning"
	EventTypeInfo    EventType = "info"
)

// Event 事件结构
type Event struct {
	Type      EventType              `json:"type"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// Client SSE 客户端
type Client struct {
	ID      string
	UserID  uint
	Channel chan Event
}

// EventManager 事件管理器
type EventManager struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

var (
	manager *EventManager
	once    sync.Once
)

// GetEventManager 获取全局事件管理器实例（单例）
func GetEventManager() *EventManager {
	once.Do(func() {
		manager = &EventManager{
			clients: make(map[string]*Client),
		}
	})
	return manager
}

// AddClient 添加客户端
func (em *EventManager) AddClient(client *Client) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.clients[client.ID] = client
}

// RemoveClient 移除客户端
func (em *EventManager) RemoveClient(clientID string) {
	em.mu.Lock()
	defer em.mu.Unlock()
	if client, ok := em.clients[clientID]; ok {
		close(client.Channel)
		delete(em.clients, clientID)
	}
}

// SendToUser 发送事件给特定用户
func (em *EventManager) SendToUser(userID uint, event Event) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	event.Timestamp = time.Now()

	fmt.Printf("[SendToUser] 尝试发送给 userID=%d, 当前客户端总数=%d\n", userID, len(em.clients))
	sentCount := 0
	for clientID, client := range em.clients {
		fmt.Printf("[SendToUser] 检查客户端 %s, userID=%d, 匹配=%v\n", clientID, client.UserID, client.UserID == userID)
		if client.UserID == userID {
			select {
			case client.Channel <- event:
				sentCount++
				fmt.Printf("[SendToUser] ✓ 成功发送给客户端 %s (userID=%d)\n", clientID, userID)
			case <-time.After(time.Second):
				fmt.Printf("[SendToUser] ✗ 发送超时，客户端 %s (userID=%d)\n", clientID, userID)
			}
		}
	}
	fmt.Printf("[SendToUser] 完成发送给 userID=%d, 成功发送到 %d 个客户端\n", userID, sentCount)
}

// Broadcast 广播事件给所有客户端
func (em *EventManager) Broadcast(event Event) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	event.Timestamp = time.Now()

	for _, client := range em.clients {
		select {
		case client.Channel <- event:
		case <-time.After(time.Second):
			// 发送超时，跳过
		}
	}
}

// GetClientCount 获取当前连接的客户端数量
func (em *EventManager) GetClientCount() int {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return len(em.clients)
}

// ========== 便捷方法 ==========

// NotifySuccess 发送成功通知
func NotifySuccess(userID uint, message string, data map[string]interface{}) {
	GetEventManager().SendToUser(userID, Event{
		Type:    EventTypeSuccess,
		Message: message,
		Data:    data,
	})
}

// NotifyError 发送错误通知
func NotifyError(userID uint, message string, data map[string]interface{}) {
	GetEventManager().SendToUser(userID, Event{
		Type:    EventTypeError,
		Message: message,
		Data:    data,
	})
}

// NotifyWarning 发送警告通知
func NotifyWarning(userID uint, message string, data map[string]interface{}) {
	GetEventManager().SendToUser(userID, Event{
		Type:    EventTypeWarning,
		Message: message,
		Data:    data,
	})
}

// NotifyInfo 发送信息通知
func NotifyInfo(userID uint, message string, data map[string]interface{}) {
	GetEventManager().SendToUser(userID, Event{
		Type:    EventTypeInfo,
		Message: message,
		Data:    data,
	})
}

// BroadcastSuccess 广播成功消息
func BroadcastSuccess(message string, data map[string]interface{}) {
	GetEventManager().Broadcast(Event{
		Type:    EventTypeSuccess,
		Message: message,
		Data:    data,
	})
}

// BroadcastError 广播错误消息
func BroadcastError(message string, data map[string]interface{}) {
	GetEventManager().Broadcast(Event{
		Type:    EventTypeError,
		Message: message,
		Data:    data,
	})
}

// HandleSSE SSE 连接处理器
func HandleSSE(c *fiber.Ctx, userID uint) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	clientID := c.Get("X-Request-ID")
	if clientID == "" {
		clientID = time.Now().Format("20060102150405") + "-" + c.IP()
	}
	fmt.Printf("[SSE] 新客户端连接: clientID=%s, userID=%d, IP=%s\n", clientID, userID, c.IP())
	client := &Client{
		ID:      clientID,
		UserID:  userID,
		Channel: make(chan Event, 100),
	}

	manager := GetEventManager()
	manager.AddClient(client)

	// 发送初始连接成功消息
	initialEvent := Event{
		Type:      EventTypeInfo,
		Message:   "连接成功",
		Timestamp: time.Now(),
	}
	eventData, _ := json.Marshal(initialEvent)
	c.Write([]byte("data: "))
	c.Write(eventData)
	c.Write([]byte("\n\n"))

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer manager.RemoveClient(clientID)
		// 心跳定时器
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		defer func() {
			fmt.Println("HandleSSE关闭")
		}()
		for {
			select {
			case event, ok := <-client.Channel:
				fmt.Printf("[HandleSSE] 客户端 %s (userID=%d) 接收到事件: type=%s, message=%s, ok=%v\n", clientID, userID, event.Type, event.Message, ok)
				if !ok {
					fmt.Printf("[HandleSSE] 客户端 %s channel 已关闭\n", clientID)
					return
				}
				eventData, err := json.Marshal(event)
				if err != nil {
					fmt.Printf("[HandleSSE] JSON 序列化失败: %v\n", err)
					continue
				}

				w.Write([]byte("data: "))
				w.Write(eventData)
				w.Write([]byte("\n\n"))
				err = w.Flush()
				if err != nil {
					fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
					return
				}

			case <-ticker.C:
				// 发送心跳
				w.Write([]byte(": heartbeat\n\n"))
				err := w.Flush()
				if err != nil {
					fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
					return
				}

			}
		}

	})

	return nil
}
