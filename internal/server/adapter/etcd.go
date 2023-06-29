package adapter

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Etcd struct {
	c *clientv3.Client
}

func NewETCDChecker(c *clientv3.Client) *Etcd {
	return &Etcd{c: c}
}

func (r *Etcd) Check(ctx context.Context) error {
	pctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if _, err := r.c.MemberList(pctx); err != nil {
		return err
	}

	return nil
}
