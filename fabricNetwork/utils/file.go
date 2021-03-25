package utils

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
)

func fileStat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

//exists file
func exists(name string) bool {
	if _, err := fileStat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//isDir file or directory
func isDir(name string) bool {
	info, err := fileStat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}

var files []*ScpFile

//GetFiles return list of files
func GetFiles(name string) ([]*ScpFile, error) {

	if !exists(name) {
		return nil, nil
	}

	if isDir(name) {
		if len(files) != 0 {
			files = files[0:0]
		}
		filepath.Walk(name, visit)
		return files, nil
	} else {
		file, err := os.Open(name)
		info, err := os.Stat(name)
		scpFile := &ScpFile{
			F:     file,
			Info:  info,
			IsDir: false,
		}
		return []*ScpFile{scpFile}, err
	}

}

func visit(path string, f os.FileInfo, err error) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		scpFile := &ScpFile{
			F:     file,
			Info:  f,
			IsDir: f.IsDir(),
		}
		files = append(files, scpFile)
	}
	return nil
}

// ReadLines read file
func ReadLines(input string) (lines []string, err error) {
	sourceAbs, err := filepath.Abs(input)
	if err != nil {
		os.Exit(1)
	}
	var (
		file   *os.File
		part   []byte
		prefix bool
	)

	if file, err = os.Open(sourceAbs); err != nil {
		return
	}

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))

	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

// ScpFile struct
type ScpFile struct {
	F     *os.File
	Info  os.FileInfo
	IsDir bool
}
