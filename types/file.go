package types

import (
	"github.com/geekymedic/neon/errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	DefFlag = os.O_CREATE | os.O_WRONLY | os.O_EXCL
	DefPerm = 0644
)

var ErrFileType = errors.NewStackError("file type error")

func NewBaseFile(abs string) *BaseFile {
	var (
		extraName string
		fileName  string
		fullName  string
		dir       string
	)
	dir, fullName = filepath.Split(abs)
	partIdx := strings.LastIndex(fullName, ".")
	if partIdx >= 0 && partIdx+1 < len(fullName) {
		extraName = fullName[partIdx+1:]
	}
	if partIdx >= 0 {
		fileName = fullName[0:partIdx]
	} else {
		fileName = fullName
	}
	return &BaseFile{
		absPath:   abs,
		dir:       filepath.Clean(dir),
		name:      fileName,
		extraName: extraName,
	}
}

// Opz add buffer for write
type BaseFile struct {
	absPath   string
	dir       string
	name      string
	extraName string
	*os.File
}

func (baseFile *BaseFile) Name() string {
	return baseFile.name
}

func (baseFile *BaseFile) ExtraName() string {
	return baseFile.extraName
}

func (baseFile *BaseFile) NodeType() NodeType {
	return NodeTypeFile
}

func (baseFile *BaseFile) BaseDir() string {
	return baseFile.dir
}

func (baseFile *BaseFile) Abs() string {
	return baseFile.absPath
}

func (baseFile *BaseFile) Split() []string {
	var absPath = path.Clean(baseFile.absPath)
	switch OsType() {
	case LinuxOs, MacOs:
		return strings.Split(absPath, Separator)
	case WindowsOs:
		return strings.Split(absPath, Separator)
	default:
		panic("unsupported os type")
	}
}

func (baseFile *BaseFile) Create(flag int, per os.FileMode) (err error) {
	if baseFile.File != nil {
		return ErrFileOpened
	}
	baseFile.File, err = os.OpenFile(baseFile.absPath, flag, per)
	return
}

func (baseFile *BaseFile) IsExist() error {
	info, err := os.Stat(baseFile.Abs())
	if err != nil {
		return err
	}
	if info.IsDir() {
		return ErrFileType
	}

	return os.ErrExist
}

func (baseFile *BaseFile) MustCreate(flag int, per os.FileMode) {
	if baseFile.File != nil {
		PanicSanity(ErrFileOpened)
	}
	err := os.MkdirAll(baseFile.BaseDir(), os.ModePerm)
	AssertNil(err)
	err = baseFile.Create(flag, per)
	AssertNil(err)
}

func (baseFile *BaseFile) Close() error {
	return baseFile.File.Close()
}

func (baseFile *BaseFile) Remove() (err error) {
	if err = baseFile.Close(); err != nil {
		return
	}
	return os.RemoveAll(baseFile.absPath)
}

func (baseFile *BaseFile) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(baseFile.File)
}

func (baseFile *BaseFile) Copy(src io.Reader) (int64, error) {
	return io.Copy(baseFile.File, src)
}

func (baseFile *BaseFile) Walk(fn filepath.WalkFunc) error {
	return filepath.Walk(baseFile.BaseDir(), fn)
}

type BaseDir struct {
	absPath string
	name    string
}

// dir must a directory
func NewBaseDir(dir string) *BaseDir {
	absPath := filepath.Clean(dir)
	_, name := filepath.Split(absPath)
	//var parts = strings.Split(filepath.Clean(dir), Separator)
	return &BaseDir{
		absPath: absPath,
		name:    name,
	}
}

func (baseDir *BaseDir) Name() string {
	return baseDir.name
}

func (baseDir *BaseDir) NodeType() NodeType {
	return NodeTypeDir
}

func (baseDir *BaseDir) BaseDir() string {
	idx := strings.LastIndex(baseDir.absPath, Separator)
	if idx >= 0 {
		return baseDir.absPath[0:idx]
	}
	return "."
}

func (baseDir *BaseDir) Abs() string {
	return baseDir.absPath
}

func (baseDir *BaseDir) Split() []string {
	var absPath = strings.Trim(path.Clean(baseDir.absPath), Separator)
	switch OsType() {
	case LinuxOs, MacOs:
		return strings.Split(absPath, Separator)
	case WindowsOs:
		// Notic
		// C:\Program\neon --> [C:, Program, neon]
		return strings.Split(absPath, Separator)
	default:
		panic("unsupported os type")
	}
}

func (baseDir *BaseDir) Create(per os.FileMode) (err error) {
	info, err := os.Stat(baseDir.absPath)
	if err == nil {
		if info.IsDir() {
			return
		}
		return errors.NewStackError("directory has exists, but it's type is not directory")
	}
	return os.MkdirAll(baseDir.absPath, per)
}

func (baseDir *BaseDir) IsExist() error {
	info, err := os.Stat(baseDir.Abs())
	if err != nil {
		return err
	}
	if info.IsDir() {
		return os.ErrExist
	}

	return ErrFileType
}

func (baseDir *BaseDir) Append(dirs ...string) DireOperation {
	var newDir = baseDir
	for _, dir := range dirs {
		newDir = NewBaseDir(newDir.Abs() + Separator + filepath.Clean(dir))
	}
	return newDir
}

func (baseDir *BaseDir) Remove() error {
	return os.RemoveAll(baseDir.absPath)
}

func (baseDir *BaseDir) Walk(fn filepath.WalkFunc) error {
	return filepath.Walk(baseDir.Abs(), fn)
}
