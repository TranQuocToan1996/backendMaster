package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (s *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("grpc metadata empty")
		return mtdt
	}

	log.Printf("metadata: %+v", md)
	if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) != 0 {
		mtdt.UserAgent = userAgents[0]
	}

	if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) != 0 {
		mtdt.ClientIP = clientIPs[0]
	}

	if userAgents := md.Get(userAgentHeader); len(userAgents) != 0 {
		mtdt.UserAgent = userAgents[0]
	}

	if peers, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = peers.Addr.String()
	}

	return mtdt
}
