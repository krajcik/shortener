package shortener

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
)

type FileRepository struct {
	rf *os.File
	w  *bufio.Writer
	r  *bufio.Reader
}

func NewFileRepository(filename string) (*FileRepository, error) {
	wFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, errors.New("open write file " + filename + " failed: " + err.Error())
	}
	rFile, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, errors.New("open read file " + filename + " failed: " + err.Error())
	}

	return &FileRepository{rf: rFile, w: bufio.NewWriter(wFile), r: bufio.NewReader(rFile)}, nil
}

func (f *FileRepository) Save(url *URL) error {
	jsonString, err := json.Marshal(url)
	if err != nil {
		return err
	}

	defer func(w *bufio.Writer) {
		err := w.Flush()
		if err != nil {
			panic(err)
		}
	}(f.w)
	_, err = f.rf.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	_, err = f.w.WriteString(string(jsonString) + "\n")
	if err != nil {
		return err
	}

	return nil
}

func (f *FileRepository) GetByURL(url string) (*URL, error) {
	scanner, err := f.getScanner()
	if err != nil {
		return nil, err
	}
	for scanner.Scan() {
		jsonString := scanner.Text()

		unmarshalURL := &URL{}
		if err := json.Unmarshal([]byte(jsonString), unmarshalURL); err != nil {
			panic(err)
		}

		if unmarshalURL.URL == url {
			return unmarshalURL, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, ErrNotFound
}

func (f *FileRepository) GetByShortCode(code string) (*URL, error) {
	scanner, err := f.getScanner()
	if err != nil {
		return nil, err
	}
	for scanner.Scan() {
		jsonString := scanner.Text()

		unmarshalURL := &URL{}
		jsonString = jsonString[:len(jsonString)-1]
		if err := json.Unmarshal([]byte(jsonString[:len(jsonString)-1]), unmarshalURL); err != nil {
			return nil, err
		}

		if unmarshalURL.ShortenedURL == code {
			return unmarshalURL, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, ErrNotFound
}

func (f *FileRepository) getScanner() (*bufio.Scanner, error) {
	_, err := f.rf.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f.r)
	return scanner, nil
}
