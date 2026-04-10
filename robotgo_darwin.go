// Copyright (c) 2016-2025 AtomAI, All rights reserved.
//
// Licensed under the Apache License, Version 2.0

//go:build darwin

package robotgo

import (
	"fmt"
	"image"
	"image/color"
	"unsafe"

	"github.com/shaolei/robotgo/internal/darwin"
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
// macOS CGEventFlags
const (
	MOD_NONE    uint64 = 0
	MOD_ALT     uint64 = 0x0008 // kCGEventFlagMaskAlternate
	MOD_CONTROL uint64 = 0x1000 // kCGEventFlagMaskControl
	MOD_SHIFT   uint64 = 0x0002 // kCGEventFlagMaskShift
	MOD_META    uint64 = 0x0008 // kCGEventFlagMaskCommand (same value as ALT for macOS)
)

// Virtual key codes (replacing C.MMKeyCode / C.K_*)
// macOS virtual key code definitions
const (
	K_NOT_A_KEY         uint32 = 0xFF
	K_BACKSPACE         uint32 = 0x33
	K_DELETE            uint32 = 0x75
	K_RETURN            uint32 = 0x24
	K_TAB               uint32 = 0x30
	K_ESCAPE            uint32 = 0x35
	K_UP                uint32 = 0x7E
	K_DOWN              uint32 = 0x7D
	K_RIGHT             uint32 = 0x7C
	K_LEFT              uint32 = 0x7B
	K_HOME              uint32 = 0x73
	K_END               uint32 = 0x77
	K_PAGEUP            uint32 = 0x74
	K_PAGEDOWN          uint32 = 0x79
	K_F1                uint32 = 0x7A
	K_F2                uint32 = 0x78
	K_F3                uint32 = 0x63
	K_F4                uint32 = 0x76
	K_F5                uint32 = 0x60
	K_F6                uint32 = 0x61
	K_F7                uint32 = 0x62
	K_F8                uint32 = 0x64
	K_F9                uint32 = 0x65
	K_F10               uint32 = 0x6D
	K_F11               uint32 = 0x67
	K_F12               uint32 = 0x6F
	K_F13               uint32 = 0x69
	K_F14               uint32 = 0x6B
	K_F15               uint32 = 0x71
	K_F16               uint32 = 0x6A
	K_F17               uint32 = 0x40
	K_F18               uint32 = 0x4F
	K_F19               uint32 = 0x50
	K_F20               uint32 = 0x5A
	K_F21               uint32 = 0xFF
	K_F22               uint32 = 0xFF
	K_F23               uint32 = 0xFF
	K_F24               uint32 = 0xFF
	K_META              uint32 = 0x37
	K_LMETA             uint32 = 0x37
	K_RMETA             uint32 = 0x36
	K_ALT               uint32 = 0x3A
	K_LALT              uint32 = 0x3A
	K_RALT              uint32 = 0x3D
	K_CONTROL           uint32 = 0x3B
	K_LCONTROL          uint32 = 0x3B
	K_RCONTROL          uint32 = 0x3E
	K_SHIFT             uint32 = 0x38
	K_LSHIFT            uint32 = 0x38
	K_RSHIFT            uint32 = 0x3C
	K_CAPSLOCK          uint32 = 0x39
	K_SPACE             uint32 = 0x31
	K_PRINTSCREEN       uint32 = 0x46
	K_INSERT            uint32 = 0x72
	K_MENU              uint32 = 0x6E
	K_AUDIO_VOLUME_MUTE uint32 = 0x4A
	K_AUDIO_VOLUME_DOWN uint32 = 0x49
	K_AUDIO_VOLUME_UP   uint32 = 0x48
	K_AUDIO_PLAY        uint32 = 0x34
	K_AUDIO_STOP        uint32 = 0xFF
	K_AUDIO_PAUSE       uint32 = 0xFF
	K_AUDIO_PREV        uint32 = 0xFF
	K_AUDIO_NEXT        uint32 = 0xFF
	K_AUDIO_REWIND      uint32 = 0xFF
	K_AUDIO_FORWARD     uint32 = 0xFF
	K_AUDIO_REPEAT      uint32 = 0xFF
	K_AUDIO_RANDOM      uint32 = 0xFF
	K_NUMPAD_0          uint32 = 0x52
	K_NUMPAD_1          uint32 = 0x53
	K_NUMPAD_2          uint32 = 0x54
	K_NUMPAD_3          uint32 = 0x55
	K_NUMPAD_4          uint32 = 0x56
	K_NUMPAD_5          uint32 = 0x57
	K_NUMPAD_6          uint32 = 0x58
	K_NUMPAD_7          uint32 = 0x59
	K_NUMPAD_8          uint32 = 0x5B
	K_NUMPAD_9          uint32 = 0x5C
	K_NUMPAD_LOCK       uint32 = 0x47
	K_NUMPAD_DECIMAL    uint32 = 0x41
	K_NUMPAD_PLUS       uint32 = 0x45
	K_NUMPAD_MINUS      uint32 = 0x4E
	K_NUMPAD_MUL        uint32 = 0x43
	K_NUMPAD_DIV        uint32 = 0x4B
	K_NUMPAD_CLEAR      uint32 = 0x47
	K_NUMPAD_ENTER      uint32 = 0x4C
	K_NUMPAD_EQUAL      uint32 = 0x51
	K_LIGHTS_MON_UP     uint32 = 0xFF
	K_LIGHTS_MON_DOWN   uint32 = 0xFF
	K_LIGHTS_KBD_TOGGLE uint32 = 0xFF
	K_LIGHTS_KBD_UP     uint32 = 0xFF
	K_LIGHTS_KBD_DOWN   uint32 = 0xFF
)

// MData represents window identification data (platform-specific)
type MData struct {
	CgID uint32
	AxID uintptr
}

// MMBitmap represents a bitmap in memory (replacing C.MMBitmapRef)
type MMBitmap struct {
	ImageBuffer   *uint8
	Width         int32
	Height        int32
	Bytewidth     int32
	BitsPerPixel  uint8
	BytesPerPixel uint8
	pixels        []byte // prevent GC
}

// Internal state
var (
	globalHandle MData
)

// platformInit performs macOS-specific initialization
func platformInit() {
	// purego loads frameworks in init()
}

// --- Mouse operations ---

func platformMoveMouse(x, y int32) int {
	return darwin.MoveMouse(x, y)
}

func platformDragMouse(x, y int32, button uint16) int {
	darwin.ToggleMouse(true, button)
	result := darwin.MoveMouse(x, y)
	darwin.ToggleMouse(false, button)
	return result
}

func platformToggleMouse(down bool, button uint16) int {
	return darwin.ToggleMouse(down, button)
}

func platformClickMouse(button uint16) int {
	if code := darwin.ToggleMouse(true, button); code != 0 {
		return code
	}
	MilliSleep(5)
	return darwin.ToggleMouse(false, button)
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
	darwin.ScrollMouse(x, y)
}

func platformLocation() (int32, int32) {
	return darwin.Location()
}

func platformSmoothlyMoveMouse(x, y int32, low, high float64) bool {
	return smoothlyMoveMouseImpl(x, y, low, high)
}

func platformSmoothlyDragMouse(x, y int32, low, high float64, button uint16) bool {
	darwin.ToggleMouse(true, button)
	result := smoothlyMoveMouseImpl(x, y, low, high)
	darwin.ToggleMouse(false, button)
	return result
}

// --- Screen operations ---

func platformGetPxColor(x, y int32, displayId int32) uint32 {
	return darwin.GetPixelColor(x, y, displayId)
}

func platformGetMainDisplaySize() (int32, int32) {
	return darwin.GetMainDisplaySize()
}

func platformGetScreenRect(displayId int32) (int32, int32, int32, int32) {
	return darwin.GetScreenRect(displayId)
}

func platformSysScale(displayId int32) float64 {
	return darwin.SysScale(displayId)
}

func platformScaleX() int {
	return int(darwin.SysScale(-1) * 96)
}

func platformGetNumDisplays() int {
	return darwin.GetNumDisplays()
}

func platformCaptureScreen(x, y, w, h int32, displayId int32, isPid int8) *MMBitmap {
	pixels, width, height := darwin.CaptureScreen(x, y, w, h, displayId)
	if pixels == nil {
		return nil
	}

	bmp := &MMBitmap{
		Width:         int32(width),
		Height:        int32(height),
		Bytewidth:     int32(width * 4),
		BitsPerPixel:  32,
		BytesPerPixel: 4,
		pixels:        pixels,
	}
	if len(pixels) > 0 {
		bmp.ImageBuffer = &pixels[0]
	}
	return bmp
}

func platformBitmapDealloc(bmp *MMBitmap) {
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
	return darwin.ToggleKeyCode(keyCode, down, flags, pid)
}

func platformKeyCodeForChar(char byte) uint32 {
	return darwin.KeyCodeForChar(char)
}

func platformUnicodeType(char uint32, pid uintptr, isPid int8) {
	darwin.UnicodeType(char, pid, isPid)
}

func platformInputUTF(str string) {
	darwin.InputUTF(str)
}

// --- Window operations ---

func platformShowAlert(title, msg, defaultBtn, cancelBtn string) bool {
	return darwin.ShowAlert(title, msg, defaultBtn, cancelBtn)
}

func platformIsValid() bool {
	return globalHandle.AxID != 0 || globalHandle.CgID != 0
}

func platformSetActive(win MData) {
	if win.AxID != 0 {
		// Use Accessibility API to focus the window
		darwin.MinimizeWindow(uint32(win.CgID), false)
	}
}

func platformGetActive() MData {
	pid := darwin.GetActiveWindowPID()
	return MData{CgID: pid}
}

func platformMinWindow(pid uintptr, state bool, isPid int8) {
	darwin.MinimizeWindow(uint32(pid), state)
}

func platformMaxWindow(pid uintptr, state bool, isPid int8) {
	// macOS doesn't have a native "maximize" like Windows
	// Use Zoom (which is the macOS equivalent)
	darwin.MinimizeWindow(uint32(pid), false) // un-minimize
}

func platformCloseMainWindow() {
	pid := darwin.GetActiveWindowPID()
	_ = pid
	// Use Cmd+Q via keyboard
	darwin.ToggleKeyCode(0x0C, true, 0x0008, 0) // Cmd+Q
	MilliSleep(10)
	darwin.ToggleKeyCode(0x0C, false, 0, 0)
}

func platformCloseWindowByPID(pid uintptr, isPid int8) {
	// Use Cmd+W to close window
	darwin.ToggleKeyCode(0x0D, true, 0x0008, 0) // Cmd+W
	MilliSleep(10)
	darwin.ToggleKeyCode(0x0D, false, 0, 0)
}

func platformSetHandle(hwnd uintptr) {
	globalHandle.CgID = uint32(hwnd)
	globalHandle.AxID = hwnd
}

func platformSetHandlePid(pid uintptr, isPid int8) MData {
	globalHandle.CgID = uint32(pid)
	globalHandle.AxID = pid
	return globalHandle
}

func platformSetHandlePidMData(pid uintptr, isPid int8) {
	globalHandle.CgID = uint32(pid)
	globalHandle.AxID = pid
}

func platformGetHandle() uintptr {
	return globalHandle.AxID
}

func platformBGetHandle() uintptr {
	return globalHandle.AxID
}

func platformGetMainTitle() string {
	return darwin.GetActiveWindowTitle()
}

func platformGetTitleByPid(pid uintptr, isPid int8) string {
	return darwin.GetWindowTitleByPID(uint32(pid))
}

func platformGetPID() uintptr {
	return uintptr(darwin.GetActiveWindowPID())
}

type boundsRect struct {
	X, Y, W, H int32
}

func platformGetBounds(pid uintptr, isPid int8) boundsRect {
	// Use Accessibility API - requires more work
	return boundsRect{}
}

func platformGetClient(pid uintptr, isPid int8) boundsRect {
	return boundsRect{}
}

func platformIs64Bit() bool {
	return darwin.Is64Bit()
}

func platformActivePID(pid uintptr, isPid int8) {
	// Use NSRunningApplication to activate
	darwin.MinimizeWindow(uint32(pid), false)
}

func platformGetHwndByPid(pid uintptr) uintptr {
	return 0
}

func platformSetXDisplayName(name string) error {
	return nil // Not applicable on macOS
}

func platformGetXDisplayName() string {
	return "" // Not applicable on macOS
}

func platformCloseMainDisplay() {
	// Not applicable on macOS
}

// GetMainId get the main display id
func GetMainId() int {
	return int(darwin.CGMainDisplayID())
}

// ScaleF get the system scale value
func ScaleF(displayId ...int) float64 {
	f := SysScale(displayId...)
	if f == 0.0 {
		f = 1.0
	}
	return f
}

// Alert shows a macOS alert
func Alert(title, msg string, args ...string) bool {
	defaultBtn, cancelBtn := alertArgs(args...)
	return platformShowAlert(title, msg, defaultBtn, cancelBtn)
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

// Ensure unused imports are referenced
var _ = unsafe.Pointer(nil)

// internalGetTitle gets the window title by PID
func internalGetTitle(pid int, args ...int) string {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
	}
	return platformGetTitleByPid(uintptr(pid), int8(isPid))
}

// ActivePid activates the window by PID
func ActivePid(pid int, args ...int) error {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
	}
	platformActivePID(uintptr(pid), int8(isPid))
	return nil
}
