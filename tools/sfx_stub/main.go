package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	ole32    = syscall.NewLazyDLL("ole32.dll")

	procRegisterClassExW     = user32.NewProc("RegisterClassExW")
	procCreateWindowExW      = user32.NewProc("CreateWindowExW")
	procDefWindowProcW       = user32.NewProc("DefWindowProcW")
	procShowWindow           = user32.NewProc("ShowWindow")
	procUpdateWindow         = user32.NewProc("UpdateWindow")
	procGetMessageW          = user32.NewProc("GetMessageW")
	procTranslateMessage     = user32.NewProc("TranslateMessage")
	procDispatchMessageW     = user32.NewProc("DispatchMessageW")
	procSendMessageW         = user32.NewProc("SendMessageW")
	procSetWindowTextW       = user32.NewProc("SetWindowTextW")
	procGetWindowTextW       = user32.NewProc("GetWindowTextW")
	procPostQuitMessage      = user32.NewProc("PostQuitMessage")
	procMessageBoxW          = user32.NewProc("MessageBoxW")
	procGetModuleHandleW     = kernel32.NewProc("GetModuleHandleW")
	procBeginPaint           = user32.NewProc("BeginPaint")
	procEndPaint             = user32.NewProc("EndPaint")

	procGetStockObject       = gdi32.NewProc("GetStockObject")
	procCreateFontW          = gdi32.NewProc("CreateFontW")
	procCreateSolidBrush     = gdi32.NewProc("CreateSolidBrush")
	procCreatePen            = gdi32.NewProc("CreatePen")
	procSelectObject         = gdi32.NewProc("SelectObject")
	procDeleteObject         = gdi32.NewProc("DeleteObject")
	procRectangle            = gdi32.NewProc("Rectangle")
	procRoundRect            = gdi32.NewProc("RoundRect")
	procSetBkMode            = gdi32.NewProc("SetBkMode")
	procSetTextColor         = gdi32.NewProc("SetTextColor")
	procTextOutW             = gdi32.NewProc("TextOutW")
	procMoveToEx             = gdi32.NewProc("MoveToEx")
	procLineTo               = gdi32.NewProc("LineTo")

	procSHBrowseForFolderW   = shell32.NewProc("SHBrowseForFolderW")
	procSHGetPathFromIDListW = shell32.NewProc("SHGetPathFromIDListW")
	procShellExecuteW        = shell32.NewProc("ShellExecuteW")
	procCoTaskMemFree        = ole32.NewProc("CoTaskMemFree")
)

const (
	WS_OVERLAPPEDWINDOW = 0x00CF0000
	WS_VISIBLE          = 0x10000000
	WS_CHILD            = 0x40000000
	WS_BORDER           = 0x00800000
	WS_TABSTOP          = 0x00010000
	ES_AUTOHSCROLL      = 0x0080
	BS_PUSHBUTTON       = 0x0000
	BS_DEFPUSHBUTTON    = 0x0001
	SS_LEFT             = 0x0000

	WM_DESTROY = 0x0002
	WM_PAINT   = 0x000F
	WM_COMMAND = 0x0111
	WM_SETFONT = 0x0030

	PBM_SETRANGE = 0x0401
	PBM_SETPOS   = 0x0402

	ID_INSTALL_BTN = 1001
	ID_BROWSE_BTN  = 1002
	ID_CANCEL_BTN  = 1003
	ID_HELP_BTN    = 1004
	ID_PATH_EDIT   = 1005
)

type WNDCLASSEXW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}

type PAINTSTRUCT struct {
	Hdc         uintptr
	FErase      int32
	RcPaint     [4]int32
	FRestore    int32
	FIncUpdate  int32
	RgbReserved [32]byte
}

type POINT struct {
	X, Y int32
}

type MSG struct {
	Hwnd     uintptr
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       POINT
	LPrivate uint32
}

type BROWSEINFOW struct {
	HwndOwner      uintptr
	PidlRoot       uintptr
	PszDisplayName *uint16
	LpszTitle      *uint16
	UlFlags        uint32
	Lpfn           uintptr
	LParam         uintptr
	IImage         int32
}

var (
	MagicMarker = []byte("WINRAR_SFX_PAYLOAD_V1\x00\x00")
	hwndMain    uintptr
	hwndEdit    uintptr
	hwndProg    uintptr
	hwndStatus  uintptr
	hwndInstall uintptr
	hwndBrowse  uintptr
	hwndCancel  uintptr
	hwndHelp    uintptr

	payloadZipBytes []byte
	configDestPath  = "C:\\Program Files\\VibeRAR"
	configSetupCmd  = ""
	configTitle     = "VibeRAR self-extracting archive"
	configPrompt    = "Please choose the destination folder and click Install to extract the application files."
	configSilent    = "0"
)

func utf16Ptr(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}

func neg(v int) int {
	return -v
}

func rgb(r, g, b uint8) uint32 {
	return uint32(r) | (uint32(g) << 8) | (uint32(b) << 16)
}

func loword(l uint32) uint16 {
	return uint16(l & 0xffff)
}

func drawVibeRARBooks(hdc uintptr, x, y, w, h int) {
	hWhiteBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(255, 255, 255)))
	hBorderPen, _, _ := procCreatePen.Call(0, 1, uintptr(rgb(205, 215, 225)))
	oldBrush, _, _ := procSelectObject.Call(hdc, hWhiteBrush)
	oldPen, _, _ := procSelectObject.Call(hdc, hBorderPen)

	procRoundRect.Call(hdc, uintptr(x), uintptr(y), uintptr(x+w), uintptr(y+h), 8, 8)

	hHeaderBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(240, 245, 252)))
	procSelectObject.Call(hdc, hHeaderBrush)
	procRoundRect.Call(hdc, uintptr(x+1), uintptr(y+1), uintptr(x+w-1), uintptr(y+26), 6, 6)

	procSetBkMode.Call(hdc, 1)
	procSetTextColor.Call(hdc, uintptr(rgb(20, 70, 160)))
	hTitleFont, _, _ := procCreateFontW.Call(
		uintptr(neg(12)), 0, 0, 0, 700, 0, 0, 0,
		0, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(utf16Ptr("Segoe UI"))),
	)
	oldFont, _, _ := procSelectObject.Call(hdc, hTitleFont)
	titleStr := utf16Ptr("VibeRAR")
	procTextOutW.Call(hdc, uintptr(x+(w-52)/2), uintptr(y+5), uintptr(unsafe.Pointer(titleStr)), 7)

	procSelectObject.Call(hdc, oldFont)
	procDeleteObject.Call(hTitleFont)
	procDeleteObject.Call(hHeaderBrush)

	hShadowBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(225, 230, 236)))
	procSelectObject.Call(hdc, hShadowBrush)
	hNullPen, _, _ := procGetStockObject.Call(8)
	procSelectObject.Call(hdc, hNullPen)
	procRoundRect.Call(hdc, uintptr(x+10), uintptr(y+130), uintptr(x+w-10), uintptr(y+142), 8, 8)
	procDeleteObject.Call(hShadowBrush)

	drawSingleBook := func(bx, by, bw, bh int, spineCol, coverCol, ribCol uint32) {
		hCoverBrush, _, _ := procCreateSolidBrush.Call(uintptr(coverCol))
		procSelectObject.Call(hdc, hCoverBrush)
		procRoundRect.Call(hdc, uintptr(bx), uintptr(by), uintptr(bx+bw), uintptr(by+bh), 4, 4)
		procDeleteObject.Call(hCoverBrush)

		hSpineBrush, _, _ := procCreateSolidBrush.Call(uintptr(spineCol))
		procSelectObject.Call(hdc, hSpineBrush)
		procRoundRect.Call(hdc, uintptr(bx), uintptr(by), uintptr(bx+16), uintptr(by+bh), 4, 4)
		procDeleteObject.Call(hSpineBrush)

		hPageBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(250, 246, 238)))
		procSelectObject.Call(hdc, hPageBrush)
		procRectangle.Call(hdc, uintptr(bx+bw-10), uintptr(by+2), uintptr(bx+bw-2), uintptr(by+bh-2))
		procDeleteObject.Call(hPageBrush)

		hRibPen, _, _ := procCreatePen.Call(0, 2, uintptr(ribCol))
		procSelectObject.Call(hdc, hRibPen)
		procMoveToEx.Call(hdc, uintptr(bx+2), uintptr(by+5), 0)
		procLineTo.Call(hdc, uintptr(bx+14), uintptr(by+5))
		procMoveToEx.Call(hdc, uintptr(bx+2), uintptr(by+bh-6), 0)
		procLineTo.Call(hdc, uintptr(bx+14), uintptr(by+bh-6))
		procDeleteObject.Call(hRibPen)
	}

	drawSingleBook(x+12, y+102, w-24, 24, rgb(130, 20, 30), rgb(185, 35, 45), rgb(235, 195, 60))
	drawSingleBook(x+11, y+74, w-24, 24, rgb(20, 95, 40), rgb(35, 145, 60), rgb(235, 195, 60))
	drawSingleBook(x+10, y+46, w-24, 24, rgb(15, 55, 125), rgb(28, 85, 185), rgb(235, 195, 60))

	procSetTextColor.Call(hdc, uintptr(rgb(255, 225, 75)))
	hRarFont, _, _ := procCreateFontW.Call(
		uintptr(neg(10)), 0, 0, 0, 800, 0, 0, 0,
		0, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(utf16Ptr("Arial"))),
	)
	procSelectObject.Call(hdc, hRarFont)
	rarStr := utf16Ptr("RAR")
	procTextOutW.Call(hdc, uintptr(x+34), uintptr(y+51), uintptr(unsafe.Pointer(rarStr)), 3)
	procDeleteObject.Call(hRarFont)

	strapX := x + 44
	strapW := 13
	strapTop := y + 38
	strapBot := y + 130

	hStrapBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(105, 50, 20)))
	procSelectObject.Call(hdc, hStrapBrush)
	procSelectObject.Call(hdc, hNullPen)
	procRectangle.Call(hdc, uintptr(strapX), uintptr(strapTop), uintptr(strapX+strapW), uintptr(strapBot))
	procDeleteObject.Call(hStrapBrush)

	buckleY := y + 80
	buckleW := 17
	buckleH := 14
	buckleX := strapX - (buckleW-strapW)/2

	hBuckleBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(240, 200, 45)))
	procSelectObject.Call(hdc, hBuckleBrush)
	procRoundRect.Call(hdc, uintptr(buckleX), uintptr(buckleY), uintptr(buckleX+buckleW), uintptr(buckleY+buckleH), 3, 3)
	procDeleteObject.Call(hBuckleBrush)

	hHoleBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(65, 30, 10)))
	procSelectObject.Call(hdc, hHoleBrush)
	procRectangle.Call(hdc, uintptr(buckleX+3), uintptr(buckleY+2), uintptr(buckleX+buckleW-3), uintptr(buckleY+buckleH-2))
	procDeleteObject.Call(hHoleBrush)

	procSetTextColor.Call(hdc, uintptr(rgb(110, 120, 135)))
	hSubFont, _, _ := procCreateFontW.Call(
		uintptr(neg(9)), 0, 0, 0, 400, 0, 0, 0,
		0, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(utf16Ptr("Segoe UI"))),
	)
	procSelectObject.Call(hdc, hSubFont)
	subStr := utf16Ptr("SFX Archive")
	procTextOutW.Call(hdc, uintptr(x+(w-54)/2), uintptr(y+148), uintptr(unsafe.Pointer(subStr)), 11)
	procDeleteObject.Call(hSubFont)

	procSelectObject.Call(hdc, oldBrush)
	procSelectObject.Call(hdc, oldPen)
	procDeleteObject.Call(hWhiteBrush)
	procDeleteObject.Call(hBorderPen)
}

func wndProc(hwnd uintptr, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		wmId := loword(uint32(wParam))
		if wmId == ID_INSTALL_BTN {
			go performExtraction()
		} else if wmId == ID_BROWSE_BTN {
			browseFolder()
		} else if wmId == ID_CANCEL_BTN {
			procPostQuitMessage.Call(0)
		} else if wmId == ID_HELP_BTN {
			procMessageBoxW.Call(
				hwnd,
				uintptr(unsafe.Pointer(utf16Ptr("VibeRAR Self-Extracting Archive:\n\nClick 'Install' to extract all contents to the specified destination folder.\nClick 'Browse' to choose a different folder."))),
				uintptr(unsafe.Pointer(utf16Ptr("VibeRAR SFX Help"))),
				0x40,
			)
		}
		return 0

	case WM_PAINT:
		var ps PAINTSTRUCT
		hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		drawVibeRARBooks(hdc, 380, 15, 105, 175)
		procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		return 0

	case WM_DESTROY:
		procPostQuitMessage.Call(0)
		return 0
	}

	ret, _, _ := procDefWindowProcW.Call(hwnd, uintptr(msg), wParam, lParam)
	return ret
}

func browseFolder() {
	var bi BROWSEINFOW
	bi.HwndOwner = hwndMain
	bi.LpszTitle = utf16Ptr("Select destination folder for extraction:")
	bi.UlFlags = 0x00000041

	pidl, _, _ := procSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(&bi)))
	if pidl != 0 {
		var pathBuf [260]uint16
		procSHGetPathFromIDListW.Call(pidl, uintptr(unsafe.Pointer(&pathBuf[0])))
		procCoTaskMemFree.Call(pidl)
		newPath := syscall.UTF16ToString(pathBuf[:])
		if strings.TrimSpace(newPath) != "" {
			procSetWindowTextW.Call(hwndEdit, uintptr(unsafe.Pointer(utf16Ptr(strings.TrimSpace(newPath)))))
		}
	}
}

func performExtraction() {
	procSetWindowTextW.Call(hwndInstall, uintptr(unsafe.Pointer(utf16Ptr("Extracting..."))))

	var buf [512]uint16
	procGetWindowTextW.Call(hwndEdit, uintptr(unsafe.Pointer(&buf[0])), 512)
	dest := syscall.UTF16ToString(buf[:])
	if strings.TrimSpace(dest) == "" {
		dest = configDestPath
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		procMessageBoxW.Call(hwndMain, uintptr(unsafe.Pointer(utf16Ptr("Failed to create destination folder: "+err.Error()))), uintptr(unsafe.Pointer(utf16Ptr("VibeRAR SFX Error"))), 0x10)
		procSetWindowTextW.Call(hwndInstall, uintptr(unsafe.Pointer(utf16Ptr("Install"))))
		return
	}

	if len(payloadZipBytes) == 0 {
		procMessageBoxW.Call(hwndMain, uintptr(unsafe.Pointer(utf16Ptr("No archive payload embedded in this executable stub."))), uintptr(unsafe.Pointer(utf16Ptr("VibeRAR SFX"))), 0x40)
		procSetWindowTextW.Call(hwndInstall, uintptr(unsafe.Pointer(utf16Ptr("Install"))))
		return
	}

	zipReader, err := zip.NewReader(bytes.NewReader(payloadZipBytes), int64(len(payloadZipBytes)))
	if err != nil {
		procMessageBoxW.Call(hwndMain, uintptr(unsafe.Pointer(utf16Ptr("Failed to read archive payload: "+err.Error()))), uintptr(unsafe.Pointer(utf16Ptr("VibeRAR SFX Error"))), 0x10)
		procSetWindowTextW.Call(hwndInstall, uintptr(unsafe.Pointer(utf16Ptr("Install"))))
		return
	}

	totalFiles := len(zipReader.File)
	for i, file := range zipReader.File {
		statusText := fmt.Sprintf("Extracting: %s", file.Name)
		procSetWindowTextW.Call(hwndStatus, uintptr(unsafe.Pointer(utf16Ptr(statusText))))

		outPath := filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(outPath, file.Mode())
			continue
		}

		os.MkdirAll(filepath.Dir(outPath), 0755)
		outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err == nil {
			rc, err := file.Open()
			if err == nil {
				io.Copy(outFile, rc)
				rc.Close()
			}
			outFile.Close()
		}

		progress := int(((i + 1) * 100) / totalFiles)
		procSendMessageW.Call(hwndProg, PBM_SETPOS, uintptr(progress), 0)
	}

	procSetWindowTextW.Call(hwndStatus, uintptr(unsafe.Pointer(utf16Ptr("Extraction completed successfully!"))))
	procMessageBoxW.Call(hwndMain, uintptr(unsafe.Pointer(utf16Ptr("All files extracted successfully to:\n"+dest))), uintptr(unsafe.Pointer(utf16Ptr("VibeRAR SFX"))), 0x40)

	if configSetupCmd != "" {
		procShellExecuteW.Call(
			0,
			uintptr(unsafe.Pointer(utf16Ptr("open"))),
			uintptr(unsafe.Pointer(utf16Ptr(configSetupCmd))),
			0,
			uintptr(unsafe.Pointer(utf16Ptr(dest))),
			1,
		)
	}

	procPostQuitMessage.Call(0)
}

func main() {
	runtime.LockOSThread()

	exePath, err := os.Executable()
	if err == nil {
		data, err := os.ReadFile(exePath)
		if err == nil {
			idx := bytes.Index(data, MagicMarker)
			if idx != -1 {
				payload := data[idx+len(MagicMarker):]
				configIdx := bytes.LastIndex(payload, []byte(";--- VibeRAR SFX Script Configuration ---"))
				if configIdx != -1 {
					scriptText := string(payload[configIdx:])
					payloadZipBytes = payload[:configIdx]

					for _, line := range strings.Split(scriptText, "\n") {
						line = strings.TrimSpace(line)
						if strings.HasPrefix(line, "Path=") {
							configDestPath = strings.TrimPrefix(line, "Path=")
							configDestPath = os.ExpandEnv(configDestPath)
						} else if strings.HasPrefix(line, "Setup=") {
							configSetupCmd = strings.TrimPrefix(line, "Setup=")
						} else if strings.HasPrefix(line, "Title=") {
							configTitle = strings.TrimPrefix(line, "Title=")
						} else if strings.HasPrefix(line, "Text=") {
							configPrompt = strings.TrimPrefix(line, "Text=")
						} else if strings.HasPrefix(line, "Silent=") {
							configSilent = strings.TrimPrefix(line, "Silent=")
						}
					}
				} else {
					payloadZipBytes = payload
				}
			}
		}
	}

	isSilent := configSilent == "1" || configSilent == "2"
	for _, arg := range os.Args[1:] {
		if arg == "-s" || arg == "/s" || arg == "--silent" {
			isSilent = true
		}
	}

	if isSilent {
		dest := configDestPath
		if len(payloadZipBytes) > 0 {
			_ = os.MkdirAll(dest, 0755)
			if zipReader, err := zip.NewReader(bytes.NewReader(payloadZipBytes), int64(len(payloadZipBytes))); err == nil {
				for _, file := range zipReader.File {
					outPath := filepath.Join(dest, file.Name)
					if file.FileInfo().IsDir() {
						_ = os.MkdirAll(outPath, file.Mode())
						continue
					}
					_ = os.MkdirAll(filepath.Dir(outPath), 0755)
					if outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode()); err == nil {
						if rc, err := file.Open(); err == nil {
							_, _ = io.Copy(outFile, rc)
							rc.Close()
						}
						outFile.Close()
					}
				}
			}
		}
		if configSetupCmd != "" {
			procShellExecuteW.Call(
				0,
				uintptr(unsafe.Pointer(utf16Ptr("open"))),
				uintptr(unsafe.Pointer(utf16Ptr(configSetupCmd))),
				0,
				uintptr(unsafe.Pointer(utf16Ptr(dest))),
				1,
			)
		}
		return
	}

	className := utf16Ptr("VibeRAR_SFX_Wizard")
	hInstance, _, _ := procGetModuleHandleW.Call(0)

	var wc WNDCLASSEXW
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.Style = 0x0003
	wc.LpfnWndProc = syscall.NewCallback(wndProc)
	wc.HInstance = hInstance
	wc.HbrBackground = 16
	wc.LpszClassName = className
	procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	windowTitle := utf16Ptr(configTitle)
	hwnd, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowTitle)),
		WS_OVERLAPPEDWINDOW&^0x00050000,
		150, 150, 510, 285,
		0, 0, hInstance, 0,
	)
	hwndMain = hwnd

	hFont, _, _ := procCreateFontW.Call(
		uintptr(neg(12)), 0, 0, 0, 400, 0, 0, 0,
		0, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(utf16Ptr("Segoe UI"))),
	)
	hFontBold, _, _ := procCreateFontW.Call(
		uintptr(neg(12)), 0, 0, 0, 700, 0, 0, 0,
		0, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(utf16Ptr("Segoe UI"))),
	)

	hwndDesc, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr(configPrompt))),
		WS_CHILD|WS_VISIBLE|SS_LEFT,
		20, 15, 345, 36,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hwndDesc, WM_SETFONT, hFont, 1)

	hwndDestLbl, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("Destination folder:"))),
		WS_CHILD|WS_VISIBLE|SS_LEFT,
		20, 58, 200, 18,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hwndDestLbl, WM_SETFONT, hFont, 1)

	hwndEditBox, _, _ := procCreateWindowExW.Call(
		0x00000200,
		uintptr(unsafe.Pointer(utf16Ptr("EDIT"))),
		uintptr(unsafe.Pointer(utf16Ptr(configDestPath))),
		WS_CHILD|WS_VISIBLE|WS_BORDER|WS_TABSTOP|ES_AUTOHSCROLL,
		20, 78, 255, 23,
		hwnd, uintptr(ID_PATH_EDIT), hInstance, 0,
	)
	hwndEdit = hwndEditBox
	procSendMessageW.Call(hwndEdit, WM_SETFONT, hFont, 1)

	hwndBrowseBtn, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Browse..."))),
		WS_CHILD|WS_VISIBLE|WS_TABSTOP|BS_PUSHBUTTON,
		285, 78, 80, 24,
		hwnd, uintptr(ID_BROWSE_BTN), hInstance, 0,
	)
	hwndBrowse = hwndBrowseBtn
	procSendMessageW.Call(hwndBrowse, WM_SETFONT, hFont, 1)

	hwndStat, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("Ready to extract files."))),
		WS_CHILD|WS_VISIBLE|SS_LEFT,
		20, 115, 345, 18,
		hwnd, 0, hInstance, 0,
	)
	hwndStatus = hwndStat
	procSendMessageW.Call(hwndStatus, WM_SETFONT, hFont, 1)

	hwndProgress, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr("msctls_progress32"))),
		uintptr(unsafe.Pointer(utf16Ptr(""))),
		WS_CHILD|WS_VISIBLE|WS_BORDER,
		20, 135, 345, 18,
		hwnd, 0, hInstance, 0,
	)
	hwndProg = hwndProgress
	procSendMessageW.Call(hwndProg, PBM_SETRANGE, 0, uintptr((100<<16)))

	hwndInst, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Install"))),
		WS_CHILD|WS_VISIBLE|WS_TABSTOP|BS_DEFPUSHBUTTON,
		205, 205, 85, 26,
		hwnd, uintptr(ID_INSTALL_BTN), hInstance, 0,
	)
	hwndInstall = hwndInst
	procSendMessageW.Call(hwndInstall, WM_SETFONT, hFontBold, 1)

	hwndCanc, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Cancel"))),
		WS_CHILD|WS_VISIBLE|WS_TABSTOP|BS_PUSHBUTTON,
		300, 205, 85, 26,
		hwnd, uintptr(ID_CANCEL_BTN), hInstance, 0,
	)
	hwndCancel = hwndCanc
	procSendMessageW.Call(hwndCancel, WM_SETFONT, hFont, 1)

	hwndHlp, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Help"))),
		WS_CHILD|WS_VISIBLE|WS_TABSTOP|BS_PUSHBUTTON,
		395, 205, 85, 26,
		hwnd, uintptr(ID_HELP_BTN), hInstance, 0,
	)
	hwndHelp = hwndHlp
	procSendMessageW.Call(hwndHelp, WM_SETFONT, hFont, 1)

	procShowWindow.Call(hwnd, 1)
	procUpdateWindow.Call(hwnd)

	var msg MSG
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 || int32(ret) == -1 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}
