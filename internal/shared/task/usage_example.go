package task

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// UsageExample 使用示例（手动初始化方式）
// 注意：推荐使用 Wire 依赖注入，参考 wire_integration_example.go
func UsageExample(db *gorm.DB) {
	// ========== 第一步：初始化组件 ==========

	// 创建任务仓库
	repo := NewTaskRepo(db)

	// 自动迁移数据库表
	err := repo.AutoMigrate()
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 获取全局事件总线单例
	eventBus := GetEventBus()

	// 创建任务执行器注册表
	registry := NewTaskExecutorRegistry()

	// ========== 第二步：注册任务执行器（支持依赖注入）==========

	// 方式1: 注册不需要依赖的执行器
	//registry.Register(tasks2.NewExampleTaskExecutor())

	// 方式2: 注册需要依赖的执行器
	// 假设你有一个 FileRepo
	// fileRepo := repo.NewFileRepo(db)
	// registry.Register(NewFileProcessExecutor(fileRepo))

	// 方式3: 注册需要多个依赖的执行器
	// lessonRepo := repo.NewLessonRepo(db)
	// registry.Register(NewAudioProcessExecutor(fileRepo, lessonRepo))

	// ========== 第三步：创建并启动任务管理器 ==========

	manager := NewTaskManager(repo, eventBus, registry)
	manager.Start()

	// 确保程序退出时停止任务管理器
	defer manager.Stop()

	// ========== 第四步：订阅任务事件（可选） ==========

	// 监听任务开始事件
	eventBus.Subscribe(EventTypeTaskStart, func(event TaskEvent) {
		log.Printf("📢 任务开始执行: id=%d, type=%s, user_id=%d",
			event.TaskID, event.Job.Type, event.Job.UserID)

		// 这里可以执行外部逻辑
		// 例如：发送 WebSocket 通知给用户
		// 例如：记录审计日志
		// 例如：触发其他业务流程
	})

	// 监听任务成功事件
	eventBus.Subscribe(EventTypeTaskSuccess, func(event TaskEvent) {
		log.Printf("✅ 任务执行成功: id=%d, result=%v",
			event.TaskID, event.Data["result"])

		// 这里可以执行外部逻辑
		// 例如：更新业务状态
		// 例如：发送成功通知
		// 例如：触发后续任务
	})

	// 监听任务失败事件
	eventBus.Subscribe(EventTypeTaskFailed, func(event TaskEvent) {
		log.Printf("❌ 任务执行失败: id=%d, error=%v",
			event.TaskID, event.Data["error"])

		// 这里可以执行外部逻辑
		// 例如：发送告警通知
		// 例如：记录错误日志
		// 例如：触发补偿逻辑
	})

	// 监听任务进度更新（如果任务执行器支持）
	eventBus.Subscribe(EventTypeTaskProgress, func(event TaskEvent) {
		log.Printf("🔄 任务进度更新: id=%d, progress=%v",
			event.TaskID, event.Data["progress"])

		// 这里可以执行外部逻辑
		// 例如：更新进度条
		// 例如：发送进度通知
	})

	// ========== 第五步：在业务代码中提交任务 ==========

	// 示例1: 提交单个任务
	err = eventBus.SubmitTask(SubmitTaskRequest{
		Type:   "example_task",
		UserID: 1,
		Payload: map[string]interface{}{
			"message": "处理订单数据",
			"count":   100,
		},
		MaxParallel: 1, // 最多同时执行1个
		Priority:    0, // 默认优先级
		MaxAttempts: 3, // 失败后最多重试3次
	})
	if err != nil {
		log.Printf("提交任务失败: %v", err)
	}

	// 示例2: 提交多个任务（批量处理）
	for i := 1; i <= 5; i++ {
		err = eventBus.SubmitTask(SubmitTaskRequest{
			Type:   "example_task",
			UserID: 1,
			Payload: map[string]interface{}{
				"message": "批量任务",
				"count":   i * 10,
			},
			MaxParallel: 3, // 最多同时执行3个
			Priority:    i, // 优先级递增
			MaxAttempts: 3,
		})
		if err != nil {
			log.Printf("提交任务失败: %v", err)
		}
	}

	// 示例3: 提交高优先级任务
	err = eventBus.SubmitTask(SubmitTaskRequest{
		Type:   "example_task",
		UserID: 2,
		Payload: map[string]interface{}{
			"message": "紧急任务",
			"count":   999,
		},
		MaxParallel: 1,
		Priority:    100, // 高优先级，会优先执行
		MaxAttempts: 5,
	})
	if err != nil {
		log.Printf("提交任务失败: %v", err)
	}

	// ========== 第六步：查询任务状态 ==========

	// 等待一段时间让任务执行
	time.Sleep(3 * time.Second)

	// 查询用户的任务列表
	tasks, total, err := repo.GetTasksByUserID(1, "", 10, 0)
	if err != nil {
		log.Printf("查询任务失败: %v", err)
	} else {
		log.Printf("用户1的任务总数: %d", total)
		for _, t := range tasks {
			log.Printf("  任务: id=%d, type=%s, status=%s, created_at=%v",
				t.ID, t.Type, t.Status, t.CreatedAt)
		}
	}

	// ========== 第七步：取消任务（如果需要） ==========

	// 假设要取消任务ID为10的任务
	// err = manager.CancelTask(10)
	// if err != nil {
	//     log.Printf("取消任务失败: %v", err)
	// }

	// ========== 第八步：保持程序运行 ==========

	// 在实际应用中，程序会持续运行
	// 这里为了演示，等待一段时间
	log.Println("任务系统运行中...")
	time.Sleep(30 * time.Second)

	log.Println("示例结束")
}

// BusinessExample 业务代码示例
// 参数 eventBus 可以通过以下方式获取：
//  1. Wire 注入：在 Handler/Service 构造函数中注入
//  2. 全局单例：eventBus := task.GetEventBus()
func BusinessExample(eventBus *EventBus) {
	// 在你的业务代码中，只需要调用 SubmitTask 即可
	// 完全不需要关心任务如何存储、如何调度、如何执行
	//
	// 如果没有注入 eventBus，也可以直接使用全局单例：
	// eventBus := task.GetEventBus()

	// 例1: 用户上传文件后，提交文件处理任务
	err := eventBus.SubmitTask(SubmitTaskRequest{
		Type:   "file_process",
		UserID: 123,
		Payload: map[string]interface{}{
			"file_id":   "abc123",
			"file_path": "/uploads/document.pdf",
			"action":    "convert_to_text",
		},
		MaxParallel: 5, // 允许同时处理5个文件
	})
	if err != nil {
		log.Printf("提交文件处理任务失败: %v", err)
	}

	// 例2: 定时任务，每天凌晨生成报表
	err = eventBus.SubmitTask(SubmitTaskRequest{
		Type:   "generate_report",
		UserID: 0, // 系统任务
		Payload: map[string]interface{}{
			"report_type": "daily_sales",
			"date":        time.Now().Format("2006-01-02"),
		},
		MaxParallel: 1, // 同时只能生成1个报表
		Priority:    10,
	})
	if err != nil {
		log.Printf("提交报表生成任务失败: %v", err)
	}

	// 例3: AI 对话任务
	err = eventBus.SubmitTask(SubmitTaskRequest{
		Type:   "ai_chat",
		UserID: 456,
		Payload: map[string]interface{}{
			"conversation_id": "conv_789",
			"message":         "请帮我分析这段代码",
			"model":           "gpt-4",
		},
		MaxParallel: 10, // 允许10个并发 AI 请求
		Priority:    5,
	})
	if err != nil {
		log.Printf("提交 AI 对话任务失败: %v", err)
	}
}
