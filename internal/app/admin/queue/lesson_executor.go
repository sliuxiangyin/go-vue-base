package queue

import (
	"context"
	sharedEvent "databaseAi/internal/shared/event"
	"databaseAi/internal/shared/models"
	protos "databaseAi/internal/shared/proto"
	"databaseAi/internal/shared/repo"
	"databaseAi/internal/shared/task"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

const LessonTask = "lesson_task"

type LessonPayload struct {
	AudioURL string `json:"audio_url"`
	UserID   uint   `json:"count"`
	ID       uint   `json:"id"`
}
type LessonTaskResult struct {
	Message string    `json:"message"`
	At      time.Time `json:"at"`
}
type LessonExecutor struct {
	lessonRepo   *repo.LessonRepo
	audioRepo    *repo.AudioRepo
	semanticRepo *repo.SemanticRepo
}

func NewLessonExecutor(lessonRepo *repo.LessonRepo, audioRepo *repo.AudioRepo, semanticRepo *repo.SemanticRepo) *LessonExecutor {
	return &LessonExecutor{
		lessonRepo:   lessonRepo,
		audioRepo:    audioRepo,
		semanticRepo: semanticRepo,
	}
}

func (e *LessonExecutor) GetType() string {
	return LessonTask
}

// Run 执行任务
func (e *LessonExecutor) Run(ctx context.Context, job *task.Job) (interface{}, error) {

	// 解析任务参数
	var payload LessonPayload
	if err := task.UnmarshalPayload(job, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	transcribeInfo, err := e.audioRepo.Transcribe(payload.AudioURL, "en", "small")
	if err != nil {
		// 发送错误通知给用户
		sharedEvent.BroadcastError("音频转录失败", map[string]interface{}{
			"lesson_id": payload.ID,
			"error":     err.Error(),
		},
		)
		return nil, err
	}
	lesson, err := e.lessonRepo.GetByID(payload.ID)
	if err != nil {
		return nil, err
	}
	lesson.ContentEN = transcribeInfo.Text

	var wait sync.WaitGroup

	levels := []repo.ChunkLevel{repo.ChunkLevel1, repo.ChunkLevel2, repo.ChunkLevel3}

	semanticMap := make(models.SemanticMapChunks)

	// 并行处理三个等级的语义切分
	for _, level := range levels {
		wait.Add(1)
		level := level // 避免闭包问题
		go func(currentLevel repo.ChunkLevel) {
			defer wait.Done()
			// 调用 AI 进行语义切分
			chunks, sematicErr := e.semanticRepo.AnalyzeSemanticChunks(ctx, lesson.ContentEN, currentLevel)
			if sematicErr != nil {
				log.Printf("语义切分失败 (level=%d): %v", currentLevel, sematicErr)
				return
			}
			// 将 chunks 和 transcribeInfo.Words 匹配处理成 models.SemanticChunk
			semanticChunks := e.matchChunksWithTimestamps(chunks, transcribeInfo.Words)
			// 存储到 semanticMap
			semanticMap[int(currentLevel)] = semanticChunks

		}(level)
	}
	wait.Wait()

	// 保存语义切分结果
	lesson.SemanticJSON = semanticMap

	// 保存 Word Timestamps
	lesson.WordTimestampJSON = e.convertProtoWordsToModel(transcribeInfo.Words)

	// 更新数据库
	if err := e.lessonRepo.Update(lesson); err != nil {
		// 发送错误通知给用户
		sharedEvent.BroadcastError("保存课程数据失败", map[string]interface{}{
			"lesson_id": payload.ID,
			"error":     err.Error(),
		})
		return nil, err
	}

	// 发送成功通知给用户
	sharedEvent.BroadcastSuccess("音频转录完成", map[string]interface{}{
		"lesson_id": payload.ID,
		"language":  transcribeInfo.Language,
		"text":      transcribeInfo.Text,
	})

	return LessonTaskResult{
		Message: fmt.Sprintf("%s", "ok"),
		At:      time.Now(),
	}, nil
}

// matchChunksWithTimestamps 将语义切分结果与时间戳匹配
func (e *LessonExecutor) matchChunksWithTimestamps(chunks []string, words []*protos.WordTimestamp) []models.SemanticChunk {
	if len(chunks) == 0 || len(words) == 0 {
		return []models.SemanticChunk{}
	}

	var semanticChunks []models.SemanticChunk
	wordIndex := 0

	// 遍历每个语义块
	for _, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}

		// 将语义块分解为单词（简单按空格分割）
		chunkWords := strings.Fields(chunk)
		if len(chunkWords) == 0 {
			continue
		}

		// 在 words 中查找匹配的单词序列
		startIdx, endIdx := e.findWordSequence(words, wordIndex, chunkWords)

		if startIdx >= 0 && endIdx >= 0 && startIdx < len(words) && endIdx < len(words) {
			// 创建 SemanticChunk
			semanticChunk := models.SemanticChunk{
				Text:  chunk,
				Start: words[startIdx].Start,
				End:   words[endIdx].End,
			}
			semanticChunks = append(semanticChunks, semanticChunk)

			// 更新单词索引
			wordIndex = endIdx + 1
		} else {
			// 如果找不到匹配，尝试从当前位置继续
			log.Printf("警告: 无法匹配语义块: %s", chunk)
			wordIndex++
		}

		// 防止超出范围
		if wordIndex >= len(words) {
			break
		}
	}

	return semanticChunks
}

// findWordSequence 在 words 中查找与 chunkWords 匹配的单词序列
func (e *LessonExecutor) findWordSequence(words []*protos.WordTimestamp, startIdx int, chunkWords []string) (int, int) {
	if startIdx >= len(words) {
		return -1, -1
	}

	// 尝试在后续的 words 中寻找匹配
	for i := startIdx; i < len(words); i++ {
		if e.matchSequence(words, i, chunkWords) {
			return i, i + len(chunkWords) - 1
		}
	}

	return -1, -1
}

// matchSequence 检查从指定位置开始是否匹配 chunkWords
func (e *LessonExecutor) matchSequence(words []*protos.WordTimestamp, startIdx int, chunkWords []string) bool {
	if startIdx+len(chunkWords) > len(words) {
		return false
	}

	for i, chunkWord := range chunkWords {
		// 正规化比较（转小写、去除标点）
		wordText := strings.ToLower(strings.Trim(words[startIdx+i].Word, ".,!?;:"))
		chunkWordText := strings.ToLower(strings.Trim(chunkWord, ".,!?;:"))

		if wordText != chunkWordText {
			return false
		}
	}

	return true
}

// convertProtoWordsToModel 将 protobuf 的 WordTimestamp 转换为 model 的 WordTimestamp
func (e *LessonExecutor) convertProtoWordsToModel(protoWords []*protos.WordTimestamp) models.WordTimestamps {
	if protoWords == nil {
		return nil
	}

	modelWords := make(models.WordTimestamps, 0, len(protoWords))

	for _, pw := range protoWords {
		modelWord := models.WordTimestamp{
			Word:       pw.Word,
			Start:      pw.Start,
			End:        pw.End,
			Confidence: pw.Confidence,
		}
		modelWords = append(modelWords, modelWord)
	}

	return modelWords
}
