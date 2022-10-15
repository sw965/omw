package omw

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

func ListDir(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, err
	}
	result := make([]string, len(files))
	for i, file := range files {
		result[i] = file.Name()
	}
	return result, nil
}

func ReadText(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ReadTextLines(filePath string) ([]string, error) {
	fp, err := os.Open(filePath)
	if err != nil {
		return []string{}, err
	}
	defer fp.Close()

	result := make([]string, 0)
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result, nil
}

func WriteText(filePath string, data string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(([]byte)(data))
	return nil
}

func WriteTextLines(filePath string, data []string) error {
	strData := ""
	for _, ele := range data {
		strData += ele + "\n"
	}
	strData = strings.TrimRight(strData, "\n")
	return WriteText(filePath, strData)
}
