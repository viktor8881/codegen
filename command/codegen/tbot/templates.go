package tbot

import (
	"github.com/viktor8881/codegen/command/codegen"
)

const TmplTbotEndpointFile = codegen.TmplCodeGeneratorPhrase + `
package tbot

import (
	"context"
	"github.com/viktor8881/service-utilities/tbot"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
	"{{.PackageName}}/generated"
)

{{- range .Endpoints }}

{{ . }}
{{- end }}

`

const TmplTbotEndpoint = `
func {{.Name}}(
	ctx context.Context,
	t *tbot.Bot,
	decodeFn tbot.DecodePayloadFunc,
	serviceFn func(ctx context.Context, in *generated.{{.InputRequest}}) (*generated.{{.OutputResponse}}, error),
	encodeFn tbot.EncodeResponseFunc,
	errorHandlerFn tbot.ErrorHandlerFunc,
	logger *zap.Logger,
	middlewares ...tbot.Middleware,
) {
	t.AddCommandHandler(
		"{{.Url}}",	
		(*generated.{{.InputRequest}})(nil),
		decodeFn,
		func(c telebot.Context, in any) (any, error) {
			return serviceFn(ctx, in.(*generated.{{.InputRequest}}))
		},
		encodeFn,
		errorHandlerFn,
		logger,	
		append(middlewares, tbot.LoggerMiddleware(logger))...,
	)
}
`

const TmplAddCodeToRouterFile = `
tbotgenerated.{{.Name}}(
		ctx,
		newTbot,
		tbot.DecodePayload,
		{{.ServiceNameToLower}}.NewService().{{.ServiceMethod}},	
		tbot.EncodeResponse,
		tbot.ErrorHandler,	
		logger,
	)
`
