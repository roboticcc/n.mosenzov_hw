package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrIsDir                 = errors.New("file is a directory")
	ErrInvalidLimit          = errors.New("limit can't be a negative number")
	ErrEmptyPaths            = errors.New("from/to path not specified")
	ErrSamePath              = errors.New("source and destination paths cannot be the same")
)

type progressWriter struct {
	copied *int64
	total  int64
}

func (pw *progressWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	*pw.copied += int64(n)
	percent := float64(*pw.copied) / float64(pw.total) * 100
	fmt.Printf("\rCopying: %.2f%%", percent)
	return n, nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrEmptyPaths
	}
	if fromPath == toPath {
		return ErrSamePath
	}

	fi, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return ErrIsDir
	}
	if !fi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if offset > fi.Size() {
		return ErrOffsetExceedsFileSize
	}
	if limit < 0 {
		return ErrInvalidLimit
	}

	fromF, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := fromF.Close(); err != nil {
			log.Printf("failed to close source file %s: %v", fromPath, err)
		}
	}()

	_, err = fromF.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	toF, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := toF.Close(); err != nil {
			log.Printf("failed to close destination file %s: %v", toPath, err)
		}
	}()

	remaining := fi.Size() - offset
	if limit == 0 || limit > remaining {
		limit = remaining
	}

	var copied int64
	pw := &progressWriter{
		copied: &copied,
		total:  limit,
	}
	reader := io.TeeReader(fromF, pw)

	_, err = io.CopyN(toF, reader, limit)
	if err != nil {
		return err
	}

	return nil
}
