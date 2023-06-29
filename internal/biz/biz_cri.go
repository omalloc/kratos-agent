package biz

import "github.com/go-kratos/kratos/v2/log"

// docker / cri info.
//
// features
// - list container
// - check container health
// - report container info to mallard2-agent
type CRIUsecase struct {
	log *log.Helper
}

func NewCRIUsecase(logger log.Logger) *CRIUsecase {
	return &CRIUsecase{log: log.NewHelper(logger)}
}
