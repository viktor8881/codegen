package http

import "github.com/viktor8881/codegen/command/codegen"

const TmplServerEndpointFile = codegen.TmplCodeGeneratorPhrase + `
package server

import (
	"context"
	"github.com/viktor8881/service-utilities/http/server"
	"go.uber.org/zap"
	"{{.PackageName}}/generated"
)


{{- range .Endpoints }}

{{ . }}
{{- end }}

`

const TmplServerEndpoint = `
func {{.Name}}(
	t *server.Transport,
	decodeFn server.DecodeRequestFunc,
	serviceFn func(ctx context.Context, in *generated.{{.InputRequest}}) (*generated.{{.OutputResponse}}, error),
	encodeFn server.EncodeResponseFunc,
	errorHandlerFn server.ErrorHandlerFunc,
	logger *zap.Logger,	
	middlewares ...server.Middleware,
) {
	t.AddEndpoint(
		"{{.Url}}",
		"{{.Method}}",
		(*generated.{{.InputRequest}})(nil),
		decodeFn,
		func(ctx context.Context, in interface{}) (interface{}, error) {
			return serviceFn(ctx, in.(*generated.{{.InputRequest}}))
		},
		encodeFn,
		errorHandlerFn,
		logger,	
		append(middlewares, server.LoggerMiddleware(logger))...,
	)
}
`

const TmplClientEndpointFile = codegen.TmplCodeGeneratorPhrase + `
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/viktor8881/service-utilities/http/client"
	"{{.PackageName}}/generated"
)

type Client struct {
	client *client.Client
}

func NewClient(client *client.Client) *Client {
	return &Client{client: client}
}

{{- range .Endpoints }}

{{ . }}
{{- end }}
`

const TmplClientEndpoint = `
func (c *Client){{.Name}}(
	ctx context.Context, in *generated.{{.InputRequest}}) (*generated.{{.OutputResponse}}, error) {
	var dest generated.{{.OutputResponse}}

	resp, err := c.client.{{toCamelCase .Method}}(ctx, "{{.Url}}", in, nil)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if err := json.NewDecoder(resp.Body).Decode(&dest); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &dest, nil
}
`

const TmplAddCodeToRouterFile = `
generated.{{.Name}}(
		tr,
		server.DecodeRequest,
		{{.ServiceNameToLower}}.NewService().{{.ServiceMethod}},
		server.EncodeResponse,
		server.ErrorHandler,
		logger,	
	)
`
