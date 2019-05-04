package os

import "os"

func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func Remove(name string) error {
	return os.Remove(name)
}

func OpenFile(name string, flag int, perm FileMode) (*File, error) {
	file, err := os.OpenFile(name, flag, perm)
	return &File{file: file}, err
}

func Open(name string) (*File, error) {
	file, err := os.Open(name)
	return &File{file: file}, err
}

func Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func Create(name string) (*File, error) {
	file, err := os.Create(name)
	return &File{file: file}, err
}

func Truncate(name string, size int64) error {
	return os.Truncate(name, size)
}

func Mkdir(name string, perm FileMode) error {
	return os.Mkdir(name, perm)
}

func MkdirAll(path string, perm FileMode) error {
	return os.MkdirAll(path, perm)
}

func IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func IsExist(err error) bool {
	return os.IsExist(err)
}

func Stat(name string) (FileInfo, error) {
	return os.Stat(name)
}

func TempDir() string {
	return os.TempDir()
}

func Exit(code int) {
	os.Exit(code)
}

func Getpid() int {
	return os.Getpid()
}
