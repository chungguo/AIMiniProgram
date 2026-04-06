package service

import (
	"context"

	"modellens/backend/internal/repository"
	pb "modellens/proto/paper/v1"
	commonpb "modellens/proto/common/v1"

	"trpc.group/trpc-go/trpc-go/log"
)

// PaperService 实现
type PaperService struct {
	repo repository.PaperRepository
}

// NewPaperService 创建服务实例
func NewPaperService() *PaperService {
	return &PaperService{
		repo: repository.NewPaperRepository(),
	}
}

// ListPapers 获取论文列表
func (s *PaperService) ListPapers(ctx context.Context, req *pb.ListPapersRequest) (*pb.ListPapersResponse, error) {
	log.Debugf("ListPapers called with page=%d, limit=%d", 
		req.Pagination.Page, req.Pagination.Limit)

	filter := &repository.PaperFilter{}
	if req.Filter != nil {
		filter.Search = req.Filter.Search
		filter.Author = req.Filter.Author
	}

	papers, total, err := s.repo.List(ctx, filter, req.Pagination.Page, req.Pagination.Limit)
	if err != nil {
		return &pb.ListPapersResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	pbPapers := make([]*pb.Paper, len(papers))
	for i, p := range papers {
		pbPapers[i] = convertToProtoPaper(p)
	}

	return &pb.ListPapersResponse{
		Success: true,
		Message: "success",
		Data:    pbPapers,
		Pagination: &commonpb.Pagination{
			Page:       req.Pagination.Page,
			Limit:      req.Pagination.Limit,
			Total:      int32(total),
			TotalPages: (int32(total) + req.Pagination.Limit - 1) / req.Pagination.Limit,
		},
	}, nil
}

// GetPaper 获取单个论文
func (s *PaperService) GetPaper(ctx context.Context, req *pb.GetPaperRequest) (*pb.GetPaperResponse, error) {
	log.Debugf("GetPaper called with id=%s", req.Id)

	paper, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return &pb.GetPaperResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.GetPaperResponse{
		Success: true,
		Message: "success",
		Data:    convertToProtoPaper(paper),
	}, nil
}

// GetLatestPapers 获取最新论文
func (s *PaperService) GetLatestPapers(ctx context.Context, req *commonpb.PaginationRequest) (*pb.GetLatestPapersResponse, error) {
	log.Debug("GetLatestPapers called")

	papers, err := s.repo.GetLatest(ctx, req.Limit)
	if err != nil {
		return &pb.GetLatestPapersResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	pbPapers := make([]*pb.Paper, len(papers))
	for i, p := range papers {
		pbPapers[i] = convertToProtoPaper(p)
	}

	return &pb.GetLatestPapersResponse{
		Success: true,
		Message: "success",
		Data:    pbPapers,
	}, nil
}

// convertToProtoPaper 转换为 proto Paper
func convertToProtoPaper(p *repository.Paper) *pb.Paper {
	return &pb.Paper{
		Id:         p.ID,
		Title:      p.Title,
		TitleCn:    p.TitleCn,
		Author:     p.Author,
		Abstract:   p.Abstract,
		AbstractCn: p.AbstractCn,
		SubmitAt:   p.SubmitAt,
	}
}
