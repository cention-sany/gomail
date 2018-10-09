package gomail

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"
)

type FileError struct {
	Err      error
	Filename string
}

func (f *FileError) Error() string {
	return fmt.Sprintf("%v", f.Err)
}

type FileSender struct {
	Err         error
	Dir, Prefix string
	counter     int
	sync.Mutex
}

func (f *FileSender) generateFileName() string {
	s := path.Join(f.Dir, fmt.Sprintf("%s-%x-%x-%d.eml", f.Prefix, rand.Int(),
		time.Now().Unix(), f.counter))
	f.counter++
	return s
}

func (f *FileSender) Send(from string, to []string, msg io.WriterTo) error {
	f.Lock()
	defer f.Unlock()
	fn := f.generateFileName()
	file, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err = msg.WriteTo(file); err != nil {
		file.Close()
		return err
	}
	if err = file.Close(); err == nil {
		if f.Err == nil {
			return nil
		}
		err = f.Err
	}
	return &FileError{Err: err, Filename: fn}
}

func (f *FileSender) Close() error {
	return nil
}
