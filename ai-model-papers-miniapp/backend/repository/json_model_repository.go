package repository

import (
	"ai-model-papers-backend/models"
	"encoding/json"
	"os"
	"strings"
	"sync"
)

// JSONModelRepository JSON文件实现的模型仓库
type JSONModelRepository struct {
	data     models.ModelsData
	dataPath string
	mu       sync.RWMutex
}

// NewJSONModelRepository 创建JSON模型仓库
func NewJSONModelRepository(dataPath string) (*JSONModelRepository, error) {
	repo := &JSONModelRepository{
		dataPath: dataPath,
	}
	
	if err := repo.load(); err != nil {
		return nil, err
	}
	
	return repo, nil
}

// load 从JSON文件加载数据
func (r *JSONModelRepository) load() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	data, err := os.ReadFile(r.dataPath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, &r.data)
}

// GetAll 获取所有模型
func (r *JSONModelRepository) GetAll() ([]models.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	models := make([]models.Model, len(r.data.Models))
	copy(models, r.data.Models)
	return models, nil
}

// GetByID 根据ID获取模型
func (r *JSONModelRepository) GetByID(id string) (*models.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	for _, model := range r.data.Models {
		if model.ID == id {
			return &model, nil
		}
	}
	return nil, nil
}

// GetByProvider 根据提供商筛选
func (r *JSONModelRepository) GetByProvider(providerID string) ([]models.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []models.Model
	for _, model := range r.data.Models {
		if model.ProviderID == providerID {
			result = append(result, model)
		}
	}
	return result, nil
}

// GetByCapability 根据能力筛选
func (r *JSONModelRepository) GetByCapability(capability string) ([]models.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []models.Model
	for _, model := range r.data.Models {
		switch capability {
		case "text":
			if model.Capabilities.Text {
				result = append(result, model)
			}
		case "image":
			if model.Capabilities.Image {
				result = append(result, model)
			}
		case "audio":
			if model.Capabilities.Audio {
				result = append(result, model)
			}
		case "video":
			if model.Capabilities.Video {
				result = append(result, model)
			}
		case "file":
			if model.Capabilities.File {
				result = append(result, model)
			}
		}
	}
	return result, nil
}

// Search 搜索模型
func (r *JSONModelRepository) Search(keyword string) ([]models.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	keyword = strings.ToLower(keyword)
	var result []models.Model
	
	for _, model := range r.data.Models {
		if strings.Contains(strings.ToLower(model.Name), keyword) ||
			strings.Contains(strings.ToLower(model.Provider), keyword) ||
			strings.Contains(strings.ToLower(model.Description), keyword) {
			result = append(result, model)
		}
	}
	return result, nil
}

// GetProviders 获取所有提供商
func (r *JSONModelRepository) GetProviders() ([]models.Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	providers := make([]models.Provider, len(r.data.Providers))
	copy(providers, r.data.Providers)
	return providers, nil
}

// GetComparisonCategories 获取对比类别
func (r *JSONModelRepository) GetComparisonCategories() ([]models.ComparisonCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	categories := make([]models.ComparisonCategory, len(r.data.ComparisonCategories))
	copy(categories, r.data.ComparisonCategories)
	return categories, nil
}
