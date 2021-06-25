package cache

import (
	"context"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Item ...
type Item struct {
	CreateTime time.Time
	Filename   string
}

// Storage ...
type Storage struct {
	config *Config
	logger *logrus.Logger
	Items  map[string]Item
	fc     chan func()
}

// New ...
func New(logger *logrus.Logger, config *Config) *Storage {
	return &Storage{
		config: config,
		logger: logger,
		Items:  make(map[string]Item),
		fc:     make(chan func()),
	}
}

// Run ...
func (s *Storage) Run(ctx context.Context) error {
	err := os.MkdirAll(s.config.Dir, 0750)
	if err != nil {
		return err
	}
	err = s.checkStorage()
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case f := <-s.fc:
				f()
			case <-ctx.Done():
				return
			}

		}
	}()
	return nil
}

// GetCacheReader ...
func (s *Storage) GetCacheReader(key string) (io.ReadCloser, error) {
	type ans struct {
		r   io.ReadCloser
		err error
	}

	ansc := make(chan ans)

	s.fc <- func() {
		item, ok := s.Items[key]
		if ok {
			// if corresponding file does not exist, then check the storage
			if _, err := os.Stat(item.Filename); os.IsNotExist(err) {
				s.logger.Warn("Битая ссылка на файл хранилище кэшей. Запускаю проверку хранилища.")
				err := s.checkStorage()
				if err != nil {
					ansc <- ans{nil, err}
					return
				}
				ansc <- ans{nil, nil}
				return
			}
			f, err := os.Open(item.Filename)
			if err != nil {
				ansc <- ans{nil, err}
				return
			}
			ansc <- ans{f, nil}
			return
		}

		ansc <- ans{nil, nil}
		return
	}

	res := <-ansc
	return res.r, res.err
}

// CacheData ...
func (s *Storage) CacheData(key string, data []byte) {

	s.fc <- func() {
		_, ok := s.Items[key]
		if ok {
			s.logger.Info("Данные уже есть в кэше. Ничего не делаем.")
			return
		}

		if len(s.Items) >= s.config.Size {
			s.removeOldestItem()
		}

		t := time.Now()
		filename := path.Join(s.config.Dir, strings.ReplaceAll(key, string(os.PathSeparator), "_"))

		f, err := os.Create(filename)
		if err != nil {
			s.logger.Error("Ошибка создания кэш-файла: ", err)
			return
		}
		defer f.Close()

		_, err = f.Write(data)
		if err != nil {
			s.logger.Error("Не смог записать данные в файл: ", err)
			return
		}

		s.Items[key] = Item{
			CreateTime: t,
			Filename:   filename,
		}
	}
}

// removeOldestItem ...
func (s *Storage) removeOldestItem() {
	if len(s.Items) < s.config.Size {
		return
	}
	minTime := time.Now()
	minKey := ""
	for key, item := range s.Items {
		if minTime.After(item.CreateTime) {
			minTime = item.CreateTime
			minKey = key
		}
	}
	err := os.Remove(s.Items[minKey].Filename)
	if err != nil {
		s.logger.Error("Не смог удалить файл: ", err)
	}
	delete(s.Items, minKey)
}

// checkStorage checks if there are any broken references between files and filenames in the Storage structure
// and eliminates them if any found.
func (s *Storage) checkStorage() error {
	s.logger.Info("Start checking the cache storage...")
	for key, item := range s.Items {
		if _, err := os.Stat(item.Filename); os.IsNotExist(err) {
			s.logger.Info("Файл не существует. Удаляю соответствующую ссылку из map")
			delete(s.Items, key)
		}
	}
	files, err := os.ReadDir(s.config.Dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		filename := path.Join(s.config.Dir, f.Name())
		rem := true
		for _, item := range s.Items {
			if item.Filename == filename {
				rem = false
			}
		}
		if rem {
			err := os.Remove(filename)
			if err != nil {
				s.logger.Info("Не удалось удалить файл: ", err)
				return err
			}
		}
	}
	return nil
}
