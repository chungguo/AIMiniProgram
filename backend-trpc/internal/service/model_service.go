package service

import (
	"context"
	"fmt"

	"aiminiprogram/backend-trpc/internal/repository"
	pb "aiminiprogram/proto/model/v1"
	commonpb "aiminiprogram/proto/common/v1"

	"trpc.group/trpc-go/trpc-go/log"
)

// ModelService 实现
type ModelService struct {
	repo repository.ModelRepository
}

// NewModelService 创建服务实例
func NewModelService() *ModelService {
	return &ModelService{
		repo: repository.NewModelRepository(),
	}
}

// ListModels 获取模型列表
func (s *ModelService) ListModels(ctx context.Context, req *pb.ListModelsRequest) (*pb.ListModelsResponse, error) {
	log.Debugf("ListModels called with page=%d, limit=%d", 
		req.Pagination.Page, req.Pagination.Limit)

	// 转换筛选条件
	filter := &repository.ModelFilter{}
	if req.Filter != nil {
		filter.Family = req.Filter.Family
		filter.Search = req.Filter.Search
		if req.Filter.HasAttachment != nil {
			filter.HasAttachment = req.Filter.HasAttachment
		}
		if req.Filter.HasReasoning != nil {
			filter.HasReasoning = req.Filter.HasReasoning
		}
		if req.Filter.HasToolCall != nil {
			filter.HasToolCall = req.Filter.HasToolCall
		}
		if req.Filter.OpenWeights != nil {
			filter.OpenWeights = req.Filter.OpenWeights
		}
		filter.MinContext = req.Filter.MinContext
		filter.MaxCostInput = req.Filter.MaxCostInput
	}

	// 调用仓库层
	models, total, err := s.repo.List(ctx, filter, req.Pagination.Page, req.Pagination.Limit)
	if err != nil {
		log.Errorf("Failed to list models: %v", err)
		return &pb.ListModelsResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 转换为 proto 消息
	pbModels := make([]*pb.Model, len(models))
	for i, m := range models {
		pbModels[i] = convertToProtoModel(m)
	}

	return &pb.ListModelsResponse{
		Success: true,
		Message: "success",
		Data:    pbModels,
		Pagination: &commonpb.Pagination{
			Page:       req.Pagination.Page,
			Limit:      req.Pagination.Limit,
			Total:      int32(total),
			TotalPages: (int32(total) + req.Pagination.Limit - 1) / req.Pagination.Limit,
		},
	}, nil
}

// GetModel 获取单个模型
func (s *ModelService) GetModel(ctx context.Context, req *pb.GetModelRequest) (*pb.GetModelResponse, error) {
	log.Debugf("GetModel called with id=%s", req.Id)

	model, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		log.Errorf("Failed to get model: %v", err)
		return &pb.GetModelResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.GetModelResponse{
		Success: true,
		Message: "success",
		Data:    convertToProtoModel(model),
	}, nil
}

// ListFamilies 获取家族列表
func (s *ModelService) ListFamilies(ctx context.Context, req *commonpb.PaginationRequest) (*pb.ListFamiliesResponse, error) {
	log.Debug("ListFamilies called")

	families, err := s.repo.GetFamilies(ctx)
	if err != nil {
		log.Errorf("Failed to list families: %v", err)
		return &pb.ListFamiliesResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.ListFamiliesResponse{
		Success:  true,
		Message:  "success",
		Families: families,
	}, nil
}

// GetFamilyModels 获取家族模型
func (s *ModelService) GetFamilyModels(ctx context.Context, req *pb.GetFamilyModelsRequest) (*pb.GetFamilyModelsResponse, error) {
	log.Debugf("GetFamilyModels called with family=%s", req.Family)

	models, err := s.repo.GetByFamily(ctx, req.Family)
	if err != nil {
		log.Errorf("Failed to get family models: %v", err)
		return &pb.GetFamilyModelsResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	pbModels := make([]*pb.Model, len(models))
	for i, m := range models {
		pbModels[i] = convertToProtoModel(m)
	}

	return &pb.GetFamilyModelsResponse{
		Success: true,
		Message: "success",
		Data:    pbModels,
	}, nil
}

// GetComparisonCategories 获取对比类别
func (s *ModelService) GetComparisonCategories(ctx context.Context, req *commonpb.PaginationRequest) (*pb.GetComparisonCategoriesResponse, error) {
	log.Debug("GetComparisonCategories called")

	categories := getComparisonCategories()

	return &pb.GetComparisonCategoriesResponse{
		Success:    true,
		Message:    "success",
		Categories: categories,
	}, nil
}

// CompareModels 对比模型
func (s *ModelService) CompareModels(ctx context.Context, req *pb.CompareModelsRequest) (*pb.CompareModelsResponse, error) {
	log.Debugf("CompareModels called with ids=%v", req.Ids)

	if len(req.Ids) < 2 {
		return &pb.CompareModelsResponse{
			Success: false,
			Message: "至少需要 2 个模型进行对比",
		}, nil
	}

	if len(req.Ids) > 5 {
		return &pb.CompareModelsResponse{
			Success: false,
			Message: "最多只能对比 5 个模型",
		}, nil
	}

	models := make([]*repository.Model, 0, len(req.Ids))
	for _, id := range req.Ids {
		model, err := s.repo.GetByID(ctx, id)
		if err != nil {
			log.Warnf("Model not found: %s", id)
			continue
		}
		models = append(models, model)
	}

	if len(models) < 2 {
		return &pb.CompareModelsResponse{
			Success: false,
			Message: "有效的模型数量不足 2 个",
		}, nil
	}

	pbModels := make([]*pb.Model, len(models))
	for i, m := range models {
		pbModels[i] = convertToProtoModel(m)
	}

	return &pb.CompareModelsResponse{
		Success:             true,
		Message:             "success",
		Models:              pbModels,
		ComparisonCategories: getComparisonCategories(),
	}, nil
}

// convertToProtoModel 转换为 proto Model
func convertToProtoModel(m *repository.Model) *pb.Model {
	modalities := make([]commonpb.Modality, len(m.ModalitiesInput))
	for i, mod := range m.ModalitiesInput {
		modalities[i] = commonpb.Modality(commonpb.Modality_value[string(mod)])
	}

	return &pb.Model{
		Id:                 m.ID,
		Name:               m.Name,
		Family:             m.Family,
		Provider:           m.Provider,
		Description:        m.Description,
		Architecture:       m.Architecture,
		Knowledge:          m.Knowledge,
		ReleaseDate:        m.ReleaseDate,
		LastUpdated:        m.LastUpdated,
		Attachment:         m.Attachment,
		Reasoning:          m.Reasoning,
		ToolCall:          m.ToolCall,
		StructuredOutput:  m.StructuredOutput,
		Temperature:        m.Temperature,
		OpenWeights:        m.OpenWeights,
		ModalitiesInput:   modalities,
		ModalitiesOutput:  m.ModalitiesOutput,
		CostInput:         m.CostInput,
		CostOutput:        m.CostOutput,
		CostReasoning:     m.CostReasoning,
		CostCacheRead:     m.CostCacheRead,
		CostCacheWrite:    m.CostCacheWrite,
		CostInputAudio:    m.CostInputAudio,
		CostOutputAudio:   m.CostOutputAudio,
		LimitContext:      int32(m.LimitContext),
		LimitInput:        int32(m.LimitInput),
		LimitOutput:       int32(m.LimitOutput),
		InterleavedField:  m.InterleavedField,
	}
}

// getComparisonCategories 获取对比类别定义
func getComparisonCategories() []*pb.ComparisonCategory {
	return []*pb.ComparisonCategory{
		{
			Key:  "basic",
			Name: "基本信息",
			Items: []*pb.ComparisonItem{
				{Key: "name", Name: "名称", Type: "string"},
				{Key: "family", Name: "家族", Type: "string"},
				{Key: "releaseDate", Name: "发布日期", Type: "string"},
			},
		},
		{
			Key:  "capabilities",
			Name: "能力特性",
			Items: []*pb.ComparisonItem{
				{Key: "reasoning", Name: "推理", Type: "boolean"},
				{Key: "toolCall", Name: "工具调用", Type: "boolean"},
				{Key: "attachment", Name: "附件", Type: "boolean"},
				{Key: "openWeights", Name: "开源", Type: "boolean"},
			},
		},
		{
			Key:  "limits",
			Name: "限制",
			Items: []*pb.ComparisonItem{
				{Key: "limitContext", Name: "上下文", Type: "number", Unit: "tokens"},
				{Key: "limitInput", Name: "最大输入", Type: "number", Unit: "tokens"},
				{Key: "limitOutput", Name: "最大输出", Type: "number", Unit: "tokens"},
			},
		},
		{
			Key:  "pricing",
			Name: "定价",
			Items: []*pb.ComparisonItem{
				{Key: "costInput", Name: "输入", Type: "number", Unit: "$"},
				{Key: "costOutput", Name: "输出", Type: "number", Unit: "$"},
			},
		},
	}
}
