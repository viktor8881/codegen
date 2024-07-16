package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
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

func GenerateHttpServerFile(dirName string) error {
	tmpl, err := template.New("gen_structure").Parse(tmplServerEndpointFile)
	if err != nil {
		fmt.Println("err read template: ", err)
		return err
	}

	data := struct {
		Endpoints []string
		Models    []string
	}{}

	endpoints, err := GenerateHttpServerEndpoints(dirName)
	if err != nil {
		fmt.Println("err call GenerateHttpServerEndpoints: ", err)
		return err
	}

	fmt.Println("find GenerateHttpServerEndpoints: ", len(endpoints))
	for _, e := range endpoints {
		data.Endpoints = append(data.Endpoints, e)
	}

	// Создание директории, если она не существует
	if err := os.MkdirAll("./generated"+dirName, os.ModePerm); err != nil {
		fmt.Println("err make dir "+dirName+": ", err)
		return err
	}

	if err := copyFile("./contracts"+dirName+"/models.go", "./generated"+dirName+"/models.go"); err != nil {
		log.Fatalf("failed to copy file: %v", err)
	}

	f, err := os.Create("./generated" + dirName + "/endpoints.go")
	if err != nil {
		fmt.Println("err create file ./generated"+dirName+"/endpoints.go: ", err)
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		fmt.Println("err call tmpl.Execute: ", err)
		return err
	}

	log.Println("File created successfully")
	return nil
}

func GenerateHttpServerEndpoints(dirName string) ([]string, error) {
	// Читаем JSON-файл
	content, err := os.ReadFile("./contracts" + dirName + "/endpoints.json")
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

	tmpl, err := template.New("gen_structure").Parse(tmplServerEndpoint)
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

func ParceFile(filePath string) ([]string, error) {
	fset := token.NewFileSet()

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

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return destinationFile.Sync()
}
