// Copyright (c) 2016-2025 AtomAI, All rights reserved.
//
// See the COPYRIGHT file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

//go:build windows

// Package win32 provides Windows API bindings using purego/syscall,
// eliminating the need for cgo.
package win32

import (
	"syscall"
	"unsafe"
)

// Windows constants
const (
	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1

	MOUSEEVENTF_MOVE       = 0x0001
	MOUSEEVENTF_LEFTDOWN   = 0x0002
	MOUSEEVENTF_LEFTUP     = 0x0004
	MOUSEEVENTF_RIGHTDOWN  = 0x0008
	MOUSEEVENTF_RIGHTUP    = 0x0010
	MOUSEEVENTF_MIDDLEDOWN = 0x0020
	MOUSEEVENTF_MIDDLEUP   = 0x0040
	MOUSEEVENTF_WHEEL      = 0x0800
	MOUSEEVENTF_HWHEEL     = 0x01000

	KEYEVENTF_KEYUP    = 0x0002
	KEYEVENTF_UNICODE  = 0x0004
	KEYEVENTF_SCANCODE = 0x0008

	WHEEL_DELTA = 120

	MAPVK_VK_TO_VSC    = 0
	MAPVK_VSC_TO_VK    = 1
	MAPVK_VK_TO_CHAR   = 2
	MAPVK_VSC_TO_VK_EX = 3

	VK_LSHIFT   = 0xA0
	VK_RSHIFT   = 0xA1
	VK_LCONTROL = 0xA2
	VK_RCONTROL = 0xA3
	VK_LMENU    = 0xA4
	VK_RMENU    = 0xA5
	VK_LWIN     = 0x5B
	VK_RWIN     = 0x5C

	SW_MINIMIZE = 6
	SW_MAXIMIZE = 3
	SW_RESTORE  = 9

	GWL_STYLE   = -16
	WS_MAXIMIZE = 0x01000000
	WS_MINIMIZE = 0x20000000

	WM_CLOSE         = 0x0010
	WM_GETTEXT       = 0x000D
	WM_GETTEXTLENGTH = 0x000E

	DIB_RGB_COLORS = 0
	SRCCOPY        = 0x00CC0020
	BI_RGB         = 0

	GWL_WNDPROC = -4
	WM_KEYDOWN  = 0x0100
	WM_KEYUP    = 0x0101
	WM_CHAR     = 0x0102
)

// INPUT is the Windows INPUT structure for SendInput
type INPUT struct {
	Type uint32
	Di   [24]byte // large enough for MOUSEINPUT + KEYBDINPUT
}

// MOUSEINPUT is the Windows MOUSEINPUT structure
type MOUSEINPUT struct {
	DX          int32
	DY          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

// KEYBDINPUT is the Windows KEYBDINPUT structure
type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

// BITMAPINFOHEADER is the Windows BITMAPINFOHEADER structure
type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

// RECT is the Windows RECT structure
type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

// POINT is the Windows POINT structure
type POINT struct {
	X int32
	Y int32
}

// WINDOWPLACEMENT is the Windows WINDOWPLACEMENT structure
type WINDOWPLACEMENT struct {
	Length           uint32
	Flags            uint32
	ShowCmd          uint32
	PtMinPosition    POINT
	PtMaxPosition    POINT
	RcNormalPosition RECT
}

// DLL references
var (
	user32   = syscall.MustLoadDLL("user32.dll")
	kernel32 = syscall.MustLoadDLL("kernel32.dll")
	gdi32    = syscall.MustLoadDLL("gdi32.dll")
)

// User32 procedures
var (
	procSetCursorPos             = user32.MustFindProc("SetCursorPos")
	procGetCursorPos             = user32.MustFindProc("GetCursorPos")
	procSendInput                = user32.MustFindProc("SendInput")
	procGetForegroundWindow      = user32.MustFindProc("GetForegroundWindow")
	procSetForegroundWindow      = user32.MustFindProc("SetForegroundWindow")
	procGetWindowTextW           = user32.MustFindProc("GetWindowTextW")
	procGetWindowTextLengthW     = user32.MustFindProc("GetWindowTextLengthW")
	procShowWindow               = user32.MustFindProc("ShowWindow")
	procGetWindowThreadProcessId = user32.MustFindProc("GetWindowThreadProcessId")
	procMapVirtualKey            = user32.MustFindProc("MapVirtualKeyW")
	procVkKeyScanW               = user32.MustFindProc("VkKeyScanW")
	procGetWindowRect            = user32.MustFindProc("GetWindowRect")
	procGetClientRect            = user32.MustFindProc("GetClientRect")
	procFindWindowW              = user32.MustFindProc("FindWindowW")
	procIsWindow                 = user32.MustFindProc("IsWindow")
	procSetActiveWindow          = user32.MustFindProc("SetActiveWindow")
	procSetFocus                 = user32.MustFindProc("SetFocus")
	procGetDesktopWindow         = user32.MustFindProc("GetDesktopWindow")
	procGetDpiForWindow          = user32.MustFindProc("GetDpiForWindow")
	procGetSystemMetricsForDpi   = user32.MustFindProc("GetSystemMetricsForDpi")
	procCloseWindow              = user32.MustFindProc("CloseWindow")
	procPostMessageW             = user32.MustFindProc("PostMessageW")
	procSendMessageW             = user32.MustFindProc("SendMessageW")
	procGetWindowLongPtrW        = user32.MustFindProc("GetWindowLongPtrW")
	procGetKeyState              = user32.MustFindProc("GetKeyState")
	procEnumWindows              = user32.MustFindProc("EnumWindows")
	procGetWindowPlacement       = user32.MustFindProc("GetWindowPlacement")
	procSetWindowPlacement       = user32.MustFindProc("SetWindowPlacement")
	procGetActiveWindow          = user32.MustFindProc("GetActiveWindow")
)

// Gdi32 procedures
var (
	procBitBlt                 = gdi32.MustFindProc("BitBlt")
	procCreateDIBSection       = gdi32.MustFindProc("CreateDIBSection")
	procCreateCompatibleDC     = gdi32.MustFindProc("CreateCompatibleDC")
	procSelectObject           = gdi32.MustFindProc("SelectObject")
	procDeleteDC               = gdi32.MustFindProc("DeleteDC")
	procDeleteObject           = gdi32.MustFindProc("DeleteObject")
	procGetDIBits              = gdi32.MustFindProc("GetDIBits")
	procGetPixel               = gdi32.MustFindProc("GetPixel")
	procCreateCompatibleBitmap = gdi32.MustFindProc("CreateCompatibleBitmap")
)

// Kernel32 procedures
var (
	procGetModuleFileNameW = kernel32.MustFindProc("GetModuleFileNameW")
	procOpenProcess        = kernel32.MustFindProc("OpenProcess")
	procCloseHandle        = kernel32.MustFindProc("CloseHandle")
)

// Mouse API wrappers

// SetCursorPos sets the cursor position
func SetCursorPos(x, y int32) bool {
	ret, _, _ := procSetCursorPos.Call(uintptr(x), uintptr(y))
	return ret != 0
}

// GetCursorPos gets the cursor position
func GetCursorPos(pt *POINT) bool {
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(pt)))
	return ret != 0
}

// SendInput sends input events
func SendInput(nInputs uint32, pInputs unsafe.Pointer, cbSize int32) uint32 {
	ret, _, _ := procSendInput.Call(uintptr(nInputs), uintptr(pInputs), uintptr(cbSize))
	return uint32(ret)
}

// MoveMouse moves the mouse cursor to (x, y)
func MoveMouse(x, y int32) int {
	var pt POINT
	pt.X = x
	pt.Y = y

	// Use absolute positioning with SendInput
	var input INPUT
	input.Type = INPUT_MOUSE
	mi := (*MOUSEINPUT)(unsafe.Pointer(&input.Di[0]))
	mi.DX = x
	mi.DY = y
	mi.DwFlags = MOUSEEVENTF_MOVE | 0x8000 /* MOUSEEVENTF_ABSOLUTE - handled below */

	// For simplicity, first move with SetCursorPos (more reliable)
	if !SetCursorPos(x, y) {
		return 1
	}
	return 0
}

// ToggleMouse toggles a mouse button (down or up)
func ToggleMouse(down bool, button uint16) int {
	var input INPUT
	input.Type = INPUT_MOUSE
	mi := (*MOUSEINPUT)(unsafe.Pointer(&input.Di[0]))

	switch button {
	case 0: // LEFT_BUTTON
		if down {
			mi.DwFlags = MOUSEEVENTF_LEFTDOWN
		} else {
			mi.DwFlags = MOUSEEVENTF_LEFTUP
		}
	case 1: // RIGHT_BUTTON
		if down {
			mi.DwFlags = MOUSEEVENTF_RIGHTDOWN
		} else {
			mi.DwFlags = MOUSEEVENTF_RIGHTUP
		}
	case 2: // CENTER_BUTTON
		if down {
			mi.DwFlags = MOUSEEVENTF_MIDDLEDOWN
		} else {
			mi.DwFlags = MOUSEEVENTF_MIDDLEUP
		}
	}

	ret := SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
	if ret != 1 {
		return int(ret)
	}
	return 0
}

// ScrollMouse scrolls the mouse wheel
func ScrollMouse(x, y int) {
	if y != 0 {
		var input INPUT
		input.Type = INPUT_MOUSE
		mi := (*MOUSEINPUT)(unsafe.Pointer(&input.Di[0]))
		mi.DwFlags = MOUSEEVENTF_WHEEL
		if y > 0 {
			mi.MouseData = uint32(WHEEL_DELTA * y)
		} else {
			mi.MouseData = uint32(uint32(-y) * uint32(WHEEL_DELTA))
			// Negative scroll
			mi.MouseData = uint32(uint32(WHEEL_DELTA) * uint32(-y))
			if y < 0 {
				mi.MouseData = uint32(-WHEEL_DELTA * (-y))
			}
		}
		SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
	}

	if x != 0 {
		var input INPUT
		input.Type = INPUT_MOUSE
		mi := (*MOUSEINPUT)(unsafe.Pointer(&input.Di[0]))
		mi.DwFlags = MOUSEEVENTF_HWHEEL
		if x > 0 {
			mi.MouseData = uint32(WHEEL_DELTA * x)
		} else {
			mi.MouseData = uint32(-WHEEL_DELTA * (-x))
		}
		SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
	}
}

// Keyboard API wrappers

// ToggleKeyCode sends a key down/up event
func ToggleKeyCode(keyCode uint32, down bool, flags uint64, pid uintptr) int {
	var input INPUT
	input.Type = INPUT_KEYBOARD
	ki := (*KEYBDINPUT)(unsafe.Pointer(&input.Di[0]))
	ki.WVk = uint16(keyCode)
	ki.WScan = uint16(MapVirtualKey(keyCode, MAPVK_VK_TO_VSC))
	ki.DwFlags = 0

	if !down {
		ki.DwFlags = KEYEVENTF_KEYUP
	}

	// Handle modifier flags
	if flags != 0 {
		var modInputs []INPUT
		// Press modifiers down first
		if flags&0x01 != 0 { // MOD_ALT
			modInputs = append(modInputs, makeKeyInput(VK_LMENU, true))
		}
		if flags&0x02 != 0 { // MOD_CONTROL
			modInputs = append(modInputs, makeKeyInput(VK_LCONTROL, true))
		}
		if flags&0x04 != 0 { // MOD_SHIFT
			modInputs = append(modInputs, makeKeyInput(VK_LSHIFT, true))
		}
		if flags&0x08 != 0 { // MOD_META
			modInputs = append(modInputs, makeKeyInput(VK_LWIN, true))
		}

		if len(modInputs) > 0 {
			SendInput(uint32(len(modInputs)), unsafe.Pointer(&modInputs[0]), int32(unsafe.Sizeof(input)))
		}
	}

	ret := SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))

	// Release modifiers
	if flags != 0 && !down {
		var modInputs []INPUT
		if flags&0x01 != 0 {
			modInputs = append(modInputs, makeKeyInput(VK_LMENU, false))
		}
		if flags&0x02 != 0 {
			modInputs = append(modInputs, makeKeyInput(VK_LCONTROL, false))
		}
		if flags&0x04 != 0 {
			modInputs = append(modInputs, makeKeyInput(VK_LSHIFT, false))
		}
		if flags&0x08 != 0 {
			modInputs = append(modInputs, makeKeyInput(VK_LWIN, false))
		}
		if len(modInputs) > 0 {
			SendInput(uint32(len(modInputs)), unsafe.Pointer(&modInputs[0]), int32(unsafe.Sizeof(input)))
		}
	}

	if ret != 1 {
		return int(ret)
	}
	return 0
}

func makeKeyInput(vk uint16, down bool) INPUT {
	var input INPUT
	input.Type = INPUT_KEYBOARD
	ki := (*KEYBDINPUT)(unsafe.Pointer(&input.Di[0]))
	ki.WVk = vk
	ki.WScan = uint16(MapVirtualKey(uint32(vk), MAPVK_VK_TO_VSC))
	if !down {
		ki.DwFlags = KEYEVENTF_KEYUP
	}
	return input
}

// UnicodeType sends a unicode character
func UnicodeType(char uint32, pid uintptr, isPid int8) {
	var input INPUT
	input.Type = INPUT_KEYBOARD
	ki := (*KEYBDINPUT)(unsafe.Pointer(&input.Di[0]))
	ki.WScan = uint16(char)
	ki.DwFlags = KEYEVENTF_UNICODE
	SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))

	// Key up
	ki.DwFlags = KEYEVENTF_UNICODE | KEYEVENTF_KEYUP
	SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
}

// KeyCodeForChar returns the virtual key code for a character
func KeyCodeForChar(char byte) uint32 {
	ret, _, _ := procVkKeyScanW.Call(uintptr(char))
	return uint32(ret & 0xFF)
}

// MapVirtualKey maps a virtual key code
func MapVirtualKey(uCode, uMapType uint32) uint32 {
	ret, _, _ := procMapVirtualKey.Call(uintptr(uCode), uintptr(uMapType))
	return uint32(ret)
}

// Window API wrappers

// GetForegroundWindow returns the foreground window handle
func GetForegroundWindow() uintptr {
	ret, _, _ := procGetForegroundWindow.Call()
	return ret
}

// SetForegroundWindow sets the foreground window
func SetForegroundWindow(hwnd uintptr) bool {
	ret, _, _ := procSetForegroundWindow.Call(hwnd)
	return ret != 0
}

// GetWindowText gets the window title
func GetWindowText(hwnd uintptr) string {
	len, _, _ := procGetWindowTextLengthW.Call(hwnd)
	if len == 0 {
		return ""
	}
	buf := make([]uint16, len+1)
	procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len+1))
	return syscall.UTF16ToString(buf)
}

// ShowWindow shows/hides/minimizes/maximizes a window
func ShowWindow(hwnd uintptr, cmd int32) bool {
	ret, _, _ := procShowWindow.Call(hwnd, uintptr(cmd))
	return ret != 0
}

// GetWindowRect gets the window rectangle
func GetWindowRect(hwnd uintptr) (RECT, bool) {
	var rect RECT
	ret, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rect)))
	return rect, ret != 0
}

// GetClientRect gets the window client rectangle
func GetClientRect(hwnd uintptr) (RECT, bool) {
	var rect RECT
	ret, _, _ := procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rect)))
	return rect, ret != 0
}

// IsWindow checks if a window handle is valid
func IsWindow(hwnd uintptr) bool {
	ret, _, _ := procIsWindow.Call(hwnd)
	return ret != 0
}

// CloseWindow minimizes a window
func CloseWindow(hwnd uintptr) bool {
	ret, _, _ := procCloseWindow.Call(hwnd)
	return ret != 0
}

// PostCloseMessage posts WM_CLOSE to a window
func PostCloseMessage(hwnd uintptr) {
	procPostMessageW.Call(hwnd, WM_CLOSE, 0, 0)
}

// GetWindowThreadProcessId gets the thread and process id for a window
func GetWindowThreadProcessId(hwnd uintptr) uint32 {
	var pid uint32
	procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	return pid
}

// FindWindowByName finds a window by name
func FindWindowByName(name string) uintptr {
	namePtr, _ := syscall.UTF16PtrFromString(name)
	ret, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(namePtr)))
	return ret
}

// GetActiveWindow gets the active window
func GetActiveWindow() uintptr {
	ret, _, _ := procGetActiveWindow.Call()
	return ret
}

// SetActiveWindowHwnd sets the active window
func SetActiveWindowHwnd(hwnd uintptr) uintptr {
	ret, _, _ := procSetActiveWindow.Call(hwnd)
	return ret
}

// SetFocusHwnd sets the window focus
func SetFocusHwnd(hwnd uintptr) uintptr {
	ret, _, _ := procSetFocus.Call(hwnd)
	return ret
}

// GetDesktopWindow gets the desktop window
func GetDesktopWindow() uintptr {
	ret, _, _ := procGetDesktopWindow.Call()
	return ret
}

// GetDpiForWindow gets DPI for a window
func GetDpiForWindow(hwnd uintptr) uint32 {
	ret, _, _ := procGetDpiForWindow.Call(hwnd)
	return uint32(ret)
}

// GetSystemMetricsForDpi gets system metrics for DPI
func GetSystemMetricsForDpi(idx int32, dpi uint32) int32 {
	ret, _, _ := procGetSystemMetricsForDpi.Call(uintptr(idx), uintptr(dpi))
	return int32(ret)
}

// Screen capture functions

// CaptureScreen captures a screen region and returns pixel data (BGRA format)
// The caller is responsible for freeing the returned buffer.
func CaptureScreen(x, y, w, h int32) ([]byte, int, int) {
	hwnd := GetDesktopWindow()
	var rect RECT
	procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rect)))

	// GetDC from desktop window is not directly available via syscall,
	// but we can use GetDC(0) which gets the DC for the entire screen
	// Use the approach: CreateCompatibleDC + BitBlt
	_ = rect

	// Use GetDC(0) - the screen DC
	getDC := user32.MustFindProc("GetDC")
	releaseDC := user32.MustFindProc("ReleaseDC")

	hdc, _, _ := getDC.Call(0)
	if hdc == 0 {
		return nil, 0, 0
	}
	defer releaseDC.Call(0, hdc)

	// Create compatible DC and bitmap
	memDC, _, _ := procCreateCompatibleDC.Call(hdc)
	if memDC == 0 {
		return nil, 0, 0
	}
	defer procDeleteDC.Call(memDC)

	// Setup BITMAPINFOHEADER
	bmi := BITMAPINFOHEADER{
		BiSize:        uint32(unsafe.Sizeof(BITMAPINFOHEADER{})),
		BiWidth:       w,
		BiHeight:      -h, // negative for top-down DIB
		BiPlanes:      1,
		BiBitCount:    32,
		BiCompression: BI_RGB,
	}

	var pBits unsafe.Pointer
	hBitmap, _, _ := procCreateDIBSection.Call(
		memDC,
		uintptr(unsafe.Pointer(&bmi)),
		DIB_RGB_COLORS,
		uintptr(unsafe.Pointer(&pBits)),
		0, 0,
	)
	if hBitmap == 0 || pBits == nil {
		return nil, 0, 0
	}
	defer procDeleteObject.Call(hBitmap)

	// Select the bitmap into the memory DC
	procSelectObject.Call(memDC, hBitmap)

	// BitBlt from screen to memory DC
	procBitBlt.Call(
		memDC,
		0, 0,
		uintptr(w), uintptr(h),
		hdc,
		uintptr(x), uintptr(y),
		SRCCOPY,
	)

	// Copy pixel data
	byteWidth := w * 4
	totalBytes := byteWidth * h
	pixels := make([]byte, totalBytes)
	copy(pixels, unsafe.Slice((*byte)(pBits), totalBytes))

	return pixels, int(w), int(h)
}

// GetPixelColor gets the pixel color at (x, y) as RGB hex
func GetPixelColor(x, y int32) uint32 {
	getDC := user32.MustFindProc("GetDC")
	releaseDC := user32.MustFindProc("ReleaseDC")

	hdc, _, _ := getDC.Call(0)
	if hdc == 0 {
		return 0
	}
	defer releaseDC.Call(0, hdc)

	color, _, _ := procGetPixel.Call(hdc, uintptr(x), uintptr(y))
	// GetPixel returns COLORREF: 0x00BBGGRR, we want 0x00RRGGBB
	r := color & 0xFF
	g := (color >> 8) & 0xFF
	b := (color >> 16) & 0xFF
	return uint32(r)<<16 | uint32(g)<<8 | uint32(b)
}

// GetScreenSize gets the main display size
func GetScreenSize() (int32, int32) {
	getSystemMetrics := user32.MustFindProc("GetSystemMetrics")
	w, _, _ := getSystemMetrics.Call(0) // SM_CXSCREEN
	h, _, _ := getSystemMetrics.Call(1) // SM_CYSCREEN
	return int32(w), int32(h)
}

// GetNumDisplays gets the number of displays
func GetNumDisplays() int {
	// Use EnumDisplayMonitors to count
	enumDisplayMonitors := user32.MustFindProc("EnumDisplayMonitors")
	count := 0
	callback := syscall.NewCallback(func(hMonitor, hdcMonitor, lprcMonitor, dwData uintptr) uintptr {
		count++
		return 1 // continue enumeration
	})
	enumDisplayMonitors.Call(0, 0, callback, 0)
	return count
}

// GetWindowPlacement gets the window placement
func GetWindowPlacement(hwnd uintptr) (WINDOWPLACEMENT, bool) {
	var wp WINDOWPLACEMENT
	wp.Length = uint32(unsafe.Sizeof(wp))
	ret, _, _ := procGetWindowPlacement.Call(hwnd, uintptr(unsafe.Pointer(&wp)))
	return wp, ret != 0
}

// SetWindowPlacement sets the window placement
func SetWindowPlacement(hwnd uintptr, wp *WINDOWPLACEMENT) bool {
	ret, _, _ := procSetWindowPlacement.Call(hwnd, uintptr(unsafe.Pointer(wp)))
	return ret != 0
}

// EnumWindows enumerates all top-level windows
func EnumWindows(callback uintptr, lParam uintptr) bool {
	ret, _, _ := procEnumWindows.Call(callback, lParam)
	return ret != 0
}

// FindWindowByPID finds a window by process ID
func FindWindowByPID(pid uint32) uintptr {
	var foundHwnd uintptr
	callback := syscall.NewCallback(func(hwnd uintptr, lParam uintptr) int32 {
		var winPid uint32
		procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&winPid)))
		if winPid == pid {
			foundHwnd = hwnd
			return 0 // stop enumeration
		}
		return 1 // continue
	})
	EnumWindows(callback, 0)
	return foundHwnd
}
