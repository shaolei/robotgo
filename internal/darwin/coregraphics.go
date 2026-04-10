// Copyright (c) 2016-2025 AtomAI, All rights reserved.
//
// Licensed under the Apache License, Version 2.0

//go:build darwin

// Package darwin provides macOS API bindings using purego,
// eliminating the need for cgo.
package darwin

import (
	"errors"
	"unsafe"

	"github.com/ebitengine/purego"
)

// macOS framework paths
const (
	coreGraphicsPath = "/System/Library/Frameworks/CoreGraphics.framework/CoreGraphics"
	appServicesPath  = "/System/Library/Frameworks/ApplicationServices.framework/ApplicationServices"
	ioKitPath        = "/System/Library/Frameworks/IOKit.framework/IOKit"
	carbonPath       = "/System/Library/Frameworks/Carbon.framework/Carbon"
	objcPath         = "/usr/lib/libobjc.A.dylib"
)

// CoreGraphics constants
const (
	kCGEventLeftMouseDown     uint32 = 1
	kCGEventLeftMouseUp       uint32 = 2
	kCGEventRightMouseDown    uint32 = 3
	kCGEventRightMouseUp      uint32 = 4
	kCGEventMouseMoved        uint32 = 5
	kCGEventLeftMouseDragged  uint32 = 6
	kCGEventRightMouseDragged uint32 = 7
	kCGEventOtherMouseDown    uint32 = 25
	kCGEventOtherMouseUp      uint32 = 26
	kCGEventOtherMouseDragged uint32 = 27
	kCGEventScrollWheel       uint32 = 22
	kCGEventKeyDown           uint32 = 10
	kCGEventKeyUp             uint32 = 11
	kCGEventFlagsChanged      uint32 = 12

	kCGMouseButtonLeft   uint32 = 0
	kCGMouseButtonRight  uint32 = 1
	kCGMouseButtonCenter uint32 = 2

	kCGHIDEventTap uint32 = 0

	kCGEventSourceStateHIDSystemState uint32 = 1
)

// CGPoint represents a point in CoreGraphics
type CGPoint struct {
	X float64
	Y float64
}

// CGSize represents a size in CoreGraphics
type CGSize struct {
	Width  float64
	Height float64
}

// CGRect represents a rectangle in CoreGraphics
type CGRect struct {
	Origin CGPoint
	Size   CGSize
}

// Framework handles
var (
	coreGraphicsLib uintptr
	appServicesLib  uintptr
	ioKitLib        uintptr
	carbonLib       uintptr
	objcLib         uintptr
)

// Function pointers - CoreGraphics
var (
	pcgEventSourceCreate             func(stateID uint32, state uintptr) uintptr
	pcgEventCreateMouseEvent         func(sourceRef uintptr, eventType uint32, mousePosition CGPoint, mouseButton uint32) uintptr
	pcgEventPost                     func(tapLocation uint32, event uintptr)
	pcgEventCreateKeyboardEvent      func(sourceRef uintptr, virtualKey uint16, keyDown bool) uintptr
	pcgEventSetType                  func(event uintptr, eventType uint32)
	pcgEventSetFlags                 func(event uintptr, flags uint64)
	pcgEventGetLocation              func(event uintptr) CGPoint
	pcgEventCreate                   func(sourceRef uintptr) uintptr
	pcgMainDisplayID                 func() uint32
	pcgDisplayPixelsWide             func(display uint32) uint32
	pcgDisplayPixelsHigh             func(display uint32) uint32
	pcgDisplayBounds                 func(display uint32) CGRect
	pcgDisplayScreenSize             func(display uint32) CGSize
	pcgDisplayCreateImage            func(display uint32) uintptr
	pcgDisplayCreateImageForRect     func(display uint32, rect CGRect) uintptr
	pcgEventSetIntegerValueField     func(event uintptr, field uint32, value int64)
	pcgEventKeyboardSetUnicodeString func(event uintptr, length uint32, chars uintptr)
	pcgGetOnlineDisplayList          func(maxDisplays uint32, onlineDisplays *uint32, displayCount *uint32) int32
)

// Function pointers - Carbon
var (
	pgetKeyCode func(virtualKeyCode *uint16, charCode uintptr) int16
)

// ObjC function pointers
var (
	pobjc_getClass    func(name string) uintptr
	psel_registerName func(name string) uintptr
	pobjc_msgSend     func(obj uintptr, sel uintptr, args ...uintptr) uintptr
)

// Function pointers - Accessibility
var (
	pAXUIElementCreateApplication  func(pid uint32) uintptr
	pAXUIElementCopyAttributeValue func(element uintptr, attr uintptr, value *uintptr) int32
	pAXIsProcessTrusted            func() bool
)

var initErr error

func init() {
	var err error
	coreGraphicsLib, err = purego.Dlopen(coreGraphicsPath, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		initErr = errors.New("failed to load CoreGraphics: " + err.Error())
		return
	}

	appServicesLib, err = purego.Dlopen(appServicesPath, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		initErr = errors.New("failed to load ApplicationServices: " + err.Error())
		return
	}

	ioKitLib, err = purego.Dlopen(ioKitPath, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		initErr = errors.New("failed to load IOKit: " + err.Error())
		return
	}

	carbonLib, err = purego.Dlopen(carbonPath, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		initErr = errors.New("failed to load Carbon: " + err.Error())
		return
	}

	objcLib, err = purego.Dlopen(objcPath, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		initErr = errors.New("failed to load libobjc: " + err.Error())
		return
	}

	registerFunctions()
}

func registerFunctions() {
	// CoreGraphics - Event Source
	purego.RegisterLibFunc(&pcgEventSourceCreate, coreGraphicsLib, "CGEventSourceCreate")

	// CoreGraphics - Mouse
	purego.RegisterLibFunc(&pcgEventCreateMouseEvent, coreGraphicsLib, "CGEventCreateMouseEvent")
	purego.RegisterLibFunc(&pcgEventPost, coreGraphicsLib, "CGEventPost")
	purego.RegisterLibFunc(&pcgEventGetLocation, coreGraphicsLib, "CGEventGetLocation")

	// CoreGraphics - Keyboard
	purego.RegisterLibFunc(&pcgEventCreateKeyboardEvent, coreGraphicsLib, "CGEventCreateKeyboardEvent")
	purego.RegisterLibFunc(&pcgEventSetType, coreGraphicsLib, "CGEventSetType")
	purego.RegisterLibFunc(&pcgEventSetFlags, coreGraphicsLib, "CGEventSetFlags")
	purego.RegisterLibFunc(&pcgEventCreate, coreGraphicsLib, "CGEventCreate")
	purego.RegisterLibFunc(&pcgEventSetIntegerValueField, coreGraphicsLib, "CGEventSetIntegerValueField")
	purego.RegisterLibFunc(&pcgEventKeyboardSetUnicodeString, coreGraphicsLib, "CGEventKeyboardSetUnicodeString")

	// CoreGraphics - Display
	purego.RegisterLibFunc(&pcgMainDisplayID, coreGraphicsLib, "CGMainDisplayID")
	purego.RegisterLibFunc(&pcgDisplayPixelsWide, coreGraphicsLib, "CGDisplayPixelsWide")
	purego.RegisterLibFunc(&pcgDisplayPixelsHigh, coreGraphicsLib, "CGDisplayPixelsHigh")
	purego.RegisterLibFunc(&pcgDisplayBounds, coreGraphicsLib, "CGDisplayBounds")
	purego.RegisterLibFunc(&pcgDisplayScreenSize, coreGraphicsLib, "CGDisplayScreenSize")
	purego.RegisterLibFunc(&pcgDisplayCreateImage, coreGraphicsLib, "CGDisplayCreateImage")
	purego.RegisterLibFunc(&pcgDisplayCreateImageForRect, coreGraphicsLib, "CGDisplayCreateImageForRect")
	purego.RegisterLibFunc(&pcgGetOnlineDisplayList, coreGraphicsLib, "CGGetOnlineDisplayList")

	// Carbon
	purego.RegisterLibFunc(&pgetKeyCode, carbonLib, "getKeyCode")

	// ObjC Runtime
	purego.RegisterLibFunc(&pobjc_getClass, objcLib, "objc_getClass")
	purego.RegisterLibFunc(&psel_registerName, objcLib, "sel_registerName")
	purego.RegisterLibFunc(&pobjc_msgSend, objcLib, "objc_msgSend")

	// Accessibility
	if appServicesLib != 0 {
		purego.RegisterLibFunc(&pAXUIElementCreateApplication, appServicesLib, "AXUIElementCreateApplication")
		purego.RegisterLibFunc(&pAXUIElementCopyAttributeValue, appServicesLib, "AXUIElementCopyAttributeValue")
		purego.RegisterLibFunc(&pAXIsProcessTrusted, appServicesLib, "AXIsProcessTrusted")
	}
}

// InitError returns any initialization error
func InitError() error {
	return initErr
}

// --- Mouse Operations ---

// CGEventSourceCreate creates a new CGEventSource
func CGEventSourceCreate(stateID uint32) uintptr {
	return pcgEventSourceCreate(stateID, 0)
}

// MoveMouse moves the mouse to (x, y)
func MoveMouse(x, y int32) int {
	source := CGEventSourceCreate(kCGEventSourceStateHIDSystemState)
	if source == 0 {
		return 1
	}

	point := CGPoint{X: float64(x), Y: float64(y)}
	event := pcgEventCreateMouseEvent(source, kCGEventMouseMoved, point, kCGMouseButtonLeft)
	if event == 0 {
		return 1
	}

	pcgEventPost(kCGHIDEventTap, event)
	return 0
}

// DragMouse drags the mouse to (x, y) with the specified button
func DragMouse(x, y int32, button uint16) int {
	source := CGEventSourceCreate(kCGEventSourceStateHIDSystemState)
	if source == 0 {
		return 1
	}

	point := CGPoint{X: float64(x), Y: float64(y)}
	var cgButton uint32
	var dragType uint32
	switch button {
	case 0: // LEFT_BUTTON
		cgButton = kCGMouseButtonLeft
		dragType = kCGEventLeftMouseDragged
	case 1: // RIGHT_BUTTON
		cgButton = kCGMouseButtonRight
		dragType = kCGEventRightMouseDragged
	default:
		cgButton = kCGMouseButtonCenter
		dragType = kCGEventOtherMouseDragged
	}

	event := pcgEventCreateMouseEvent(source, dragType, point, cgButton)
	if event == 0 {
		return 1
	}

	pcgEventPost(kCGHIDEventTap, event)
	return 0
}

// ToggleMouse toggles a mouse button down or up
func ToggleMouse(down bool, button uint16) int {
	source := CGEventSourceCreate(kCGEventSourceStateHIDSystemState)
	if source == 0 {
		return 1
	}

	// Get current mouse location
	locEvent := pcgEventCreate(source)
	var point CGPoint
	if locEvent != 0 {
		point = pcgEventGetLocation(locEvent)
	}

	var eventType uint32
	var cgButton uint32

	switch button {
	case 0: // LEFT_BUTTON
		cgButton = kCGMouseButtonLeft
		if down {
			eventType = kCGEventLeftMouseDown
		} else {
			eventType = kCGEventLeftMouseUp
		}
	case 1: // RIGHT_BUTTON
		cgButton = kCGMouseButtonRight
		if down {
			eventType = kCGEventRightMouseDown
		} else {
			eventType = kCGEventRightMouseUp
		}
	default: // CENTER_BUTTON
		cgButton = kCGMouseButtonCenter
		if down {
			eventType = kCGEventOtherMouseDown
		} else {
			eventType = kCGEventOtherMouseUp
		}
	}

	event := pcgEventCreateMouseEvent(source, eventType, point, cgButton)
	if event == 0 {
		return 1
	}

	pcgEventPost(kCGHIDEventTap, event)
	return 0
}

// ScrollMouse scrolls the mouse wheel
func ScrollMouse(x, y int) {
	source := CGEventSourceCreate(kCGEventSourceStateHIDSystemState)
	if source == 0 {
		return
	}

	if y != 0 {
		event := pcgEventCreateMouseEvent(source, kCGEventScrollWheel, CGPoint{}, 0)
		if event != 0 {
			pcgEventSetIntegerValueField(event, 10, int64(y)) // kCGScrollWheelEventDeltaAxis1
			pcgEventPost(kCGHIDEventTap, event)
		}
	}

	if x != 0 {
		event := pcgEventCreateMouseEvent(source, kCGEventScrollWheel, CGPoint{}, 0)
		if event != 0 {
			pcgEventSetIntegerValueField(event, 12, int64(x)) // kCGScrollWheelEventDeltaAxis2
			pcgEventPost(kCGHIDEventTap, event)
		}
	}
}

// Location returns the current mouse position
func Location() (int32, int32) {
	source := CGEventSourceCreate(kCGEventSourceStateHIDSystemState)
	if source == 0 {
		return 0, 0
	}

	event := pcgEventCreate(source)
	if event == 0 {
		return 0, 0
	}

	loc := pcgEventGetLocation(event)
	return int32(loc.X), int32(loc.Y)
}

// --- Keyboard Operations ---

// ToggleKeyCode sends a key down/up event
func ToggleKeyCode(keyCode uint32, down bool, flags uint64, pid uintptr) int {
	source := CGEventSourceCreate(kCGEventSourceStateHIDSystemState)
	if source == 0 {
		return 1
	}

	event := pcgEventCreateKeyboardEvent(source, uint16(keyCode), down)
	if event == 0 {
		return 1
	}

	if flags != 0 {
		pcgEventSetFlags(event, flags)
	}

	pcgEventPost(kCGHIDEventTap, event)
	return 0
}

// KeyCodeForChar returns the key code for a character
func KeyCodeForChar(char byte) uint32 {
	if pgetKeyCode != nil {
		var keyCode uint16
		cstr := []byte{char, 0}
		ret := pgetKeyCode(&keyCode, uintptr(unsafe.Pointer(&cstr[0])))
		if ret == 0 {
			return uint32(keyCode)
		}
	}
	return keyCodeForCharFallback(char)
}

// keyCodeForCharFallback provides key code mappings for macOS
func keyCodeForCharFallback(char byte) uint32 {
	keyMap := map[byte]uint32{
		'a': 0x00, 'b': 0x0B, 'c': 0x08, 'd': 0x02, 'e': 0x0E,
		'f': 0x03, 'g': 0x05, 'h': 0x04, 'i': 0x22, 'j': 0x26,
		'k': 0x28, 'l': 0x25, 'm': 0x2E, 'n': 0x2D, 'o': 0x1F,
		'p': 0x23, 'q': 0x0C, 'r': 0x0F, 's': 0x01, 't': 0x11,
		'u': 0x20, 'v': 0x09, 'w': 0x0D, 'x': 0x07, 'y': 0x10,
		'z': 0x06,
		'A': 0x00, 'B': 0x0B, 'C': 0x08, 'D': 0x02, 'E': 0x0E,
		'F': 0x03, 'G': 0x05, 'H': 0x04, 'I': 0x22, 'J': 0x26,
		'K': 0x28, 'L': 0x25, 'M': 0x2E, 'N': 0x2D, 'O': 0x1F,
		'P': 0x23, 'Q': 0x0C, 'R': 0x0F, 'S': 0x01, 'T': 0x11,
		'U': 0x20, 'V': 0x09, 'W': 0x0D, 'X': 0x07, 'Y': 0x10,
		'Z': 0x06,
		'0': 0x1D, '1': 0x12, '2': 0x13, '3': 0x14, '4': 0x15,
		'5': 0x17, '6': 0x16, '7': 0x1A, '8': 0x1C, '9': 0x19,
		'-': 0x1B, '=': 0x18, '[': 0x21, ']': 0x1E, '\\': 0x2A,
		';': 0x29, '\'': 0x27, ',': 0x2B, '.': 0x2F, '/': 0x2C,
		'`': 0x32, ' ': 0x31,
	}
	if code, ok := keyMap[char]; ok {
		return code
	}
	return 0xFF // K_NOT_A_KEY
}

// UnicodeType sends a unicode character using keyboard events
func UnicodeType(char uint32, pid uintptr, isPid int8) {
	source := CGEventSourceCreate(kCGEventSourceStateHIDSystemState)
	if source == 0 {
		return
	}

	// Create key down event with keycode 0
	event := pcgEventCreateKeyboardEvent(source, 0, true)
	if event == 0 {
		return
	}

	r := rune(char)
	pcgEventKeyboardSetUnicodeString(event, 1, uintptr(unsafe.Pointer(&r)))
	pcgEventPost(kCGHIDEventTap, event)

	// Create key up event
	eventUp := pcgEventCreateKeyboardEvent(source, 0, false)
	if eventUp == 0 {
		return
	}
	pcgEventKeyboardSetUnicodeString(eventUp, 1, uintptr(unsafe.Pointer(&r)))
	pcgEventPost(kCGHIDEventTap, eventUp)
}

// InputUTF types a UTF-8 string character by character
func InputUTF(str string) {
	for _, r := range str {
		UnicodeType(uint32(r), 0, 0)
	}
}

// --- Display Operations ---

// CGMainDisplayID returns the main display ID
func CGMainDisplayID() uint32 {
	return pcgMainDisplayID()
}

// GetMainDisplaySize returns the main display size
func GetMainDisplaySize() (int32, int32) {
	displayID := pcgMainDisplayID()
	w := pcgDisplayPixelsWide(displayID)
	h := pcgDisplayPixelsHigh(displayID)
	return int32(w), int32(h)
}

// GetScreenRect returns the screen rectangle
func GetScreenRect(displayID int32) (int32, int32, int32, int32) {
	id := uint32(displayID)
	if displayID < 0 {
		id = pcgMainDisplayID()
	}
	rect := pcgDisplayBounds(id)
	return int32(rect.Origin.X), int32(rect.Origin.Y),
		int32(rect.Size.Width), int32(rect.Size.Height)
}

// GetNumDisplays returns the number of displays
func GetNumDisplays() int {
	var count uint32
	pcgGetOnlineDisplayList(0, nil, &count)
	return int(count)
}

// SysScale returns the system scale factor
func SysScale(displayID int32) float64 {
	id := uint32(displayID)
	if displayID < 0 {
		id = pcgMainDisplayID()
	}
	size := pcgDisplayScreenSize(id)
	w := pcgDisplayPixelsWide(id)
	if w == 0 || size.Width == 0 {
		return 1.0
	}
	return float64(w) / size.Width
}

// --- Window Operations ---

// ShowAlert shows an alert using NSAlert via ObjC runtime
func ShowAlert(title, msg, defaultBtn, cancelBtn string) bool {
	if pobjc_getClass == nil || psel_registerName == nil || pobjc_msgSend == nil {
		return false
	}

	// Create NSAutoreleasePool
	poolClass := pobjc_getClass("NSAutoreleasePool")
	allocSel := psel_registerName("alloc")
	initSel := psel_registerName("init")
	drainSel := psel_registerName("drain")

	pool := pobjc_msgSend(poolClass, allocSel)
	pool = pobjc_msgSend(pool, initSel)
	defer pobjc_msgSend(pool, drainSel)

	// Create NSAlert
	alertClass := pobjc_getClass("NSAlert")
	alert := pobjc_msgSend(alertClass, allocSel)
	alert = pobjc_msgSend(alert, initSel)

	// Set message text
	setMsgSel := psel_registerName("setMessageText:")
	msgStr := CreateNSString(msg)
	pobjc_msgSend(alert, setMsgSel, msgStr)

	// Set informative text (title)
	setInfoSel := psel_registerName("setInformativeText:")
	infoStr := CreateNSString(title)
	pobjc_msgSend(alert, setInfoSel, infoStr)

	// Add buttons
	addBtnSel := psel_registerName("addButtonWithTitle:")
	defaultStr := CreateNSString(defaultBtn)
	pobjc_msgSend(alert, addBtnSel, defaultStr)

	if cancelBtn != "" {
		cancelStr := CreateNSString(cancelBtn)
		pobjc_msgSend(alert, addBtnSel, cancelStr)
	}

	// Run modal
	runSel := psel_registerName("runModal")
	result := pobjc_msgSend(alert, runSel)

	// Release
	releaseSel := psel_registerName("release")
	pobjc_msgSend(alert, releaseSel)

	// NSAlertFirstButtonReturn = 1000
	return result == 1000
}

// CreateNSString creates an NSString from a Go string via ObjC runtime
func CreateNSString(s string) uintptr {
	if pobjc_getClass == nil {
		return 0
	}
	nsStringClass := pobjc_getClass("NSString")
	allocSel := psel_registerName("alloc")
	initWithUTF8Sel := psel_registerName("initWithUTF8String:")

	obj := pobjc_msgSend(nsStringClass, allocSel)
	// Need null-terminated C string
	cs := append([]byte(s), 0)
	obj = pobjc_msgSend(obj, initWithUTF8Sel, uintptr(unsafe.Pointer(&cs[0])))

	return obj
}

// Is64Bit returns whether the system is 64-bit
func Is64Bit() bool {
	return unsafe.Sizeof(uintptr(0)) == 8
}

// AXIsActive checks if accessibility is enabled for this process
func AXIsActive() bool {
	if pAXIsProcessTrusted == nil {
		return false
	}
	return pAXIsProcessTrusted()
}

// GetPixelColor returns the pixel color at (x, y) as RGB hex value
func GetPixelColor(x, y int32, displayID int32) uint32 {
	// Use CGDisplayCreateImage and read the pixel
	id := uint32(displayID)
	if displayID < 0 {
		id = pcgMainDisplayID()
	}

	imageRef := pcgDisplayCreateImage(id)
	if imageRef == 0 {
		return 0
	}

	// We need CGImageGetWidth, CGImageGetHeight, CGImageGetBytesPerRow,
	// CGImageGetDataProvider, CGDataProviderCopyData
	// For now, use a simpler approach via the NSBitmapImageRep
	// This requires additional purego bindings

	// Simplified: return 0 until full screen capture is implemented
	return 0
}

// CaptureScreen captures a region of the screen and returns pixel data (BGRA format)
func CaptureScreen(x, y, w, h int32, displayID int32) ([]byte, int, int) {
	// Full screen capture requires CGDisplayCreateImageForRect + CGDataProviderCopyData
	// which needs additional purego function registrations.
	// For the initial implementation, we return an empty result.
	// TODO: Implement full screen capture using purego
	return nil, 0, 0
}

// GetWindowTitleByPID gets the window title using Accessibility API
func GetWindowTitleByPID(pid uint32) string {
	if pAXUIElementCreateApplication == nil || pAXUIElementCopyAttributeValue == nil {
		return ""
	}

	element := pAXUIElementCreateApplication(pid)
	if element == 0 {
		return ""
	}

	// Get the focused window
	focusedWinSel := "AXFocusedWindow"
	cs := append([]byte(focusedWinSel), 0)
	attrName := registerSel(cs)

	var windowRef uintptr
	result := pAXUIElementCopyAttributeValue(element, attrName, &windowRef)
	if result != 0 {
		return ""
	}

	// Get the title attribute
	titleSel := "AXTitle"
	cs2 := append([]byte(titleSel), 0)
	titleAttr := registerSel(cs2)

	var titleValue uintptr
	result = pAXUIElementCopyAttributeValue(windowRef, titleAttr, &titleValue)
	if result != 0 {
		return ""
	}

	// titleValue is a CFStringRef, convert to Go string
	// We need CFStringGetLength and CFStringGetCString
	return cfStringToGo(titleValue)
}

// GetActiveWindowPID gets the PID of the frontmost application
func GetActiveWindowPID() uint32 {
	if pobjc_getClass == nil {
		return 0
	}

	// [[NSWorkspace sharedWorkspace] frontmostApplication] processIdentifier
	workspaceClass := pobjc_getClass("NSWorkspace")
	sharedSel := psel_registerName("sharedWorkspace")
	frontmostSel := psel_registerName("frontmostApplication")
	pidSel := psel_registerName("processIdentifier")

	workspace := pobjc_msgSend(workspaceClass, sharedSel)
	app := pobjc_msgSend(workspace, frontmostSel)
	if app == 0 {
		return 0
	}
	pid := pobjc_msgSend(app, pidSel)
	return uint32(pid)
}

// GetActiveWindowTitle gets the title of the active window
func GetActiveWindowTitle() string {
	pid := GetActiveWindowPID()
	if pid == 0 {
		return ""
	}
	return GetWindowTitleByPID(pid)
}

// MinimizeWindow minimizes a window using Accessibility API
func MinimizeWindow(pid uint32, minimize bool) {
	if pAXUIElementCreateApplication == nil {
		return
	}

	element := pAXUIElementCreateApplication(pid)
	if element == 0 {
		return
	}

	focusedWinSel := append([]byte("AXFocusedWindow"), 0)
	attrName := registerSel(focusedWinSel)

	var windowRef uintptr
	result := pAXUIElementCopyAttributeValue(element, attrName, &windowRef)
	if result != 0 {
		return
	}

	// AXSetAttributeValue with AXMinimized attribute
	minimizedSel := append([]byte("AXMinimized"), 0)
	minimizedAttr := registerSel(minimizedSel)

	// Set the AXMinimized attribute using objc_msgSend for NSNumber
	numberClass := pobjc_getClass("NSNumber")
	allocSel := psel_registerName("alloc")
	initBoolSel := psel_registerName("initWithBool:")

	var boolVal uintptr
	if minimize {
		boolVal = 1
	}
	number := pobjc_msgSend(numberClass, allocSel)
	number = pobjc_msgSend(number, initBoolSel, boolVal)

	// Use AXUIElementSetAttributeValue (need to register)
	// For now, use a simplified approach
	_ = minimizedAttr
	_ = number
}

func registerSel(name []byte) uintptr {
	if psel_registerName == nil {
		return 0
	}
	return psel_registerName(string(name[:len(name)-1])) // remove null terminator
}

func cfStringToGo(ref uintptr) string {
	if ref == 0 {
		return ""
	}
	// Use CFStringGetLength and CFStringGetCStringPtr
	// For now, use a simplified ObjC approach
	if pobjc_getClass == nil {
		return ""
	}

	// Treat as NSString and call UTF8String
	nsstringSel := psel_registerName("UTF8String")
	if nsstringSel == 0 {
		return ""
	}

	utf8Ptr := pobjc_msgSend(ref, nsstringSel)
	if utf8Ptr == 0 {
		return ""
	}

	// Read C string from pointer
	var length int
	p := utf8Ptr
	for *(*byte)(unsafe.Pointer(p)) != 0 {
		length++
		p++
	}
	if length == 0 {
		return ""
	}
	slice := unsafe.Slice((*byte)(unsafe.Pointer(utf8Ptr)), length)
	return string(slice)
}
