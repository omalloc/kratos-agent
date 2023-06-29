package biz_test

import (
	"context"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"

	"github.com/omalloc/kratos-agent/internal/biz"
	"github.com/omalloc/kratos-agent/internal/conf"
	"github.com/omalloc/kratos-agent/pkg/cluster"
)

func loadConfig() *conf.Bootstrap {

	c := config.New(
		config.WithSource(
			file.NewSource("../../configs"),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	return &bc
}

func TestGetClusters(t *testing.T) {
	bc := loadConfig()

	clis, err := cluster.NewClients(log.GetLogger(), bc)
	if err != nil {
		t.Fatal(err)
	}

	clusterUsecase := biz.NewClusterUsecase(log.GetLogger(), clis)

	if clusters, err := clusterUsecase.GetClusters(context.TODO()); err == nil {
		for _, c := range clusters {
			println("cluster name:", c.Name, "\tmembers:", strings.Join(lo.Map(c.Members, func(item *biz.ClusterMember, _ int) string { return item.Name }), ","))
		}
	}
}

func TestGetServices(t *testing.T) {
	bc := loadConfig()

	clis, err := cluster.NewClients(log.GetLogger(), bc)
	if err != nil {
		t.Fatal(err)
	}

	clusterUsecase := biz.NewClusterUsecase(log.GetLogger(), clis)

	if services, err := clusterUsecase.GetServices(context.TODO()); err == nil {
		for _, s := range services {
			println("cluster name:", s.Name, "\tservice count:", len(s.Services))
			for _, m := range s.Services {
				println("service name:", m.Name, "\tversion:", m.Version, "\tendpoints:", strings.Join(m.Endpoints, ","))
			}
		}
	}
}
