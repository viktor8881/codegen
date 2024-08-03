package codegen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
