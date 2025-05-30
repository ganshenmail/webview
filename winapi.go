package main

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows API constants
const (
	// Window style constants
	GWL_STYLE      = -16
	WS_MAXIMIZEBOX = 0x00010000
	
	// System metrics indices
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
	
	// Window position flags
	SWP_NOSIZE         = 0x0001
	SWP_NOMOVE         = 0x0002
	SWP_NOZORDER       = 0x0004
	SWP_NOACTIVATE     = 0x0010
	SWP_FRAMECHANGED   = 0x0020
	SWP_SHOWWINDOW     = 0x0040
	SWP_NOOWNERZORDER  = 0x0200
	
	// Message box types
	MB_OK              = 0x00000000
	MB_ICONINFORMATION = 0x00000040
)

// WindowRect represents a rectangle structure used by Windows API
type WindowRect struct {
	Left, Top, Right, Bottom int
}

// Windows API DLL and procedure references
var (
	user32                = syscall.NewLazyDLL("user32.dll")
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	procFindWindow        = user32.NewProc("FindWindowW")
	procShowWindow        = user32.NewProc("ShowWindow")
	procGetWindowRect     = user32.NewProc("GetWindowRect")
	procGetSystemMetrics  = user32.NewProc("GetSystemMetrics")
	procSetWindowPos      = user32.NewProc("SetWindowPos")
	procGetDpiForSystem   = user32.NewProc("GetDpiForSystem")
	procGetDpiForWindow   = user32.NewProc("GetDpiForWindow")
	procSetWindowLongW    = user32.NewProc("SetWindowLongW")
)

// SetWindowLong sets a window attribute
func SetWindowLong(hwnd unsafe.Pointer, nIndex int, dwNewLong float64) (uintptr, uintptr, error) {
	return procSetWindowLongW.Call(uintptr(hwnd), uintptr(nIndex), uintptr(dwNewLong))
}

// ShowWindow sets the specified window's show state
func ShowWindow(hwnd unsafe.Pointer, show int) error {
	ret, _, err := procShowWindow.Call(uintptr(hwnd), uintptr(show))
	if ret == 0 {
		return err
	}
	return nil
}

// SetWindowPos changes the size, position, and Z order of a window
func SetWindowPos(hwnd unsafe.Pointer, hwndInsertAfter uintptr, x, y, cx, cy int, flags uint) error {
	ret, _, err := procSetWindowPos.Call(
		uintptr(hwnd),
		hwndInsertAfter,
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		uintptr(flags),
	)
	if ret == 0 {
		return err
	}
	return nil
}

// GetWindowRect retrieves the dimensions of the bounding rectangle of the specified window
func GetWindowRect(hwnd unsafe.Pointer) (WindowRect, error) {
	var rect WindowRect
	ret, _, err := procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&rect)))
	if ret == 0 {
		return rect, err
	}
	return rect, nil
}

// GetSystemMetrics retrieves the specified system metric
func GetSystemMetrics(nIndex int) (int, error) {
	ret, _, err := procGetSystemMetrics.Call(uintptr(nIndex))
	if ret == 0 {
		return 0, err
	}
	return int(ret), nil
}

// GetDpiForSystem retrieves the system DPI
func GetDpiForSystem() (uint, error) {
	ret, _, err := procGetDpiForSystem.Call()
	if ret == 0 {
		return 0, err
	}
	return uint(ret), nil
}

// GetDpiForWindow retrieves the DPI for the specified window
func GetDpiForWindow(hwnd unsafe.Pointer) (uint, error) {
	ret, _, err := procGetDpiForWindow.Call(uintptr(hwnd))
	if ret == 0 {
		return 0, err
	}
	return uint(ret), nil
}

// MessageBox displays a message box
func MessageBox(title string, message string, utype uint32) error {
	hwnd := windows.HWND(0)
	ptitle, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return err
	}
	
	pmessage, err := syscall.UTF16PtrFromString(message)
	if err != nil {
		return err
	}
	
	ret, err := windows.MessageBox(hwnd, pmessage, ptitle, utype)
	if ret == 0 {
		return err
	}
	return nil
}