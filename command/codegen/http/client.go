package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"
)

func GenerateHttpClientFile(dirName string) error {
	tmpl, err := template.New("gen_structure").Parse(tmplClientEndpointFile)
	if err != nil {
		fmt.Println("err read template: ", err)
		return err
	}

	data := struct {
		Endpoints []string
		Models    []string
	}{}

	endpoints, err := GenerateHttpClientEndpoints(dirName)
	if err != nil {
		fmt.Println("err call GenerateHttpClientEndpoints: ", err)
		return err
	}

	fmt.Println("find GenerateHttpClientEndpoints: ", len(endpoints))
	for _, e := range endpoints {
		data.Endpoints = append(data.Endpoints, e)
	}

	// Создание директории, если она не существует
	if err := os.MkdirAll("./generated"+dirName, os.ModePerm); err != nil {
		fmt.Println("err make dir "+dirName+": ", err)
		return err
	}

	// copy models.go
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

func GenerateHttpClientEndpoints(dirName string) ([]string, error) {
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

	tmpl, err := template.New("gen_structure").Parse(tmplClientEndpoint)
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
