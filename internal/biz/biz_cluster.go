package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"
	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/connectivity"

	"github.com/omalloc/kratos-agent/pkg/cluster"
)

var (
	ErrEtcdClientNotReady = errors.New("etcd client not ready")
)

// etcd clusters info.
//
// features
// - get etcd cluster info
// - get all service by microservice namespace
// - get all election leader key

type Microservice struct {
	Key       string            `json:"key"`
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []string          `json:"endpoints"`
	CreatedAt int64             `json:"created_at"`
}

type ClusterMicroservice struct {
	Name     string          `json:"name"`
	Services []*Microservice `json:"services"`
}

type ClusterMember struct {
	ID         uint64   `json:"id"`
	Name       string   `json:"name"`
	PeerURLs   []string `json:"peer_urls"`
	ClientURLs []string `json:"client_urls"`
	IsLearner  bool     `json:"is_learner"`
}

type Cluster struct {
	Name    string           `json:"name"`    // 集群名称
	Healthy bool             `json:"healthy"` // 是否健康
	Members []*ClusterMember `json:"members"` // 集群成员
}

type ClusterUsecase struct {
	log  *log.Helper
	clis []*cluster.Client
}

func NewClusterUsecase(logger log.Logger, clis []*cluster.Client) *ClusterUsecase {
	return &ClusterUsecase{
		log:  log.NewHelper(logger),
		clis: clis,
	}
}

func (r *ClusterUsecase) GetClusters(ctx context.Context) ([]*Cluster, error) {
	ctx, span := tracer.Start(ctx,
		fmt.Sprintf("func %s.%s", "ClusterUsecase", "GetClusters"),
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer func() {
		if err := recover(); err != nil {
			span.RecordError(err.(error))
			span.SetStatus(codes.Error, err.(error).Error())
		}

		span.End()
	}()

	var (
		clusterInfo = make([]*Cluster, 0, len(r.clis))
		wg          = sync.WaitGroup{}
	)
	wg.Add(len(r.clis))

	for _, cli := range r.clis {
		go func(client *cluster.Client) {
			defer wg.Done()

			currCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()

			members, err := getMembers(currCtx, client)
			if err != nil {
				clusterInfo = append(clusterInfo, &Cluster{
					Name:    client.Name,
					Healthy: false,
				})
				return
			}

			clusterInfo = append(clusterInfo, &Cluster{
				Name:    client.Name,
				Healthy: true,
				Members: lo.Map(members, func(item *pb.Member, _ int) *ClusterMember {
					return &ClusterMember{
						ID:         item.ID,
						Name:       item.Name,
						PeerURLs:   item.PeerURLs,
						ClientURLs: item.ClientURLs,
						IsLearner:  item.IsLearner,
					}
				}),
			})
		}(cli)
	}

	wg.Wait()

	return clusterInfo, nil
}

func (r *ClusterUsecase) GetServices(ctx context.Context) ([]*ClusterMicroservice, error) {
	var (
		cm = make([]*ClusterMicroservice, 0, len(r.clis))
	)
	for _, cli := range r.clis {
		srvs, err := getServices(ctx, cli)
		if err != nil {
			continue
		}

		cm = append(cm, &ClusterMicroservice{
			Name:     cli.Name,
			Services: srvs,
		})
	}

	return cm, nil
}

func (r *ClusterUsecase) Ping(ctx context.Context) error {
	for _, cli := range r.clis {
		if cli.ActiveConnection().GetState() != connectivity.Ready {
			return ErrEtcdClientNotReady
		}
	}
	return nil
}

func getMembers(ctx context.Context, cli *cluster.Client) ([]*pb.Member, error) {
	ctx, span := tracer.Start(ctx,
		fmt.Sprintf("func %s.%s", "ClusterUsecase", "getMembers"),
		trace.WithSpanKind(trace.SpanKindServer),
	)

	resp, err := cli.MemberList(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}

		span.End()
	}()

	span.SetAttributes(
		attribute.String("cluster.name", cli.Name),
		attribute.String("cluster.members.size", fmt.Sprintf("%d", len(resp.Members))),
	)

	return resp.Members, nil
}

func getServices(ctx context.Context, cli *cluster.Client) ([]*Microservice, error) {
	resp, err := cli.Get(ctx, "/microservices", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var services = make([]*Microservice, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var v Microservice
		if err := json.Unmarshal(kv.Value, &v); err != nil {
			continue
		}
		v.Key = string(kv.Key)
		services = append(services, &v)
	}

	return services, nil
}
