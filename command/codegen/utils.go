package codegen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"unicode"
)

func CreateDirIfNeed(dirname string) error {
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

func CopyFile(src, dst string, firstLine string) error {
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

	if firstLine != "" {
		_, err = destinationFile.WriteString(firstLine)
		if err != nil {
			fmt.Println("Ошибка записи строки в файл назначения:", err)
			return err
		}
	}

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return destinationFile.Sync()
}

var (
	packageName     string
	packageNameInit bool
	packageNameMu   sync.Mutex
)

func GetPackageName() (string, error) {
	packageNameMu.Lock()
	defer packageNameMu.Unlock()

	if packageNameInit {
		return packageName, nil
	}

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
			packageName = strings.TrimSpace(strings.TrimPrefix(line, "module"))
			packageNameInit = true
			return packageName, nil
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

	packageName = filepath.Base(dir)
	packageNameInit = true
	return packageName, nil
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func CreateInnerFiles(e Endpoint, tmplRouterCode *template.Template) error {
	tmplFile, err := template.New("gen_logic_service").Parse(TmplLogicServiceFile)
	if err != nil {
		return err
	}

	tmplService, err := template.New("gen_logic_service_endpoints").Parse(TmplLogicServiceEndpoint)
	if err != nil {
		return err
	}

	packageName, err := GetPackageName()
	if err != nil {
		fmt.Println("err getPackageName: ", err)
	}

	dirname := "./inner/" + strings.ToLower(e.ServiceName)
	err = CreateDirIfNeed(dirname)
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

	err = addMethodToLogicServiceFileIfNeed(logicServiceFileName, tmplService, tmplRouterCode, e)
	if err != nil {
		fmt.Println("err addMethodToLogicServiceFileIfNeed: ", err)
		return err
	}

	return nil
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

func addMethodToLogicServiceFileIfNeed(fileName string, tmplService, tmplRouterCode *template.Template, e Endpoint) error {
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
		//tmplRouterCode := template.Must(template.New("service").Parse(TmplAddCodeToRouterFile))
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
