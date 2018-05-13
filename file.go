package until

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// GetIntSliceFromFile 从文件中读取[]int结构
// 用\n作为行分隔符, "splitString"作为列分隔符
func GetIntSliceFromFile(file, splitString string) ([]int, error) {
	s := make([]int, 0)
	f, err := os.Open(file)
	if err != nil {
		return s, err
	}
	defer f.Close()

	// 读取文件到buffer里边
	buf := bufio.NewReader(f)
	for {
		// 按照换行读取每一行
		l, err := buf.ReadString('\n')
		// 跳过空行
		if l == "\n" {
			continue
		}

		lineSplit := strings.SplitN(l, splitString, 1024)
		for _, v := range lineSplit {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			value, _ := strconv.Atoi(v)
			s = append(s, value)
		}
		if err != nil {
			break
		}
	}
	return s, nil
}

// BytesToFile 字节数组写入到文件
func BytesToFile(data []byte, filepath string) error {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}
