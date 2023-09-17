package plugins

import (
	"context"
	"net/http"

	"github.com/ZiXian92/my-helper-bot/plugins/interfaces"
	"github.com/ZiXian92/my-helper-bot/plugins/proto"
	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type WebHandlerPlugin struct {
	plugin.Plugin
	Impl interfaces.WebHandler
}

func (p *WebHandlerPlugin) GRPCServer(b *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterWebHandlerServer(s, &WebHandlerGRPCServer{Impl: p.Impl})
	return nil
}

func (p *WebHandlerPlugin) GRPCClient(ctx context.Context, b *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &WebHandlerGRPCClient{client: proto.NewWebHandlerClient(c)}, nil
}

type WebHandlerGRPCServer struct {
	proto.UnimplementedWebHandlerServer
	Impl interfaces.WebHandler
}

func (s *WebHandlerGRPCServer) GetEndpoints(ctx context.Context, msg *proto.Empty) (*proto.WebEndpointList, error) {
	goEndpoints := s.Impl.GetEndpoints()
	grpcEndpoints := []*proto.WebEndpoint{}
	for _, ep := range goEndpoints {
		grpcEndpoints = append(grpcEndpoints, &proto.WebEndpoint{
			Name:    ep.Name,
			Methods: ep.Methods,
			Path:    ep.Path,
		})
	}
	return &proto.WebEndpointList{Endpoints: grpcEndpoints}, nil
}

func (s *WebHandlerGRPCServer) HandleRequest(ctx context.Context, r *proto.WebRequest) (*proto.WebResponse, error) {
	processedReqHeaders := map[string][]string{}
	for k, vList := range r.Headers {
		processedReqHeaders[k] = vList.Values
	}
	processedQueryParams := map[string][]string{}
	for k, vList := range r.QueryParams {
		processedQueryParams[k] = vList.Values
	}
	res := s.Impl.HandleRequest(interfaces.WebRequest{
		EndPointName: r.EndpointName,
		Headers:      processedReqHeaders,
		URIParams:    r.UriParams,
		QueryParams:  processedQueryParams,
		Body:         r.Body,
	})

	protobufResHeaders := map[string]*proto.StringList{}
	for k, vList := range res.Headers {
		protobufResHeaders[k] = &proto.StringList{Values: vList}
	}
	return &proto.WebResponse{
		Code:    int32(res.Code),
		Headers: protobufResHeaders,
		Body:    res.Body,
	}, nil
}

type WebHandlerGRPCClient struct {
	client proto.WebHandlerClient
}

func (c *WebHandlerGRPCClient) GetEndpoints() []interfaces.WebEndpoint {
	rawEpList, err := c.client.GetEndpoints(context.Background(), &proto.Empty{})
	if err != nil {
		return nil
	}
	endpoints := []interfaces.WebEndpoint{}
	for _, rawEp := range rawEpList.Endpoints {
		endpoints = append(endpoints, interfaces.WebEndpoint{
			Name:    rawEp.Name,
			Methods: rawEp.Methods,
			Path:    rawEp.Path,
		})
	}
	return endpoints
}

func (c *WebHandlerGRPCClient) HandleRequest(r interfaces.WebRequest) interfaces.WebResponse {
	protobufReqHeaders := map[string]*proto.StringList{}
	for k, vList := range r.Headers {
		protobufReqHeaders[k] = &proto.StringList{Values: vList}
	}
	protobufReqQueryParams := map[string]*proto.StringList{}
	for k, vList := range r.QueryParams {
		protobufReqQueryParams[k] = &proto.StringList{Values: vList}
	}
	protobufRes, err := c.client.HandleRequest(context.Background(), &proto.WebRequest{
		EndpointName: r.EndPointName,
		Headers:      protobufReqHeaders,
		UriParams:    r.URIParams,
		QueryParams:  protobufReqQueryParams,
		Body:         r.Body,
	})
	if err != nil {
		return interfaces.WebResponse{
			Code: http.StatusInternalServerError,
			Body: []byte("Server error processing request."),
		}
	}

	processedResHeaders := map[string][]string{}
	for k, vList := range protobufRes.Headers {
		processedResHeaders[k] = vList.Values
	}
	return interfaces.WebResponse{
		Code:    int(protobufRes.Code),
		Headers: processedResHeaders,
		Body:    protobufRes.Body,
	}
}
