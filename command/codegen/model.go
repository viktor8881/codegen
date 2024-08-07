package codegen

type Endpoint struct {
	Name           string
	Description    string
	Url            string // url for http, command for tbot
	Method         string
	ServiceName    string
	ServiceMethod  string
	InputRequest   string
	OutputResponse string
}
