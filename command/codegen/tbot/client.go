package tbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/viktor8881/codegen/command/codegen"
	"log"
	"os"
	"text/template"
)

func GenerateTbotClientFile(dirName string) error {
	tmpl, err := template.New("gen_structure").Parse(TmplTbotEndpointFile)
	if err != nil {
		fmt.Println("err read template: ", err)
		return err
	}

	packageName, err := codegen.GetPackageName()
	if err != nil {
		fmt.Println("err getPackageName: ", err)
	}

	data := struct {
		Endpoints   []string
		PackageName string
	}{
		PackageName: packageName,
	}

	endpoints, err := GenerateTbotEndpoints(dirName)
	if err != nil {
		fmt.Println("err call GenerateTbotEndpoints: ", err)
		return err
	}

	fmt.Println("find GenerateTbotEndpoints: ", len(endpoints))
	for _, e := range endpoints {
		data.Endpoints = append(data.Endpoints, e)
	}

	// Создание директории, если она не существует
	if err := os.MkdirAll("./generated"+dirName, os.ModePerm); err != nil {
		fmt.Println("err make dir "+dirName+": ", err)
		return err
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

func GenerateTbotEndpoints(dirName string) ([]string, error) {
	// Читаем JSON-файл
	content, err := os.ReadFile("./contracts" + dirName + "/endpoints.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	// Десериализуем JSON
	var endpoints []codegen.Endpoint
	err = json.Unmarshal(content, &endpoints)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}

	tmpl, err := template.New("gen_structure").Parse(TmplTbotEndpoint)
	if err != nil {
		return nil, err
	}

	tmplRouterCode := template.Must(template.New("service").Parse(TmplAddCodeToRouterFile))

	res := make([]string, 0, len(endpoints))
	for _, e := range endpoints {
		var buf bytes.Buffer

		if err := tmpl.Execute(&buf, e); err != nil {
			return nil, err
		}
		res = append(res, buf.String())

		err = codegen.CreateInnerFiles(e, tmplRouterCode)
		if err != nil {
			fmt.Println("err create inner files: ", err)
			return nil, err
		}
	}

	return res, nil
}
