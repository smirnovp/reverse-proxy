package cache

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

// Item ...
type Item struct {
	CreateTime time.Time
	Filename   string
}

// Storage ...
type Storage struct {
	config *Config
	Items  map[string]Item
	fc     chan func()
}

// New ...
func New(config *Config) *Storage {
	return &Storage{
		config: config,
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
			f, err := os.Open(path.Join(s.config.Dir, item.Filename))
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
			fmt.Println("Данные уже есть в кэше. Ничего не делаем.")
			return
		}

		if len(s.Items) >= s.config.Size {
			s.RemoveItem()
		}

		t := time.Now()
		filename := path.Clean(strings.ReplaceAll(key, string(os.PathSeparator), "_"))

		f, err := os.Create(path.Join(s.config.Dir, filename))
		if err != nil {
			log.Println("Ошибка создания кэш-файла: ", err)
			return
		}
		defer f.Close()

		_, err = f.Write(data)
		if err != nil {
			log.Println("Не смог записать данные в файл: ", err)
			return
		}

		s.Items[key] = Item{
			CreateTime: t,
			Filename:   filename,
		}
	}
}

// RemoveItem ...
func (s *Storage) RemoveItem() {

}
