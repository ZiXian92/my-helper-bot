package main

import (
	"net/http"

	"github.com/ZiXian92/my-helper-bot/plugins"
	"github.com/ZiXian92/my-helper-bot/plugins/interfaces"
	"github.com/hashicorp/go-plugin"
)

type healthcheck struct{}

func (hc *healthcheck) GetEndpoints() []interfaces.WebEndpoint {
	return []interfaces.WebEndpoint{
		{
			Name:    "healthcheck",
			Methods: []string{"GET"},
			Path:    "/healthz",
		},
	}
}

func (hc *healthcheck) HandleRequest(r interfaces.WebRequest) interfaces.WebResponse {
	return interfaces.WebResponse{
		Code: http.StatusOK,
		Body: []byte("OK"),
	}
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			plugins.PluginMapKeyWebHandler: &plugins.WebHandlerPlugin{Impl: &healthcheck{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
