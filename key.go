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

package robotgo

import (
	"errors"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"github.com/shaolei/robotgo/clipboard"
)

// Defining a bunch of constants.
const (
	// KeyA define key "a"
	KeyA = "a"
	KeyB = "b"
	KeyC = "c"
	KeyD = "d"
	KeyE = "e"
	KeyF = "f"
	KeyG = "g"
	KeyH = "h"
	KeyI = "i"
	KeyJ = "j"
	KeyK = "k"
	KeyL = "l"
	KeyM = "m"
	KeyN = "n"
	KeyO = "o"
	KeyP = "p"
	KeyQ = "q"
	KeyR = "r"
	KeyS = "s"
	KeyT = "t"
	KeyU = "u"
	KeyV = "v"
	KeyW = "w"
	KeyX = "x"
	KeyY = "y"
	KeyZ = "z"
	//
	CapA = "A"
	CapB = "B"
	CapC = "C"
	CapD = "D"
	CapE = "E"
	CapF = "F"
	CapG = "G"
	CapH = "H"
	CapI = "I"
	CapJ = "J"
	CapK = "K"
	CapL = "L"
	CapM = "M"
	CapN = "N"
	CapO = "O"
	CapP = "P"
	CapQ = "Q"
	CapR = "R"
	CapS = "S"
	CapT = "T"
	CapU = "U"
	CapV = "V"
	CapW = "W"
	CapX = "X"
	CapY = "Y"
	CapZ = "Z"
	//
	Key0      = "0"
	Key1      = "1"
	Key2      = "2"
	Key3      = "3"
	Key4      = "4"
	Key5      = "5"
	Key6      = "6"
	Key7      = "7"
	Key8      = "8"
	Key9      = "9"
	KeyGrave  = "`"
	KeyQuoter = '"'
	KeyQuote  = "'"

	// Backspace backspace key string
	Backspace = "backspace"
	Delete    = "delete"
	Enter     = "enter"
	Tab       = "tab"
	Esc       = "esc"
	Escape    = "escape"
	Up        = "up"    // Up arrow key
	Down      = "down"  // Down arrow key
	Right     = "right" // Right arrow key
	Left      = "left"  // Left arrow key
	Home      = "home"
	End       = "end"
	Pageup    = "pageup"
	Pagedown  = "pagedown"

	F1  = "f1"
	F2  = "f2"
	F3  = "f3"
	F4  = "f4"
	F5  = "f5"
	F6  = "f6"
	F7  = "f7"
	F8  = "f8"
	F9  = "f9"
	F10 = "f10"
	F11 = "f11"
	F12 = "f12"
	F13 = "f13"
	F14 = "f14"
	F15 = "f15"
	F16 = "f16"
	F17 = "f17"
	F18 = "f18"
	F19 = "f19"
	F20 = "f20"
	F21 = "f21"
	F22 = "f22"
	F23 = "f23"
	F24 = "f24"

	Cmd         = "cmd"  // is the "win" key for windows
	Lcmd        = "lcmd" // left command
	Rcmd        = "rcmd" // right command
	Alt         = "alt"
	Lalt        = "lalt"
	Ralt        = "ralt"
	Ctrl        = "ctrl"
	Lctrl       = "lctrl"
	Rctrl       = "rctrl"
	Control     = "control"
	Shift       = "shift"
	Lshift      = "lshift"
	Rshift      = "rshift"
	Capslock    = "capslock"
	Space       = "space"
	Print       = "print"
	Printscreen = "printscreen"
	Insert      = "insert"
	Menu        = "menu"

	AudioMute    = "audio_mute"
	AudioVolDown = "audio_vol_down"
	AudioVolUp   = "audio_vol_up"
	AudioPlay    = "audio_play"
	AudioStop    = "audio_stop"
	AudioPause   = "audio_pause"
	AudioPrev    = "audio_prev"
	AudioNext    = "audio_next"
	AudioRewind  = "audio_rewind"
	AudioForward = "audio_forward"
	AudioRepeat  = "audio_repeat"
	AudioRandom  = "audio_random"

	Num0    = "num0"
	Num1    = "num1"
	Num2    = "num2"
	Num3    = "num3"
	Num4    = "num4"
	Num5    = "num5"
	Num6    = "num6"
	Num7    = "num7"
	Num8    = "num8"
	Num9    = "num9"
	NumLock = "num_lock"

	NumDecimal = "num."
	NumPlus    = "num+"
	NumMinus   = "num-"
	NumMul     = "num*"
	NumDiv     = "num/"
	NumClear   = "num_clear"
	NumEnter   = "num_enter"
	NumEqual   = "num_equal"

	LightsMonUp     = "lights_mon_up"
	LightsMonDown   = "lights_mon_down"
	LightsKbdToggle = "lights_kbd_toggle"
	LightsKbdUp     = "lights_kbd_up"
	LightsKbdDown   = "lights_kbd_down"
)

// keyNames define a map of key names to key code
var keyNames = map[string]uint32{
	"backspace": K_BACKSPACE,
	"delete":    K_DELETE,
	"enter":     K_RETURN,
	"tab":       K_TAB,
	"esc":       K_ESCAPE,
	"escape":    K_ESCAPE,
	"up":        K_UP,
	"down":      K_DOWN,
	"right":     K_RIGHT,
	"left":      K_LEFT,
	"home":      K_HOME,
	"end":       K_END,
	"pageup":    K_PAGEUP,
	"pagedown":  K_PAGEDOWN,
	//
	"f1":  K_F1,
	"f2":  K_F2,
	"f3":  K_F3,
	"f4":  K_F4,
	"f5":  K_F5,
	"f6":  K_F6,
	"f7":  K_F7,
	"f8":  K_F8,
	"f9":  K_F9,
	"f10": K_F10,
	"f11": K_F11,
	"f12": K_F12,
	"f13": K_F13,
	"f14": K_F14,
	"f15": K_F15,
	"f16": K_F16,
	"f17": K_F17,
	"f18": K_F18,
	"f19": K_F19,
	"f20": K_F20,
	"f21": K_F21,
	"f22": K_F22,
	"f23": K_F23,
	"f24": K_F24,
	//
	"cmd":         K_META,
	"lcmd":        K_LMETA,
	"rcmd":        K_RMETA,
	"command":     K_META,
	"alt":         K_ALT,
	"lalt":        K_LALT,
	"ralt":        K_RALT,
	"ctrl":        K_CONTROL,
	"lctrl":       K_LCONTROL,
	"rctrl":       K_RCONTROL,
	"control":     K_CONTROL,
	"shift":       K_SHIFT,
	"lshift":      K_LSHIFT,
	"rshift":      K_RSHIFT,
	"right_shift": K_RSHIFT,
	"capslock":    K_CAPSLOCK,
	"space":       K_SPACE,
	"print":       K_PRINTSCREEN,
	"printscreen": K_PRINTSCREEN,
	"insert":      K_INSERT,
	"menu":        K_MENU,

	"audio_mute":     K_AUDIO_VOLUME_MUTE,
	"audio_vol_down": K_AUDIO_VOLUME_DOWN,
	"audio_vol_up":   K_AUDIO_VOLUME_UP,
	"audio_play":     K_AUDIO_PLAY,
	"audio_stop":     K_AUDIO_STOP,
	"audio_pause":    K_AUDIO_PAUSE,
	"audio_prev":     K_AUDIO_PREV,
	"audio_next":     K_AUDIO_NEXT,
	"audio_rewind":   K_AUDIO_REWIND,
	"audio_forward":  K_AUDIO_FORWARD,
	"audio_repeat":   K_AUDIO_REPEAT,
	"audio_random":   K_AUDIO_RANDOM,

	"num0":     K_NUMPAD_0,
	"num1":     K_NUMPAD_1,
	"num2":     K_NUMPAD_2,
	"num3":     K_NUMPAD_3,
	"num4":     K_NUMPAD_4,
	"num5":     K_NUMPAD_5,
	"num6":     K_NUMPAD_6,
	"num7":     K_NUMPAD_7,
	"num8":     K_NUMPAD_8,
	"num9":     K_NUMPAD_9,
	"num_lock": K_NUMPAD_LOCK,

	// todo: removed
	"numpad_0":    K_NUMPAD_0,
	"numpad_1":    K_NUMPAD_1,
	"numpad_2":    K_NUMPAD_2,
	"numpad_3":    K_NUMPAD_3,
	"numpad_4":    K_NUMPAD_4,
	"numpad_5":    K_NUMPAD_5,
	"numpad_6":    K_NUMPAD_6,
	"numpad_7":    K_NUMPAD_7,
	"numpad_8":    K_NUMPAD_8,
	"numpad_9":    K_NUMPAD_9,
	"numpad_lock": K_NUMPAD_LOCK,

	"num.":      K_NUMPAD_DECIMAL,
	"num+":      K_NUMPAD_PLUS,
	"num-":      K_NUMPAD_MINUS,
	"num*":      K_NUMPAD_MUL,
	"num/":      K_NUMPAD_DIV,
	"num_clear": K_NUMPAD_CLEAR,
	"num_enter": K_NUMPAD_ENTER,
	"num_equal": K_NUMPAD_EQUAL,

	"lights_mon_up":     K_LIGHTS_MON_UP,
	"lights_mon_down":   K_LIGHTS_MON_DOWN,
	"lights_kbd_toggle": K_LIGHTS_KBD_TOGGLE,
	"lights_kbd_up":     K_LIGHTS_KBD_UP,
	"lights_kbd_down":   K_LIGHTS_KBD_DOWN,
}

// CmdCtrl If the operating system is macOS, return the key string "cmd",
// otherwise return the key string "ctrl
func CmdCtrl() string {
	if runtime.GOOS == "darwin" {
		return "cmd"
	}
	return "ctrl"
}

// tapKeyCode sends a key press and release to the active application
func tapKeyCode(code uint32, flags uint64, pid uintptr) {
	platformToggleKeyCode(code, true, flags, pid)
	MilliSleep(3)
	platformToggleKeyCode(code, false, flags, pid)
}

var keyErr = errors.New("Invalid key flag specified.")

func checkKeyCodes(k string) (key uint32, err error) {
	if k == "" {
		return
	}

	if len(k) == 1 {
		key = platformKeyCodeForChar(k[0])
		if key == K_NOT_A_KEY {
			err = keyErr
			return
		}
		return
	}

	if v, ok := keyNames[k]; ok {
		key = v
		if key == K_NOT_A_KEY {
			err = keyErr
			return
		}
	}
	return
}

func checkKeyFlags(f string) uint64 {
	m := map[string]uint64{
		"alt":    MOD_ALT,
		"ralt":   MOD_ALT,
		"lalt":   MOD_ALT,
		"cmd":    MOD_META,
		"rcmd":   MOD_META,
		"lcmd":   MOD_META,
		"ctrl":   MOD_CONTROL,
		"rctrl":  MOD_CONTROL,
		"lctrl":  MOD_CONTROL,
		"shift":  MOD_SHIFT,
		"rshift": MOD_SHIFT,
		"lshift": MOD_SHIFT,
		"none":   MOD_NONE,
	}

	if v, ok := m[f]; ok {
		return v
	}
	return 0
}

func getFlagsFromValue(value []string) uint64 {
	var flags uint64
	if len(value) <= 0 {
		return flags
	}

	for i := 0; i < len(value); i++ {
		f := checkKeyFlags(value[i])
		flags = flags | f
	}
	return flags
}

func upKeyArr(keyArr []string, pid int) {
	for i := 0; i < len(keyArr); i++ {
		key1, _ := checkKeyCodes(keyArr[i])
		platformToggleKeyCode(key1, false, MOD_NONE, uintptr(pid))
	}
}

func keyTaps(k string, keyArr []string, pid int) error {
	flags := getFlagsFromValue(keyArr)
	key, err := checkKeyCodes(k)
	if err != nil {
		return err
	}

	tapKeyCode(key, flags, uintptr(pid))
	MilliSleep(KeySleep)
	upKeyArr(keyArr, pid)
	return nil
}

func getKeyDown(keyArr []string) (bool, []string) {
	if len(keyArr) <= 0 {
		keyArr = append(keyArr, "down")
	}

	down := true
	if keyArr[0] == "up" {
		down = false
	}

	if keyArr[0] == "up" || keyArr[0] == "down" {
		keyArr = keyArr[1:]
	}
	return down, keyArr
}

func keyTogglesB(k string, down bool, keyArr []string, pid int) error {
	flags := getFlagsFromValue(keyArr)
	key, err := checkKeyCodes(k)
	if err != nil {
		return err
	}

	platformToggleKeyCode(key, down, flags, uintptr(pid))
	MilliSleep(KeySleep)
	if !down {
		upKeyArr(keyArr, pid)
	}
	return nil
}

func keyToggles(k string, keyArr []string, pid int) error {
	down, keyArr1 := getKeyDown(keyArr)
	return keyTogglesB(k, down, keyArr1, pid)
}

/*
 __  ___  ___________    ____ .______     ______        ___      .______       _______
|  |/  / |   ____\   \  /   / |   _  \   /  __  \      /   \     |   _  \     |       \
|  '  /  |  |__   \   \/   /  |  |_)  | |  |  |  |    /  ^  \    |  |_)  |    |  .--.  |
|    <   |   __|   \_    _/   |   _  <  |  |  |  |   /  /_\  \   |      /     |  |  |  |
|  .  \  |  |____    |  |     |  |_)  | |  `--'  |  /  _____  \  |  |\  \----.|  '--'  |
|__|\__\ |_______|   |__|     |______/   \______/  /__/     \__\ | _| `._____||_______|

*/

// ToInterfaces convert []string to []interface{}
func ToInterfaces(fields []string) []interface{} {
	res := make([]interface{}, 0, len(fields))
	for _, s := range fields {
		res = append(res, s)
	}
	return res
}

// ToStrings convert []interface{} to []string
func ToStrings(fields []interface{}) []string {
	res := make([]string, 0, len(fields))
	for _, s := range fields {
		res = append(res, s.(string))
	}
	return res
}

// toErr converts a string to a Go error
func toErr(str string) error {
	if str == "" {
		return nil
	}
	return errors.New(str)
}

func appendShift(key string, len1 int, args ...interface{}) (string, []interface{}) {
	if len(key) > 0 && unicode.IsUpper([]rune(key)[0]) {
		args = append(args, "shift")
	}

	key = strings.ToLower(key)
	if _, ok := Special[key]; ok {
		key = Special[key]
		if len(args) <= len1 {
			args = append(args, "shift")
		}
	}

	return key, args
}

// KeyTap taps the keyboard code;
//
// See keys supported:
//
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#keys
//
// Examples:
//
//	robotgo.KeySleep = 100 // 100 millisecond
//	robotgo.KeyTap("a")
//	robotgo.KeyTap("i", "alt", "command")
//
//	arr := []string{"alt", "command"}
//	robotgo.KeyTap("i", arr)
//
//	robotgo.KeyTap("k", pid int)
func KeyTap(key string, args ...interface{}) error {
	var keyArr []string
	key, args = appendShift(key, 0, args...)

	pid := 0
	if len(args) > 0 {
		if reflect.TypeOf(args[0]) == reflect.TypeOf(keyArr) {
			keyArr = args[0].([]string)
		} else {
			if reflect.TypeOf(args[0]) == reflect.TypeOf(pid) {
				pid = args[0].(int)
				keyArr = ToStrings(args[1:])
			} else {
				keyArr = ToStrings(args)
			}
		}
	}

	return keyTaps(key, keyArr, pid)
}

func getToggleArgs(args ...interface{}) (pid int, keyArr []string) {
	if len(args) > 0 && reflect.TypeOf(args[0]) == reflect.TypeOf(pid) {
		pid = args[0].(int)
		keyArr = ToStrings(args[1:])
	} else {
		keyArr = ToStrings(args)
	}
	return
}

// KeyToggle toggles the keyboard, if there not have args default is "down"
//
// See keys:
//
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#keys
//
// Examples:
//
//	robotgo.KeyToggle("a")
//	robotgo.KeyToggle("a", "up")
//
//	robotgo.KeyToggle("a", "up", "alt", "cmd")
//	robotgo.KeyToggle("k", pid int)
func KeyToggle(key string, args ...interface{}) error {
	key, args = appendShift(key, 1, args...)
	pid, keyArr := getToggleArgs(args...)
	return keyToggles(key, keyArr, pid)
}

// KeyPress press key string
func KeyPress(key string, args ...interface{}) error {
	err := KeyDown(key, args...)
	if err != nil {
		return err
	}

	MilliSleep(1 + rand.Intn(3))
	return KeyUp(key, args...)
}

// KeyDown press down a key
func KeyDown(key string, args ...interface{}) error {
	return KeyToggle(key, args...)
}

// KeyUp press up a key
func KeyUp(key string, args ...interface{}) error {
	arr := []interface{}{"up"}
	arr = append(arr, args...)
	return KeyToggle(key, arr...)
}

// ReadAll read string from clipboard
func ReadAll() (string, error) {
	return clipboard.ReadAll()
}

// WriteAll write string to clipboard
func WriteAll(text string) error {
	return clipboard.WriteAll(text)
}

// CharCodeAt char code at utf-8
func CharCodeAt(s string, n int) rune {
	i := 0
	for _, r := range s {
		if i == n {
			return r
		}
		i++
	}
	return 0
}

// UnicodeType tap the uint32 unicode
func UnicodeType(str uint32, args ...int) {
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}
	isPid := 0
	if len(args) > 1 {
		isPid = args[1]
	}
	platformUnicodeType(str, uintptr(pid), int8(isPid))
}

// ToUC trans string to unicode []string
func ToUC(text string) []string {
	var uc []string
	for _, r := range text {
		textQ := strconv.QuoteToASCII(string(r))
		textUnQ := textQ[1 : len(textQ)-1]
		st := strings.Replace(textUnQ, "\\u", "U", -1)
		if st == "\\\\" {
			st = "\\"
		}
		if st == `\"` {
			st = `"`
		}
		uc = append(uc, st)
	}
	return uc
}

func inputUTF(str string) {
	platformInputUTF(str)
}

// TypeStr tap a string
//
// Deprecated: use the Type()
func TypeStr(str string, args ...int) {
	Type(str, args...)
}

// Type type a string (supported UTF-8)
//
// robotgo.Type(string: "The string to send", int: pid, "milli_sleep time", "x11 option")
//
// Examples:
//
//	robotgo.Type("abc@123, Hi galaxy, こんにちは")
//	robotgo.Type("To be or not to be, this is questions.", pid int)
func Type(str string, args ...int) {
	var tm, tm1 = 0, 7

	if len(args) > 1 {
		tm = args[1]
	}
	if len(args) > 2 {
		tm1 = args[2]
	}
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}

	if runtime.GOOS == "linux" {
		strUc := ToUC(str)
		for i := 0; i < len(strUc); i++ {
			ru := []rune(strUc[i])
			if len(ru) <= 1 {
				ustr := uint32(CharCodeAt(strUc[i], 0))
				UnicodeType(ustr, pid)
			} else {
				inputUTF(strUc[i])
				MilliSleep(tm1)
			}
			MilliSleep(tm)
		}
		return
	}

	for i := 0; i < len([]rune(str)); i++ {
		ustr := uint32(CharCodeAt(str, i))
		UnicodeType(ustr, pid)
		MilliSleep(tm)
	}
	MilliSleep(KeySleep)
}

// PasteStr paste a string
//
// Deprecated: use the Paste()
func PasteStr(str string) error {
	return Paste(str)
}

// Paste paste a string (supported UTF-8),
// write the string to clipboard and tap `cmd + v`
func Paste(str string) error {
	err := clipboard.WriteAll(str)
	if err != nil {
		return err
	}
	return CmdV()
}

// CmdV tap key command + v or control + v
func CmdV() error {
	if runtime.GOOS == "darwin" {
		return KeyTap("v", "command")
	}
	return KeyTap("v", "control")
}

// TypeStrDelay type string width delay
//
// Deprecated: use the TypeDelay()
func TypeStrDelay(str string, delay int) {
	TypeDelay(str, delay)
}

// TypeDelay type string with delayed
// And you can use robotgo.KeySleep = 100 to delayed not this function
func TypeDelay(str string, delay int) {
	TypeStr(str)
	MilliSleep(delay)
}

// SetDelay sets the key and mouse delay
// robotgo.SetDelay(100) option the robotgo.KeySleep and robotgo.MouseSleep = d
func SetDelay(d ...int) {
	v := 10
	if len(d) > 0 {
		v = d[0]
	}
	KeySleep = v
	MouseSleep = v
}
