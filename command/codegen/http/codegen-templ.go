package http

const tmplServerEndpointFile = `package server

import (
	"context"
	customHttp "github.com/viktor8881/service-utilities/http/server"
	"go.uber.org/zap"
	"net/http"
)


{{- range .Endpoints }}

{{ . }}
{{- end }}

`

const tmplServerEndpoint = `
func {{.Name}}(
	t *customHttp.Transport,
	handlerFn func(ctx context.Context, in *{{.InputRequest}}) (*{{.OutputResponse}}, error),
	logger *zap.Logger,
	fErrorHandler func(w http.ResponseWriter, r *http.Request, err error, logger *zap.Logger),
) {
	t.AddEndpoint(
		"{{.Url}}",
		"{{.Method}}",
		(*{{.InputRequest}})(nil),
		func(ctx context.Context, in interface{}) (interface{}, error) {
			return handlerFn(ctx, in.(*{{.InputRequest}}))
		},
		logger,
		fErrorHandler,
	)
}
`

const tmplClientEndpointFile = `package client

import (
	"context"
	"encoding/json"
	simpleClient "github.com/viktor8881/service-utilities/http/client"	
)

{{- range .Endpoints }}

{{ . }}
{{- end }}
`

const tmplClientEndpoint = `
func {{.Name}}(
	ctx context.Context, client *simpleClient.SimpleClient, in any) (*{{.OutputResponse}}, error) {
	var dest {{.OutputResponse}}

	url, err := simpleClient.BuildURL("{{.Url}}", in)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if err := json.NewDecoder(resp.Body).Decode(&dest); err != nil {
		return nil, err
	}

	return &dest, nil
}
`

const tmplLogicServiceFile = `package {{.ServiceNameToLower}}

import (
	"context"
	"errors"
	simplehttp "github.com/viktor8881/service-utilities/http/client"
	generated "{{.PackageName}}/generated/http/server"
)

type {{.ServiceName}}Service struct {}
`

const tmplLogicServiceEndpoint = `
func {{.ServiceMethod}}(ctx context.Context, in *generated.{{.InputRequest}}) (*generated.{{.OutputResponse}}, error) {
	return *generated.{{.OutputResponse}}{}, errors.New("not implemented")
}
`
