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

package robotgo

import (
	"fmt"
	"image"
	"image/color"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/shaolei/robotgo/internal/win32"
)

// Mouse button constants (replacing C.MMMouseButton)
const (
	LEFT_BUTTON   uint16 = 0
	RIGHT_BUTTON  uint16 = 1
	CENTER_BUTTON uint16 = 2
	WHEEL_DOWN    uint16 = 3
	WHEEL_UP      uint16 = 4
	WHEEL_LEFT    uint16 = 5
	WHEEL_RIGHT   uint16 = 6
)

// Key modifier flags (replacing C.MMKeyFlags / C.MOD_*)
const (
	MOD_NONE    uint64 = 0
	MOD_ALT     uint64 = 0x01
	MOD_CONTROL uint64 = 0x02
	MOD_SHIFT   uint64 = 0x04
	MOD_META    uint64 = 0x08
)

// Virtual key codes (replacing C.MMKeyCode / C.K_*)
const (
	K_NOT_A_KEY         uint32 = 0xFF
	K_BACKSPACE         uint32 = 0x08
	K_DELETE            uint32 = 0x2E
	K_RETURN            uint32 = 0x0D
	K_TAB               uint32 = 0x09
	K_ESCAPE            uint32 = 0x1B
	K_UP                uint32 = 0x26
	K_DOWN              uint32 = 0x28
	K_RIGHT             uint32 = 0x27
	K_LEFT              uint32 = 0x25
	K_HOME              uint32 = 0x24
	K_END               uint32 = 0x23
	K_PAGEUP            uint32 = 0x21
	K_PAGEDOWN          uint32 = 0x22
	K_F1                uint32 = 0x70
	K_F2                uint32 = 0x71
	K_F3                uint32 = 0x72
	K_F4                uint32 = 0x73
	K_F5                uint32 = 0x74
	K_F6                uint32 = 0x75
	K_F7                uint32 = 0x76
	K_F8                uint32 = 0x77
	K_F9                uint32 = 0x78
	K_F10               uint32 = 0x79
	K_F11               uint32 = 0x7A
	K_F12               uint32 = 0x7B
	K_F13               uint32 = 0x7C
	K_F14               uint32 = 0x7D
	K_F15               uint32 = 0x7E
	K_F16               uint32 = 0x7F
	K_F17               uint32 = 0x80
	K_F18               uint32 = 0x81
	K_F19               uint32 = 0x82
	K_F20               uint32 = 0x83
	K_F21               uint32 = 0x84
	K_F22               uint32 = 0x85
	K_F23               uint32 = 0x86
	K_F24               uint32 = 0x87
	K_META              uint32 = 0x5B // VK_LWIN
	K_LMETA             uint32 = 0x5B
	K_RMETA             uint32 = 0x5C
	K_ALT               uint32 = 0xA4 // VK_LMENU
	K_LALT              uint32 = 0xA4
	K_RALT              uint32 = 0xA5
	K_CONTROL           uint32 = 0xA2 // VK_LCONTROL
	K_LCONTROL          uint32 = 0xA2
	K_RCONTROL          uint32 = 0xA3
	K_SHIFT             uint32 = 0xA0 // VK_LSHIFT
	K_LSHIFT            uint32 = 0xA0
	K_RSHIFT            uint32 = 0xA1
	K_CAPSLOCK          uint32 = 0x14
	K_SPACE             uint32 = 0x20
	K_PRINTSCREEN       uint32 = 0x2C
	K_INSERT            uint32 = 0x2D
	K_MENU              uint32 = 0x5D
	K_AUDIO_VOLUME_MUTE uint32 = 0xAD
	K_AUDIO_VOLUME_DOWN uint32 = 0xAE
	K_AUDIO_VOLUME_UP   uint32 = 0xAF
	K_AUDIO_PLAY        uint32 = 0xB3
	K_AUDIO_STOP        uint32 = 0xB2
	K_AUDIO_PAUSE       uint32 = 0x13
	K_AUDIO_PREV        uint32 = 0xB1
	K_AUDIO_NEXT        uint32 = 0xB0
	K_AUDIO_REWIND      uint32 = 0xFF // Not standard on Windows
	K_AUDIO_FORWARD     uint32 = 0xFF
	K_AUDIO_REPEAT      uint32 = 0xFF
	K_AUDIO_RANDOM      uint32 = 0xFF
	K_NUMPAD_0          uint32 = 0x60
	K_NUMPAD_1          uint32 = 0x61
	K_NUMPAD_2          uint32 = 0x62
	K_NUMPAD_3          uint32 = 0x63
	K_NUMPAD_4          uint32 = 0x64
	K_NUMPAD_5          uint32 = 0x65
	K_NUMPAD_6          uint32 = 0x66
	K_NUMPAD_7          uint32 = 0x67
	K_NUMPAD_8          uint32 = 0x68
	K_NUMPAD_9          uint32 = 0x69
	K_NUMPAD_LOCK       uint32 = 0x90
	K_NUMPAD_DECIMAL    uint32 = 0x6E
	K_NUMPAD_PLUS       uint32 = 0x6B
	K_NUMPAD_MINUS      uint32 = 0x6D
	K_NUMPAD_MUL        uint32 = 0x6A
	K_NUMPAD_DIV        uint32 = 0x6F
	K_NUMPAD_CLEAR      uint32 = 0x0C
	K_NUMPAD_ENTER      uint32 = 0x0D
	K_NUMPAD_EQUAL      uint32 = 0x0D
	K_LIGHTS_MON_UP     uint32 = 0xFF
	K_LIGHTS_MON_DOWN   uint32 = 0xFF
	K_LIGHTS_KBD_TOGGLE uint32 = 0xFF
	K_LIGHTS_KBD_UP     uint32 = 0xFF
	K_LIGHTS_KBD_DOWN   uint32 = 0xFF
)

// MData represents window identification data (platform-specific)
type MData struct {
	HWnd  uintptr
	Title [512]uint16
}

// MMBitmap represents a bitmap in memory (replacing C.MMBitmapRef)
type MMBitmap struct {
	ImageBuffer   *uint8
	pixels        []byte // underlying pixel data (prevents GC)
	Width         int32
	Height        int32
	Bytewidth     int32
	BitsPerPixel  uint8
	BytesPerPixel uint8
}

// Internal state
var (
	globalHandle MData
)

// platformInit performs Windows-specific initialization
func platformInit() {
	// No special initialization needed for Windows
}

// --- Mouse operations ---

func platformMoveMouse(x, y int32) int {
	if !win32.SetCursorPos(x, y) {
		return 1
	}
	return 0
}

func platformDragMouse(x, y int32, button uint16) int {
	// Press the button first
	win32.ToggleMouse(true, button)
	// Move
	result := platformMoveMouse(x, y)
	// Release the button
	win32.ToggleMouse(false, button)
	return result
}

func platformToggleMouse(down bool, button uint16) int {
	return win32.ToggleMouse(down, button)
}

func platformClickMouse(button uint16) int {
	if code := win32.ToggleMouse(true, button); code != 0 {
		return code
	}
	MilliSleep(5)
	return win32.ToggleMouse(false, button)
}

func platformDoubleClick(button uint16, count int) int {
	for i := 0; i < count; i++ {
		if code := platformClickMouse(button); code != 0 {
			return code
		}
		MilliSleep(10)
	}
	return 0
}

func platformScrollMouse(x, y int) {
	win32.ScrollMouse(x, y)
}

func platformLocation() (int32, int32) {
	var pt win32.POINT
	win32.GetCursorPos(&pt)
	return pt.X, pt.Y
}

func platformSmoothlyMoveMouse(x, y int32, low, high float64) bool {
	// Pure Go implementation of smooth mouse movement
	return smoothlyMoveMouseImpl(x, y, low, high)
}

func platformSmoothlyDragMouse(x, y int32, low, high float64, button uint16) bool {
	win32.ToggleMouse(true, button)
	result := smoothlyMoveMouseImpl(x, y, low, high)
	win32.ToggleMouse(false, button)
	return result
}

// --- Screen operations ---

func platformGetPxColor(x, y int32, displayId int32) uint32 {
	return win32.GetPixelColor(x, y)
}

func platformGetMainDisplaySize() (int32, int32) {
	return win32.GetScreenSize()
}

func platformGetScreenRect(displayId int32) (int32, int32, int32, int32) {
	w, h := win32.GetScreenSize()
	return 0, 0, w, h
}

func platformSysScale(displayId int32) float64 {
	var hwnd uintptr
	if displayId >= 0 {
		hwnd = uintptr(displayId)
	} else {
		hwnd = win32.GetForegroundWindow()
	}
	dpi := win32.GetDpiForWindow(hwnd)
	return float64(dpi) / 96.0
}

func platformScaleX() int {
	hwnd := win32.GetForegroundWindow()
	dpi := win32.GetDpiForWindow(hwnd)
	return int(dpi)
}

func platformGetNumDisplays() int {
	return win32.GetNumDisplays()
}

func platformCaptureScreen(x, y, w, h int32, displayId int32, isPid int8) *MMBitmap {
	pixels, width, height := win32.CaptureScreen(x, y, w, h)
	if pixels == nil {
		return nil
	}

	bmp := &MMBitmap{
		Width:         int32(width),
		Height:        int32(height),
		Bytewidth:     int32(width * 4),
		BitsPerPixel:  32,
		BytesPerPixel: 4,
	}
	if len(pixels) > 0 {
		bmp.ImageBuffer = &pixels[0]
		// Store the slice to prevent GC
		bmp.pixels = pixels
	}

	return bmp
}

func platformBitmapDealloc(bmp *MMBitmap) {
	// Go GC will handle the memory
	bmp.ImageBuffer = nil
	bmp.pixels = nil
}

func platformCreateMMBitmap(imageBuffer *uint8, width, height, bytewidth int32, bitsPerPixel, bytesPerPixel uint8) *MMBitmap {
	return &MMBitmap{
		ImageBuffer:   imageBuffer,
		Width:         width,
		Height:        height,
		Bytewidth:     bytewidth,
		BitsPerPixel:  bitsPerPixel,
		BytesPerPixel: bytesPerPixel,
	}
}

// --- Keyboard operations ---

func platformToggleKeyCode(keyCode uint32, down bool, flags uint64, pid uintptr) int {
	return win32.ToggleKeyCode(keyCode, down, flags, pid)
}

func platformKeyCodeForChar(char byte) uint32 {
	return win32.KeyCodeForChar(char)
}

func platformUnicodeType(char uint32, pid uintptr, isPid int8) {
	win32.UnicodeType(char, pid, isPid)
}

func platformInputUTF(str string) {
	// Input UTF-8 string using UnicodeType
	for _, r := range str {
		win32.UnicodeType(uint32(r), 0, 0)
		MilliSleep(7)
	}
}

// --- Window operations ---

func platformShowAlert(title, msg, defaultBtn, cancelBtn string) bool {
	// Use Windows MessageBoxW
	messageBoxW := syscall.MustLoadDLL("user32.dll").MustFindProc("MessageBoxW")

	titlePtr, _ := syscall.UTF16PtrFromString(title)
	msgPtr, _ := syscall.UTF16PtrFromString(msg)

	var flags uint32 = 0x01 // MB_OK
	if cancelBtn != "" {
		flags = 0x11 // MB_OKCANCEL
	}

	ret, _, _ := messageBoxW.Call(0, uintptr(unsafe.Pointer(msgPtr)), uintptr(unsafe.Pointer(titlePtr)), uintptr(flags))
	// IDOK = 1, IDCANCEL = 2
	return ret == 1
}

func platformIsValid() bool {
	hwnd := globalHandle.HWnd
	if hwnd == 0 {
		hwnd = win32.GetForegroundWindow()
	}
	return win32.IsWindow(hwnd)
}

func platformSetActive(win MData) {
	if win.HWnd != 0 {
		win32.SetForegroundWindow(win.HWnd)
		win32.SetActiveWindowHwnd(win.HWnd)
	}
}

func platformGetActive() MData {
	hwnd := win32.GetForegroundWindow()
	return MData{HWnd: hwnd}
}

func platformMinWindow(pid uintptr, state bool, isPid int8) {
	var hwnd uintptr
	if isPid == 1 {
		hwnd = win32.FindWindowByPID(uint32(pid))
	} else {
		hwnd = pid
	}
	if hwnd != 0 {
		if state {
			win32.ShowWindow(hwnd, win32.SW_MINIMIZE)
		} else {
			win32.ShowWindow(hwnd, win32.SW_RESTORE)
		}
	}
}

func platformMaxWindow(pid uintptr, state bool, isPid int8) {
	var hwnd uintptr
	if isPid == 1 {
		hwnd = win32.FindWindowByPID(uint32(pid))
	} else {
		hwnd = pid
	}
	if hwnd != 0 {
		if state {
			win32.ShowWindow(hwnd, win32.SW_MAXIMIZE)
		} else {
			win32.ShowWindow(hwnd, win32.SW_RESTORE)
		}
	}
}

func platformCloseMainWindow() {
	hwnd := win32.GetForegroundWindow()
	if hwnd != 0 {
		win32.PostCloseMessage(hwnd)
	}
}

func platformCloseWindowByPID(pid uintptr, isPid int8) {
	var hwnd uintptr
	if isPid == 1 {
		hwnd = win32.FindWindowByPID(uint32(pid))
	} else {
		hwnd = pid
	}
	if hwnd != 0 {
		win32.PostCloseMessage(hwnd)
	}
}

func platformSetHandle(hwnd uintptr) {
	globalHandle.HWnd = hwnd
}

func platformSetHandlePid(pid uintptr, isPid int8) MData {
	var hwnd uintptr
	if isPid == 1 {
		hwnd = win32.FindWindowByPID(uint32(pid))
	} else {
		hwnd = pid
	}
	globalHandle.HWnd = hwnd
	return globalHandle
}

func platformSetHandlePidMData(pid uintptr, isPid int8) {
	var hwnd uintptr
	if isPid == 1 {
		hwnd = win32.FindWindowByPID(uint32(pid))
	} else {
		hwnd = pid
	}
	globalHandle.HWnd = hwnd
}

func platformGetHandle() uintptr {
	return globalHandle.HWnd
}

func platformBGetHandle() uintptr {
	return globalHandle.HWnd
}

func platformGetMainTitle() string {
	hwnd := win32.GetForegroundWindow()
	return win32.GetWindowText(hwnd)
}

func platformGetTitleByPid(pid uintptr, isPid int8) string {
	var hwnd uintptr
	if isPid == 1 {
		hwnd = win32.FindWindowByPID(uint32(pid))
	} else {
		hwnd = pid
	}
	if hwnd == 0 {
		return ""
	}
	return win32.GetWindowText(hwnd)
}

func platformGetPID() uintptr {
	hwnd := win32.GetForegroundWindow()
	return uintptr(win32.GetWindowThreadProcessId(hwnd))
}

type boundsRect struct {
	X, Y, W, H int32
}

func platformGetBounds(pid uintptr, isPid int8) boundsRect {
	var hwnd uintptr
	if isPid == 1 {
		hwnd = win32.FindWindowByPID(uint32(pid))
	} else {
		hwnd = pid
	}
	if hwnd == 0 {
		return boundsRect{}
	}
	rect, ok := win32.GetWindowRect(hwnd)
	if !ok {
		return boundsRect{}
	}
	return boundsRect{
		X: rect.Left,
		Y: rect.Top,
		W: rect.Right - rect.Left,
		H: rect.Bottom - rect.Top,
	}
}

func platformGetClient(pid uintptr, isPid int8) boundsRect {
	var hwnd uintptr
	if isPid == 1 {
		hwnd = win32.FindWindowByPID(uint32(pid))
	} else {
		hwnd = pid
	}
	if hwnd == 0 {
		return boundsRect{}
	}
	rect, ok := win32.GetClientRect(hwnd)
	if !ok {
		return boundsRect{}
	}
	return boundsRect{
		X: rect.Left,
		Y: rect.Top,
		W: rect.Right - rect.Left,
		H: rect.Bottom - rect.Top,
	}
}

func platformIs64Bit() bool {
	return uintptr(0) != 4 // Always true on 64-bit Windows, check at runtime
}

func platformActivePID(pid uintptr, isPid int8) {
	var hwnd uintptr
	if isPid == 1 {
		hwnd = win32.FindWindowByPID(uint32(pid))
	} else {
		hwnd = pid
	}
	if hwnd != 0 {
		win32.SetForegroundWindow(hwnd)
	}
}

func platformGetHwndByPid(pid uintptr) uintptr {
	return win32.FindWindowByPID(uint32(pid))
}

func platformSetXDisplayName(name string) error {
	return nil // Not applicable on Windows
}

func platformGetXDisplayName() string {
	return "" // Not applicable on Windows
}

func platformCloseMainDisplay() {
	// Not applicable on Windows
}

// --- FindWindow wraps win32.FindWindowByName ---
func FindWindow(name string) uintptr {
	return win32.FindWindowByName(name)
}

// --- GetHWND gets foreground window hwnd ---
func GetHWND() uintptr {
	return win32.GetForegroundWindow()
}

// --- SendInput wraps win32.SendInput ---
func SendInput(nInputs uint32, pInputs unsafe.Pointer, cbSize int32) uint32 {
	return win32.SendInput(nInputs, pInputs, cbSize)
}

// --- SendMsg sends a message to a window ---
func SendMsg(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := syscall.MustLoadDLL("user32.dll").MustFindProc("SendMessageW").Call(hwnd, uintptr(msg), wParam, lParam)
	return ret
}

// --- SetActiveWindow sets active window ---
func SetActiveWindow(hwnd uintptr) uintptr {
	return win32.SetActiveWindowHwnd(hwnd)
}

// --- SetFocus sets window focus ---
func SetFocus(hwnd uintptr) uintptr {
	return win32.SetFocusHwnd(hwnd)
}

// --- SetForeg sets foreground window ---
func SetForeg(hwnd uintptr) bool {
	return win32.SetForegroundWindow(hwnd)
}

// --- GetMain gets active window ---
func GetMain() uintptr {
	return win32.GetActiveWindow()
}

// --- GetMainId gets the main display id ---
func GetMainId() int {
	return int(GetMain())
}

// --- ScaleF gets the system scale value ---
func ScaleF(displayId ...int) (f float64) {
	if len(displayId) > 0 && displayId[0] != -1 {
		if displayId[0] >= 0 {
			dpi := win32.GetDpiForWindow(uintptr(displayId[0]))
			f = float64(dpi) / 96.0
		}
		if displayId[0] == -2 {
			f = float64(win32.GetDpiForWindow(win32.GetDesktopWindow())) / 96.0
		}
	} else {
		hwnd := win32.GetForegroundWindow()
		f = float64(win32.GetDpiForWindow(hwnd)) / 96.0
	}
	if f == 0.0 {
		f = 1.0
	}
	return f
}

// --- GetDesktopWindow gets desktop window hwnd ---
func GetDesktopWindow() uintptr {
	return win32.GetDesktopWindow()
}

// --- GetMainDPI gets the display dpi ---
func GetMainDPI() int {
	return int(win32.GetDpiForWindow(win32.GetForegroundWindow()))
}

// --- GetDPI gets the window dpi ---
func GetDPI(hwnd uintptr) uint32 {
	return win32.GetDpiForWindow(hwnd)
}

// --- GetSysDPI gets system metrics for DPI ---
func GetSysDPI(idx int32, dpi uint32) int32 {
	return win32.GetSystemMetricsForDpi(idx, dpi)
}

// --- DisplaysNum gets the count of displays ---
func DisplaysNum() int {
	return platformGetNumDisplays()
}

// --- Alert shows a Windows message box ---
func Alert(title, msg string, args ...string) bool {
	defaultBtn, cancelBtn := alertArgs(args...)
	return platformShowAlert(title, msg, defaultBtn, cancelBtn)
}

// --- Window management for darwin||windows build tag ---
func GetBounds(pid int, args ...int) (int, int, int, int) {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
	}
	b := platformGetBounds(uintptr(pid), int8(isPid))
	return int(b.X), int(b.Y), int(b.W), int(b.H)
}

func GetClient(pid int, args ...int) (int, int, int, int) {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
	}
	b := platformGetClient(uintptr(pid), int8(isPid))
	return int(b.X), int(b.Y), int(b.W), int(b.H)
}

func internalGetTitle(pid int, args ...int) string {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
	}
	return platformGetTitleByPid(uintptr(pid), int8(isPid))
}

func ActivePid(pid int, args ...int) error {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
	}
	platformActivePID(uintptr(pid), int8(isPid))
	return nil
}

// platformIs64Bit implementation
func init() {
	// Detect 64-bit at runtime
	_ = runtime.GOARCH // just to reference runtime
}

// padHex converts uint32 color to hex string
func padHex(hex uint32) string {
	return fmt.Sprintf("%06x", hex)
}

// colorHexToRGB converts hex color to RGB
func colorHexToRGB(hex uint32) (r, g, b uint8) {
	r = uint8((hex >> 16) & 0xFF)
	g = uint8((hex >> 8) & 0xFF)
	b = uint8(hex & 0xFF)
	return
}

// colorRGBToHex converts RGB to hex
func colorRGBToHex(r, g, b uint8) uint32 {
	return uint32(r)<<16 | uint32(g)<<8 | uint32(b)
}

// Bitmap to image conversion
func platformToImage(bmp *MMBitmap) image.Image {
	if bmp == nil || bmp.ImageBuffer == nil {
		return nil
	}

	img := image.NewRGBA(image.Rect(0, 0, int(bmp.Width), int(bmp.Height)))
	pixels := bmp.pixels
	if pixels == nil {
		return img
	}

	for y := 0; y < int(bmp.Height); y++ {
		for x := 0; x < int(bmp.Width); x++ {
			i := (y*int(bmp.Width) + x) * 4
			if i+3 < len(pixels) {
				// BGRA -> RGBA
				img.SetRGBA(x, y, color.RGBA{
					R: pixels[i+2],
					G: pixels[i+1],
					B: pixels[i],
					A: pixels[i+3],
				})
			}
		}
	}
	return img
}
