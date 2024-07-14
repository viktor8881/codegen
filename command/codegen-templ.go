package command

const tmplStr = `
package generated

import (
	"context"
	customHttp "github.com/viktor8881/service-utilities/http/server"
	"go.uber.org/zap"
	"net/http"
)


{{- range .Endpoints }}

{{ . }}
{{- end }}

{{- range .Models }}

{{ . }}
{{- end }}
`

const tmplEndpoint = `
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
