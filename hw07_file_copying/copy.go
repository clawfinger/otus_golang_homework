package main

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const BufferLen = 5120

func Copy(fromPath, toPath string, offset, limit int64) error {

	count := 100
	// create and start new bar
	bar := pb.StartNew(count)

	// finish bar
	defer bar.Finish()
	source, err := os.OpenFile(fromPath, os.O_RDONLY, os.ModeAppend)

	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(toPath)
	if err != nil {
		return err
	}

	fi, err := source.Stat()

	if err != nil {
		return err
	}

	fileSize := fi.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = source.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	var reader io.Reader
	if limit == 0 {
		reader = source
	} else {
		reader = io.LimitReader(source, limit)
	}

	totalWriten := 0

	writeBuffer := bufio.NewWriter(dest)
	readBuffer := make([]byte, BufferLen)
	for {
		read, err := reader.Read(readBuffer)

		if err == io.EOF {
			break
		}
		writeBuffer.Write(readBuffer[:read])
		totalWriten += writeBuffer.Buffered()
		writeBuffer.Flush()
		bar.Set(totalWriten / int(fileSize) * 100)
	}

	return nil
}
