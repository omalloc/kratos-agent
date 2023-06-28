package biz

import "github.com/go-kratos/kratos/v2/log"

type CRIUsecase struct {
	log *log.Helper
}

func NewCRIUsecase(logger log.Logger) *CRIUsecase {
	return &CRIUsecase{log: log.NewHelper(logger)}
}
