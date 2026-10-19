package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

//go:embed payload.exe
var embeddedPayload []byte

//go:embed appicon.ico
var embeddedIcon []byte

//go:embed viberar_shell.dll
var embeddedShellDll []byte

//go:embed VibeRAR.identity.msix
var embeddedIdentityMsix []byte

//go:embed VibeRAR.cer
var embeddedCer []byte

//go:embed Assets/StoreLogo.png
var embeddedStoreLogo []byte

//go:embed Assets/Square150x150Logo.png
var embeddedSquare150 []byte

//go:embed Assets/Square44x44Logo.png
var embeddedSquare44 []byte

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	ole32    = syscall.NewLazyDLL("ole32.dll")
	advapi32 = syscall.NewLazyDLL("advapi32.dll")

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
	procLoadIconW            = user32.NewProc("LoadIconW")
	procGetModuleHandleW     = kernel32.NewProc("GetModuleHandleW")
	procSetWindowPos         = user32.NewProc("SetWindowPos")
	procInvalidateRect       = user32.NewProc("InvalidateRect")
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
	procSHChangeNotify       = shell32.NewProc("SHChangeNotify")
	procCoTaskMemFree        = ole32.NewProc("CoTaskMemFree")

	procRegCreateKeyExW      = advapi32.NewProc("RegCreateKeyExW")
	procRegSetValueExW       = advapi32.NewProc("RegSetValueExW")
	procRegCloseKey          = advapi32.NewProc("RegCloseKey")
)

const (
	WS_OVERLAPPEDWINDOW = 0x00CF0000
	WS_VISIBLE          = 0x10000000
	WS_CHILD            = 0x40000000
	WS_BORDER           = 0x00800000
	WS_TABSTOP          = 0x00010000
	ES_AUTOHSCROLL      = 0x0080
	ES_MULTILINE        = 0x0004
	ES_READONLY         = 0x0800
	WS_VSCROLL          = 0x00200000
	BS_PUSHBUTTON       = 0x0000
	BS_DEFPUSHBUTTON    = 0x0001
	BS_CHECKBOX         = 0x0002
	BS_AUTOCHECKBOX     = 0x0003
	BS_GROUPBOX         = 0x0007
	SS_LEFT             = 0x0000

	WM_DESTROY  = 0x0002
	WM_PAINT    = 0x000F
	WM_COMMAND  = 0x0111
	WM_SETFONT  = 0x0030
	WM_SETICON  = 0x0080

	BM_GETCHECK = 0x00F0
	BM_SETCHECK = 0x00F1
	BST_CHECKED = 1

	PBM_SETRANGE = 0x0401
	PBM_SETPOS   = 0x0402

	HKEY_CURRENT_USER = 0x80000001
	KEY_WRITE         = 0x00020006
	REG_SZ            = 1

	// Control IDs
	ID_INSTALL_BTN   = 2001
	ID_HELP_BTN      = 2002
	ID_CANCEL_BTN    = 2003
	ID_BROWSE_BTN    = 2004
	ID_PATH_EDIT     = 2005

	ID_STEP2_OK      = 3001
	ID_STEP2_CANCEL  = 3002
	ID_STEP2_HELP    = 3003
	ID_TOGGLE_ALL    = 3004

	ID_STEP3_DONE    = 4001
	ID_STEP3_RUN     = 4002
	ID_STEP3_ORDER   = 4003
	ID_STEP3_HELP    = 4004
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
	hwndMain       uintptr
	hInstance      uintptr
	hFontUi        uintptr
	hFontBold      uintptr
	currentStep    = 1
	installPath    = "C:\\Program Files\\VibeRAR"
	step1Controls  []uintptr
	step2Controls  []uintptr
	step3Controls  []uintptr

	// Step 1 handles
	hwndPathEdit   uintptr
	hwndProgBar    uintptr
	hwndStatusText uintptr
	hwndInstBtn    uintptr

	// Step 2 handles
	assocCheckboxes []uintptr
	chkDesktop      uintptr
	chkStartMenu    uintptr
	chkShellInteg   uintptr
	chkCascaded     uintptr
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

func setRegString(subKey, valueName, value string) error {
	var hKey uintptr
	var disp uint32
	subKeyPtr := utf16Ptr(subKey)
	ret, _, _ := procRegCreateKeyExW.Call(
		HKEY_CURRENT_USER,
		uintptr(unsafe.Pointer(subKeyPtr)),
		0, 0, 0,
		KEY_WRITE,
		0,
		uintptr(unsafe.Pointer(&hKey)),
		uintptr(unsafe.Pointer(&disp)),
	)
	if ret != 0 {
		return fmt.Errorf("RegCreateKeyExW error: %d", ret)
	}
	defer procRegCloseKey.Call(hKey)

	valBytes, _ := syscall.UTF16FromString(value)
	var valNamePtr *uint16
	if valueName != "" {
		valNamePtr = utf16Ptr(valueName)
	}
	procRegSetValueExW.Call(
		hKey,
		uintptr(unsafe.Pointer(valNamePtr)),
		0,
		REG_SZ,
		uintptr(unsafe.Pointer(&valBytes[0])),
		uintptr(len(valBytes)*2),
	)
	return nil
}

func createShortcut(shortcutPath, targetPath, workingDir, description string) {
	vbs := fmt.Sprintf(`Set w = CreateObject("WScript.Shell")
Set s = w.CreateShortcut("%s")
s.TargetPath = "%s"
s.WorkingDirectory = "%s"
s.Description = "%s"
s.Save
`, shortcutPath, targetPath, workingDir, description)

	tmpVbs := filepath.Join(os.TempDir(), "viberar_sc.vbs")
	if err := os.WriteFile(tmpVbs, []byte(vbs), 0644); err == nil {
		cmd := exec.Command("wscript.exe", "//B", "//Nologo", tmpVbs)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000, // CREATE_NO_WINDOW
		}
		_ = cmd.Run()
		_ = os.Remove(tmpVbs)
	}
}

func drawVibeRARBooks(hdc uintptr, x, y, w, h int) {
	hWhiteBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(255, 255, 255)))
	hBorderPen, _, _ := procCreatePen.Call(0, 1, uintptr(rgb(205, 215, 225)))
	oldBrush, _, _ := procSelectObject.Call(hdc, hWhiteBrush)
	oldPen, _, _ := procSelectObject.Call(hdc, hBorderPen)

	// Outer card frame
	procRoundRect.Call(hdc, uintptr(x), uintptr(y), uintptr(x+w), uintptr(y+h), 8, 8)

	// Top title bar banner
	hHeaderBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(240, 245, 252)))
	procSelectObject.Call(hdc, hHeaderBrush)
	procRoundRect.Call(hdc, uintptr(x+1), uintptr(y+1), uintptr(x+w-1), uintptr(y+28), 6, 6)

	procSetBkMode.Call(hdc, 1) // TRANSPARENT
	procSetTextColor.Call(hdc, uintptr(rgb(20, 70, 160)))
	hTitleFont, _, _ := procCreateFontW.Call(
		uintptr(neg(13)), 0, 0, 0, 700, 0, 0, 0,
		0, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(utf16Ptr("Segoe UI"))),
	)
	oldFont, _, _ := procSelectObject.Call(hdc, hTitleFont)
	titleStr := utf16Ptr("VibeRAR")
	procTextOutW.Call(hdc, uintptr(x+(w-56)/2), uintptr(y+6), uintptr(unsafe.Pointer(titleStr)), 7)

	procSelectObject.Call(hdc, oldFont)
	procDeleteObject.Call(hTitleFont)
	procDeleteObject.Call(hHeaderBrush)

	// Shadow under the books
	hShadowBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(225, 230, 236)))
	procSelectObject.Call(hdc, hShadowBrush)
	hNullPen, _, _ := procGetStockObject.Call(8) // NULL_PEN
	procSelectObject.Call(hdc, hNullPen)
	procRoundRect.Call(hdc, uintptr(x+14), uintptr(y+162), uintptr(x+w-12), uintptr(y+174), 10, 10)
	procDeleteObject.Call(hShadowBrush)

	// Single book renderer
	drawSingleBook := func(bx, by, bw, bh int, spineCol, coverCol, ribCol uint32) {
		hCoverBrush, _, _ := procCreateSolidBrush.Call(uintptr(coverCol))
		procSelectObject.Call(hdc, hCoverBrush)
		procRoundRect.Call(hdc, uintptr(bx), uintptr(by), uintptr(bx+bw), uintptr(by+bh), 4, 4)
		procDeleteObject.Call(hCoverBrush)

		// Darker left spine
		hSpineBrush, _, _ := procCreateSolidBrush.Call(uintptr(spineCol))
		procSelectObject.Call(hdc, hSpineBrush)
		procRoundRect.Call(hdc, uintptr(bx), uintptr(by), uintptr(bx+20), uintptr(by+bh), 4, 4)
		procDeleteObject.Call(hSpineBrush)

		// Page edge
		hPageBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(250, 246, 238)))
		procSelectObject.Call(hdc, hPageBrush)
		procRectangle.Call(hdc, uintptr(bx+bw-12), uintptr(by+3), uintptr(bx+bw-3), uintptr(by+bh-3))
		procDeleteObject.Call(hPageBrush)

		// Spine ribs (gold foil lines)
		hRibPen, _, _ := procCreatePen.Call(0, 2, uintptr(ribCol))
		procSelectObject.Call(hdc, hRibPen)
		procMoveToEx.Call(hdc, uintptr(bx+2), uintptr(by+6), 0)
		procLineTo.Call(hdc, uintptr(bx+18), uintptr(by+6))
		procMoveToEx.Call(hdc, uintptr(bx+2), uintptr(by+bh-7), 0)
		procLineTo.Call(hdc, uintptr(bx+18), uintptr(by+bh-7))
		procDeleteObject.Call(hRibPen)
	}

	// 1. Bottom Book (Crimson Red)
	drawSingleBook(x+18, y+128, w-34, 28, rgb(130, 20, 30), rgb(185, 35, 45), rgb(235, 195, 60))

	// 2. Middle Book (Emerald Green)
	drawSingleBook(x+16, y+94, w-34, 28, rgb(20, 95, 40), rgb(35, 145, 60), rgb(235, 195, 60))

	// 3. Top Book (Royal Blue)
	drawSingleBook(x+14, y+60, w-34, 28, rgb(15, 55, 125), rgb(28, 85, 185), rgb(235, 195, 60))

	// "RAR" text on top book cover
	procSetTextColor.Call(hdc, uintptr(rgb(255, 225, 75))) // Gold
	hRarFont, _, _ := procCreateFontW.Call(
		uintptr(neg(12)), 0, 0, 0, 800, 0, 0, 0,
		0, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(utf16Ptr("Arial"))),
	)
	procSelectObject.Call(hdc, hRarFont)
	rarStr := utf16Ptr("RAR")
	procTextOutW.Call(hdc, uintptr(x+40), uintptr(y+66), uintptr(unsafe.Pointer(rarStr)), 3)
	procDeleteObject.Call(hRarFont)

	// Leather Belt wrapping the 3 books
	strapX := x + 52
	strapW := 15
	strapTop := y + 50
	strapBot := y + 162

	hStrapBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(105, 50, 20)))
	procSelectObject.Call(hdc, hStrapBrush)
	procSelectObject.Call(hdc, hNullPen)
	procRectangle.Call(hdc, uintptr(strapX), uintptr(strapTop), uintptr(strapX+strapW), uintptr(strapBot))
	procDeleteObject.Call(hStrapBrush)

	// Stitching on belt
	hStitchPen, _, _ := procCreatePen.Call(0, 1, uintptr(rgb(220, 180, 80)))
	procSelectObject.Call(hdc, hStitchPen)
	for sy := strapTop + 3; sy < strapBot-3; sy += 6 {
		procMoveToEx.Call(hdc, uintptr(strapX+1), uintptr(sy), 0)
		procLineTo.Call(hdc, uintptr(strapX+2), uintptr(sy))
		procMoveToEx.Call(hdc, uintptr(strapX+strapW-2), uintptr(sy), 0)
		procLineTo.Call(hdc, uintptr(strapX+strapW-1), uintptr(sy))
	}
	procDeleteObject.Call(hStitchPen)

	// Gold Buckle
	buckleY := y + 102
	buckleW := 21
	buckleH := 16
	buckleX := strapX - (buckleW-strapW)/2

	hBuckleBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(240, 200, 45)))
	procSelectObject.Call(hdc, hBuckleBrush)
	procRoundRect.Call(hdc, uintptr(buckleX), uintptr(buckleY), uintptr(buckleX+buckleW), uintptr(buckleY+buckleH), 4, 4)
	procDeleteObject.Call(hBuckleBrush)

	// Inner leather cutout
	hHoleBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(65, 30, 10)))
	procSelectObject.Call(hdc, hHoleBrush)
	procRectangle.Call(hdc, uintptr(buckleX+4), uintptr(buckleY+3), uintptr(buckleX+buckleW-4), uintptr(buckleY+buckleH-3))
	procDeleteObject.Call(hHoleBrush)

	// Buckle pin
	hProngPen, _, _ := procCreatePen.Call(0, 2, uintptr(rgb(255, 230, 90)))
	procSelectObject.Call(hdc, hProngPen)
	procMoveToEx.Call(hdc, uintptr(buckleX+buckleW/2), uintptr(buckleY+2), 0)
	procLineTo.Call(hdc, uintptr(buckleX+buckleW/2), uintptr(buckleY+2))
	procDeleteObject.Call(hProngPen)

	// Subtitle text under books
	procSetTextColor.Call(hdc, uintptr(rgb(110, 120, 135)))
	hSubFont, _, _ := procCreateFontW.Call(
		uintptr(neg(10)), 0, 0, 0, 400, 0, 0, 0,
		0, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(utf16Ptr("Segoe UI"))),
	)
	procSelectObject.Call(hdc, hSubFont)
	subStr := utf16Ptr("v7.00 Pro")
	procTextOutW.Call(hdc, uintptr(x+(w-54)/2), uintptr(y+182), uintptr(unsafe.Pointer(subStr)), 9)
	procDeleteObject.Call(hSubFont)

	// Restore original GDI objects
	procSelectObject.Call(hdc, oldBrush)
	procSelectObject.Call(hdc, oldPen)
	procDeleteObject.Call(hWhiteBrush)
	procDeleteObject.Call(hBorderPen)
}

func showStep1() {
	currentStep = 1
	for _, h := range step2Controls {
		procShowWindow.Call(h, 0)
	}
	for _, h := range step3Controls {
		procShowWindow.Call(h, 0)
	}

	procSetWindowTextW.Call(hwndMain, uintptr(unsafe.Pointer(utf16Ptr("VibeRAR 7.00"))))
	procSetWindowPos.Call(hwndMain, 0, 150, 120, 540, 365, 0x0004)

	for _, h := range step1Controls {
		procShowWindow.Call(h, 1)
	}
	procInvalidateRect.Call(hwndMain, 0, 1)
}

func showStep2() {
	currentStep = 2
	for _, h := range step1Controls {
		procShowWindow.Call(h, 0)
	}
	for _, h := range step3Controls {
		procShowWindow.Call(h, 0)
	}

	procSetWindowTextW.Call(hwndMain, uintptr(unsafe.Pointer(utf16Ptr("VibeRAR Setup"))))
	procSetWindowPos.Call(hwndMain, 0, 150, 100, 530, 435, 0x0004)

	for _, h := range step2Controls {
		procShowWindow.Call(h, 1)
	}
	procInvalidateRect.Call(hwndMain, 0, 1)
}

func showStep3() {
	currentStep = 3
	for _, h := range step1Controls {
		procShowWindow.Call(h, 0)
	}
	for _, h := range step2Controls {
		procShowWindow.Call(h, 0)
	}

	procSetWindowTextW.Call(hwndMain, uintptr(unsafe.Pointer(utf16Ptr("VibeRAR Setup"))))
	procSetWindowPos.Call(hwndMain, 0, 180, 150, 510, 275, 0x0004)

	for _, h := range step3Controls {
		procShowWindow.Call(h, 1)
	}
	procInvalidateRect.Call(hwndMain, 0, 1)
}

func browseFolder() {
	// Native Win32 Folder Picker (ZERO console windows / ZERO cmd.exe)
	var bi BROWSEINFOW
	bi.HwndOwner = hwndMain
	bi.LpszTitle = utf16Ptr("Select destination folder for VibeRAR installation:")
	bi.UlFlags = 0x00000041 // BIF_RETURNONLYFSDIRS | BIF_NEWDIALOGSTYLE

	pidl, _, _ := procSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(&bi)))
	if pidl != 0 {
		var pathBuf [260]uint16
		procSHGetPathFromIDListW.Call(pidl, uintptr(unsafe.Pointer(&pathBuf[0])))
		procCoTaskMemFree.Call(pidl)
		newPath := syscall.UTF16ToString(pathBuf[:])
		if strings.TrimSpace(newPath) != "" {
			installPath = strings.TrimSpace(newPath)
			procSetWindowTextW.Call(hwndPathEdit, uintptr(unsafe.Pointer(utf16Ptr(installPath))))
		}
	}
}

func performInstallation() {
	procSetWindowTextW.Call(hwndInstBtn, uintptr(unsafe.Pointer(utf16Ptr("Installing..."))))

	var buf [512]uint16
	procGetWindowTextW.Call(hwndPathEdit, uintptr(unsafe.Pointer(&buf[0])), 512)
	dest := syscall.UTF16ToString(buf[:])
	if strings.TrimSpace(dest) != "" {
		installPath = strings.TrimSpace(dest)
	}

	go func() {
		if err := os.MkdirAll(installPath, 0755); err != nil {
			userProfile := os.Getenv("USERPROFILE")
			installPath = filepath.Join(userProfile, "AppData", "Local", "Programs", "VibeRAR")
			os.MkdirAll(installPath, 0755)
		}

		stages := []struct {
			progress int
			status   string
			delay    time.Duration
		}{
			{15, "Extracting: VibeRAR.exe", 150 * time.Millisecond},
			{35, "Extracting: VibeRAR_Portable.exe", 200 * time.Millisecond},
			{55, "Extracting: compression_engine.dll & assets", 180 * time.Millisecond},
			{75, "Configuring hardware benchmark engine...", 200 * time.Millisecond},
			{90, "Setting up Windows shell integration...", 220 * time.Millisecond},
			{100, "Installation completed successfully!", 150 * time.Millisecond},
		}

		for _, s := range stages {
			procSendMessageW.Call(hwndProgBar, PBM_SETPOS, uintptr(s.progress), 0)
			procSetWindowTextW.Call(hwndStatusText, uintptr(unsafe.Pointer(utf16Ptr(s.status))))
			time.Sleep(s.delay)
		}

		// Write payload files
		targetExe := filepath.Join(installPath, "VibeRAR.exe")
		targetPortable := filepath.Join(installPath, "VibeRAR_Portable.exe")
		targetIco := filepath.Join(installPath, "appicon.ico")

		if len(embeddedPayload) > 0 {
			_ = os.WriteFile(targetExe, embeddedPayload, 0755)
			_ = os.WriteFile(targetPortable, embeddedPayload, 0755)
		} else {
			cwd, _ := os.Getwd()
			srcExe := filepath.Join(cwd, "VibeRAR.exe")
			if data, err := os.ReadFile(srcExe); err == nil {
				_ = os.WriteFile(targetExe, data, 0755)
				_ = os.WriteFile(targetPortable, data, 0755)
			}
		}

		if len(embeddedIcon) > 0 {
			_ = os.WriteFile(targetIco, embeddedIcon, 0644)
		} else {
			cwd, _ := os.Getwd()
			srcIco := filepath.Join(cwd, "appicon.ico")
			if data, err := os.ReadFile(srcIco); err == nil {
				_ = os.WriteFile(targetIco, data, 0644)
			}
		}

		// Write Windows 11 Modern Context Menu (Sparse Package & COM DLL)
		targetShell := filepath.Join(installPath, "viberar_shell.dll")
		targetMsix := filepath.Join(installPath, "VibeRAR.identity.msix")
		targetCer := filepath.Join(installPath, "VibeRAR.cer")
		assetsDir := filepath.Join(installPath, "Assets")
		_ = os.MkdirAll(assetsDir, 0755)

		if len(embeddedShellDll) > 0 {
			_ = os.WriteFile(targetShell, embeddedShellDll, 0755)
		}
		if len(embeddedIdentityMsix) > 0 {
			_ = os.WriteFile(targetMsix, embeddedIdentityMsix, 0644)
		}
		if len(embeddedCer) > 0 {
			_ = os.WriteFile(targetCer, embeddedCer, 0644)
		}
		if len(embeddedStoreLogo) > 0 {
			_ = os.WriteFile(filepath.Join(assetsDir, "StoreLogo.png"), embeddedStoreLogo, 0644)
		}
		if len(embeddedSquare150) > 0 {
			_ = os.WriteFile(filepath.Join(assetsDir, "Square150x150Logo.png"), embeddedSquare150, 0644)
		}
		if len(embeddedSquare44) > 0 {
			_ = os.WriteFile(filepath.Join(assetsDir, "Square44x44Logo.png"), embeddedSquare44, 0644)
		}

		time.Sleep(200 * time.Millisecond)
		showStep2()
	}()
}

func applyStep2AndContinue() {
	chkDeskRet, _, _ := procSendMessageW.Call(chkDesktop, BM_GETCHECK, 0, 0)
	isDesktop := chkDeskRet == BST_CHECKED

	chkMenuRet, _, _ := procSendMessageW.Call(chkStartMenu, BM_GETCHECK, 0, 0)
	isStartMenu := chkMenuRet == BST_CHECKED

	chkShellRet, _, _ := procSendMessageW.Call(chkShellInteg, BM_GETCHECK, 0, 0)
	isShellInteg := chkShellRet == BST_CHECKED

	targetExe := filepath.Join(installPath, "VibeRAR.exe")
	iconIco := filepath.Join(installPath, "appicon.ico")

	if isDesktop {
		desktopDir := filepath.Join(os.Getenv("USERPROFILE"), "Desktop")
		shortcutPath := filepath.Join(desktopDir, "VibeRAR.lnk")
		createShortcut(shortcutPath, targetExe, installPath, "VibeRAR - High Performance Archive Manager")
	}

	if isStartMenu {
		startMenuDir := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs")
		shortcutPath := filepath.Join(startMenuDir, "VibeRAR.lnk")
		createShortcut(shortcutPath, targetExe, installPath, "VibeRAR - High Performance Archive Manager")
	}

	// Register in Windows App Paths using pure Win32 Registry APIs (ZERO command prompt)
	regSubKey := "Software\\Microsoft\\Windows\\CurrentVersion\\App Paths\\VibeRAR.exe"
	_ = setRegString(regSubKey, "", targetExe)
	_ = setRegString(regSubKey, "Path", installPath)

	if isShellInteg {
		// 1. Direct Context Menus on Directory (Folders)
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Add", "", "Add to archive with VibeRAR...")
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Add", "MUIVerb", "Add to archive with VibeRAR...")
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Add", "Icon", iconIco)
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Add\\command", "", fmt.Sprintf("\"%s\" --add \"%%1\"", targetExe))

		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Open", "", "Open with VibeRAR")
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Open", "MUIVerb", "Open with VibeRAR")
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Open", "Icon", iconIco)
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Open\\command", "", fmt.Sprintf("\"%s\" --open \"%%1\"", targetExe))

		// 2. Direct Context Menus on Folder
		_ = setRegString("Software\\Classes\\Folder\\shell\\VibeRAR.Add", "", "Add to archive with VibeRAR...")
		_ = setRegString("Software\\Classes\\Folder\\shell\\VibeRAR.Add", "MUIVerb", "Add to archive with VibeRAR...")
		_ = setRegString("Software\\Classes\\Folder\\shell\\VibeRAR.Add", "Icon", iconIco)
		_ = setRegString("Software\\Classes\\Folder\\shell\\VibeRAR.Add\\command", "", fmt.Sprintf("\"%s\" --add \"%%1\"", targetExe))

		_ = setRegString("Software\\Classes\\Folder\\shell\\VibeRAR.Open", "", "Open with VibeRAR")
		_ = setRegString("Software\\Classes\\Folder\\shell\\VibeRAR.Open", "MUIVerb", "Open with VibeRAR")
		_ = setRegString("Software\\Classes\\Folder\\shell\\VibeRAR.Open", "Icon", iconIco)
		_ = setRegString("Software\\Classes\\Folder\\shell\\VibeRAR.Open\\command", "", fmt.Sprintf("\"%s\" --open \"%%1\"", targetExe))

		// 3. Directory Background (Empty space)
		_ = setRegString("Software\\Classes\\Directory\\Background\\shell\\VibeRAR", "", "Open VibeRAR here")
		_ = setRegString("Software\\Classes\\Directory\\Background\\shell\\VibeRAR", "MUIVerb", "Open VibeRAR here")
		_ = setRegString("Software\\Classes\\Directory\\Background\\shell\\VibeRAR", "Icon", iconIco)
		_ = setRegString("Software\\Classes\\Directory\\Background\\shell\\VibeRAR\\command", "", fmt.Sprintf("\"%s\" --open \"%%V\"", targetExe))

		// 4. All Files (*)
		_ = setRegString("Software\\Classes\\*\\shell\\VibeRAR.Add", "", "Add to archive with VibeRAR...")
		_ = setRegString("Software\\Classes\\*\\shell\\VibeRAR.Add", "MUIVerb", "Add to archive with VibeRAR...")
		_ = setRegString("Software\\Classes\\*\\shell\\VibeRAR.Add", "Icon", iconIco)
		_ = setRegString("Software\\Classes\\*\\shell\\VibeRAR.Add\\command", "", fmt.Sprintf("\"%s\" --add \"%%1\"", targetExe))

		_ = setRegString("Software\\Classes\\*\\shell\\VibeRAR.Open", "", "Open with VibeRAR")
		_ = setRegString("Software\\Classes\\*\\shell\\VibeRAR.Open", "MUIVerb", "Open with VibeRAR")
		_ = setRegString("Software\\Classes\\*\\shell\\VibeRAR.Open", "Icon", iconIco)
		_ = setRegString("Software\\Classes\\*\\shell\\VibeRAR.Open\\command", "", fmt.Sprintf("\"%s\" --open \"%%1\"", targetExe))

		// 5. Shared Context Menu Definition
		menuSubKey := "Software\\Classes\\VibeRAR.ContextMenu\\shell"
		_ = setRegString(menuSubKey+"\\01_add", "MUIVerb", "Add to archive...")
		_ = setRegString(menuSubKey+"\\01_add", "Icon", iconIco)
		_ = setRegString(menuSubKey+"\\01_add\\command", "", fmt.Sprintf("\"%s\" --add \"%%1\"", targetExe))

		_ = setRegString(menuSubKey+"\\02_fast_rar", "MUIVerb", "Add to archive (Fast .rar)")
		_ = setRegString(menuSubKey+"\\02_fast_rar", "Icon", iconIco)
		_ = setRegString(menuSubKey+"\\02_fast_rar\\command", "", fmt.Sprintf("\"%s\" --add-fast-rar \"%%1\"", targetExe))

		_ = setRegString(menuSubKey+"\\03_fast_zip", "MUIVerb", "Add to archive (Fast .zip)")
		_ = setRegString(menuSubKey+"\\03_fast_zip", "Icon", iconIco)
		_ = setRegString(menuSubKey+"\\03_fast_zip\\command", "", fmt.Sprintf("\"%s\" --add-fast-zip \"%%1\"", targetExe))

		_ = setRegString(menuSubKey+"\\04_extract", "MUIVerb", "Extract files...")
		_ = setRegString(menuSubKey+"\\04_extract", "Icon", iconIco)
		_ = setRegString(menuSubKey+"\\04_extract\\command", "", fmt.Sprintf("\"%s\" --extract \"%%1\"", targetExe))

		_ = setRegString(menuSubKey+"\\05_extract_here", "MUIVerb", "Extract Here")
		_ = setRegString(menuSubKey+"\\05_extract_here", "Icon", iconIco)
		_ = setRegString(menuSubKey+"\\05_extract_here\\command", "", fmt.Sprintf("\"%s\" --extract-here \"%%1\"", targetExe))

		_ = setRegString(menuSubKey+"\\06_open", "MUIVerb", "Open with VibeRAR")
		_ = setRegString(menuSubKey+"\\06_open", "Icon", iconIco)
		_ = setRegString(menuSubKey+"\\06_open\\command", "", fmt.Sprintf("\"%s\" --open \"%%1\"", targetExe))

		// 6. Hook Cascaded Submenu on Directory, Folder, and All Files (*)
		for _, baseKey := range []string{"Software\\Classes\\Directory\\shell\\VibeRAR", "Software\\Classes\\Folder\\shell\\VibeRAR", "Software\\Classes\\*\\shell\\VibeRAR"} {
			_ = setRegString(baseKey, "", "")
			_ = setRegString(baseKey, "MUIVerb", "VibeRAR")
			_ = setRegString(baseKey, "Icon", iconIco)
			_ = setRegString(baseKey, "ExtendedSubCommandsKey", "VibeRAR.ContextMenu")
		}

		// 7. Register Modern Windows 11 Context Menu (MSIX Sparse Package Integration)
		go registerModernWindows11ContextMenu(installPath)
	}

	// Register File Associations
	archiveProgID := "Software\\Classes\\VibeRAR.Archive"
	_ = setRegString(archiveProgID, "", "VibeRAR Archive")
	_ = setRegString(archiveProgID+"\\DefaultIcon", "", fmt.Sprintf("\"%s\",0", iconIco))
	_ = setRegString(archiveProgID+"\\shell\\open\\command", "", fmt.Sprintf("\"%s\" --open \"%%1\"", targetExe))
	_ = setRegString(archiveProgID+"\\shell\\VibeRAR.Extract", "", "Extract files...")
	_ = setRegString(archiveProgID+"\\shell\\VibeRAR.Extract", "Icon", iconIco)
	_ = setRegString(archiveProgID+"\\shell\\VibeRAR.Extract\\command", "", fmt.Sprintf("\"%s\" --extract \"%%1\"", targetExe))
	_ = setRegString(archiveProgID+"\\shell\\VibeRAR.ExtractHere", "", "Extract Here")
	_ = setRegString(archiveProgID+"\\shell\\VibeRAR.ExtractHere", "Icon", iconIco)
	_ = setRegString(archiveProgID+"\\shell\\VibeRAR.ExtractHere\\command", "", fmt.Sprintf("\"%s\" --extract-here \"%%1\"", targetExe))

	extList := []string{"rar", "zip", "7z", "arj", "bz2", "cab", "gz", "iso", "jar", "lz", "lzh", "tar", "uue", "xz", "z", "001"}
	for i, chk := range assocCheckboxes {
		state, _, _ := procSendMessageW.Call(chk, BM_GETCHECK, 0, 0)
		if state == BST_CHECKED && i < len(extList) {
			ext := "." + extList[i]
			_ = setRegString("Software\\Classes\\"+ext, "", "VibeRAR.Archive")
		}
	}

	// Notify Windows Explorer of association & shell verb changes
	procSHChangeNotify.Call(uintptr(0x08000000), uintptr(0x0000|0x1000), 0, 0)

	showStep3()
}

func runPowerShellHidden(psCmd string) {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", psCmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
	_ = cmd.Run()
}

func registerModernWindows11ContextMenu(installDir string) {
	cerPath := filepath.Join(installDir, "VibeRAR.cer")
	msixPath := filepath.Join(installDir, "VibeRAR.identity.msix")

	if _, err := os.Stat(msixPath); err != nil {
		return
	}

	// 1. Import certificate into TrustedPeople store (silent)
	if _, err := os.Stat(cerPath); err == nil {
		cmdCert := fmt.Sprintf(`Import-Certificate -FilePath '%s' -CertStoreLocation 'Cert:\CurrentUser\TrustedPeople' -ErrorAction SilentlyContinue; Import-Certificate -FilePath '%s' -CertStoreLocation 'Cert:\LocalMachine\TrustedPeople' -ErrorAction SilentlyContinue`, cerPath, cerPath)
		runPowerShellHidden(cmdCert)
	}

	// 2. Register Sparse Package with External Location in user context
	cmdAppx := fmt.Sprintf(`Add-AppxPackage -Path '%s' -ExternalLocation '%s' -ErrorAction SilentlyContinue`, msixPath, installDir)
	runPowerShellHidden(cmdAppx)

	// 3. Notify Windows Explorer
	procSHChangeNotify.Call(uintptr(0x08000000), uintptr(0x0000|0x1000), 0, 0)
}


func launchApplication() {
	targetExe := filepath.Join(installPath, "VibeRAR.exe")
	if _, err := os.Stat(targetExe); os.IsNotExist(err) {
		targetExe = filepath.Join(installPath, "VibeRAR_Portable.exe")
	}

	// Launch via ShellExecuteW (ZERO cmd.exe / NO command prompt)
	procShellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr("open"))),
		uintptr(unsafe.Pointer(utf16Ptr(targetExe))),
		0, 0, 1,
	)
	procPostQuitMessage.Call(0)
}

func toggleAllAssociations() {
	if len(assocCheckboxes) == 0 {
		return
	}
	firstState, _, _ := procSendMessageW.Call(assocCheckboxes[0], BM_GETCHECK, 0, 0)
	newState := uintptr(BST_CHECKED)
	if firstState == BST_CHECKED {
		newState = 0
	}
	for _, h := range assocCheckboxes {
		procSendMessageW.Call(h, BM_SETCHECK, newState, 0)
	}
}

func wndProc(hwnd uintptr, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		wmId := loword(uint32(wParam))
		switch wmId {
		// Step 1
		case ID_INSTALL_BTN:
			performInstallation()
		case ID_BROWSE_BTN:
			browseFolder()
		case ID_CANCEL_BTN, ID_STEP2_CANCEL:
			procPostQuitMessage.Call(0)
		case ID_HELP_BTN, ID_STEP2_HELP, ID_STEP3_HELP:
			procMessageBoxW.Call(
				hwnd,
				uintptr(unsafe.Pointer(utf16Ptr("VibeRAR v1.0\nA fast, modern desktop archiver.\n\nCreated by Ahmed Samir\nOfficial Site: https://ahmedsamir.com\nGitHub: https://github.com/.../viberar\nsponsored by https://www.0webhost.com/"))),
				uintptr(unsafe.Pointer(utf16Ptr("About VibeRAR"))),
				0x40,
			)

		// Step 2
		case ID_TOGGLE_ALL:
			toggleAllAssociations()
		case ID_STEP2_OK:
			applyStep2AndContinue()

		// Step 3
		case ID_STEP3_RUN:
			launchApplication()
		case ID_STEP3_DONE:
			launchApplication()
		case ID_STEP3_ORDER:
			procMessageBoxW.Call(
				hwnd,
				uintptr(unsafe.Pointer(utf16Ptr("VibeRAR is 100% FREE forever!\n\nNo license key or purchase order required. Enjoy unlimited features!"))),
				uintptr(unsafe.Pointer(utf16Ptr("Evaluation Notice"))),
				0x40,
			)
		}
		return 0

	case WM_PAINT:
		var ps PAINTSTRUCT
		hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))

		if currentStep == 1 {
			drawVibeRARBooks(hdc, 395, 15, 115, 235)
		} else if currentStep == 3 {
			drawVibeRARBooks(hdc, 18, 15, 115, 205)
		}

		procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		return 0

	case WM_DESTROY:
		procPostQuitMessage.Call(0)
		return 0
	}

	ret, _, _ := procDefWindowProcW.Call(hwnd, uintptr(msg), wParam, lParam)
	return ret
}

func main() {
	runtime.LockOSThread()

	className := utf16Ptr("VibeRAR_Installer_Wizard")
	hInstanceVal, _, _ := procGetModuleHandleW.Call(0)
	hInstance = hInstanceVal

	var wc WNDCLASSEXW
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.Style = 0x0003 // CS_HREDRAW | CS_VREDRAW
	wc.LpfnWndProc = syscall.NewCallback(wndProc)
	wc.HInstance = hInstance
	wc.HbrBackground = 16 // COLOR_BTNFACE + 1
	wc.LpszClassName = className
	hAppIcon, _, _ := procLoadIconW.Call(hInstance, uintptr(1))
	if hAppIcon != 0 {
		wc.HIcon = hAppIcon
		wc.HIconSm = hAppIcon
	}
	procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	// Main window
	hwnd, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(utf16Ptr("VibeRAR 7.00"))),
		WS_OVERLAPPEDWINDOW&^0x00050000, // No Maximize / Resize
		150, 120, 540, 365,
		0, 0, hInstance, 0,
	)
	hwndMain = hwnd
	if hAppIcon != 0 {
		procSendMessageW.Call(hwndMain, WM_SETICON, 0, hAppIcon)
		procSendMessageW.Call(hwndMain, WM_SETICON, 1, hAppIcon)
	}

	// Create Modern Segoe UI Fonts
	hFontUi, _, _ = procCreateFontW.Call(
		uintptr(neg(12)), 0, 0, 0, 400, 0, 0, 0,
		0, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(utf16Ptr("Segoe UI"))),
	)
	hFontBold, _, _ = procCreateFontW.Call(
		uintptr(neg(12)), 0, 0, 0, 700, 0, 0, 0,
		0, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(utf16Ptr("Segoe UI"))),
	)

	// ==========================================
	// STEP 1 CONTROLS (VibeRAR 7.00 Installation)
	// ==========================================
	hInfo, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("Please note that VibeRAR is free software with unlimited features."))),
		WS_CHILD|SS_LEFT, 20, 15, 365, 18,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hInfo, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hInfo)

	// EULA / License edit box
	licenseText := `VibeRAR 7.00 Pro (64-bit)
High-Performance Archive Manager & Compression Suite
===================================================
Software Creator: Ahmed Samir Hany
Official Website: https://ahmedsamir.com

Sponsored Web Hosting We Honor:
https://www.0webhost.com/ - Buy hosting at cheap prices
===================================================

License Agreement:
VibeRAR is distributed as 100% FREE and UNLIMITED software.
You may use this software without any trial expiration, registration prompts, or feature limitations for personal, commercial, and educational purposes.

Key Capabilities:
• Ultra-fast decompression of RAR (RAR4, RAR5, RAR7), ZIP, 7Z, TAR, GZ, ISO, CAB, BZ2, XZ, and Z formats.
• High-ratio compression to ZIP64 and 7z with full multi-threading and AES-256 encryption.
• Dual-density toolbars (Large 48x36 and Small 24x24) matching authentic ergonomics.
• Integrated hardware multi-threaded CPU benchmark utility.
• Self-extracting archive (SFX) creation engine with custom destination paths and setup scripts.
• Fully self-contained offline architecture compatible with Windows 7, 8, 10, and 11.

Click "Install" to begin copying VibeRAR to your system, or click "Browse" to choose an alternative folder.`

	hEula, _, _ := procCreateWindowExW.Call(
		0x00000200, uintptr(unsafe.Pointer(utf16Ptr("EDIT"))),
		uintptr(unsafe.Pointer(utf16Ptr(licenseText))),
		WS_CHILD|WS_BORDER|ES_MULTILINE|ES_READONLY|WS_VSCROLL,
		20, 38, 365, 140,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hEula, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hEula)

	// Destination folder label
	hDestLbl, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("Destination folder"))),
		WS_CHILD|SS_LEFT, 20, 190, 200, 18,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hDestLbl, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hDestLbl)

	// Destination path edit
	hwndPathEdit, _, _ = procCreateWindowExW.Call(
		0x00000200, uintptr(unsafe.Pointer(utf16Ptr("EDIT"))),
		uintptr(unsafe.Pointer(utf16Ptr(installPath))),
		WS_CHILD|WS_BORDER|WS_TABSTOP|ES_AUTOHSCROLL,
		20, 210, 275, 24,
		hwnd, uintptr(ID_PATH_EDIT), hInstance, 0,
	)
	procSendMessageW.Call(hwndPathEdit, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hwndPathEdit)

	// Browse button
	hBrowse, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Browse..."))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		305, 210, 80, 24,
		hwnd, uintptr(ID_BROWSE_BTN), hInstance, 0,
	)
	procSendMessageW.Call(hBrowse, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hBrowse)

	// Status text (hidden until installing)
	hwndStatusText, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr(""))),
		WS_CHILD|SS_LEFT, 20, 240, 365, 16,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hwndStatusText, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hwndStatusText)

	// Progress bar
	hwndProgBar, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("msctls_progress32"))),
		uintptr(unsafe.Pointer(utf16Ptr(""))),
		WS_CHILD|WS_BORDER,
		20, 258, 365, 16,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hwndProgBar, PBM_SETRANGE, 0, uintptr((100<<16)))
	step1Controls = append(step1Controls, hwndProgBar)

	// Bottom Buttons
	hwndInstBtn, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Install"))),
		WS_CHILD|WS_TABSTOP|BS_DEFPUSHBUTTON,
		240, 285, 85, 26,
		hwnd, uintptr(ID_INSTALL_BTN), hInstance, 0,
	)
	procSendMessageW.Call(hwndInstBtn, WM_SETFONT, hFontBold, 1)
	step1Controls = append(step1Controls, hwndInstBtn)

	hHelp1, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Help"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		335, 285, 85, 26,
		hwnd, uintptr(ID_HELP_BTN), hInstance, 0,
	)
	procSendMessageW.Call(hHelp1, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hHelp1)

	hCancel1, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Cancel"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		430, 285, 85, 26,
		hwnd, uintptr(ID_CANCEL_BTN), hInstance, 0,
	)
	procSendMessageW.Call(hCancel1, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hCancel1)

	// ==========================================
	// STEP 2 CONTROLS (VibeRAR Setup Options)
	// ==========================================
	hGrpAssoc, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Associate VibeRAR with"))),
		WS_CHILD|BS_GROUPBOX,
		15, 10, 485, 155,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hGrpAssoc, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, hGrpAssoc)

	exts := []string{
		"RAR", "ZIP", "7Z", "ARJ",
		"BZ2", "CAB", "GZ", "ISO",
		"JAR", "LZ", "LZH", "TAR",
		"UUE", "XZ", "Z", "001",
	}
	startX := 30
	startY := 32
	colW := 85
	rowH := 24
	for i, ext := range exts {
		col := i / 4
		row := i % 4
		cx := startX + (col * colW)
		cy := startY + (row * rowH)
		chk, _, _ := procCreateWindowExW.Call(
			0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
			uintptr(unsafe.Pointer(utf16Ptr(ext))),
			WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
			uintptr(cx), uintptr(cy), 75, 20,
			hwnd, 0, hInstance, 0,
		)
		procSendMessageW.Call(chk, WM_SETFONT, hFontUi, 1)
		procSendMessageW.Call(chk, BM_SETCHECK, BST_CHECKED, 0)
		assocCheckboxes = append(assocCheckboxes, chk)
		step2Controls = append(step2Controls, chk)
	}

	// Toggle all button
	hToggleAll, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Toggle all"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		380, 125, 105, 25,
		hwnd, uintptr(ID_TOGGLE_ALL), hInstance, 0,
	)
	procSendMessageW.Call(hToggleAll, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, hToggleAll)

	// Group: Interface
	hGrpIntf, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Interface"))),
		WS_CHILD|BS_GROUPBOX,
		15, 175, 235, 140,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hGrpIntf, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, hGrpIntf)

	chkDesktop, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Add VibeRAR to Desktop"))),
		WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
		25, 200, 215, 22,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(chkDesktop, WM_SETFONT, hFontUi, 1)
	procSendMessageW.Call(chkDesktop, BM_SETCHECK, BST_CHECKED, 0)
	step2Controls = append(step2Controls, chkDesktop)

	chkStartMenu, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Add VibeRAR to Start Menu"))),
		WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
		25, 230, 215, 22,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(chkStartMenu, WM_SETFONT, hFontUi, 1)
	procSendMessageW.Call(chkStartMenu, BM_SETCHECK, BST_CHECKED, 0)
	step2Controls = append(step2Controls, chkStartMenu)

	chkGroup, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Create VibeRAR program group"))),
		WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
		25, 260, 215, 22,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(chkGroup, WM_SETFONT, hFontUi, 1)
	procSendMessageW.Call(chkGroup, BM_SETCHECK, BST_CHECKED, 0)
	step2Controls = append(step2Controls, chkGroup)

	// Group: Shell integration
	hGrpShell, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Shell integration"))),
		WS_CHILD|BS_GROUPBOX,
		265, 175, 235, 140,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hGrpShell, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, hGrpShell)

	chkShellInteg, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Integrate VibeRAR into shell"))),
		WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
		275, 200, 215, 22,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(chkShellInteg, WM_SETFONT, hFontUi, 1)
	procSendMessageW.Call(chkShellInteg, BM_SETCHECK, BST_CHECKED, 0)
	step2Controls = append(step2Controls, chkShellInteg)

	chkCascaded, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Cascaded context menus"))),
		WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
		275, 230, 215, 22,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(chkCascaded, WM_SETFONT, hFontUi, 1)
	procSendMessageW.Call(chkCascaded, BM_SETCHECK, BST_CHECKED, 0)
	step2Controls = append(step2Controls, chkCascaded)

	// Bottom buttons for Step 2
	hOk2, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("OK"))),
		WS_CHILD|WS_TABSTOP|BS_DEFPUSHBUTTON,
		225, 345, 85, 26,
		hwnd, uintptr(ID_STEP2_OK), hInstance, 0,
	)
	procSendMessageW.Call(hOk2, WM_SETFONT, hFontBold, 1)
	step2Controls = append(step2Controls, hOk2)

	hHelp2, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Help"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		320, 345, 85, 26,
		hwnd, uintptr(ID_STEP2_HELP), hInstance, 0,
	)
	procSendMessageW.Call(hHelp2, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, hHelp2)

	hCancel2, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Cancel"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		415, 345, 85, 26,
		hwnd, uintptr(ID_STEP2_CANCEL), hInstance, 0,
	)
	procSendMessageW.Call(hCancel2, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, hCancel2)

	// ==========================================
	// STEP 3 CONTROLS (VibeRAR Setup Finished)
	// ==========================================
	hDoneTitle, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("VibeRAR has been successfully installed!"))),
		WS_CHILD|SS_LEFT,
		150, 25, 330, 24,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hDoneTitle, WM_SETFONT, hFontBold, 1)
	step3Controls = append(step3Controls, hDoneTitle)

	hDoneDesc, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("VibeRAR v1.0\nA fast, modern desktop archiver.\n\nCreated by Ahmed Samir\nOfficial Site: https://ahmedsamir.com\nGitHub: https://github.com/.../viberar\nsponsored by https://www.0webhost.com/\n\nDestination: C:\\Program Files\\VibeRAR"))),
		WS_CHILD|SS_LEFT,
		150, 52, 330, 120,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hDoneDesc, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hDoneDesc)

	// Action buttons for Step 3
	hDoneBtn, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Done"))),
		WS_CHILD|WS_TABSTOP|BS_DEFPUSHBUTTON,
		145, 185, 75, 26,
		hwnd, uintptr(ID_STEP3_DONE), hInstance, 0,
	)
	procSendMessageW.Call(hDoneBtn, WM_SETFONT, hFontBold, 1)
	step3Controls = append(step3Controls, hDoneBtn)

	hRunBtn, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Run VibeRAR"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		230, 185, 95, 26,
		hwnd, uintptr(ID_STEP3_RUN), hInstance, 0,
	)
	procSendMessageW.Call(hRunBtn, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hRunBtn)

	hHelp3, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Help"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		335, 185, 75, 26,
		hwnd, uintptr(ID_STEP3_HELP), hInstance, 0,
	)
	procSendMessageW.Call(hHelp3, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hHelp3)

	hOrderBtn, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Order..."))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		420, 185, 65, 26,
		hwnd, uintptr(ID_STEP3_ORDER), hInstance, 0,
	)
	procSendMessageW.Call(hOrderBtn, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hOrderBtn)

	// Initially show Step 1
	showStep1()
	procShowWindow.Call(hwnd, 1)
	procUpdateWindow.Call(hwnd)

	// Message Loop
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
