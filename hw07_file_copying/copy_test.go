package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestFile(t *testing.T, content []byte) (string, func()) {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "testfile_*.txt")
	require.NoError(t, err)

	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)

	return tmpfile.Name(), func() {
		os.Remove(tmpfile.Name())
	}
}

func TestCopy(t *testing.T) {
	t.Run("successful copy with offset and limit", func(t *testing.T) {
		srcPath, cleanupSrc := createTestFile(t, []byte("Hello, this is a test file for copying!"))
		defer cleanupSrc()

		dstPath, cleanupDst := createTestFile(t, nil)
		defer cleanupDst()

		offset := int64(7)
		limit := int64(4)
		err := Copy(srcPath, dstPath, offset, limit)
		require.NoError(t, err)

		result, err := os.ReadFile(dstPath)
		require.NoError(t, err)
		require.Equal(t, []byte("this"), result, "copied content mismatch")
	})

	t.Run("empty paths should return error", func(t *testing.T) {
		err := Copy("", "somepath.txt", 0, 0)
		require.True(t, errors.Is(err, ErrEmptyPaths), "expected ErrEmptyPaths")

		err = Copy("somepath.txt", "", 0, 0)
		require.True(t, errors.Is(err, ErrEmptyPaths), "expected ErrEmptyPaths")
	})

	t.Run("same source and destination paths", func(t *testing.T) {
		path, cleanup := createTestFile(t, []byte("test"))
		defer cleanup()

		err := Copy(path, path, 0, 0)
		require.True(t, errors.Is(err, ErrSameFiles), "expected ErrSameFiles")
	})

	t.Run("non-existent source file", func(t *testing.T) {
		dstPath, cleanupDst := createTestFile(t, nil)
		defer cleanupDst()

		err := Copy("testdata/nonexistent.txt", dstPath, 0, 0)
		require.Error(t, err)
		require.False(t, errors.Is(err, ErrEmptyPaths), "should not be ErrEmptyPaths")
	})

	t.Run("source is a directory", func(t *testing.T) {
		dirPath := filepath.Join("testdata", "testdir")
		err := os.MkdirAll(dirPath, 0o755)
		require.NoError(t, err)
		defer os.RemoveAll(dirPath)

		dstPath, cleanupDst := createTestFile(t, nil)
		defer cleanupDst()

		err = Copy(dirPath, dstPath, 0, 0)
		require.True(t, errors.Is(err, ErrIsDir), "expected ErrIsDir")
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		srcPath, cleanupSrc := createTestFile(t, []byte("short"))
		defer cleanupSrc()

		dstPath, cleanupDst := createTestFile(t, nil)
		defer cleanupDst()

		err := Copy(srcPath, dstPath, 10, 0)
		require.True(t, errors.Is(err, ErrOffsetExceedsFileSize), "expected ErrOffsetExceedsFileSize")
	})

	t.Run("negative limit", func(t *testing.T) {
		srcPath, cleanupSrc := createTestFile(t, []byte("test"))
		defer cleanupSrc()

		dstPath, cleanupDst := createTestFile(t, nil)
		defer cleanupDst()

		err := Copy(srcPath, dstPath, 0, -1)
		require.True(t, errors.Is(err, ErrInvalidLimit), "expected ErrInvalidLimit")
	})

	t.Run("full file copy with limit exceeding size", func(t *testing.T) {
		content := []byte("full file copy test")
		srcPath, cleanupSrc := createTestFile(t, content)
		defer cleanupSrc()

		dstPath, cleanupDst := createTestFile(t, nil)
		defer cleanupDst()

		err := Copy(srcPath, dstPath, 0, 1000)
		require.NoError(t, err)

		result, err := os.ReadFile(dstPath)
		require.NoError(t, err)
		require.Equal(t, content, result, "full file should be copied")
	})

	t.Run("copy from testdata file", func(t *testing.T) {
		srcPath := filepath.Join("testdata", "sample.txt")
		err := os.WriteFile(srcPath, []byte("sample content"), 0o644)
		require.NoError(t, err)
		defer os.Remove(srcPath)

		dstPath, cleanupDst := createTestFile(t, nil)
		defer cleanupDst()

		err = Copy(srcPath, dstPath, 0, 0)
		require.NoError(t, err)

		result, err := os.ReadFile(dstPath)
		require.NoError(t, err)
		require.Equal(t, []byte("sample content"), result, "content should match")
	})
}
