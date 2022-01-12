package grpc

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"log"
	"net"

	structpb "github.com/golang/protobuf/ptypes/struct"
	manager_pluginpb "github.com/xellos00/silver-bentonville/dist/proto/dsrv/api/node_manager/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
	grpcPort = 9091
)

type grpcService struct {
	manager_pluginpb.UnimplementedManagerPluginServer
}

func (s *grpcService) Init(ctx context.Context, in *emptypb.Empty) (*manager_pluginpb.PluginInfo, error) {
	var metaData = &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"Protocol": &structpb.Value{
				Kind: &structpb.Value_StringValue{
					StringValue: "Near",
				},
			},
			"Service": &structpb.Value{
				Kind: &structpb.Value_StringValue{
					StringValue: "Node Manager",
				},
			},
		},
	}
	resp := manager_pluginpb.PluginInfo{
		Metadata: metaData,
	}
	return &resp, nil
}

func (s *grpcService) Verify(ctx context.Context, in *emptypb.Empty) (*manager_pluginpb.VerifyInfo, error) {
	resp := manager_pluginpb.VerifyInfo{
		VerifyMsg: "Verified",
	}

	return &resp, nil
}

func (s *grpcService) Execute(ctx context.Context, in *manager_pluginpb.ExecuteRequest) (*manager_pluginpb.ExecuteResponse, error) {

	//TODO
	// 이곳을 기점으로 어떻게 플러그인 기준을 어떻게 MVC 기준으로 가져갈지 생각할 수 있게 한다.?
	// 아니면 그것은 자율에 맞긴다.

	c, b := exec.Command("ls", "-al"), new(strings.Builder)
	c.Stdout = b
	c.Run()

	return &manager_pluginpb.ExecuteResponse{
		State:        manager_pluginpb.ExecuteResponse_SUCCESS,
		Message:      b.String(),
		ResourceType: "Go",
	}, nil
}

// StartServer try to start grpc service.
func StartServer() error {
	s := grpc.NewServer()
	serv := grpcService{}

	manager_pluginpb.RegisterManagerPluginServer(s, &serv)
	reflection.Register(s)

	addr := fmt.Sprintf(":%d", grpcPort)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
		return err
	}

	log.Println("listen ", addr)

	go func() {
		if err := s.Serve(l); err != nil {
			log.Panic(err)
		}
	}()

	return nil
}
