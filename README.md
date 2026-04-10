# Robotgo (Fork)

> **This is a fork of [go-vgo/robotgo](https://github.com/go-vgo/robotgo) maintained at [shaolei/robotgo](https://github.com/shaolei/robotgo).**
>
> Module path: `github.com/shaolei/robotgo` (NOT `github.com/go-vgo/robotgo`)

<!-- <img align="right" src="https://raw.githubusercontent.com/go-vgo/robotgo/master/logo.jpg"> -->

[![Go Report Card](https://goreportcard.com/badge/github.com/shaolei/robotgo)](https://goreportcard.com/report/github.com/shaolei/robotgo)
[![GoDoc](https://pkg.go.dev/badge/github.com/shaolei/robotgo?status.svg)](https://pkg.go.dev/github.com/shaolei/robotgo?tab=doc)
[![GitHub release](https://img.shields.io/github/release/shaolei/robotgo.svg)](https://github.com/shaolei/robotgo/releases/latest)

> Golang Desktop Automation, auto test and AI Computer Use. <br>
> Control the mouse, keyboard, read the screen, process, Window Handle, image and bitmap and global event listener.

RobotGo supports Mac, Windows, and Linux (X11); and robotgo supports arm64 and x86-amd64.

---

## Version & Compatibility

### Current Version: **v1.0.0**

This fork is based on the upstream `go-vgo/robotgo` but introduces **breaking changes** that are incompatible with the upstream API. If you are migrating from the upstream version, please read the following carefully.

### Key Differences from Upstream

| Feature | Upstream (`go-vgo/robotgo`) | This Fork (`shaolei/robotgo`) |
|---|---|---|
| **Build System** | CGo (requires C compiler: GCC/MinGW) | Pure Go (syscall + purego, **no C compiler needed**) |
| **Module Path** | `github.com/go-vgo/robotgo` | `github.com/shaolei/robotgo` |
| **Keycode Type** | `uint16` | `uint32` (supports X11 XF86 multimedia keys > 0xFFFF) |
| **Wheel Constants** | `WheelDown`, `WheelUp`, `WheelLeft`, `WheelRight` | `WHEEL_DOWN`, `WHEEL_UP`, `WHEEL_LEFT`, `WHEEL_RIGHT` |
| **Color Type** | Custom `colorRGBA` struct | Standard `image/color.RGBA` |
| **Platform Files** | `robotgo_mac.go`, `robotgo_win.go`, `robotgo_x11.go` (CGo) | `robotgo_darwin.go`, `robotgo_windows.go`, `robotgo_linux.go` (pure Go) |
| **Bitmap Memory** | `unsafe.Pointer` arithmetic | `unsafe.Slice` (Go 1.17+) |
| **Requirements** | GCC + CGo enabled | Go only (no GCC/CGo required) |

### Breaking Changes (Migration Required)

If you are migrating from `github.com/go-vgo/robotgo`, you MUST update your code:

#### 1. Change Import Path

```go
// Before
import "github.com/go-vgo/robotgo"

// After
import "github.com/shaolei/robotgo"
```

#### 2. Keycode Type Change (uint16 → uint32)

```go
// Before
var code uint16 = robotgo.K_Enter

// After
var code uint32 = robotgo.K_Enter
```

This change was necessary because X11 KeySym values for XF86 multimedia keys (e.g., `0x1008FF12` for XF86AudioMute) exceed the `uint16` range (max 0xFFFF).

#### 3. Wheel Constant Names

```go
// Before
robotgo.CheckMouse(robotgo.WheelDown, ...)

// After
robotgo.CheckMouse(robotgo.WHEEL_DOWN, ...)
```

#### 4. No C Compiler Required

This fork eliminates the CGo dependency entirely:
- **Windows**: Uses `golang.org/x/sys/windows` syscall
- **macOS**: Uses `github.com/ebitengine/purego` for dynamic library loading
- **Linux**: Uses X11/XCB Go bindings (`github.com/jezek/xgb`)

You no longer need GCC, MinGW, or Xcode Command Line Tools to build this library.

### What Remains Compatible

The following APIs are **fully compatible** with the upstream version (no code changes needed beyond import path):

- Mouse control: `Move`, `Click`, `Toggle`, `Scroll`, `DragSmooth`, `MoveSmooth`
- Keyboard control: `KeyTap`, `KeyToggle`, `Type`, `UnicodeType`
- Screen capture: `CaptureScreen`, `CaptureImg`, `GetPixelColor`, `GetScreenSize`
- Bitmap operations: `ToImage`, `Save`, `FreeBitmap`
- Clipboard: `ReadAll`, `WriteAll`
- Process management: `FindIds`, `ActivePid`, `Kill`, `PidExists`
- Window handle: `GetTitle`, `GetHandle`, `SetActive`, `MinWindow`, `MaxWindow`, `CloseWindow`

### Not Yet Implemented

The following features from upstream are **not yet available** in this fork:

- **Global event hook** (`gohook` integration) — requires CGo, not yet ported
- **Bitmap find** (`bitmap.Find`) — requires the `vcaesar/bitmap` CGo library
- **OCR** (`GetText`) — available but requires `gosseract` CGo binding and tesseract

---

## Contents

- [Version & Compatibility](#version--compatibility)
- [Docs](#docs)
- [Requirements](#requirements)
- [Installation](#installation)
- [Update](#update)
- [Examples](#examples)
- [Type Conversion and keys](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [Cross-Compiling](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [Authors](#authors)
- [License](#license)

## Docs

- [GoDoc](https://pkg.go.dev/github.com/shaolei/robotgo) <br>
- [API Docs](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (Deprecated, not updated)

## Requirements:

**This fork requires only Go (1.21+). No C compiler is needed.**

### ALL:

```
Golang 1.21+
```

#### For MacOS:

```
brew install go
```

Privacy setting: add Screen Recording and Accessibility under:
`System Settings > Privacy & Security > Accessibility, Screen & System Audio Recording`.

#### For Windows:

```
winget install Golang.go
```

No MinGW or GCC required. This fork uses pure Go syscall.

#### For Linux:

```
# Go
sudo snap install go --classic

# X11 (required for keyboard/mouse/screen control)
sudo apt install libx11-dev xorg-dev libxtst-dev

# Clipboard
sudo apt install xsel xclip
```

## Installation:

With Go module support (Go 1.11+), just import:

```go
import "github.com/shaolei/robotgo"
```

Otherwise, to install the robotgo package, run the command:

```
go get github.com/shaolei/robotgo
```

## Update:

```
go get -u github.com/shaolei/robotgo
```

## [Examples:](https://github.com/shaolei/robotgo/blob/master/examples)

#### [Mouse](https://github.com/shaolei/robotgo/blob/master/examples/mouse/main.go)

```Go
package main

import (
  "fmt"
  "github.com/shaolei/robotgo"
)

func main() {
  robotgo.MouseSleep = 300

  robotgo.Move(100, 100)
  fmt.Println(robotgo.Location())
  robotgo.Move(100, -200) // multi screen supported
  robotgo.MoveSmooth(120, -150)
  fmt.Println(robotgo.Location())

  robotgo.ScrollDir(10, "up")
  robotgo.ScrollDir(20, "right")

  robotgo.Scroll(0, -10)
  robotgo.Scroll(100, 0)

  robotgo.MilliSleep(100)
  robotgo.ScrollSmooth(-10, 6)

  robotgo.Move(10, 20)
  robotgo.MoveRelative(0, -10)
  robotgo.DragSmooth(10, 10)

  robotgo.Click("wheelRight")
  robotgo.Click("left", true)
  robotgo.MoveSmooth(100, 200, 1.0, 10.0)

  robotgo.Toggle("left")
  robotgo.Toggle("left", "up")
}
```

#### [Keyboard](https://github.com/shaolei/robotgo/blob/master/examples/key/main.go)

```Go
package main

import (
  "fmt"

  "github.com/shaolei/robotgo"
)

func main() {
  robotgo.Type("Hello World")
  robotgo.Type("だんしゃり", 0, 1)

  robotgo.KeySleep = 100
  robotgo.KeyTap("enter")
  robotgo.KeyTap("i", "alt", "cmd")

  arr := []string{"alt", "cmd"}
  robotgo.KeyTap("i", arr)

  robotgo.MilliSleep(100)
  robotgo.KeyToggle("a")
  robotgo.KeyToggle("a", "up")

  robotgo.WriteAll("Test")
  text, err := robotgo.ReadAll()
  if err == nil {
    fmt.Println(text)
  }
}
```

#### [Screen](https://github.com/shaolei/robotgo/blob/master/examples/screen/main.go)

```Go
package main

import (
  "fmt"
  "strconv"

  "github.com/shaolei/robotgo"
)

func main() {
  x, y := robotgo.Location()
  fmt.Println("pos: ", x, y)

  color := robotgo.GetPixelColor(100, 200)
  fmt.Println("color---- ", color)

  sx, sy := robotgo.GetScreenSize()
  fmt.Println("get screen size: ", sx, sy)

  bit := robotgo.CaptureScreen(10, 10, 30, 30)
  defer robotgo.FreeBitmap(bit)

  img := robotgo.ToImage(bit)
  robotgo.Save(img, "test.png")

  num := robotgo.DisplaysNum()
  for i := 0; i < num; i++ {
    robotgo.DisplayID = i
    img1, _ := robotgo.CaptureImg()
    path1 := "save_" + strconv.Itoa(i)
    robotgo.Save(img1, path1+".png")
    robotgo.SaveJpeg(img1, path1+".jpeg", 50)

    img2, _ := robotgo.CaptureImg(10, 10, 20, 20)
    robotgo.Save(img2, "test_"+strconv.Itoa(i)+".png")

    x, y, w, h := robotgo.GetDisplayBounds(i)
    img3, err := robotgo.CaptureImg(x, y, w, h)
    fmt.Println("Capture error: ", err)
    robotgo.Save(img3, path1+"_1.png")
  }
}
```

#### [Window](https://github.com/shaolei/robotgo/blob/master/examples/window/main.go)

```Go
package main

import (
  "fmt"

  "github.com/shaolei/robotgo"
)

func main() {
  fpid, err := robotgo.FindIds("Google")
  if err == nil {
    fmt.Println("pids... ", fpid)

    if len(fpid) > 0 {
      robotgo.ActivePid(fpid[0])

      robotgo.Kill(fpid[0])
    }
  }

  robotgo.ActiveName("chrome")

  isExist, err := robotgo.PidExists(100)
  if err == nil && isExist {
    fmt.Println("pid exists is", isExist)

    robotgo.Kill(100)
  }

  abool := robotgo.Alert("test", "robotgo")
  if abool {
    fmt.Println("ok@@@ ", "ok")
  }

  title := robotgo.GetTitle()
  fmt.Println("title@@@ ", title)
}
```

## Authors

- [Original author: Evans (vcaesar)](https://github.com/vcaesar)
- [Fork maintainer: shaolei](https://github.com/shaolei)
- [Original Maintainers](https://github.com/orgs/go-vgo/people)

## Upstream Reference

- Original repository: [go-vgo/robotgo](https://github.com/go-vgo/robotgo)
- [RobotGo-Pro](https://github.com/vcaesar/robotgo-pro) — JavaScript, Python, Lua versions, tech support, and newest features (Wayland support, etc.)

## License

Robotgo is primarily distributed under the terms of "the Apache License (Version 2.0)", with portions covered by various BSD-like licenses.

See [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/shaolei/robotgo/blob/master/LICENSE).
