package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/viktor8881/codegen/command/codegen"
	"github.com/viktor8881/codegen/command/codegen/http"
	"github.com/viktor8881/codegen/command/codegen/tbot"
	"log"
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
	err := prepare()
	if err != nil {
		log.Fatalf("failed to prepare: %v", err)
	}

	// gen server
	if codegen.FileExists("./contracts/http/server/endpoints.json") {
		err = http.GenerateHttpServerFile("/http/server/")
		if err != nil {
			fmt.Println("err GenerateHttpServerFile: ", err)
			return
		}
	}

	// gen client
	if codegen.FileExists("./contracts/http/client/endpoints.json") {
		err = http.GenerateHttpClientFile("/http/client/")
		if err != nil {
			fmt.Println("err GenerateHttpClientFile: ", err)
			return
		}
	}

	// gen tbot
	if codegen.FileExists("./contracts/tbot/endpoints.json") {
		err = tbot.GenerateTbotClientFile("/tbot")
		if err != nil {
			fmt.Println("err GenerateTbotClientFile: ", err)
			return
		}
	}
}

func prepare() error {
	err := codegen.CreateDirIfNeed("./generated")
	if err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	// copy models.go
	if err := codegen.CopyFile("./contracts/models.go", "./generated/models.go", codegen.TmplCodeGeneratorPhrase); err != nil {
		return fmt.Errorf("failed to copy models.go: %w", err)
	}

	return nil
}
