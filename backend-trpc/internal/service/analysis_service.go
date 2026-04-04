package service

import (
	"context"

	"aiminiprogram/backend-trpc/internal/repository"
	pb "aiminiprogram/proto/analysis/v1"
	commonpb "aiminiprogram/proto/common/v1"

	"trpc.group/trpc-go/trpc-go/log"
)

// AnalysisService 实现
type AnalysisService struct {
	repo repository.AnalysisRepository
}

// NewAnalysisService 创建服务实例
func NewAnalysisService() *AnalysisService {
	return &AnalysisService{
		repo: repository.NewAnalysisRepository(),
	}
}

// ListArtificialAnalysis 获取评测数据列表
func (s *AnalysisService) ListArtificialAnalysis(ctx context.Context, req *pb.ListArtificialAnalysisRequest) (*pb.ListArtificialAnalysisResponse, error) {
	log.Debug("ListArtificialAnalysis called")

	analysisList, total, err := s.repo.List(ctx, req.Pagination.Page, req.Pagination.Limit)
	if err != nil {
		return &pb.ListArtificialAnalysisResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	pbAnalysis := make([]*pb.ArtificialAnalysis, len(analysisList))
	for i, a := range analysisList {
		pbAnalysis[i] = convertToProtoAnalysis(a)
	}

	return &pb.ListArtificialAnalysisResponse{
		Success: true,
		Message: "success",
		Data:    pbAnalysis,
		Pagination: &commonpb.Pagination{
			Page:       req.Pagination.Page,
			Limit:      req.Pagination.Limit,
			Total:      int32(total),
			TotalPages: (int32(total) + req.Pagination.Limit - 1) / req.Pagination.Limit,
		},
	}, nil
}

// GetArtificialAnalysisBySlug 通过 slug 获取评测数据
func (s *AnalysisService) GetArtificialAnalysisBySlug(ctx context.Context, req *pb.GetArtificialAnalysisBySlugRequest) (*pb.GetArtificialAnalysisBySlugResponse, error) {
	log.Debugf("GetArtificialAnalysisBySlug called with slug=%s", req.Slug)

	analysis, err := s.repo.GetBySlug(ctx, req.Slug)
	if err != nil {
		return &pb.GetArtificialAnalysisBySlugResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.GetArtificialAnalysisBySlugResponse{
		Success: true,
		Message: "success",
		Data:    convertToProtoAnalysis(analysis),
	}, nil
}

// GetModelWithAnalysis 获取模型与评测数据
func (s *AnalysisService) GetModelWithAnalysis(ctx context.Context, req *pb.GetModelWithAnalysisRequest) (*pb.GetModelWithAnalysisResponse, error) {
	log.Debugf("GetModelWithAnalysis called with model_id=%s", req.ModelId)

	// TODO: 实现组合查询
	return &pb.GetModelWithAnalysisResponse{
		Success: false,
		Message: "not implemented",
	}, nil
}

// convertToProtoAnalysis 转换为 proto Analysis
func convertToProtoAnalysis(a *repository.ArtificialAnalysis) *pb.ArtificialAnalysis {
	return &pb.ArtificialAnalysis{
		Id:                                a.ID,
		Slug:                              a.Slug,
		ModelCreator:                      a.ModelCreator,
		ArtificialAnalysisIntelligenceIndex: a.ArtificialAnalysisIntelligenceIndex,
		ArtificialAnalysisCodingIndex:     a.ArtificialAnalysisCodingIndex,
		ArtificialAnalysisMathIndex:       a.ArtificialAnalysisMathIndex,
		MmluPro:                           a.MmluPro,
		Gpqa:                              a.Gpqa,
		Hle:                               a.Hle,
		Livecodebench:                     a.Livecodebench,
		Scicode:                           a.Scicode,
		Math_500:                          a.Math500,
		Aime:                              a.Aime,
		Aime_25:                           a.Aime25,
		Ifbench:                           a.Ifbench,
		Lcr:                               a.Lcr,
		TerminalbenchHard:                 a.TerminalbenchHard,
		Tau2:                              a.Tau2,
		Price_1m_blended_3_to_1:           a.Price1mBlended31,
		Price_1m_input_tokens:             a.Price1mInputTokens,
		Price_1m_output_tokens:            a.Price1mOutputTokens,
		Median_output_tokens_per_second:   a.MedianOutputTokensPerSecond,
		Median_time_to_first_token_seconds: a.MedianTimeToFirstTokenSeconds,
		Median_time_to_first_answer_token: a.MedianTimeToFirstAnswerToken,
	}
}
