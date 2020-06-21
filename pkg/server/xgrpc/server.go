// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xgrpc

import (
	"context"
	"net"

	"github.com/douyu/jupiter/pkg"

	"github.com/douyu/jupiter/pkg/server"
	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	*grpc.Server
	*Config
}

func newServer(config *Config) *Server {
	if len(config.streamInterceptors) > 0 {
		config.serverOptions = append(config.serverOptions, grpc.StreamInterceptor(StreamInterceptorChain(config.streamInterceptors...)))
	}
	if len(config.unaryInterceptors) > 0 {
		config.serverOptions = append(config.serverOptions, grpc.UnaryInterceptor(UnaryInterceptorChain(config.unaryInterceptors...)))
	}
	newServer := grpc.NewServer(config.serverOptions...)
	return &Server{Server: newServer, Config: config}
}

// Serve implements server.Serve interface.
func (s *Server) Serve() error {
	l, err := net.Listen(s.Network, s.Address())
	if err != nil {
		return err
	}
	s.Port = l.Addr().(*net.TCPAddr).Port
	err = s.Server.Serve(l)
	if err == grpc.ErrServerStopped {
		return nil
	}
	return err
}

// Stop implements server.Server interface
// it will terminate echo server immediately
func (s *Server) Stop() error {
	s.Server.Stop()
	return nil
}

// GracefulStop implements server.Server interface
// it will stop echo server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	s.Server.GracefulStop()
	return nil
}

// Info returns server info, used by governor and consumer balancer
// TODO(gorexlv): implements government protocol with juno
func (s *Server) Info() *server.ServiceInfo {
	return &server.ServiceInfo{
		Name:      pkg.Name(),
		Scheme:    "grpc",
		IP:        s.Host,
		Port:      s.Port,
		Weight:    0.0,
		Enable:    true,
		Healthy:   true,
		Metadata:  map[string]string{},
		Region:    "",
		Zone:      "",
		GroupName: "",
	}
}
