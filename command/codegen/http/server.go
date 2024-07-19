package http

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	// Создание директории, если она не существует
	if err := os.MkdirAll("./generated"+dirName, os.ModePerm); err != nil {
		fmt.Println("err make dir "+dirName+": ", err)
		return err
	}

	endpoints, err := GenerateHttpServerEndpoints(dirName)
	if err != nil {
		fmt.Println("err call GenerateHttpServerEndpoints: ", err)
		return err
	}

	fmt.Println("find GenerateHttpServerEndpoints: ", len(endpoints))
	for _, e := range endpoints {
		data.Endpoints = append(data.Endpoints, e)
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

	log.Println("File is generated successfully")
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

	tmplLSFile, err := template.New("gen_logic_service").Parse(tmplLogicServiceFile)
	if err != nil {
		return nil, err
	}

	tmplLSEndpoint, err := template.New("gen_logic_service_endpoints").Parse(tmplLogicServiceEndpoint)
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

		err = createInnerFiles(e, tmplLSFile, tmplLSEndpoint)
		if err != nil {
			fmt.Println("err create inner files: ", err)
			return nil, err
		}
	}

	return res, nil
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

	_, err = destinationFile.WriteString(tmplCodeGeneratorPhrase)
	if err != nil {
		fmt.Println("Ошибка записи строки в файл назначения:", err)
		return err
	}

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return destinationFile.Sync()
}

func createInnerFiles(e Endpoint, tmplFile, tmplService *template.Template) error {
	packageName, err := getPackageName()
	if err != nil {
		fmt.Println("err getPackageName: ", err)
	}

	dirname := "./inner/" + strings.ToLower(e.ServiceName)
	err = createDirIfNeed(dirname)
	if err != nil {
		fmt.Println("err create dir "+dirname+": ", err)
		return err
	}

	data := struct {
		PackageName        string
		ServiceName        string
		ServiceNameToLower string
	}{
		PackageName:        packageName,
		ServiceName:        e.ServiceName,
		ServiceNameToLower: strings.ToLower(e.ServiceName),
	}

	logicServiceFileName := dirname + "/logic_service.go"
	err = createLogicServiceFileIfNeed(logicServiceFileName, tmplFile, data)
	if err != nil {
		fmt.Println("err createLogicServiceFileIfNeed: ", err)
		return err
	}

	err = addMethodToLogicServiceFileIfNeed(logicServiceFileName, tmplService, e)
	if err != nil {
		fmt.Println("err addMethodToLogicServiceFileIfNeed: ", err)
		return err
	}

	return nil
}

func listFunctionByFileName(filePath string) (map[string]string, error) {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if unicode.IsUpper(rune(x.Name.Name[0])) {
				result[x.Name.Name] = x.Name.Name
			}
		}
		return true
	})

	return result, err
}

func createLogicServiceFileIfNeed(fileName string, tmplFile *template.Template, dstr any) error {
	fStat, err := os.Stat(fileName)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if fStat == nil {
		f, err := os.Create(fileName)
		if err != nil {
			fmt.Println("err create file "+fileName, err)
			return err
		}
		defer f.Close()

		err = tmplFile.Execute(f, dstr)
		if err != nil {
			fmt.Println("err call tmpl.Execute: ", err)
			return err
		}
	}

	return nil
}

func createDirIfNeed(dirname string) error {
	dState, err := os.Stat(dirname)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if dState == nil {
		if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func addMethodToLogicServiceFileIfNeed(fileName string, tmplService *template.Template, e Endpoint) error {
	fNames, err := listFunctionByFileName(fileName)
	if err != nil {
		fmt.Println("err list function by file name: ", err)
		return err
	}

	if _, ok := fNames[e.ServiceMethod]; !ok {
		// add method code to file
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("err open file "+fileName, err)
			return err
		}

		err = tmplService.Execute(f, e)
		if err != nil {
			fmt.Println("err call tmpl.Execute: ", err)
			return err
		}

		fmt.Println("Add next code to your router file:")
		tmplRouterCode := template.Must(template.New("service").Parse(tmplAddCodeToRouterFile))
		var buf bytes.Buffer
		data := struct {
			Name               string
			ServiceName        string
			ServiceNameToLower string
			ServiceMethod      string
		}{
			Name:               e.Name,
			ServiceName:        e.ServiceName,
			ServiceNameToLower: strings.ToLower(e.ServiceName),
			ServiceMethod:      e.ServiceMethod,
		}
		err = tmplRouterCode.Execute(&buf, data)
		if err != nil {
			fmt.Println("err call tmplRouterCode.Execute: ", err)
			return err
		}

		fmt.Println(buf.String())
	}

	return nil
}

func getPackageName() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module") {
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module"))
			return moduleName, nil
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return "", err
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка получения текущей директории:", err)
		return "", err
	}

	return filepath.Base(dir), nil

}
