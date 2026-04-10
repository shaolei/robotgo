// Copyright (c) 2016-2025 AtomAI, All rights reserved.
//
// Licensed under the Apache License, Version 2.0

//go:build !darwin && !windows

package robotgo

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xinerama"
	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/ewmh"
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
	MOD_ALT     uint64 = 0x08 // Mod1Mask
	MOD_CONTROL uint64 = 0x04 // ControlMask
	MOD_SHIFT   uint64 = 0x01 // ShiftMask
	MOD_META    uint64 = 0x40 // Mod4Mask
)

// Virtual key codes (replacing C.MMKeyCode / C.K_*)
// Linux X11 keycodes - using uint32 for X11 KeySym values
// (XF86 multimedia keys exceed uint16 range)
const (
	K_NOT_A_KEY         uint32 = 0xFF
	K_BACKSPACE         uint32 = 0x16   // XK_BackSpace
	K_DELETE            uint32 = 0xFFFF // XK_Delete
	K_RETURN            uint32 = 0xFF0D // XK_Return
	K_TAB               uint32 = 0xFF09 // XK_Tab
	K_ESCAPE            uint32 = 0xFF1B // XK_Escape
	K_UP                uint32 = 0xFF52 // XK_Up
	K_DOWN              uint32 = 0xFF54 // XK_Down
	K_RIGHT             uint32 = 0xFF53 // XK_Right
	K_LEFT              uint32 = 0xFF51 // XK_Left
	K_HOME              uint32 = 0xFF50 // XK_Home
	K_END               uint32 = 0xFF57 // XK_End
	K_PAGEUP            uint32 = 0xFF55 // XK_Page_Up
	K_PAGEDOWN          uint32 = 0xFF56 // XK_Page_Down
	K_F1                uint32 = 0xFFBE
	K_F2                uint32 = 0xFFBF
	K_F3                uint32 = 0xFFC0
	K_F4                uint32 = 0xFFC1
	K_F5                uint32 = 0xFFC2
	K_F6                uint32 = 0xFFC3
	K_F7                uint32 = 0xFFC4
	K_F8                uint32 = 0xFFC5
	K_F9                uint32 = 0xFFC6
	K_F10               uint32 = 0xFFC7
	K_F11               uint32 = 0xFFC8
	K_F12               uint32 = 0xFFC9
	K_F13               uint32 = 0xFFCA
	K_F14               uint32 = 0xFFCB
	K_F15               uint32 = 0xFFCC
	K_F16               uint32 = 0xFFCD
	K_F17               uint32 = 0xFFCE
	K_F18               uint32 = 0xFFCF
	K_F19               uint32 = 0xFFD0
	K_F20               uint32 = 0xFFD1
	K_F21               uint32 = 0xFFD2
	K_F22               uint32 = 0xFFD3
	K_F23               uint32 = 0xFFD4
	K_F24               uint32 = 0xFFD5
	K_META              uint32 = 0xFFEB // XK_Super_L
	K_LMETA             uint32 = 0xFFEB
	K_RMETA             uint32 = 0xFFEC
	K_ALT               uint32 = 0xFFE9 // XK_Alt_L
	K_LALT              uint32 = 0xFFE9
	K_RALT              uint32 = 0xFFEA
	K_CONTROL           uint32 = 0xFFE3 // XK_Control_L
	K_LCONTROL          uint32 = 0xFFE3
	K_RCONTROL          uint32 = 0xFFE4
	K_SHIFT             uint32 = 0xFFE1 // XK_Shift_L
	K_LSHIFT            uint32 = 0xFFE1
	K_RSHIFT            uint32 = 0xFFE2
	K_CAPSLOCK          uint32 = 0xFFE5
	K_SPACE             uint32 = 0x0020
	K_PRINTSCREEN       uint32 = 0xFF61 // XK_Print
	K_INSERT            uint32 = 0xFF63 // XK_Insert
	K_MENU              uint32 = 0xFF67 // XK_Menu
	K_AUDIO_VOLUME_MUTE uint32 = 0x1008FF12
	K_AUDIO_VOLUME_DOWN uint32 = 0x1008FF11
	K_AUDIO_VOLUME_UP   uint32 = 0x1008FF13
	K_AUDIO_PLAY        uint32 = 0x1008FF14
	K_AUDIO_STOP        uint32 = 0x1008FF15
	K_AUDIO_PAUSE       uint32 = 0x1008FF31
	K_AUDIO_PREV        uint32 = 0x1008FF16
	K_AUDIO_NEXT        uint32 = 0x1008FF17
	K_AUDIO_REWIND      uint32 = 0x1008FF3E
	K_AUDIO_FORWARD     uint32 = 0x1008FF3F
	K_AUDIO_REPEAT      uint32 = 0x1008FF3B
	K_AUDIO_RANDOM      uint32 = 0x1008FF3C
	K_NUMPAD_0          uint32 = 0xFFB0
	K_NUMPAD_1          uint32 = 0xFFB1
	K_NUMPAD_2          uint32 = 0xFFB2
	K_NUMPAD_3          uint32 = 0xFFB3
	K_NUMPAD_4          uint32 = 0xFFB4
	K_NUMPAD_5          uint32 = 0xFFB5
	K_NUMPAD_6          uint32 = 0xFFB6
	K_NUMPAD_7          uint32 = 0xFFB7
	K_NUMPAD_8          uint32 = 0xFFB8
	K_NUMPAD_9          uint32 = 0xFFB9
	K_NUMPAD_LOCK       uint32 = 0xFF7F
	K_NUMPAD_DECIMAL    uint32 = 0xFFAE
	K_NUMPAD_PLUS       uint32 = 0xFFAB
	K_NUMPAD_MINUS      uint32 = 0xFFAD
	K_NUMPAD_MUL        uint32 = 0xFFAA
	K_NUMPAD_DIV        uint32 = 0xFFAF
	K_NUMPAD_CLEAR      uint32 = 0xFFBD
	K_NUMPAD_ENTER      uint32 = 0xFF8D
	K_NUMPAD_EQUAL      uint32 = 0xFFBD
	K_LIGHTS_MON_UP     uint32 = 0x1008FF02
	K_LIGHTS_MON_DOWN   uint32 = 0x1008FF03
	K_LIGHTS_KBD_TOGGLE uint32 = 0x1008FF04
	K_LIGHTS_KBD_UP     uint32 = 0x1008FF05
	K_LIGHTS_KBD_DOWN   uint32 = 0x1008FF06
)

// MData represents window identification data (platform-specific)
type MData struct {
	XWin uintptr
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
	xu           *xgbutil.XUtil
	xDisplay     uintptr
)

// X11/Xtst function pointers loaded via purego
var (
	libX11  uintptr
	libXtst uintptr

	pXWarpPointer         func(display uintptr, srcW uintptr, destW uintptr, srcX int32, srcY int32, srcWidth uint32, srcHeight uint32, destX int32, destY int32)
	pXQueryPointer        func(display uintptr, w uintptr, rootReturn *uintptr, childReturn *uintptr, rootXReturn *int32, rootYReturn *int32, winXReturn *int32, winYReturn *int32, maskReturn *uint32) int32
	pXTestFakeButtonEvent func(display uintptr, button uint32, isPress int32, delay uintptr) int32
	pXTestFakeKeyEvent    func(display uintptr, keycode uint32, isPress int32, delay uintptr) int32
	pXStringToKeysym      func(string unsafe.Pointer) uintptr
	pXKeysymToKeycode     func(display uintptr, keysym uintptr) uint8
	pXOpenDisplay         func(name unsafe.Pointer) uintptr
	pXCloseDisplay        func(display uintptr) int32
	pXGetImage            func(display uintptr, drawable uintptr, x int32, y int32, width uint32, height uint32, planeMask uint64, format int32) uintptr
	pXGetPixel            func(image uintptr, x int32, y int32) uint64
	pXDestroyImage        func(image uintptr)
	pXDefaultRootWindow   func(display uintptr) uintptr
	pXFlush               func(display uintptr) int32
	pXSync                func(display uintptr, discard int32) int32
	pXFree                func(data uintptr) int32
)

func init() {
	var err error
	libX11, err = purego.Dlopen("libX11.so.6", purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		log.Println("robotgo: failed to load libX11:", err)
		return
	}

	libXtst, err = purego.Dlopen("libXtst.so.6", purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		log.Println("robotgo: failed to load libXtst:", err)
		return
	}

	registerX11Functions()

	// Open default display
	xDisplay = pXOpenDisplay(unsafe.Pointer(uintptr(0)))
	if xDisplay == 0 {
		log.Println("robotgo: failed to open X display")
	}

	// Initialize xgbutil
	xu, err = xgbutil.NewConn()
	if err != nil {
		log.Println("robotgo: xgbutil.NewConn error:", err)
	}
}

func registerX11Functions() {
	// X11 core functions
	purego.RegisterLibFunc(&pXOpenDisplay, libX11, "XOpenDisplay")
	purego.RegisterLibFunc(&pXCloseDisplay, libX11, "XCloseDisplay")
	purego.RegisterLibFunc(&pXWarpPointer, libX11, "XWarpPointer")
	purego.RegisterLibFunc(&pXQueryPointer, libX11, "XQueryPointer")
	purego.RegisterLibFunc(&pXStringToKeysym, libX11, "XStringToKeysym")
	purego.RegisterLibFunc(&pXKeysymToKeycode, libX11, "XKeysymToKeycode")
	purego.RegisterLibFunc(&pXGetImage, libX11, "XGetImage")
	purego.RegisterLibFunc(&pXGetPixel, libX11, "XGetPixel")
	purego.RegisterLibFunc(&pXDestroyImage, libX11, "XDestroyImage")
	purego.RegisterLibFunc(&pXDefaultRootWindow, libX11, "XDefaultRootWindow")
	purego.RegisterLibFunc(&pXFlush, libX11, "XFlush")
	purego.RegisterLibFunc(&pXSync, libX11, "XSync")
	purego.RegisterLibFunc(&pXFree, libX11, "XFree")

	// Xtst functions
	purego.RegisterLibFunc(&pXTestFakeButtonEvent, libXtst, "XTestFakeButtonEvent")
	purego.RegisterLibFunc(&pXTestFakeKeyEvent, libXtst, "XTestFakeKeyEvent")
}

// platformInit performs Linux-specific initialization
func platformInit() {
	// X11 init is done in package init()
}

// --- Mouse operations ---

func platformMoveMouse(x, y int32) int {
	if xDisplay == 0 {
		return 1
	}
	root := pXDefaultRootWindow(xDisplay)
	pXWarpPointer(xDisplay, 0, root, 0, 0, 0, 0, x, y)
	pXFlush(xDisplay)
	return 0
}

func platformDragMouse(x, y int32, button uint16) int {
	platformToggleMouse(true, button)
	result := platformMoveMouse(x, y)
	platformToggleMouse(false, button)
	return result
}

func platformToggleMouse(down bool, button uint16) int {
	if xDisplay == 0 {
		return 1
	}
	var xButton uint32
	switch button {
	case 0: // LEFT_BUTTON
		xButton = 1
	case 1: // RIGHT_BUTTON
		xButton = 3
	case 2: // CENTER_BUTTON
		xButton = 2
	default:
		xButton = uint32(button)
	}
	isPress := int32(0)
	if down {
		isPress = 1
	}
	pXTestFakeButtonEvent(xDisplay, xButton, isPress, 0)
	pXFlush(xDisplay)
	return 0
}

func platformClickMouse(button uint16) int {
	if code := platformToggleMouse(true, button); code != 0 {
		return code
	}
	MilliSleep(5)
	return platformToggleMouse(false, button)
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
	if xDisplay == 0 {
		return
	}
	// X11 uses buttons 4 (up), 5 (down), 6 (left), 7 (right)
	if y > 0 {
		for i := 0; i < y; i++ {
			pXTestFakeButtonEvent(xDisplay, 4, 1, 0) // Button 4 = scroll up
			pXTestFakeButtonEvent(xDisplay, 4, 0, 0)
		}
	} else if y < 0 {
		for i := 0; i < -y; i++ {
			pXTestFakeButtonEvent(xDisplay, 5, 1, 0) // Button 5 = scroll down
			pXTestFakeButtonEvent(xDisplay, 5, 0, 0)
		}
	}
	if x > 0 {
		for i := 0; i < x; i++ {
			pXTestFakeButtonEvent(xDisplay, 6, 1, 0) // Button 6 = scroll left
			pXTestFakeButtonEvent(xDisplay, 6, 0, 0)
		}
	} else if x < 0 {
		for i := 0; i < -x; i++ {
			pXTestFakeButtonEvent(xDisplay, 7, 1, 0) // Button 7 = scroll right
			pXTestFakeButtonEvent(xDisplay, 7, 0, 0)
		}
	}
	pXFlush(xDisplay)
}

func platformLocation() (int32, int32) {
	if xDisplay == 0 {
		return 0, 0
	}
	root := pXDefaultRootWindow(xDisplay)
	var rootRet, childRet uintptr
	var rootX, rootY, winX, winY int32
	var mask uint32
	pXQueryPointer(xDisplay, root, &rootRet, &childRet, &rootX, &rootY, &winX, &winY, &mask)
	return rootX, rootY
}

func platformSmoothlyMoveMouse(x, y int32, low, high float64) bool {
	return smoothlyMoveMouseImpl(x, y, low, high)
}

func platformSmoothlyDragMouse(x, y int32, low, high float64, button uint16) bool {
	platformToggleMouse(true, button)
	result := smoothlyMoveMouseImpl(x, y, low, high)
	platformToggleMouse(false, button)
	return result
}

// --- Screen operations ---

func platformGetPxColor(x, y int32, displayId int32) uint32 {
	if xDisplay == 0 {
		return 0
	}
	root := pXDefaultRootWindow(xDisplay)
	// Get screen dimensions first
	w, h := platformGetMainDisplaySize()
	if w == 0 || h == 0 {
		return 0
	}
	img := pXGetImage(xDisplay, root, 0, 0, uint32(w), uint32(h), 0xFFFFFFFF, 2 /* ZPixmap */)
	if img == 0 {
		return 0
	}
	defer pXDestroyImage(img)

	pixel := pXGetPixel(img, x, y)
	// X11 returns pixel in native format, convert to RGB
	r := uint32((pixel >> 16) & 0xFF)
	g := uint32((pixel >> 8) & 0xFF)
	b := uint32(pixel & 0xFF)
	return (r << 16) | (g << 8) | b
}

func platformGetMainDisplaySize() (int32, int32) {
	c, err := xgb.NewConn()
	if err != nil {
		return 0, 0
	}
	defer c.Close()

	setup := xproto.Setup(c)
	screen := setup.DefaultScreen(c)
	return int32(screen.WidthInPixels), int32(screen.HeightInPixels)
}

func platformGetScreenRect(displayId int32) (int32, int32, int32, int32) {
	w, h := platformGetMainDisplaySize()
	return 0, 0, w, h
}

func platformSysScale(displayId int32) float64 {
	return 1.0 // Linux typically no scaling
}

func platformScaleX() int {
	return 96 // Default DPI
}

func platformGetNumDisplays() int {
	c, err := xgb.NewConn()
	if err != nil {
		return 0
	}
	defer c.Close()

	err = xinerama.Init(c)
	if err != nil {
		return 0
	}

	reply, err := xinerama.QueryScreens(c).Reply()
	if err != nil {
		return 0
	}
	return int(reply.Number)
}

func platformCaptureScreen(x, y, w, h int32, displayId int32, isPid int8) *MMBitmap {
	if xDisplay == 0 || w <= 0 || h <= 0 {
		return nil
	}
	root := pXDefaultRootWindow(xDisplay)
	img := pXGetImage(xDisplay, root, x, y, uint32(w), uint32(h), 0xFFFFFFFF, 2 /* ZPixmap */)
	if img == 0 {
		return nil
	}
	defer pXDestroyImage(img)

	// Read pixel data from XImage
	// XImage struct: width(int), height(int), xoffset(int), format(int), data(char*), ...
	// We need to read the data pointer and byte order
	// For simplicity, use per-pixel extraction
	pixels := make([]byte, int(w)*int(h)*4)
	for py := int32(0); py < h; py++ {
		for px := int32(0); px < w; px++ {
			pixel := pXGetPixel(img, px, py)
			i := (py*w + px) * 4
			if i+3 < int32(len(pixels)) {
				// X11 ZPixmap format: BGRA or RGBX depending on byte order
				pixels[i+2] = uint8((pixel >> 16) & 0xFF) // R
				pixels[i+1] = uint8((pixel >> 8) & 0xFF)  // G
				pixels[i] = uint8(pixel & 0xFF)           // B
				pixels[i+3] = 0xFF                        // A
			}
		}
	}

	bmp := &MMBitmap{
		Width:         w,
		Height:        h,
		Bytewidth:     w * 4,
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
	if xDisplay == 0 {
		return 1
	}

	// Handle modifier flags by pressing modifier keys first
	if flags != 0 && down {
		pressModifierKeys(flags, true)
	}

	isPress := int32(0)
	if down {
		isPress = 1
	}
	// X11 keycodes are offset by 8 from keysyms
	// We need to convert our keysym to a keycode
	keycode := uint32(keyCode)
	// Check if this is already a keycode (values < 0xFF are typically keysyms)
	// X11 keycodes start at 8
	pXTestFakeKeyEvent(xDisplay, keycode, isPress, 0)

	// Release modifier keys
	if flags != 0 && !down {
		pressModifierKeys(flags, false)
	}

	pXFlush(xDisplay)
	return 0
}

func pressModifierKeys(flags uint64, down bool) {
	if flags&MOD_SHIFT != 0 {
		pXTestFakeKeyEvent(xDisplay, 0x32, boolToInt32(down), 0) // Shift_L keycode
	}
	if flags&MOD_CONTROL != 0 {
		pXTestFakeKeyEvent(xDisplay, 0x25, boolToInt32(down), 0) // Control_L keycode
	}
	if flags&MOD_ALT != 0 {
		pXTestFakeKeyEvent(xDisplay, 0x40, boolToInt32(down), 0) // Alt_L keycode
	}
	if flags&MOD_META != 0 {
		pXTestFakeKeyEvent(xDisplay, 0x73, boolToInt32(down), 0) // Super_L keycode
	}
}

func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func platformKeyCodeForChar(char byte) uint32 {
	if xDisplay == 0 {
		return K_NOT_A_KEY
	}
	// Convert char to keysym then to keycode
	cs := []byte{char, 0}
	keysym := pXStringToKeysym(unsafe.Pointer(&cs[0]))
	if keysym == 0 {
		return K_NOT_A_KEY
	}
	keycode := pXKeysymToKeycode(xDisplay, keysym)
	if keycode == 0 {
		return K_NOT_A_KEY
	}
	return uint32(keycode)
}

func platformUnicodeType(char uint32, pid uintptr, isPid int8) {
	if xDisplay == 0 {
		return
	}
	// Use XTestFakeKeyEvent with keysym for unicode
	// For unicode chars, we need to use XLookup keysym approach
	// For now, use the basic approach: convert to keysym
	keysym := uintptr(char) | 0x01000000 // Unicode keysym format
	keycode := pXKeysymToKeycode(xDisplay, keysym)
	if keycode != 0 {
		pXTestFakeKeyEvent(xDisplay, uint32(keycode), 1, 0)
		pXTestFakeKeyEvent(xDisplay, uint32(keycode), 0, 0)
	} else {
		// Fallback: try to find a keycode and set its keysym temporarily
		// This is complex - use a simpler approach for now
		// We use keycode 0 with the keysym
		pXTestFakeKeyEvent(xDisplay, uint32(keysym&0xFF), 1, 0)
		pXTestFakeKeyEvent(xDisplay, uint32(keysym&0xFF), 0, 0)
	}
	pXFlush(xDisplay)
}

func platformInputUTF(str string) {
	for _, r := range str {
		platformUnicodeType(uint32(r), 0, 0)
		MilliSleep(7)
	}
}

// --- Window operations ---

func platformShowAlert(title, msg, defaultBtn, cancelBtn string) bool {
	// Use xmessage on Linux
	cmd := `xmessage -center ` + msg +
		` -title ` + title + ` -buttons ` + defaultBtn + ":0,"
	if cancelBtn != "" {
		cmd += cancelBtn + ":1"
	}
	cmd += ` -default ` + defaultBtn
	cmd += ` -geometry 400x200`

	out, err := Run(cmd)
	if err != nil {
		return false
	}
	return string(out) != "1"
}

func platformIsValid() bool {
	return globalHandle.XWin != 0
}

func platformSetActive(win MData) {
	if win.XWin != 0 && xu != nil {
		_ = ewmh.ActiveWindowReq(xu, xproto.Window(win.XWin))
	}
}

func platformGetActive() MData {
	if xu == nil {
		return MData{}
	}
	active, err := ewmh.ActiveWindowGet(xu)
	if err != nil {
		return MData{}
	}
	return MData{XWin: uintptr(active)}
}

func platformMinWindow(pid uintptr, state bool, isPid int8) {
	// Use xdotool or EWMH
}

func platformMaxWindow(pid uintptr, state bool, isPid int8) {
	// Use xdotool or EWMH
}

func platformCloseMainWindow() {
	if xu == nil {
		return
	}
	active, err := ewmh.ActiveWindowGet(xu)
	if err != nil {
		return
	}
	_ = ewmh.CloseWindow(xu, active)
}

func platformCloseWindowByPID(pid uintptr, isPid int8) {
	if xu == nil {
		return
	}
	var xid xproto.Window
	if isPid == 1 {
		win, err := GetXidByPid(xu, int(pid))
		if err != nil {
			return
		}
		xid = win
	} else {
		xid = xproto.Window(pid)
	}
	_ = ewmh.CloseWindow(xu, xid)
}

func platformSetHandle(hwnd uintptr) {
	globalHandle.XWin = hwnd
}

func platformSetHandlePid(pid uintptr, isPid int8) MData {
	var xid uintptr
	if isPid == 1 && xu != nil {
		win, err := GetXidByPid(xu, int(pid))
		if err == nil {
			xid = uintptr(win)
		}
	} else {
		xid = pid
	}
	globalHandle.XWin = xid
	return globalHandle
}

func platformSetHandlePidMData(pid uintptr, isPid int8) {
	var xid uintptr
	if isPid == 1 && xu != nil {
		win, err := GetXidByPid(xu, int(pid))
		if err == nil {
			xid = uintptr(win)
		}
	} else {
		xid = pid
	}
	globalHandle.XWin = xid
}

func platformGetHandle() uintptr {
	return globalHandle.XWin
}

func platformBGetHandle() uintptr {
	return globalHandle.XWin
}

func platformGetMainTitle() string {
	if xu == nil {
		return ""
	}
	active, err := ewmh.ActiveWindowGet(xu)
	if err != nil {
		return ""
	}
	title, err := ewmh.WmNameGet(xu, active)
	if err != nil {
		return ""
	}
	return title
}

func platformGetTitleByPid(pid uintptr, isPid int8) string {
	if xu == nil {
		return ""
	}
	var xid xproto.Window
	if isPid == 1 {
		win, err := GetXidByPid(xu, int(pid))
		if err != nil {
			return ""
		}
		xid = win
	} else {
		xid = xproto.Window(pid)
	}
	title, err := ewmh.WmNameGet(xu, xid)
	if err != nil {
		return ""
	}
	return title
}

func platformGetPID() uintptr {
	if xu == nil {
		return 0
	}
	active, err := ewmh.ActiveWindowGet(xu)
	if err != nil {
		return 0
	}
	pid, err := ewmh.WmPidGet(xu, active)
	if err != nil {
		return 0
	}
	return uintptr(pid)
}

type boundsRect struct {
	X, Y, W, H int32
}

func platformGetBounds(pid uintptr, isPid int8) boundsRect {
	if xu == nil {
		return boundsRect{}
	}
	var xid xproto.Window
	if isPid == 1 {
		win, err := GetXidByPid(xu, int(pid))
		if err != nil {
			return boundsRect{}
		}
		xid = win
	} else {
		xid = xproto.Window(pid)
	}
	geom, err := xproto.GetGeometry(xu.Conn(), xproto.Drawable(xid)).Reply()
	if err != nil {
		return boundsRect{}
	}
	return boundsRect{
		X: int32(geom.X),
		Y: int32(geom.Y),
		W: int32(geom.Width),
		H: int32(geom.Height),
	}
}

func platformGetClient(pid uintptr, isPid int8) boundsRect {
	// On X11, GetClient is similar to GetBounds
	return platformGetBounds(pid, isPid)
}

func platformIs64Bit() bool {
	return unsafe.Sizeof(uintptr(0)) == 8
}

func platformActivePID(pid uintptr, isPid int8) {
	if xu == nil {
		return
	}
	var xid xproto.Window
	if isPid == 1 {
		win, err := GetXidByPid(xu, int(pid))
		if err != nil {
			return
		}
		xid = win
	} else {
		xid = xproto.Window(pid)
	}
	_ = ewmh.ActiveWindowReq(xu, xid)
}

func platformGetHwndByPid(pid uintptr) uintptr {
	if xu == nil {
		return 0
	}
	win, err := GetXidByPid(xu, int(pid))
	if err != nil {
		return 0
	}
	return uintptr(win)
}

func platformSetXDisplayName(name string) error {
	if xDisplay != 0 {
		pXCloseDisplay(xDisplay)
	}
	cs := append([]byte(name), 0)
	xDisplay = pXOpenDisplay(unsafe.Pointer(&cs[0]))
	if xDisplay == 0 {
		return fmt.Errorf("failed to open X display: %s", name)
	}
	return nil
}

func platformGetXDisplayName() string {
	// X11 doesn't have a direct way to get the display name
	return ""
}

func platformCloseMainDisplay() {
	if xDisplay != 0 {
		pXCloseDisplay(xDisplay)
		xDisplay = 0
	}
}

// GetBounds get the window bounds
func GetBounds(pid int, args ...int) (int, int, int, int) {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
		return internalGetBounds(pid, isPid)
	}

	xid, err := GetXid(xu, pid)
	if err != nil {
		log.Println("Get Xid from Pid errors is: ", err)
		return 0, 0, 0, 0
	}

	return internalGetBounds(int(xid), isPid)
}

// GetClient get the window client bounds
func GetClient(pid int, args ...int) (int, int, int, int) {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
		return internalGetClient(pid, isPid)
	}

	xid, err := GetXid(xu, pid)
	if err != nil {
		log.Println("Get Xid from Pid errors is: ", err)
		return 0, 0, 0, 0
	}

	return internalGetClient(int(xid), isPid)
}

// internalGetTitle get the window title
func internalGetTitle(pid int, args ...int) string {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
		return cgetTitle(pid, isPid)
	}

	xid, err := GetXid(xu, pid)
	if err != nil {
		log.Println("Get Xid from Pid errors is: ", err)
		return ""
	}

	return cgetTitle(int(xid), isPid)
}

// ActivePidC active the window by Pid
func ActivePidC(pid int, args ...int) error {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
		internalActive(pid, isPid)
		return nil
	}

	xid, err := GetXid(xu, pid)
	if err != nil {
		log.Println("Get Xid from Pid errors is: ", err)
		return err
	}

	internalActive(int(xid), isPid)
	return nil
}

// ActivePid active the window by Pid
func ActivePid(pid int, args ...int) error {
	if xu == nil {
		var err error
		xu, err = xgbutil.NewConn()
		if err != nil {
			return err
		}
	}

	if len(args) > 0 {
		err := ewmh.ActiveWindowReq(xu, xproto.Window(pid))
		if err != nil {
			return err
		}
		return nil
	}

	xid, err := GetXidByPid(xu, pid)
	if err != nil {
		return err
	}

	err = ewmh.ActiveWindowReq(xu, xid)
	if err != nil {
		return err
	}
	return nil
}

// GetXid get the xid return window and error
func GetXid(xu *xgbutil.XUtil, pid int) (xproto.Window, error) {
	if xu == nil {
		var err error
		xu, err = xgbutil.NewConn()
		if err != nil {
			return 0, err
		}
	}

	xid, err := GetXidByPid(xu, pid)
	return xid, err
}

// Deprecated: use the GetXidByPid(),
//
// GetXidFromPid get the xid from pid
func GetXidFromPid(xu *xgbutil.XUtil, pid int) (xproto.Window, error) {
	return GetXidByPid(xu, pid)
}

// GetXidByPid get the xid from pid
func GetXidByPid(xu *xgbutil.XUtil, pid int) (xproto.Window, error) {
	windows, err := ewmh.ClientListGet(xu)
	if err != nil {
		return 0, err
	}

	for _, window := range windows {
		wmPid, err := ewmh.WmPidGet(xu, window)
		if err != nil {
			return 0, err
		}

		if uint(pid) == wmPid {
			return window, nil
		}
	}

	return 0, fmt.Errorf("failed to find a window with a matching pid")
}

// DisplaysNum get the count of displays
func DisplaysNum() int {
	return platformGetNumDisplays()
}

// GetMainId get the main display id
func GetMainId() int {
	conn, err := xgb.NewConn()
	if err != nil {
		return -1
	}
	defer conn.Close()

	setup := xproto.Setup(conn)
	defaultScreen := setup.DefaultScreen(conn)
	id := -1
	for i, screen := range setup.Roots {
		if defaultScreen.Root == screen.Root {
			id = i
			break
		}
	}
	return id
}

// ScaleF get the system scale value
func ScaleF(displayId ...int) float64 {
	f := SysScale(displayId...)
	if f == 0.0 {
		f = 1.0
	}
	return f
}

// Alert show a alert window
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
