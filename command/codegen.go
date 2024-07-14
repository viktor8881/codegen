package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"text/template"
	"unicode"
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
	tmpl, err := template.New("gen_structure").Parse(tmplStr)
	if err != nil {
		fmt.Println("err read template: ", err)
		return
	}

	data := struct {
		Endpoints []string
		Models    []string
	}{}

	endpoints, err := HttpEndpoints()
	if err != nil {
		fmt.Println("err call HttpEndpoints: ", err)
	}

	fmt.Println("find HttpEndpoints: ", len(endpoints))
	for _, e := range endpoints {
		data.Endpoints = append(data.Endpoints, e)
	}

	models, err := ParceFile()
	if err != nil {
		fmt.Println("err call ParceFile: ", err)
	}

	fmt.Println("find models: ", len(models))
	for _, m := range models {
		data.Models = append(data.Models, m)
	}

	// Создание директории, если она не существует
	if err := os.MkdirAll("./generated", os.ModePerm); err != nil {
		fmt.Println("err make dir ./generated: ", err)
	}

	f, err := os.Create("./generated/endpoints.go")
	if err != nil {
		fmt.Println("err create file ./generated/endpoints.go: ", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		fmt.Println("err call tmpl.Execute: ", err)
	}

	log.Println("File created successfully")
}

func HttpEndpoints() ([]string, error) {
	// Читаем JSON-файл
	content, err := os.ReadFile("./contracts/endpoints.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	// Десериализуем JSON
	var endpoints []Endpoint
	err = json.Unmarshal(content, &endpoints)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}

	tmpl, err := template.New("gen_structure").Parse(tmplEndpoint)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(endpoints))
	for _, e := range endpoints {
		var buf bytes.Buffer

		if err := tmpl.Execute(&buf, e); err != nil {
			return nil, err
		}
		res = append(res, buf.String())
	}

	return res, nil
}

func ParceFile() ([]string, error) {
	// Укажите путь к файлу, который хотите прочитать
	filePath := "./contracts/models.go"

	// Создаем новый файловый набор
	fset := token.NewFileSet()

	// Парсим файл
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	models := make([]string, 0)

	// Обходим AST и выводим все публичные структуры, имя которых заканчивается на "Request" или "Response"
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if structType, ok := x.Type.(*ast.StructType); ok {
				//fmt.Printf("Struct: %s\n", x.Name.Name)
				var buf bytes.Buffer
				printer.Fprint(&buf, fset, structType)
				//fmt.Println(buf.String())
				models = append(models, "type "+capitalize(x.Name.Name)+" "+buf.String())
			}
		case *ast.FuncDecl:
			// Выводим публичные функции, если это нужно
			if unicode.IsUpper(rune(x.Name.Name[0])) {
				fmt.Printf("Function: %s\n", x.Name.Name)
			}
		}
		return true
	})

	return models, err
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
