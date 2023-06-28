package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	pb "github.com/omalloc/kratos-agent/api/agent"
	"github.com/omalloc/kratos-agent/internal/biz"
)

type AgentService struct {
	pb.UnimplementedAgentServer

	log *log.Helper
	cri *biz.CRIUsecase
}

func NewAgentService(logger log.Logger, cri *biz.CRIUsecase) *AgentService {
	return &AgentService{
		log: log.NewHelper(logger),
		cri: cri,
	}
}

func (s *AgentService) ListService(ctx context.Context, req *pb.ListServiceRequest) (*pb.ListServiceReply, error) {
	return &pb.ListServiceReply{}, nil
}

func (s *AgentService) Check(ctx context.Context) error {
	return nil
}
