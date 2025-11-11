package repo

import (
	"databaseAi/internal/infra/database"
	"databaseAi/internal/shared/models"
)

// PhoneticDictionaryRepo 发音词典仓库层
type PhoneticDictionaryRepo struct {
	db *database.DB
}

func NewPhoneticDictionaryRepo(db *database.DB) *PhoneticDictionaryRepo {
	return &PhoneticDictionaryRepo{
		db: db,
	}
}

// CreateOrUpdate 创建或更新单词发音（如果单词已存在则更新）
func (r *PhoneticDictionaryRepo) CreateOrUpdate(word, ipa string) (*models.PhoneticDictionary, error) {
	var phonetic models.PhoneticDictionary

	// 先查找是否存在
	err := r.db.Where("word = ?", word).First(&phonetic).Error
	if err != nil {
		// 不存在，创建新记录
		phonetic = models.PhoneticDictionary{
			Word: word,
			IPA:  ipa,
		}
		if err := r.db.Create(&phonetic).Error; err != nil {
			return nil, err
		}
	} else {
		// 存在，更新 IPA
		phonetic.IPA = ipa
		if err := r.db.Save(&phonetic).Error; err != nil {
			return nil, err
		}
	}

	return &phonetic, nil
}

// BatchCreateOrUpdate 批量创建或更新发音
func (r *PhoneticDictionaryRepo) BatchCreateOrUpdate(phonetics []models.PhoneticDictionary) error {
	for _, p := range phonetics {
		if _, err := r.CreateOrUpdate(p.Word, p.IPA); err != nil {
			return err
		}
	}
	return nil
}

// GetByWord 根据单词获取发音
func (r *PhoneticDictionaryRepo) GetByWord(word string) (*models.PhoneticDictionary, error) {
	var phonetic models.PhoneticDictionary
	err := r.db.Where("word = ?", word).First(&phonetic).Error
	if err != nil {
		return nil, err
	}
	return &phonetic, nil
}

// GetByWords 批量获取单词发音
func (r *PhoneticDictionaryRepo) GetByWords(words []string) ([]models.PhoneticDictionary, error) {
	var phonetics []models.PhoneticDictionary
	err := r.db.Where("word IN ?", words).Find(&phonetics).Error
	if err != nil {
		return nil, err
	}
	return phonetics, nil
}

// List 获取所有发音词典（分页）
func (r *PhoneticDictionaryRepo) List(page, pageSize int) ([]models.PhoneticDictionary, int64, error) {
	var phonetics []models.PhoneticDictionary
	var total int64

	query := r.db.Model(&models.PhoneticDictionary{})

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("word ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&phonetics).Error; err != nil {
		return nil, 0, err
	}

	return phonetics, total, nil
}

// Delete 删除单词发音
func (r *PhoneticDictionaryRepo) Delete(id uint) error {
	return r.db.Delete(&models.PhoneticDictionary{}, id).Error
}
