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

/*
Package robotgo Go native cross-platform system automation.

No C compiler is required — all platform APIs are called via purego/syscall.

Installation:

With Go module support (Go 1.11+), just import:

	import "github.com/shaolei/robotgo"
*/
package robotgo

import (
	"errors"
	"fmt"
	"image"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"github.com/vcaesar/tt"
)

const (
	// Version get the robotgo version
	Version = "v1.00.0.1189, MT. Baker!"
)

// GetVersion get the robotgo version
func GetVersion() string {
	return Version
}

var (
	// MouseSleep set the mouse default millisecond sleep time
	MouseSleep = 0
	// KeySleep set the key default millisecond sleep time
	KeySleep = 10

	// DisplayID set the screen display id
	DisplayID = -1

	// NotPid used the hwnd not pid in windows
	NotPid bool
	// Scale option the os screen scale
	Scale bool
)

type (
	// Map a map[string]interface{}
	Map map[string]interface{}
	// CHex define CHex as rgb Hex type (uint32)
	CHex uint32
	// CBitmap define CBitmap as MMBitmap pointer type
	CBitmap *MMBitmap
	// Handle define window Handle as MData type
	Handle MData
)

// Bitmap define the go Bitmap struct
//
// The common type conversion of bitmap:
//
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#type-conversion
type Bitmap struct {
	ImgBuf        *uint8
	Width, Height int

	Bytewidth     int
	BitsPixel     uint8
	BytesPerPixel uint8
}

// Point is point struct
type Point struct {
	X int
	Y int
}

// Size is size structure
type Size struct {
	W, H int
}

// Rect is rect structure
type Rect struct {
	Point
	Size
}

// Try handler(err)
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// MilliSleep sleep tm milli second
func MilliSleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Millisecond)
}

// Sleep time.Sleep tm second
func Sleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Second)
}

// Deprecated: use the MilliSleep(),
//
// MicroSleep time microsecond sleep
func MicroSleep(tm float64) {
	time.Sleep(time.Duration(tm * float64(time.Microsecond)))
}

/*
      _______.  ______ .______       _______  _______ .__   __.
    /       | /      ||   _  \     |   ____||   ____||  \ |  |
   |   (----`|  ,----'|  |_)  |    |  |__   |  |__   |   \|  |
    \   \    |  |     |      /     |   __|  |   __|  |  . `  |
.----)   |   |  `----.|  |\  \----.|  |____ |  |____ |  |\   |
|_______/     \______|| _| `._____||_______||_______||__| \__|
*/

// ToMMRGBHex trans CHex to uint32
func ToMMRGBHex(hex CHex) uint32 {
	return uint32(hex)
}

// UintToHex trans uint32 to robotgo.CHex
func UintToHex(u uint32) CHex {
	return CHex(u)
}

// U32ToHex trans uint32 to CHex
func U32ToHex(hex uint32) CHex {
	return CHex(hex)
}

// U8ToHex trans *uint8 to CHex
func U8ToHex(hex *uint8) CHex {
	return CHex(*hex)
}

// PadHex trans CHex to string
func PadHex(hex CHex) string {
	return padHex(uint32(hex))
}

// PadHexs trans CHex to string
func PadHexs(hex CHex) string {
	return PadHex(hex)
}

// HexToRgb trans hex to rgb
func HexToRgb(hex uint32) (r, g, b uint8) {
	return colorHexToRGB(hex)
}

// RgbToHex trans rgb to hex
func RgbToHex(r, g, b uint8) uint32 {
	return colorRGBToHex(r, g, b)
}

// GetPxColor get the pixel color return CHex
func GetPxColor(x, y int, displayId ...int) CHex {
	display := displayIdx(displayId...)
	return CHex(platformGetPxColor(int32(x), int32(y), int32(display)))
}

// GetPixelColor get the pixel color return string
func GetPixelColor(x, y int, displayId ...int) string {
	return PadHex(GetPxColor(x, y, displayId...))
}

// GetLocationColor get the location pos's color
func GetLocationColor(displayId ...int) string {
	x, y := Location()
	return GetPixelColor(x, y, displayId...)
}

// IsMain is main display
func IsMain(displayId int) bool {
	return displayId == GetMainId()
}

func displayIdx(id ...int) int {
	display := -1
	if DisplayID != -1 {
		display = DisplayID
	}
	if len(id) > 0 {
		display = id[0]
	}
	return display
}

func getNumDisplays() int {
	return platformGetNumDisplays()
}

// GetHWNDByPid get the hwnd by pid
func GetHWNDByPid(pid int) int {
	return int(platformGetHwndByPid(uintptr(pid)))
}

// SysScale get the sys scale
func SysScale(displayId ...int) float64 {
	display := displayIdx(displayId...)
	return platformSysScale(int32(display))
}

// Scaled get the screen scaled return scale size
func Scaled(x int, displayId ...int) int {
	f := ScaleF(displayId...)
	return Scaled0(x, f)
}

// Scaled0 return int(x * f)
func Scaled0(x int, f float64) int {
	return int(float64(x) * f)
}

// Scaled1 return int(x / f)
func Scaled1(x int, f float64) int {
	return int(float64(x) / f)
}

// GetScreenSize get the screen size
func GetScreenSize() (int, int) {
	w, h := platformGetMainDisplaySize()
	return int(w), int(h)
}

// GetScreenRect get the screen rect (x, y, w, h)
func GetScreenRect(displayId ...int) Rect {
	display := -1
	if len(displayId) > 0 {
		display = displayId[0]
	}

	x, y, w, h := platformGetScreenRect(int32(display))

	if runtime.GOOS == "windows" {
		f := ScaleF()
		x, y, w, h = int32(Scaled0(int(x), f)), int32(Scaled0(int(y), f)),
			int32(Scaled0(int(w), f)), int32(Scaled0(int(h), f))
	}
	return Rect{
		Point{X: int(x), Y: int(y)},
		Size{W: int(w), H: int(h)},
	}
}

// GetScaleSize get the screen scale size
func GetScaleSize(displayId ...int) (int, int) {
	x, y := GetScreenSize()
	f := ScaleF(displayId...)
	return int(float64(x) * f), int(float64(y) * f)
}

// CaptureScreen capture the screen return bitmap,
// use `defer robotgo.FreeBitmap(bitmap)` to free the bitmap
//
// robotgo.CaptureScreen(x, y, w, h int)
func CaptureScreen(args ...int) CBitmap {
	var x, y, w, h int32
	displayId := -1
	if DisplayID != -1 {
		displayId = DisplayID
	}

	if len(args) > 4 {
		displayId = args[4]
	}

	if len(args) > 3 {
		x = int32(args[0])
		y = int32(args[1])
		w = int32(args[2])
		h = int32(args[3])
	} else {
		rect := GetScreenRect(displayId)
		if runtime.GOOS == "windows" {
			x = int32(rect.X)
			y = int32(rect.Y)
		}
		w = int32(rect.W)
		h = int32(rect.H)
	}

	isPid := 0
	if NotPid || len(args) > 5 {
		isPid = 1
	}

	return platformCaptureScreen(x, y, w, h, int32(displayId), int8(isPid))
}

// CaptureGo capture the screen and return bitmap(go struct)
func CaptureGo(args ...int) Bitmap {
	bit := CaptureScreen(args...)
	defer FreeBitmap(bit)
	return ToBitmap(bit)
}

// CaptureImg capture the screen and return image.Image, error
func CaptureImg(args ...int) (image.Image, error) {
	bit := CaptureScreen(args...)
	if bit == nil {
		return nil, errors.New("Capture image not found.")
	}
	defer FreeBitmap(bit)
	return ToImage(bit), nil
}

// FreeBitmap free and dealloc the bitmap
func FreeBitmap(bitmap CBitmap) {
	platformBitmapDealloc(bitmap)
}

// FreeBitmapArr free and dealloc the bitmap array
func FreeBitmapArr(bit ...CBitmap) {
	for i := 0; i < len(bit); i++ {
		FreeBitmap(bit[i])
	}
}

// ToMMBitmapRef trans CBitmap to *MMBitmap (identity)
func ToMMBitmapRef(bit CBitmap) *MMBitmap {
	return bit
}

// ToBitmap trans CBitmap to Bitmap
func ToBitmap(bit CBitmap) Bitmap {
	if bit == nil {
		return Bitmap{}
	}
	bitmap := Bitmap{
		ImgBuf:        bit.ImageBuffer,
		Width:         int(bit.Width),
		Height:        int(bit.Height),
		Bytewidth:     int(bit.Bytewidth),
		BitsPixel:     uint8(bit.BitsPerPixel),
		BytesPerPixel: uint8(bit.BytesPerPixel),
	}
	return bitmap
}

// ToCBitmap trans Bitmap to CBitmap
func ToCBitmap(bit Bitmap) CBitmap {
	return platformCreateMMBitmap(
		bit.ImgBuf,
		int32(bit.Width),
		int32(bit.Height),
		int32(bit.Bytewidth),
		uint8(bit.BitsPixel),
		uint8(bit.BytesPerPixel),
	)
}

// ToImage convert CBitmap to standard image.Image
func ToImage(bit CBitmap) image.Image {
	return ToRGBA(bit)
}

// ToRGBA convert CBitmap to standard image.RGBA
func ToRGBA(bit CBitmap) *image.RGBA {
	bmp1 := ToBitmap(bit)
	return ToRGBAGo(bmp1)
}

// ImgToCBitmap trans image.Image to CBitmap
func ImgToCBitmap(img image.Image) CBitmap {
	return ToCBitmap(ImgToBitmap(img))
}

// ByteToCBitmap trans []byte to CBitmap
func ByteToCBitmap(by []byte) CBitmap {
	img, _ := ByteToImg(by)
	return ImgToCBitmap(img)
}

// SetXDisplayName set XDisplay name (Linux)
func SetXDisplayName(name string) error {
	return platformSetXDisplayName(name)
}

// GetXDisplayName get XDisplay name (Linux)
func GetXDisplayName() string {
	return platformGetXDisplayName()
}

// CloseMainDisplay close the main X11 display
func CloseMainDisplay() {
	platformCloseMainDisplay()
}

// Deprecated: use the ScaledF(),
//
// ScaleX get the primary display horizontal DPI scale factor, drop
func ScaleX() int {
	return platformScaleX()
}

/*
.___  ___.   ______    __    __       _______. _______
|   \/   |  /  __  \  |  |  |  |     /       ||   ____|
|  \  /  | |  |  |  | |  |  |  |    |   (----`|  |__
|  |\/|  | |  |  |  | |  |  |  |     \   \    |   __|
|  |  |  | |  `--'  | |  `--'  | .----)   |   |  |____
|__|  |__|  \______/   \______/  |_______/    |_______|

*/

// CheckMouse check the mouse button
func CheckMouse(btn string) uint16 {
	m1 := map[string]uint16{
		"left":       LEFT_BUTTON,
		"center":     CENTER_BUTTON,
		"right":      RIGHT_BUTTON,
		"wheelDown":  WHEEL_DOWN,
		"wheelUp":    WHEEL_UP,
		"wheelLeft":  WHEEL_LEFT,
		"wheelRight": WHEEL_RIGHT,
	}
	if v, ok := m1[btn]; ok {
		return v
	}
	return LEFT_BUTTON
}

// MouseButtonString converts a mouse button to a readable name.
func MouseButtonString(btn uint16) string {
	m1 := map[uint16]string{
		LEFT_BUTTON:   "left",
		CENTER_BUTTON: "center",
		RIGHT_BUTTON:  "right",
		WHEEL_DOWN:    "wheelDown",
		WHEEL_UP:      "wheelUp",
		WHEEL_LEFT:    "wheelLeft",
		WHEEL_RIGHT:   "wheelRight",
	}
	if v, ok := m1[btn]; ok {
		return v
	}
	return fmt.Sprintf("button%d", btn)
}

// MoveScale calculate the os scale factor x, y
func MoveScale(x, y int, displayId ...int) (int, int) {
	if Scale || runtime.GOOS == "windows" {
		f := ScaleF()
		x, y = Scaled1(x, f), Scaled1(y, f)
	}
	return x, y
}

// Move move the mouse to (x, y)
//
// Examples:
//
//	robotgo.MouseSleep = 100  // 100 millisecond
//	robotgo.Move(10, 10)
func Move(x, y int, displayId ...int) {
	x, y = MoveScale(x, y, displayId...)
	platformMoveMouse(int32(x), int32(y))
	MilliSleep(MouseSleep)
}

// Deprecated: use the DragSmooth(),
//
// Drag drag the mouse to (x, y),
// It's not valid now, use the DragSmooth()
func Drag(x, y int, args ...string) {
	x, y = MoveScale(x, y)
	button := LEFT_BUTTON
	if len(args) > 0 {
		button = CheckMouse(args[0])
	}
	platformDragMouse(int32(x), int32(y), button)
	MilliSleep(MouseSleep)
}

// DragSmooth drag the mouse like smooth to (x, y)
//
// Examples:
//
//	robotgo.DragSmooth(10, 10)
func DragSmooth(x, y int, args ...interface{}) {
	Toggle("left")
	MilliSleep(50)
	smoothMove(x, y, true, args...)
	Toggle("left", "up")
}

func smoothMove(x, y int, drag bool, args ...interface{}) bool {
	x, y = MoveScale(x, y)

	var (
		mouseDelay = 1
		low        = 1.0
		high       = 3.0
	)

	if len(args) > 2 {
		mouseDelay = args[2].(int)
	}
	if len(args) > 1 {
		low = args[0].(float64)
		high = args[1].(float64)
	}

	var result bool
	if drag {
		result = platformSmoothlyDragMouse(int32(x), int32(y), low, high, LEFT_BUTTON)
	} else {
		result = platformSmoothlyMoveMouse(int32(x), int32(y), low, high)
	}
	MilliSleep(MouseSleep + mouseDelay)
	return result
}

// MoveSmooth move the mouse smooth,
// moves mouse to x, y human like, with the mouse button up.
//
// robotgo.MoveSmooth(x, y int, low, high float64, mouseDelay int)
//
// Examples:
//
//	robotgo.MoveSmooth(10, 10)
//	robotgo.MoveSmooth(10, 10, 1.0, 2.0)
func MoveSmooth(x, y int, args ...interface{}) bool {
	return smoothMove(x, y, false, args...)
}

// MoveArgs get the mouse relative args
func MoveArgs(x, y int) (int, int) {
	mx, my := Location()
	mx = mx + x
	my = my + y
	return mx, my
}

// MoveRelative move mouse with relative
func MoveRelative(x, y int) {
	Move(MoveArgs(x, y))
}

// MoveSmoothRelative move mouse smooth with relative
func MoveSmoothRelative(x, y int, args ...interface{}) {
	mx, my := MoveArgs(x, y)
	MoveSmooth(mx, my, args...)
}

// Location get the mouse location position return x, y
func Location() (int, int) {
	x, y := platformLocation()
	if Scale || runtime.GOOS == "windows" {
		f := ScaleF()
		x, y = int32(Scaled0(int(x), f)), int32(Scaled0(int(y), f))
	}
	return int(x), int(y)
}

// ClickV1 click the mouse button
//
// robotgo.Click(button string, double bool)
//
// Examples:
//
//	robotgo.Click() // default is left button
//	robotgo.Click("right")
//	robotgo.Click("wheelLeft")
func ClickV1(args ...interface{}) {
	var (
		button = LEFT_BUTTON
		double bool
	)
	if len(args) > 0 {
		button = CheckMouse(args[0].(string))
	}
	if len(args) > 1 {
		double = args[1].(bool)
	}

	if !double {
		platformClickMouse(button)
	} else {
		platformDoubleClick(button, 2)
	}
	MilliSleep(MouseSleep)
}

// Click click the mouse button and return error
//
// robotgo.Click(button string, double bool)
//
// Examples:
//
//	err := robotgo.Click() // default is left button
//	err := robotgo.Click("right")
func Click(args ...interface{}) error {
	var (
		button = LEFT_BUTTON
		double bool
		count  int
	)

	if len(args) > 0 {
		btn, ok := args[0].(string)
		if !ok {
			return errors.New("first argument must be a button string")
		}
		button = CheckMouse(btn)
	}

	if len(args) > 1 {
		dbl, ok := args[1].(bool)
		if !ok {
			return errors.New("second argument must be a bool indicating double click")
		}
		double = dbl
	}
	if len(args) > 2 {
		count = args[2].(int)
	}

	defer MilliSleep(MouseSleep)
	if !double {
		if code := platformToggleMouse(true, button); code != 0 {
			return formatClickError(code, button, "down", count)
		}
		MilliSleep(5)
		code := platformToggleMouse(false, button)
		return formatClickError(code, button, "up", count)
	}

	code := platformDoubleClick(button, 2)
	return formatClickError(code, button, "double", 2)
}

// MultiClick performs multiple clicks and returns error
//
// robotgo.MultiClick(button string, count int)
func MultiClick(button string, count int, click ...bool) error {
	if count < 1 {
		return nil
	}
	defer MilliSleep(MouseSleep)

	if runtime.GOOS == "darwin" && len(click) <= 0 {
		btn := CheckMouse(button)
		code := platformDoubleClick(btn, count)
		return formatClickError(code, btn, "down", count)
	}

	for i := 0; i < count; i++ {
		if err := Click(button, false, i+1); err != nil {
			return err
		}
	}
	return nil
}

func formatClickError(code int, button uint16, stage string, count int) error {
	if code == 0 {
		return nil
	}
	btnName := MouseButtonString(button)
	detail := ""

	switch runtime.GOOS {
	case "windows":
		if code != 0 {
			detail = syscall.Errno(code).Error()
		}
	case "darwin":
		cgErrors := map[int]string{
			0:    "kCGErrorSuccess",
			1000: "kCGErrorFailure",
			1001: "kCGErrorIllegalArgument",
			1002: "kCGErrorInvalidConnection",
			1003: "kCGErrorInvalidContext",
			1004: "kCGErrorCannotComplete",
			1005: "kCGErrorNotImplemented",
			1006: "kCGErrorRangeCheck",
			1007: "kCGErrorTypeCheck",
			1008: "kCGErrorNoCurrentPoint",
			1010: "kCGErrorInvalidOperation",
		}
		if v, ok := cgErrors[code]; ok {
			detail = v
		}
	default:
		if code == 1 {
			detail = "XTestFakeButtonEvent returned false"
		}
	}

	if detail != "" {
		return fmt.Errorf("click %s failed (%s, count=%d): %s (code=%d)", stage, btnName, count, detail, code)
	}
	return fmt.Errorf("click %s failed (%s, count=%d), code=%d", stage, btnName, count, code)
}

// MoveClick move and click the mouse
//
// robotgo.MoveClick(x, y int, button string, double bool)
//
// Examples:
//
//	robotgo.MouseSleep = 100
//	robotgo.MoveClick(10, 10)
func MoveClick(x, y int, args ...interface{}) {
	Move(x, y)
	MilliSleep(50)
	Click(args...)
}

// MovesClick move smooth and click the mouse
//
// use the `robotgo.MouseSleep = 100`
func MovesClick(x, y int, args ...interface{}) {
	MoveSmooth(x, y)
	MilliSleep(50)
	Click(args...)
}

// Toggle toggle the mouse, support button:
//
//	"left", "center", "right",
//	"wheelDown", "wheelUp", "wheelLeft", "wheelRight"
//
// Examples:
//
//	robotgo.Toggle("left") // default is down
//	robotgo.Toggle("left", "up")
func Toggle(key ...interface{}) error {
	button := LEFT_BUTTON
	if len(key) > 0 {
		button = CheckMouse(key[0].(string))
	}

	down := true
	if len(key) > 1 && key[1].(string) == "up" {
		down = false
	}

	code := platformToggleMouse(down, button)
	if len(key) > 2 {
		MilliSleep(MouseSleep)
	}
	return formatClickError(code, button, "down", 1)
}

// MouseDown send mouse down event
func MouseDown(key ...interface{}) error {
	return Toggle(key...)
}

// MouseUp send mouse up event
func MouseUp(key ...interface{}) error {
	if len(key) <= 0 {
		key = append(key, "left")
	}
	return Toggle(append(key, "up")...)
}

// Scroll scroll the mouse to (x, y)
//
// robotgo.Scroll(x, y, msDelay int)
//
// Examples:
//
//	robotgo.Scroll(10, 10)
func Scroll(x, y int, args ...int) {
	var msDelay = 10
	if len(args) > 0 {
		msDelay = args[0]
	}
	platformScrollMouse(x, y)
	MilliSleep(MouseSleep + msDelay)
}

// ScrollDir scroll the mouse with direction to (x, "up")
// supported: "up", "down", "left", "right"
//
// Examples:
//
//	robotgo.ScrollDir(10, "down")
//	robotgo.ScrollDir(10, "up")
func ScrollDir(x int, direction ...interface{}) {
	d := "down"
	if len(direction) > 0 {
		d = direction[0].(string)
	}

	if d == "down" {
		Scroll(0, -x)
	}
	if d == "up" {
		Scroll(0, x)
	}
	if d == "left" {
		Scroll(x, 0)
	}
	if d == "right" {
		Scroll(-x, 0)
	}
}

// ScrollSmooth scroll the mouse smooth,
// default scroll 5 times and sleep 100 millisecond
//
// robotgo.ScrollSmooth(toy, num, sleep, tox)
//
// Examples:
//
//	robotgo.ScrollSmooth(-10)
//	robotgo.ScrollSmooth(-10, 6, 200, -10)
func ScrollSmooth(to int, args ...int) {
	i := 0
	num := 5
	if len(args) > 0 {
		num = args[0]
	}
	tm := 100
	if len(args) > 1 {
		tm = args[1]
	}
	tox := 0
	if len(args) > 2 {
		tox = args[2]
	}

	for {
		Scroll(tox, to)
		MilliSleep(tm)
		i++
		if i == num {
			break
		}
	}
	MilliSleep(MouseSleep)
}

// ScrollRelative scroll mouse with relative
//
// Examples:
//
//	robotgo.ScrollRelative(10, 10)
func ScrollRelative(x, y int, args ...int) {
	mx, my := MoveArgs(x, y)
	Scroll(mx, my, args...)
}

/*
____    __    ____  __  .__   __.  _______   ______   ____    __    ____
\   \  /  \  /   / |  | |  \ |  | |       \ /  __  \  \   \  /  \  /   /
 \   \/    \/   /  |  | |   \|  | |  .--.  |  |  |  |  \   \/    \/   /
  \            /   |  | |  . `  | |  |  |  |  |  |  |   \            /
   \    /\    /    |  | |  |\   | |  '--'  |  `--'  |    \    /\    /
    \__/  \__/     |__| |__| \__| |_______/ \______/      \__/  \__/

*/

func alertArgs(args ...string) (string, string) {
	var (
		defaultBtn = "Ok"
		cancelBtn  = "Cancel"
	)
	if len(args) > 0 {
		defaultBtn = args[0]
	}
	if len(args) > 1 {
		cancelBtn = args[1]
	}
	return defaultBtn, cancelBtn
}

func showAlert(title, msg string, args ...string) bool {
	defaultBtn, cancelBtn := alertArgs(args...)
	return platformShowAlert(title, msg, defaultBtn, cancelBtn)
}

// IsValid valid the window
func IsValid() bool {
	return platformIsValid()
}

// SetActive set the window active
func SetActive(win Handle) {
	platformSetActive(MData(win))
}

// SetActiveC set the window active (MData version)
func SetActiveC(win MData) {
	platformSetActive(win)
}

// GetActive get the active window
func GetActive() Handle {
	return Handle(platformGetActive())
}

// GetActiveC get the active window
func GetActiveC() MData {
	return platformGetActive()
}

// MinWindow set the window min
func MinWindow(pid int, args ...interface{}) {
	var (
		state = true
		isPid int
	)
	if len(args) > 0 {
		state = args[0].(bool)
	}
	if len(args) > 1 || NotPid {
		isPid = 1
	}
	platformMinWindow(uintptr(pid), state, int8(isPid))
}

// MaxWindow set the window max
func MaxWindow(pid int, args ...interface{}) {
	var (
		state = true
		isPid int
	)
	if len(args) > 0 {
		state = args[0].(bool)
	}
	if len(args) > 1 || NotPid {
		isPid = 1
	}
	platformMaxWindow(uintptr(pid), state, int8(isPid))
}

// CloseWindow close the window
func CloseWindow(args ...int) {
	if len(args) <= 0 {
		platformCloseMainWindow()
		return
	}

	var pid, isPid int
	if len(args) > 0 {
		pid = args[0]
	}
	if len(args) > 1 || NotPid {
		isPid = 1
	}
	platformCloseWindowByPID(uintptr(pid), int8(isPid))
}

// SetHandle set the window handle
func SetHandle(hwnd int) {
	platformSetHandle(uintptr(hwnd))
}

// SetHandlePid set the window handle by pid
func SetHandlePid(pid int, args ...int) {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
	}
	platformSetHandlePidMData(uintptr(pid), int8(isPid))
}

// GetHandById get handle mdata by id
func GetHandById(id int, args ...int) Handle {
	isPid := 1
	if len(args) > 0 {
		isPid = args[0]
	}
	return GetHandByPid(id, isPid)
}

// GetHandByPid get handle mdata by pid
func GetHandByPid(pid int, args ...int) Handle {
	isPid := 0
	if len(args) > 0 || NotPid {
		isPid = 1
	}
	return Handle(platformSetHandlePid(uintptr(pid), int8(isPid)))
}

// Deprecated: use the GetHandByPid(),
//
// GetHandPid get handle mdata by pid
func GetHandPid(pid int, args ...int) Handle {
	return GetHandByPid(pid, args...)
}

// GetHandByPidC get handle mdata by pid (returns MData)
func GetHandByPidC(pid int, args ...int) MData {
	isPid := 0
	if len(args) > 0 || NotPid {
		isPid = 1
	}
	return platformSetHandlePid(uintptr(pid), int8(isPid))
}

// GetHandle get the window handle
func GetHandle() int {
	return int(platformGetHandle())
}

// Deprecated: use the GetHandle(),
//
// # GetBHandle get the window handle, Wno-deprecated
//
// This function will be removed in version v1.0.0
func GetBHandle() int {
	tt.Drop("GetBHandle", "GetHandle")
	return int(platformBGetHandle())
}

func cgetTitle(pid, isPid int) string {
	return platformGetTitleByPid(uintptr(pid), int8(isPid))
}

// GetTitle get the window title return string
//
// Examples:
//
//	fmt.Println(robotgo.GetTitle())
//
//	ids, _ := robotgo.FindIds()
//	robotgo.GetTitle(ids[0])
func GetTitle(args ...int) string {
	if len(args) <= 0 {
		return platformGetMainTitle()
	}
	if len(args) > 1 {
		return internalGetTitle(args[0], args[1])
	}
	return internalGetTitle(args[0])
}

// GetPid get the process id return int32
func GetPid() int {
	return int(platformGetPID())
}

// internalGetBounds get the window bounds
func internalGetBounds(pid, isPid int) (int, int, int, int) {
	b := platformGetBounds(uintptr(pid), int8(isPid))
	return int(b.X), int(b.Y), int(b.W), int(b.H)
}

// internalGetClient get the window client bounds
func internalGetClient(pid, isPid int) (int, int, int, int) {
	b := platformGetClient(uintptr(pid), int8(isPid))
	return int(b.X), int(b.Y), int(b.W), int(b.H)
}

// Is64Bit determine whether the sys is 64bit
func Is64Bit() bool {
	return platformIs64Bit()
}

func internalActive(pid, isPid int) {
	platformActivePID(uintptr(pid), int8(isPid))
}

// ActiveName active the window by name
//
// Examples:
//
//	robotgo.ActiveName("chrome")
func ActiveName(name string) error {
	pids, err := FindIds(name)
	if err == nil && len(pids) > 0 {
		return ActivePid(pids[0])
	}
	return err
}

// GoString trans *int8 to string (replaces C.GoString)
func GoString(char *int8) string {
	if char == nil {
		return ""
	}
	var length int
	for p := unsafe.Pointer(char); *(*byte)(p) != 0; p = unsafe.Pointer(uintptr(p) + 1) {
		length++
	}
	if length == 0 {
		return ""
	}
	slice := unsafe.Slice((*byte)(unsafe.Pointer(char)), length)
	return string(slice)
}

func init() {
	platformInit()
}
