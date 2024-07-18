package http

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func GenerateUtilsErrorHandlerFile(fileName string) error {
	// Создание директории, если она не существует
	if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
		fmt.Println("err make dir "+filepath.Dir(fileName)+": ", err)
		return err
	}

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

		tmpl, err := template.New("gen_structure").Parse(tmplUtilsErrorHelperFile)
		if err != nil {
			fmt.Println("err read template: ", err)
			return err
		}

		data := struct {
			TagJsonErrorMess string
		}{
			TagJsonErrorMess: "`json:\"message\"`",
		}
		err = tmpl.Execute(f, data)
		if err != nil {
			fmt.Println("err call tmpl.Execute: ", err)
			return err
		}

		log.Println("File created successfully")
	}

	return nil
}
