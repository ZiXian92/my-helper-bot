package plugins

import "github.com/hashicorp/go-plugin"

const (
	PluginMapKeyWebHandler = "webhandler"
)

// HandshakeConfig is to be used by all plugin implementers wishing
// to communicate with this main application.
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "my-helper-bot",
	MagicCookieValue: "my-helper-bot-plugins",
}

// PluginMap is to be used by the main application.
// For plugin implementers, define your own map declaring only for keys that your plugin is implementing.
// This acts as a reference to which plugin map keys are used for which plugin type.
var PluginMap = map[string]plugin.Plugin{
	PluginMapKeyWebHandler: &WebHandlerPlugin{},
}
