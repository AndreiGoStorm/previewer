package service

import (
	"os"
	"path/filepath"
	"previewer/internal/logger"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	tempExtension := "*.txt"

	logg := logger.New("INFO")
	storage := NewStorage(logg)

	t.Run("exists upload dir", func(t *testing.T) {
		_, err := os.Stat(storage.Dir)
		require.NoError(t, err)
	})

	t.Run("read storage dir for two files", func(t *testing.T) {
		file1, _ := os.CreateTemp(storage.Dir, tempExtension)
		defer os.Remove(file1.Name())

		file2, _ := os.CreateTemp(storage.Dir, tempExtension)
		defer os.Remove(file2.Name())

		names, err := storage.ReadDirNames()
		require.NoError(t, err)

		require.Len(t, names, 2)
	})

	t.Run("delete file from storage dir", func(t *testing.T) {
		file1, _ := os.CreateTemp(storage.Dir, tempExtension)

		err := storage.DeleteFile(filepath.Base(file1.Name()))
		require.NoError(t, err)

		names, err := storage.ReadDirNames()
		require.NoError(t, err)

		require.Len(t, names, 0)
	})

	err := os.RemoveAll(storage.Dir)
	require.NoError(t, err)
}
