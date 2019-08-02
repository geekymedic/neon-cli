package types

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/geekymedic/neon/errors"
)

var (
	ErrFileOpened = errors.NewStackError("file has opened")
)

var Separator = string(filepath.Separator)

type NodeType uint8

const (
	NodeTypeDir NodeType = iota
	NodeTypeFile
)

type BaseNode interface {
	Name() string
	IsExist() error
	NodeType() NodeType
	BaseDir() string
	Abs() string
	Split() []string
}

type FileNode interface {
	BaseNode
	FileOperation
}

type FileOperation interface {
	Create(flag int, perm os.FileMode) error
	MustCreate(flag int, perm os.FileMode)
	ExtraName() string
	Remove() (err error)
	ReadAll() ([]byte, error)
	Copy(src io.Reader) (int64, error)
	Walk(walkFunc filepath.WalkFunc) error
	io.Writer
	io.Reader
	io.Closer
	io.StringWriter
	io.ReaderAt
	io.Seeker
}

type DirNode interface {
	BaseNode
	DireOperation
}

type DireOperation interface {
	Create(per os.FileMode) (err error)
	Append(...string) DireOperation
	Remove() (err error)
	Walk(fn filepath.WalkFunc) error
}

//func PanicSanity(format string, args ...interface{}) {
//	panic(fmt.Sprintf(format, args...))
//}

var PanicSanity = func(v interface{}) {
	panic(v)
}

var PanicSanityf = func(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v...))
}

var AssertNil = func(err error) {
	if err == nil {
		return
	}
	PanicSanity(err)
}

var AssertNotNil = func(v interface{}) {
	if v == nil {
		PanicSanity("It should be not nil")
	}
}
