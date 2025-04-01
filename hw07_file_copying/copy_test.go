package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("successful copy with offset and limit", func(t *testing.T) {
		srcContent := []byte("Hello, this is a test file for copying!")
		srcFile, err := os.CreateTemp("", "src_test_*.txt")
		require.NoError(t, err)
		defer os.Remove(srcFile.Name())

		_, err = srcFile.Write(srcContent)
		require.NoError(t, err)
		err = srcFile.Close()
		require.NoError(t, err)

		dstFile, err := os.CreateTemp("", "dst_test_*.txt")
		require.NoError(t, err)
		defer os.Remove(dstFile.Name())

		offset := int64(7)
		limit := int64(4)
		err = Copy(srcFile.Name(), dstFile.Name(), offset, limit)
		require.NoError(t, err)

		result, err := os.ReadFile(dstFile.Name())
		require.NoError(t, err)
		expected := []byte("this")
		require.Equal(t, expected, result, "copied content mismatch")
	})

	t.Run("empty paths should return error", func(t *testing.T) {
		err := Copy("", "somepath.txt", 0, 0)
		require.True(t, errors.Is(err, ErrEmptyPaths), "expected ErrEmptyPaths")

		err = Copy("somepath.txt", "", 0, 0)
		require.True(t, errors.Is(err, ErrEmptyPaths), "expected ErrEmptyPaths")
	})

	t.Run("same source and destination paths", func(t *testing.T) {
		path := "testdata/samepath.txt"
		err := Copy(path, path, 0, 0)
		require.True(t, errors.Is(err, ErrSamePath), "expected ErrSamePath")
	})

	t.Run("non-existent source file", func(t *testing.T) {
		err := Copy("testdata/nonexistent.txt", "testdata/output.txt", 0, 0)
		require.Error(t, err)
		require.False(t, errors.Is(err, ErrEmptyPaths), "should not be ErrEmptyPaths")
	})

	t.Run("source is a directory", func(t *testing.T) {
		err := Copy("testdata", "testdata/output.txt", 0, 0)
		require.True(t, errors.Is(err, ErrIsDir), "expected ErrIsDir")
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		srcFile, err := os.CreateTemp("", "src_test_*.txt")
		require.NoError(t, err)
		defer os.Remove(srcFile.Name())

		_, err = srcFile.Write([]byte("short"))
		require.NoError(t, err)
		err = srcFile.Close()
		require.NoError(t, err)

		dstFile, err := os.CreateTemp("", "dst_test_*.txt")
		require.NoError(t, err)
		defer os.Remove(dstFile.Name())

		err = Copy(srcFile.Name(), dstFile.Name(), 10, 0)
		require.True(t, errors.Is(err, ErrOffsetExceedsFileSize), "expected ErrOffsetExceedsFileSize")
	})

	t.Run("negative limit", func(t *testing.T) {
		srcFile, err := os.CreateTemp("", "src_test_*.txt")
		require.NoError(t, err)
		defer os.Remove(srcFile.Name())

		_, err = srcFile.Write([]byte("test"))
		require.NoError(t, err)
		err = srcFile.Close()
		require.NoError(t, err)

		dstFile, err := os.CreateTemp("", "dst_test_*.txt")
		require.NoError(t, err)
		defer os.Remove(dstFile.Name())

		err = Copy(srcFile.Name(), dstFile.Name(), 0, -1)
		require.True(t, errors.Is(err, ErrInvalidLimit), "expected ErrInvalidLimit")
	})

	t.Run("full file copy with limit exceeding size", func(t *testing.T) {
		srcFile, err := os.CreateTemp("", "src_test_*.txt")
		require.NoError(t, err)
		defer os.Remove(srcFile.Name())

		content := []byte("full file copy test")
		_, err = srcFile.Write(content)
		require.NoError(t, err)
		err = srcFile.Close()
		require.NoError(t, err)

		dstFile, err := os.CreateTemp("", "dst_test_*.txt")
		require.NoError(t, err)
		defer os.Remove(dstFile.Name())

		err = Copy(srcFile.Name(), dstFile.Name(), 0, 1000)
		require.NoError(t, err)

		result, err := os.ReadFile(dstFile.Name())
		require.NoError(t, err)
		require.Equal(t, content, result, "full file should be copied")
	})

	t.Run("copy from testdata file", func(t *testing.T) {
		srcPath := filepath.Join("testdata", "sample.txt")
		err := os.WriteFile(srcPath, []byte("sample content"), 0o644)
		require.NoError(t, err)
		defer os.Remove(srcPath)

		dstFile, err := os.CreateTemp("", "dst_test_*.txt")
		require.NoError(t, err)
		defer os.Remove(dstFile.Name())

		err = Copy(srcPath, dstFile.Name(), 0, 0)
		require.NoError(t, err)

		result, err := os.ReadFile(dstFile.Name())
		require.NoError(t, err)
		require.Equal(t, []byte("sample content"), result, "content should match")
	})
}
