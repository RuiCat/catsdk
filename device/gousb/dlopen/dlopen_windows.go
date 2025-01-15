package dlopen

import (
	"path/filepath"
	"syscall"
	"unsafe"
	_ "unsafe" // only for go:linkname
)

type syscall15Args struct {
	fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15 uintptr
	f1, f2, f3, f4, f5, f6, f7, f8                                       uintptr
	r1, r2, err                                                          uintptr
}

//go:nosplit
func Syscall_syscall15X(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15 uintptr) (r1, r2, err uintptr) {
	r1, r2, errno := syscall.SyscallN(fn, 15, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15)
	return r1, r2, uintptr(errno)
}

func GetHandle(name string) (Handle, error) {
	name, _ = filepath.Abs(name)
	handle, err := syscall.LoadLibrary(name)
	if err != nil {
		return nil, err
	}
	return &libHandle{
		Handle:  handle,
		Libname: name,
	}, nil
}

// Handle 动态库接口
type Handle interface {
	GetHandle() unsafe.Pointer
	GetSymbolPointer(symbol string) (unsafe.Pointer, error)
	Close() error
}

// LibHandle represents an open handle to a library (.so/dll)
type libHandle struct {
	Handle  syscall.Handle
	Libname string
}

func (l *libHandle) GetHandle() unsafe.Pointer {
	return unsafe.Pointer(l.Handle)
}

// GetSymbolPointer takes a symbol name and returns a pointer to the symbol.
func (l *libHandle) GetSymbolPointer(symbol string) (unsafe.Pointer, error) {
	ptr, err := syscall.GetProcAddress(l.Handle, symbol)
	return unsafe.Pointer(ptr), err
}

// Close closes a libHandle.
func (l *libHandle) Close() error {
	if l.Handle != 0 {
		err := syscall.FreeLibrary(l.Handle)
		l.Handle = 0
		return err
	}
	return nil
}
