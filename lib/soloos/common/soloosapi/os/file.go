package os

import (
	"os"
	"soloos/common/fsapi"
	"soloos/common/log"
	"soloos/common/sdfsapitypes"
)

type FileMode = os.FileMode
type FileInfo = os.FileInfo

type File struct {
	isSoloOSFile  bool
	soloOSPosixFS fsapi.PosixFS
	soloOSFdID    sdfsapitypes.FsINodeFileHandlerID
	file          *os.File
}

func (p *File) SetSoloOSFsINode(fdID sdfsapitypes.FsINodeFileHandlerID, posixFS fsapi.PosixFS) {
	p.soloOSFdID = fdID
	p.soloOSPosixFS = posixFS
	p.isSoloOSFile = true
}

func (p *File) Fd() uintptr {
	return p.file.Fd()
}

func (p *File) Name() string {
	return p.file.Name()
}

func (p *File) Stat() (FileInfo, error) {
	return p.file.Stat()
}

func (p *File) Seek(offset int64, whence int) (ret int64, err error) {
	return p.file.Seek(offset, whence)
}

func (p *File) ReadAt(b []byte, off int64) (n int, err error) {
	return p.file.ReadAt(b, off)
}

func (p *File) Read(b []byte) (n int, err error) {
	return p.file.Read(b)
}

func (p *File) Write(b []byte) (n int, err error) {
	log.Debug("fuck file Write", len(b), p.isSoloOSFile, p.soloOSFdID)
	if p.isSoloOSFile {
		var (
			fd  = p.soloOSPosixFS.FdTableGetFd(p.soloOSFdID)
			err error
		)
		err = p.soloOSPosixFS.SimpleWriteWithMem(fd.FsINodeID, b, fd.AppendPosition)
		if err != nil {
			log.Warn(err)
			return 0, err
		}

		p.soloOSPosixFS.FdTableFdAddAppendPosition(p.soloOSFdID, uint64(len(b)))
	}
	return p.file.Write(b)
}

func (p *File) WriteAt(b []byte, off int64) (n int, err error) {
	log.Debug("fuck file WriteAt", len(b), off)
	return p.file.WriteAt(b, off)
}

func (p *File) Sync() error {
	return p.file.Sync()
}

func (p *File) Truncate(size int64) error {
	return p.file.Truncate(size)
}

func (p *File) Readdir(n int) ([]FileInfo, error) {
	return p.file.Readdir(n)
}

func (p *File) Close() error {
	return p.file.Close()
}
