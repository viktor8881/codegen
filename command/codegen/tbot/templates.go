package tbot

import (
	"github.com/viktor8881/codegen/command/codegen/http"
)

const TmplTbotEndpointFile = http.TmplCodeGeneratorPhrase + `
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
	t *tbot.CustomBot,
	decodeFn tbot.DecodePayloadFunc,
	handlerFn tbog.HandlerFunc,
	encodeFn tbot.EncodeResponseFunc,
	errorHandlerFn tbog.ErrorHandlerFunc,
	logger *zap.Logger,
	middlewares ...tbot.Middleware,
) {
	t.AddCommandHandler(
		"{{.Url}}",	
		(*generated.{{.InputRequest}})(nil),
		decodeFn,
		func(c telebot.Context, in any) (any, error) {
			return handlerFn(ctx, in.(*generated.{{.InputRequest}}))
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
