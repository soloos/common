package fsapi

// _Dirent is what we send to the kernel, but we offer DirEntry and
// DirEntryList to the user.
type _Dirent struct {
	Ino     uint64
	Off     uint64
	NameLen uint32
	Typ     uint32
}
