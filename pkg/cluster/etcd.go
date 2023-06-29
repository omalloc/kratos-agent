package cluster

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/omalloc/kratos-agent/internal/conf"
)

// Client is a wrapper of etcd clientv3.Client
type Client struct {
	*clientv3.Client

	Name    string
	Enabled bool
}

func NewClients(logger log.Logger, bc *conf.Bootstrap) ([]*Client, error) {
	clog := log.NewHelper(logger)
	var clis = make([]*Client, 0, len(bc.Clusters))

	for _, cluster := range bc.Clusters {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   cluster.Endpoints,
			DialTimeout: 2 * time.Second,
		})

		if err != nil {
			clog.Warnf("etcd cluster %s connect failed: %v", cluster.Name, err)
			continue
		}

		clis = append(clis, &Client{
			Client:  client,
			Name:    cluster.Name,
			Enabled: cluster.Enabled,
		})
	}

	return clis, nil
}
