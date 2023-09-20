package service

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"

	pb "github.com/omalloc/kratos-agent/api/agent"
	"github.com/omalloc/kratos-agent/internal/biz"
)

type AgentService struct {
	pb.UnimplementedAgentServer

	log      *log.Helper
	cri      *biz.CRIUsecase
	clusters *biz.ClusterUsecase
	aa       *biz.CrontabUsecase
}

func NewAgentService(logger log.Logger, cri *biz.CRIUsecase, clusters *biz.ClusterUsecase) *AgentService {
	return &AgentService{
		log:      log.NewHelper(logger),
		cri:      cri,
		clusters: clusters,
		aa:       nil,
	}
}

func (s *AgentService) ListService(ctx context.Context, req *pb.ListServiceRequest) (*pb.ListServiceReply, error) {
	allService, err := s.clusters.GetServices(ctx)
	if err != nil {
		return nil, err
	}

	result := lo.FlatMap(allService, func(item *biz.ClusterMicroservice, _ int) []*pb.Microservice {
		return lo.Map(item.Services, func(cur *biz.Microservice, _ int) *pb.Microservice {
			return &pb.Microservice{
				Cluster:   item.Name,
				Id:        cur.ID,
				Key:       cur.Key,
				Name:      cur.Name,
				Version:   cur.Version,
				Endpoints: cur.Endpoints,
				Metadata:  cur.Metadata,
				Namespace: s.getFirstNameKey(cur.Key),
			}
		})
	})

	return &pb.ListServiceReply{
		Data: result,
	}, nil
}

func (s *AgentService) getFirstNameKey(key string) string {
	values := strings.Split(key, "/")
	// key-format: /microservices/{service-name}/{hostname}
	if len(values) > 1 {
		return values[1]
	}
	return ""
}

func (s *AgentService) ListServiceGroup(ctx context.Context, req *pb.ListServiceGroupRequest) (*pb.ListServiceGroupReply, error) {
	allService, err := s.clusters.GetServices(ctx)
	if err != nil {
		return nil, err
	}

	allData := lo.FlatMap(allService, func(item *biz.ClusterMicroservice, _ int) []*pb.Microservice {
		return lo.Map(item.Services, func(cur *biz.Microservice, _ int) *pb.Microservice {
			return &pb.Microservice{
				Cluster:   item.Name,
				Id:        cur.ID,
				Key:       cur.Key,
				Name:      cur.Name,
				Version:   cur.Version,
				Endpoints: cur.Endpoints,
				Metadata:  cur.Metadata,
			}
		})
	})

	mapData := lo.GroupBy(allData, func(item *pb.Microservice) string {
		return item.Name
	})

	results := lo.MapToSlice(mapData, func(key string, value []*pb.Microservice) *pb.MicroserviceGroup {
		var (
			keys      = make([]string, 0, len(value))
			endpoints = make([]string, 0, len(value)*2)
			clusters  = make([]string, 0, len(value))
			hostnames = make([]string, 0, len(value))
		)

		for _, service := range value {
			keys = append(keys, service.Key)
			endpoints = append(endpoints, service.Endpoints...)
			clusters = append(clusters, service.Cluster)
			hostnames = append(hostnames, service.Id)
		}

		return &pb.MicroserviceGroup{
			Name:      key,
			Keys:      keys,
			Endpoints: endpoints,
			Clusters:  lo.Uniq(clusters),
			Hostnames: lo.Uniq(hostnames),
		}
	})

	return &pb.ListServiceGroupReply{
		Data: results,
	}, nil
}

func (s *AgentService) ListCluster(ctx context.Context, req *pb.ListClusterRequest) (*pb.ListClusterReply, error) {
	clusters, err := s.clusters.GetClusters(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListClusterReply{
		Data: lo.Map(clusters, func(item *biz.Cluster, _ int) *pb.Cluster {
			return &pb.Cluster{
				Name:    item.Name,
				Healthy: item.Healthy,
				Members: lo.Map(item.Members, func(cur *biz.ClusterMember, _ int) *pb.Cluster_Member {
					return &pb.Cluster_Member{
						Name:       cur.Name,
						PeerUrls:   cur.PeerURLs,
						ClientUrls: cur.ClientURLs,
						IsLearner:  cur.IsLearner,
					}
				}),
			}
		}),
	}, nil
}

func (s *AgentService) Check(ctx context.Context) error {
	return s.clusters.Ping(ctx)
}
