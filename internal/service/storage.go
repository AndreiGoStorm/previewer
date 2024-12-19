package service

import (
	"os"
	"path/filepath"

	"github.com/AndreiGoStorm/previewer/internal/logger"
)

type Storage struct {
	logg *logger.Logger
	Dir  string
}

func NewStorage(logg *logger.Logger) *Storage {
	dir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	s := &Storage{
		logg: logg,
		Dir:  filepath.Join(dir, "uploads"),
	}

	if err := s.createDir(s.Dir); err != nil {
		logg.Error("storage new createDir: %w", err)
		panic(err)
	}

	return s
}

func (s *Storage) ReadDirNames() (names []string, err error) {
	fd, err := os.Open(s.Dir)
	if err != nil {
		s.logg.Error("storage ReadDirNames Open: %w", err)
		return nil, err
	}
	defer fd.Close()

	names, err = fd.Readdirnames(-1)
	if err != nil {
		s.logg.Error("storage ReadDirNames Readdirnames: %w", err)
		return nil, err
	}

	return names, nil
}

func (s *Storage) DeleteFile(filename string) error {
	path := s.getStorageFullPath(filename)
	_, err := os.Stat(path)
	if err == nil {
		if err = os.Remove(path); err != nil {
			s.logg.Error("storage DeleteFile Remove: %w", err)
			return err
		}
	}
	return nil
}

func (s *Storage) GetImagePath(filename string) (string, error) {
	imagePath := s.getStorageFullPath(filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		s.logg.Info("Image not fount: " + filename)
		return "", err
	}
	return imagePath, err
}

func (s *Storage) getStorageFullPath(filename string) string {
	return filepath.Join(s.Dir, filename)
}

func (s *Storage) createDir(path string) (err error) {
	_, err = os.ReadDir(path)
	if err != nil {
		err = os.MkdirAll(path, 0o755)
		if err != nil {
			return err
		}
	}
	return
}
