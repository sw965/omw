package os

import (
    "os"
	"io/ioutil"
    "bufio"    
)

func FileNames(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, err
	}
	y := make([]string, len(files))
	for i, file := range files {
		y[i] = file.Name()
	}
	return y, nil
}

func ReadLine(path string, cap_ int) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return []string{}, err
    }
    defer file.Close()

    y := make([]string, 0, cap_)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        y = append(y, line)
    }
    if err := scanner.Err(); err != nil {
        return []string{}, err
    }
    return y, nil
}