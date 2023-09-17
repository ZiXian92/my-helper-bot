package interfaces

type WebEndpoint struct {
	Name    string
	Methods []string
	Path    string
}

type WebRequest struct {
	EndPointName string
	Headers      map[string][]string
	URIParams    map[string]string
	QueryParams  map[string][]string
	Body         []byte
}

type WebResponse struct {
	Code    int
	Headers map[string][]string
	Body    []byte
}

type WebHandler interface {
	GetEndpoints() []WebEndpoint
	HandleRequest(WebRequest) WebResponse
}
