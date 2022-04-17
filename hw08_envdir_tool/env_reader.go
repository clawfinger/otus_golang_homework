package main

import (
	"bufio"
	"bytes"
	"os"
	"path"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	env := Environment{}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fileName := entry.Name()
		filePath := path.Join(dir, fileName)

		stat, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}
		if stat.Size() == 0 {
			env[fileName] = EnvValue{NeedRemove: true}
			continue
		}

		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		scanner.Scan()
		err = file.Close()
		if err != nil {
			return nil, err
		}

		content := scanner.Bytes()

		replaced := bytes.ReplaceAll(content, []byte{0}, []byte("\n"))
		content = bytes.TrimRight(replaced, " \t")

		env[fileName] = EnvValue{Value: string(content)}
	}

	return env, nil
}
