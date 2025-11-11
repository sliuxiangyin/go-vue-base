package repo

import (
	"databaseAi/internal/infra/database"
	sharedModels "databaseAi/internal/shared/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *database.DB {
	// 使用内存 SQLite 进行测试，每个测试独立数据库
	db, err := database.NewDB("file::memory:?cache=private")
	require.NoError(t, err)

	// 迁移表结构
	err = db.AutoMigrate(&sharedModels.EnglishLesson{})
	require.NoError(t, err)

	return db
}

func TestLessonRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLessonRepo(db)

	lesson := &sharedModels.EnglishLesson{
		Title:     "Forrest Gump - Opening Narration",
		AudioURL:  "/audio/forrest_gump_opening.mp3",
		Duration:  23.6,
		ContentEN: "Life was like a box of chocolates, you never know what you're gonna get.",
		ContentZH: "生活就像一盒巧克力，你永远不知道下一颗是什么味道。",
		SemanticJSON: sharedModels.SemanticChunks{
			{Text: "Life was like a box of chocolates", Start: 0.00, End: 2.50},
			{Text: "you never know what you're gonna get", Start: 2.50, End: 5.00},
		},
		WordTimestampJSON: sharedModels.WordTimestamps{
			{Word: "Life", Start: 0.00, End: 0.30, Confidence: 0.98},
			{Word: "was", Start: 0.30, End: 0.50, Confidence: 0.99},
			{Word: "like", Start: 0.50, End: 0.70, Confidence: 0.97},
		},
		PhoneticJSON: sharedModels.Phonetics{
			{Word: "Life", IPA: "/laɪf/"},
			{Word: "chocolates", IPA: "/ˈtʃɒklɪts/"},
		},
		Tags:     sharedModels.Tags{"movie", "classic", "life"},
		Level:    2,
		IsPublic: true,
	}

	err := repo.Create(lesson)
	require.NoError(t, err)
	assert.Greater(t, lesson.ID, uint(0))
}

func TestLessonRepo_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLessonRepo(db)

	// 创建测试数据
	lesson := &sharedModels.EnglishLesson{
		Title:     "Test Lesson",
		AudioURL:  "/audio/test.mp3",
		Duration:  10.0,
		ContentEN: "Hello world",
		ContentZH: "你好世界",
		Tags:      sharedModels.Tags{"test"},
		Level:     1,
		IsPublic:  true,
	}
	err := repo.Create(lesson)
	require.NoError(t, err)

	// 查询
	result, err := repo.GetByID(lesson.ID)
	require.NoError(t, err)
	assert.Equal(t, lesson.Title, result.Title)
	assert.Equal(t, lesson.AudioURL, result.AudioURL)
	assert.Equal(t, lesson.Tags, result.Tags)
}

func TestLessonRepo_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLessonRepo(db)

	// 创建多个测试课程
	lessons := []*sharedModels.EnglishLesson{
		{Title: "Lesson 1", Level: 1, IsPublic: true, AudioURL: "/audio/1.mp3", ContentEN: "Test 1", ContentZH: "测试1", Tags: sharedModels.Tags{"beginner"}},
		{Title: "Lesson 2", Level: 2, IsPublic: true, AudioURL: "/audio/2.mp3", ContentEN: "Test 2", ContentZH: "测试2", Tags: sharedModels.Tags{"intermediate"}},
		{Title: "Lesson 3", Level: 1, IsPublic: false, AudioURL: "/audio/3.mp3", ContentEN: "Test 3", ContentZH: "测试3", Tags: sharedModels.Tags{"beginner"}},
	}

	for _, lesson := range lessons {
		err := repo.Create(lesson)
		require.NoError(t, err)
	}

	// 测试分页
	result, total, err := repo.List(1, 10, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, result, 3)

	// 测试过滤 level
	level := int8(1)
	result, total, err = repo.List(1, 10, &level, nil)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)

	// 测试过滤 is_public
	isPublic := true
	result, total, err = repo.List(1, 10, nil, &isPublic)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total) // Lesson 1 和 Lesson 2 是公开的
	assert.Len(t, result, 2)
}

func TestLessonRepo_SearchByTags(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLessonRepo(db)

	// 创建测试课程
	lessons := []*sharedModels.EnglishLesson{
		{
			Title: "Movie Lesson", AudioURL: "/audio/movie.mp3", ContentEN: "Test", ContentZH: "测试",
			Tags: sharedModels.Tags{"movie", "classic"}, IsPublic: true, Level: 2,
		},
		{
			Title: "Music Lesson", AudioURL: "/audio/music.mp3", ContentEN: "Test", ContentZH: "测试",
			Tags: sharedModels.Tags{"music", "pop"}, IsPublic: true, Level: 1,
		},
	}

	for _, lesson := range lessons {
		err := repo.Create(lesson)
		require.NoError(t, err)
	}

	// 搜索包含 "movie" 标签的课程
	result, err := repo.SearchByTags([]string{"movie"})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(result), 1)
	if len(result) > 0 {
		assert.Equal(t, "Movie Lesson", result[0].Title)
	}
}

func TestLessonRepo_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLessonRepo(db)

	// 创建课程
	lesson := &sharedModels.EnglishLesson{
		Title:     "Original Title",
		AudioURL:  "/audio/test.mp3",
		ContentEN: "Test",
		ContentZH: "测试",
		Tags:      sharedModels.Tags{"test"},
		Level:     1,
		IsPublic:  true,
	}
	err := repo.Create(lesson)
	require.NoError(t, err)

	// 更新课程
	lesson.Title = "Updated Title"
	lesson.Level = 3
	err = repo.Update(lesson)
	require.NoError(t, err)

	// 验证更新
	result, err := repo.GetByID(lesson.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", result.Title)
	assert.Equal(t, int8(3), result.Level)
}

func TestLessonRepo_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLessonRepo(db)

	// 创建课程
	lesson := &sharedModels.EnglishLesson{
		Title:     "To Delete",
		AudioURL:  "/audio/test.mp3",
		ContentEN: "Test",
		ContentZH: "测试",
		Tags:      sharedModels.Tags{"test"},
		Level:     1,
		IsPublic:  true,
	}
	err := repo.Create(lesson)
	require.NoError(t, err)

	// 删除课程
	err = repo.Delete(lesson.ID)
	require.NoError(t, err)

	// 验证删除（软删除，查询不到）
	_, err = repo.GetByID(lesson.ID)
	assert.Error(t, err)
}
