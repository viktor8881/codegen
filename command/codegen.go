package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/viktor8881/codegen/command/codegen/http"
)

type Endpoint struct {
	Name           string
	Description    string
	Url            string
	Method         string
	ServiceName    string
	ServiceMethod  string
	InputRequest   string
	OutputResponse string
}

func CodeGen(cmd *cobra.Command, args []string) {
	err := http.GenerateHttpServerFile("/http/server/")
	if err != nil {
		fmt.Println("err GenerateHttpServerFile: ", err)
		return
	}

	err = http.GenerateHttpClientFile("/http/client/")
	if err != nil {
		fmt.Println("err GenerateHttpClientFile: ", err)
		return
	}
}
