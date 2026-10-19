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

//go:embed register_win11_menu.bat
var embeddedRegBat []byte

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
	SS_ETCHEDHORZ       = 0x0010

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

	// Step 1 Controls
	ID_INSTALL_BTN   = 2001
	ID_CANCEL1_BTN   = 2002
	ID_BROWSE_BTN    = 2003
	ID_PATH_EDIT     = 2004

	// Step 2 Controls
	ID_STEP2_OK      = 3001
	ID_STEP2_HELP    = 3002
	ID_TOGGLE_ALL    = 3003
	ID_CONTEXT_ITEMS = 3004

	// Step 3 Controls
	ID_STEP3_DONE     = 4001
	ID_STEP3_RUN      = 4002
	ID_STEP3_WHATSNEW = 4003
	ID_STEP3_HELP     = 4004
	ID_STEP3_LICENSE  = 4005
	ID_STEP3_ORDER    = 4006
	ID_STEP3_HOMEPAGE = 4007
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
	hFontSmall     uintptr
	hFontTitle     uintptr
	hFontBanner    uintptr

	currentStep    = 1
	installPath    = "C:\\Program Files\\VibeRAR"

	step1Controls  []uintptr
	step2Controls  []uintptr
	step3Controls  []uintptr

	// Step 1 handles
	hwndPathEdit   uintptr
	hwndProgBar    uintptr
	hwndInstBtn    uintptr

	// Step 2 handles
	assocCheckboxes []uintptr
	chkDesktop      uintptr
	chkStartMenu    uintptr
	chkProgramGroup uintptr
	chkShellInteg   uintptr
	chkCascaded     uintptr
	chkIcons        uintptr
	chkLegacy       uintptr

	// Step 3 handles
	hwndDestDisplay uintptr
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

func createShortcut(shortcutPath, targetPath, workingDir, description, iconPath string) {
	vbs := fmt.Sprintf(`Set w = CreateObject("WScript.Shell")
Set s = w.CreateShortcut("%s")
s.TargetPath = "%s"
s.WorkingDirectory = "%s"
s.Description = "%s"
s.IconLocation = "%s,0"
s.Save
`, shortcutPath, targetPath, workingDir, description, iconPath)

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

	// 1. Import certificate into TrustedPeople store using certutil (instant & universal)
	if _, err := os.Stat(cerPath); err == nil {
		cmdCertUtil := exec.Command("certutil.exe", "-user", "-addstore", "-f", "TrustedPeople", cerPath)
		cmdCertUtil.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
		_ = cmdCertUtil.Run()

		cmdCert := fmt.Sprintf(`Import-Certificate -FilePath '%s' -CertStoreLocation 'Cert:\CurrentUser\TrustedPeople' -ErrorAction SilentlyContinue`, cerPath)
		runPowerShellHidden(cmdCert)
	}

	// 2. Register COM CLSID in HKCU so Windows Explorer can instantiate the inproc COM handler
	clsidKey := `Software\Classes\CLSID\{b09a765d-05a3-4c4f-835e-5c87336ee000}`
	_ = setRegString(clsidKey, "", "VibeRAR Modern Shell Command")
	_ = setRegString(clsidKey+`\InProcServer32`, "", filepath.Join(installDir, "viberar_shell.dll"))
	_ = setRegString(clsidKey+`\InProcServer32`, "ThreadingModel", "Apartment")

	// 3. Register Sparse Package with External Location in user context
	cmdAppx := fmt.Sprintf(`Add-AppxPackage -Path '%s' -ExternalLocation '%s' -ErrorAction SilentlyContinue`, msixPath, installDir)
	runPowerShellHidden(cmdAppx)

	// 4. Notify Windows Explorer
	procSHChangeNotify.Call(uintptr(0x08000000), uintptr(0x0000|0x1000), 0, 0)
}

// Draw the iconic WinRAR 3 stacked books with belt and title in Step 1 header
func drawWinRARHeader(hdc uintptr, x, y int) {
	// 1. Draw stacked books (Blue, Green, Red)
	bx := x + 2
	by := y + 10
	bw := 48
	bh := 13

	drawBook := func(kX, kY int, spineCol, coverCol, ribCol uint32) {
		hCoverBrush, _, _ := procCreateSolidBrush.Call(uintptr(coverCol))
		procSelectObject.Call(hdc, hCoverBrush)
		procRoundRect.Call(hdc, uintptr(kX), uintptr(kY), uintptr(kX+bw), uintptr(kY+bh), 3, 3)
		procDeleteObject.Call(hCoverBrush)

		hSpineBrush, _, _ := procCreateSolidBrush.Call(uintptr(spineCol))
		procSelectObject.Call(hdc, hSpineBrush)
		procRoundRect.Call(hdc, uintptr(kX), uintptr(kY), uintptr(kX+10), uintptr(kY+bh), 3, 3)
		procDeleteObject.Call(hSpineBrush)

		hPageBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(250, 246, 238)))
		procSelectObject.Call(hdc, hPageBrush)
		procRectangle.Call(hdc, uintptr(kX+bw-6), uintptr(kY+2), uintptr(kX+bw-2), uintptr(kY+bh-2))
		procDeleteObject.Call(hPageBrush)

		hRibPen, _, _ := procCreatePen.Call(0, 1, uintptr(ribCol))
		procSelectObject.Call(hdc, hRibPen)
		procMoveToEx.Call(hdc, uintptr(kX+2), uintptr(kY+3), 0)
		procLineTo.Call(hdc, uintptr(kX+9), uintptr(kY+3))
		procMoveToEx.Call(hdc, uintptr(kX+2), uintptr(kY+bh-4), 0)
		procLineTo.Call(hdc, uintptr(kX+9), uintptr(kY+bh-4))
		procDeleteObject.Call(hRibPen)
	}

	// Bottom book (Crimson Red)
	drawBook(bx+4, by+28, rgb(130, 20, 30), rgb(190, 35, 45), rgb(240, 200, 60))
	// Middle book (Emerald Green)
	drawBook(bx+2, by+14, rgb(20, 95, 40), rgb(35, 145, 60), rgb(240, 200, 60))
	// Top book (Royal Blue)
	drawBook(bx, by, rgb(15, 55, 125), rgb(28, 85, 190), rgb(240, 200, 60))

	// Leather Strap wrapping books
	strapX := bx + 22
	hStrapBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(40, 25, 15)))
	procSelectObject.Call(hdc, hStrapBrush)
	procRectangle.Call(hdc, uintptr(strapX), uintptr(by-2), uintptr(strapX+8), uintptr(by+bh+30))
	procDeleteObject.Call(hStrapBrush)

	// Gold buckle
	hBuckleBrush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(255, 215, 0)))
	procSelectObject.Call(hdc, hBuckleBrush)
	procRoundRect.Call(hdc, uintptr(strapX-2), uintptr(by+12), uintptr(strapX+10), uintptr(by+26), 2, 2)
	procDeleteObject.Call(hBuckleBrush)

	// 2. Main Title Text: "VibeRAR 7.00"
	procSetBkMode.Call(hdc, 1) // TRANSPARENT
	procSetTextColor.Call(hdc, uintptr(rgb(10, 45, 110)))
	oldFont, _, _ := procSelectObject.Call(hdc, hFontBanner)
	titleText := utf16Ptr("VibeRAR 7.00")
	procTextOutW.Call(hdc, uintptr(x+62), uintptr(y+8), uintptr(unsafe.Pointer(titleText)), 12)

	// Sub-title: "64-bit"
	procSetTextColor.Call(hdc, uintptr(rgb(100, 110, 125)))
	procSelectObject.Call(hdc, hFontUi)
	subText := utf16Ptr("64-bit high performance archiver")
	procTextOutW.Call(hdc, uintptr(x+65), uintptr(y+35), uintptr(unsafe.Pointer(subText)), 32)

	procSelectObject.Call(hdc, oldFont)
}

func showStep1() {
	currentStep = 1
	for _, h := range step2Controls {
		procShowWindow.Call(h, 0)
	}
	for _, h := range step3Controls {
		procShowWindow.Call(h, 0)
	}
	for _, h := range step1Controls {
		procShowWindow.Call(h, 1)
	}

	procSetWindowTextW.Call(hwndMain, uintptr(unsafe.Pointer(utf16Ptr("WinRAR 7.00"))))
	procSetWindowPos.Call(hwndMain, 0, 150, 100, 510, 451, 0x0004)
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
	for _, h := range step2Controls {
		procShowWindow.Call(h, 1)
	}

	procSetWindowTextW.Call(hwndMain, uintptr(unsafe.Pointer(utf16Ptr("WinRAR Setup"))))
	procSetWindowPos.Call(hwndMain, 0, 150, 100, 510, 430, 0x0004)
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
	for _, h := range step3Controls {
		procShowWindow.Call(h, 1)
	}

	// Update destination path display in Step 3
	procSetWindowTextW.Call(hwndDestDisplay, uintptr(unsafe.Pointer(utf16Ptr(installPath))))

	procSetWindowTextW.Call(hwndMain, uintptr(unsafe.Pointer(utf16Ptr("WinRAR Setup"))))
	procSetWindowPos.Call(hwndMain, 0, 150, 100, 510, 430, 0x0004)
	procInvalidateRect.Call(hwndMain, 0, 1)
}

func browseFolder() {
	var bi BROWSEINFOW
	bi.HwndOwner = hwndMain
	bi.LpszTitle = utf16Ptr("Select destination folder for VibeRAR:")
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

		// Progress emulation exactly matching WinRAR setup
		for p := 0; p <= 100; p += 5 {
			procSendMessageW.Call(hwndProgBar, PBM_SETPOS, uintptr(p), 0)
			time.Sleep(30 * time.Millisecond)
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
		if len(embeddedRegBat) > 0 {
			_ = os.WriteFile(filepath.Join(installPath, "register_win11_menu.bat"), embeddedRegBat, 0755)
		}

		time.Sleep(100 * time.Millisecond)
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
		createShortcut(shortcutPath, targetExe, installPath, "VibeRAR - High Performance Archive Manager", iconIco)
	}

	if isStartMenu {
		startMenuDir := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs")
		shortcutPath := filepath.Join(startMenuDir, "VibeRAR.lnk")
		createShortcut(shortcutPath, targetExe, installPath, "VibeRAR - High Performance Archive Manager", iconIco)
	}

	// Register in Windows App Paths using pure Win32 Registry APIs
	regSubKey := "Software\\Microsoft\\Windows\\CurrentVersion\\App Paths\\VibeRAR.exe"
	_ = setRegString(regSubKey, "", targetExe)
	_ = setRegString(regSubKey, "Path", installPath)

	if isShellInteg {
		// Context Menus on Directory (Folders)
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Add", "", "Add to archive with VibeRAR...")
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Add", "MUIVerb", "Add to archive with VibeRAR...")
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Add", "Icon", iconIco)
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Add\\command", "", fmt.Sprintf("\"%s\" --add \"%%1\"", targetExe))

		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Open", "", "Open with VibeRAR")
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Open", "MUIVerb", "Open with VibeRAR")
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Open", "Icon", iconIco)
		_ = setRegString("Software\\Classes\\Directory\\shell\\VibeRAR.Open\\command", "", fmt.Sprintf("\"%s\" --open \"%%1\"", targetExe))

		// Context Menus on All Files (*)
		_ = setRegString("Software\\Classes\\*\\shell\\VibeRAR.Add", "", "Add to archive with VibeRAR...")
		_ = setRegString("Software\\Classes\\*\\shell\\VibeRAR.Add", "MUIVerb", "Add to archive with VibeRAR...")
		_ = setRegString("Software\\Classes\\*\\shell\\VibeRAR.Add", "Icon", iconIco)
		_ = setRegString("Software\\Classes\\*\\shell\\VibeRAR.Add\\command", "", fmt.Sprintf("\"%s\" --add \"%%1\"", targetExe))

		// Shared Cascaded Submenu
		menuSubKey := "Software\\Classes\\VibeRAR.ContextMenu\\shell"
		_ = setRegString(menuSubKey+"\\01_add", "MUIVerb", "Add to archive...")
		_ = setRegString(menuSubKey+"\\01_add", "Icon", iconIco)
		_ = setRegString(menuSubKey+"\\01_add\\command", "", fmt.Sprintf("\"%s\" --add \"%%1\"", targetExe))

		_ = setRegString(menuSubKey+"\\02_fast_rar", "MUIVerb", "Add to archive (Fast .rar)")
		_ = setRegString(menuSubKey+"\\02_fast_rar", "Icon", iconIco)
		_ = setRegString(menuSubKey+"\\02_fast_rar\\command", "", fmt.Sprintf("\"%s\" --quick-rar \"%%1\"", targetExe))

		_ = setRegString(menuSubKey+"\\03_fast_zip", "MUIVerb", "Add to archive (Fast .zip)")
		_ = setRegString(menuSubKey+"\\03_fast_zip", "Icon", iconIco)
		_ = setRegString(menuSubKey+"\\03_fast_zip\\command", "", fmt.Sprintf("\"%s\" --quick-zip \"%%1\"", targetExe))

		_ = setRegString(menuSubKey+"\\04_extract", "MUIVerb", "Extract files...")
		_ = setRegString(menuSubKey+"\\04_extract", "Icon", iconIco)
		_ = setRegString(menuSubKey+"\\04_extract\\command", "", fmt.Sprintf("\"%s\" --extract \"%%1\"", targetExe))

		_ = setRegString(menuSubKey+"\\05_extract_here", "MUIVerb", "Extract Here")
		_ = setRegString(menuSubKey+"\\05_extract_here", "Icon", iconIco)
		_ = setRegString(menuSubKey+"\\05_extract_here\\command", "", fmt.Sprintf("\"%s\" --extract-here \"%%1\"", targetExe))

		_ = setRegString(menuSubKey+"\\06_open", "MUIVerb", "Open with VibeRAR")
		_ = setRegString(menuSubKey+"\\06_open", "Icon", iconIco)
		_ = setRegString(menuSubKey+"\\06_open\\command", "", fmt.Sprintf("\"%s\" --open \"%%1\"", targetExe))

		for _, baseKey := range []string{"Software\\Classes\\Directory\\shell\\VibeRAR", "Software\\Classes\\Folder\\shell\\VibeRAR", "Software\\Classes\\*\\shell\\VibeRAR"} {
			_ = setRegString(baseKey, "", "")
			_ = setRegString(baseKey, "MUIVerb", "VibeRAR")
			_ = setRegString(baseKey, "Icon", iconIco)
			_ = setRegString(baseKey, "ExtendedSubCommandsKey", "VibeRAR.ContextMenu")
		}

		// Windows 11 Modern Context Menu (MSIX Sparse Package Integration)
		registerModernWindows11ContextMenu(installPath)
	}

	// Register File Associations
	archiveProgID := "Software\\Classes\\VibeRAR.Archive"
	_ = setRegString(archiveProgID, "", "VibeRAR Archive")
	_ = setRegString(archiveProgID+"\\DefaultIcon", "", fmt.Sprintf("\"%s\",0", iconIco))
	_ = setRegString(archiveProgID+"\\shell\\open\\command", "", fmt.Sprintf("\"%s\" --open \"%%1\"", targetExe))

	extList := []string{
		"rar", "zip", "7z", "arj", "bz2", "cab", "gz", "iso", "jar",
		"lz", "lzh", "tar", "uue", "xz", "z", "zipx", "zst", "001",
	}
	for i, chk := range assocCheckboxes {
		state, _, _ := procSendMessageW.Call(chk, BM_GETCHECK, 0, 0)
		if state == BST_CHECKED && i < len(extList) {
			ext := "." + extList[i]
			_ = setRegString("Software\\Classes\\"+ext, "", "VibeRAR.Archive")
		}
	}

	procSHChangeNotify.Call(uintptr(0x08000000), uintptr(0x0000|0x1000), 0, 0)
	showStep3()
}

func refreshWindowsExplorer() {
	runPowerShellHidden(`Stop-Process -Name explorer -Force -ErrorAction SilentlyContinue; Start-Sleep -Milliseconds 400; Start-Process explorer`)
}

func launchApplication() {
	go refreshWindowsExplorer()

	targetExe := filepath.Join(installPath, "VibeRAR.exe")
	if _, err := os.Stat(targetExe); os.IsNotExist(err) {
		targetExe = filepath.Join(installPath, "VibeRAR_Portable.exe")
	}

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

func openURL(url string) {
	procShellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr("open"))),
		uintptr(unsafe.Pointer(utf16Ptr(url))),
		0, 0, 1,
	)
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
		case ID_CANCEL1_BTN:
			procPostQuitMessage.Call(0)

		// Step 2
		case ID_TOGGLE_ALL:
			toggleAllAssociations()
		case ID_STEP2_OK:
			applyStep2AndContinue()
		case ID_STEP2_HELP:
			openURL("https://github.com/iegyi/VibeRAR#readme")
		case ID_CONTEXT_ITEMS:
			procMessageBoxW.Call(
				hwnd,
				uintptr(unsafe.Pointer(utf16Ptr("All essential context menu items are enabled by default:\n\n• Add to archive...\n• Add to <name>.zip\n• Add to <name>.rar\n• Extract files...\n• Extract Here\n• Extract to <folder>\\\n• Open with VibeRAR"))),
				uintptr(unsafe.Pointer(utf16Ptr("Context menu items"))),
				0x40,
			)

		// Step 3
		case ID_STEP3_DONE:
			procPostQuitMessage.Call(0)
		case ID_STEP3_RUN:
			launchApplication()
		case ID_STEP3_WHATSNEW:
			openURL("https://github.com/iegyi/VibeRAR/releases")
		case ID_STEP3_HELP:
			openURL("https://github.com/iegyi/VibeRAR#readme")
		case ID_STEP3_LICENSE:
			procMessageBoxW.Call(
				hwnd,
				uintptr(unsafe.Pointer(utf16Ptr("VibeRAR License Agreement\n\nVibeRAR is distributed as 100% FREE and UNLIMITED software forever.\nYou may use this software for personal, commercial, and educational purposes without purchasing a license.\n\nAuthor: Ahmed Samir Hany (https://ahmedsamir.com)\nRepository: https://github.com/iegyi/VibeRAR"))),
				uintptr(unsafe.Pointer(utf16Ptr("VibeRAR License"))),
				0x40,
			)
		case ID_STEP3_ORDER:
			procMessageBoxW.Call(
				hwnd,
				uintptr(unsafe.Pointer(utf16Ptr("VibeRAR is 100% FREE forever!\n\nNo license key or purchase order required. Enjoy unlimited features!"))),
				uintptr(unsafe.Pointer(utf16Ptr("Evaluation Notice"))),
				0x40,
			)
		case ID_STEP3_HOMEPAGE:
			openURL("https://ahmedsamir.com")
		}
		return 0

	case WM_PAINT:
		var ps PAINTSTRUCT
		hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))

		if currentStep == 1 {
			// Draw identical WinRAR header banner
			drawWinRARHeader(hdc, 6, 6)

			// Separator line under header
			hPen, _, _ := procCreatePen.Call(0, 1, uintptr(rgb(215, 220, 225)))
			oldPen, _, _ := procSelectObject.Call(hdc, hPen)
			procMoveToEx.Call(hdc, 0, 73, 0)
			procLineTo.Call(hdc, 500, 73)
			procSelectObject.Call(hdc, oldPen)
			procDeleteObject.Call(hPen)
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

	// Create main window (510x451 for exact WinRAR 7.00 size)
	hwnd, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(utf16Ptr("WinRAR 7.00"))),
		WS_OVERLAPPEDWINDOW&^0x00050000, // No Maximize / Resize
		150, 120, 510, 451,
		0, 0, hInstance, 0,
	)
	hwndMain = hwnd
	if hAppIcon != 0 {
		procSendMessageW.Call(hwndMain, WM_SETICON, 0, hAppIcon)
		procSendMessageW.Call(hwndMain, WM_SETICON, 1, hAppIcon)
	}

	// Create Standard Windows Segoe UI Fonts matching WinRAR
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
	hFontSmall, _, _ = procCreateFontW.Call(
		uintptr(neg(11)), 0, 0, 0, 400, 0, 0, 0,
		0, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(utf16Ptr("Segoe UI"))),
	)
	hFontBanner, _, _ = procCreateFontW.Call(
		uintptr(neg(20)), 0, 0, 0, 700, 0, 0, 0,
		0, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(utf16Ptr("Segoe UI"))),
	)

	// =========================================================================
	// STEP 1 CONTROLS (Identical to WinRAR 7.00 Step 1)
	// =========================================================================

	// Copyright text (top right of banner)
	hCopy1, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("Copyright © 2024-2026"))),
		WS_CHILD|SS_LEFT, 292, 10, 195, 14,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hCopy1, WM_SETFONT, hFontSmall, 1)
	step1Controls = append(step1Controls, hCopy1)

	hCopy2, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("by Ahmed Samir Hany"))),
		WS_CHILD|SS_LEFT, 292, 26, 195, 14,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hCopy2, WM_SETFONT, hFontSmall, 1)
	step1Controls = append(step1Controls, hCopy2)

	// Destination folder label: Pos (6, 78) Size: 485x14
	hDestLbl, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("&Destination folder"))),
		WS_CHILD|SS_LEFT, 6, 78, 485, 14,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hDestLbl, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hDestLbl)

	// Destination folder Edit: Pos (6, 94) Size: 405x22
	hwndPathEdit, _, _ = procCreateWindowExW.Call(
		0x00000200, uintptr(unsafe.Pointer(utf16Ptr("EDIT"))),
		uintptr(unsafe.Pointer(utf16Ptr(installPath))),
		WS_CHILD|WS_BORDER|WS_TABSTOP|ES_AUTOHSCROLL,
		6, 94, 403, 22,
		hwnd, uintptr(ID_PATH_EDIT), hInstance, 0,
	)
	procSendMessageW.Call(hwndPathEdit, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hwndPathEdit)

	// Browse button: Pos (415, 93) Size: 75x24
	hBrowse, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Bro&wse..."))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		415, 93, 75, 24,
		hwnd, uintptr(ID_BROWSE_BTN), hInstance, 0,
	)
	procSendMessageW.Call(hBrowse, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hBrowse)

	// License agreement multiline box: Pos (4, 122) Size: 486x188
	licenseText := `VibeRAR 7.00 (64-bit)
High-Performance Archive Manager & Compression Suite
===================================================
Software Author: Ahmed Samir Hany
Official Website: https://ahmedsamir.com
Repository: https://github.com/iegyi/VibeRAR

Sponsored Web Hosting:
https://www.0webhost.com/ - Buy reliable web hosting at affordable rates
===================================================

END USER LICENSE AGREEMENT (EULA):
VibeRAR is distributed as 100% FREE, FULLY FUNCTIONAL and UNLIMITED software.
You may use this software indefinitely for personal, business, and educational purposes without purchasing a license, entering serial keys, or experiencing expiration nag screens.

FEATURES & SPECIFICATIONS:
• High-performance decompression for RAR (RAR4, RAR5, RAR7), ZIP, 7Z, TAR, GZ, ISO, CAB, BZ2, XZ, and ZST.
• Multi-threaded compression to ZIP64 and RAR5 with hardware-accelerated CRC32 checksums.
• Full Windows 11 Modern Top-Level Context Menu integration (IExplorerCommand + Sparse Package).
• Instant shell operations (< 50ms): Extract Here, Extract to folder, and Fast Add.
• Multi-threaded CPU Benchmark utility for hardware evaluation.

Please click "Install" to begin copying VibeRAR to your computer.`

	hEula, _, _ := procCreateWindowExW.Call(
		0x00000200, uintptr(unsafe.Pointer(utf16Ptr("EDIT"))),
		uintptr(unsafe.Pointer(utf16Ptr(licenseText))),
		WS_CHILD|WS_BORDER|ES_MULTILINE|ES_READONLY|WS_VSCROLL,
		4, 122, 486, 188,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hEula, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hEula)

	// Progress bar: Pos (4, 317) Size: 486x14
	hwndProgBar, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("msctls_progress32"))),
		uintptr(unsafe.Pointer(utf16Ptr(""))),
		WS_CHILD|WS_BORDER,
		4, 317, 486, 14,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hwndProgBar, PBM_SETRANGE, 0, uintptr((100<<16)))
	step1Controls = append(step1Controls, hwndProgBar)

	// EULA Acceptance notice: Pos (6, 338) Size: 485x34
	hEulaNotice, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("If you agree to the END USER LICENSE AGREEMENT (EULA), please click \"Install\" to continue, or click \"Cancel\" to abort the installation."))),
		WS_CHILD|SS_LEFT, 6, 338, 485, 34,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hEulaNotice, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hEulaNotice)

	// Install button: Pos (135, 376) Size: 105x24
	hwndInstBtn, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Install"))),
		WS_CHILD|WS_TABSTOP|BS_DEFPUSHBUTTON,
		135, 376, 105, 24,
		hwnd, uintptr(ID_INSTALL_BTN), hInstance, 0,
	)
	procSendMessageW.Call(hwndInstBtn, WM_SETFONT, hFontBold, 1)
	step1Controls = append(step1Controls, hwndInstBtn)

	// Cancel button: Pos (255, 376) Size: 105x24
	hCancel1, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Cancel"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		255, 376, 105, 24,
		hwnd, uintptr(ID_CANCEL1_BTN), hInstance, 0,
	)
	procSendMessageW.Call(hCancel1, WM_SETFONT, hFontUi, 1)
	step1Controls = append(step1Controls, hCancel1)


	// =========================================================================
	// STEP 2 CONTROLS (Identical to WinRAR Setup Step 2)
	// =========================================================================

	// GroupBox 1: Associate WinRAR with: Pos (10, 8) Size: 191x225
	hGrpAssoc, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Associate VibeRAR with"))),
		WS_CHILD|BS_GROUPBOX,
		10, 8, 191, 225,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hGrpAssoc, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, hGrpAssoc)

	// Column 1 (x=24)
	col1 := []string{"&RAR", "&ZIP", "7Z", "ARJ", "BZ2", "CAB", "GZ", "ISO", "JAR"}
	for i, ext := range col1 {
		cy := 33 + (i * 19)
		chk, _, _ := procCreateWindowExW.Call(
			0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
			uintptr(unsafe.Pointer(utf16Ptr(ext))),
			WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
			24, uintptr(cy), 56, 19,
			hwnd, 0, hInstance, 0,
		)
		procSendMessageW.Call(chk, WM_SETFONT, hFontUi, 1)
		procSendMessageW.Call(chk, BM_SETCHECK, BST_CHECKED, 0)
		assocCheckboxes = append(assocCheckboxes, chk)
		step2Controls = append(step2Controls, chk)
	}

	// Column 2 (x=84)
	col2 := []string{"LZ", "LZH", "TAR", "UUE", "XZ", "Z", "ZIPX", "ZST", "001"}
	for i, ext := range col2 {
		cy := 33 + (i * 19)
		chk, _, _ := procCreateWindowExW.Call(
			0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
			uintptr(unsafe.Pointer(utf16Ptr(ext))),
			WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
			84, uintptr(cy), 56, 19,
			hwnd, 0, hInstance, 0,
		)
		procSendMessageW.Call(chk, WM_SETFONT, hFontUi, 1)
		procSendMessageW.Call(chk, BM_SETCHECK, BST_CHECKED, 0)
		assocCheckboxes = append(assocCheckboxes, chk)
		step2Controls = append(step2Controls, chk)
	}

	// Toggle all button: Pos (24, 205) Size: 162x22
	hToggleAll, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Toggle &all"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		24, 205, 162, 22,
		hwnd, uintptr(ID_TOGGLE_ALL), hInstance, 0,
	)
	procSendMessageW.Call(hToggleAll, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, hToggleAll)

	// GroupBox 2: Interface: Pos (210, 8) Size: 276x89
	hGrpUi, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Interface"))),
		WS_CHILD|BS_GROUPBOX,
		210, 8, 276, 89,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hGrpUi, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, hGrpUi)

	chkDesktop, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Add VibeRAR to &Desktop"))),
		WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
		222, 30, 255, 19,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(chkDesktop, WM_SETFONT, hFontUi, 1)
	procSendMessageW.Call(chkDesktop, BM_SETCHECK, BST_CHECKED, 0)
	step2Controls = append(step2Controls, chkDesktop)

	chkStartMenu, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Add VibeRAR to &Start Menu"))),
		WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
		222, 49, 255, 19,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(chkStartMenu, WM_SETFONT, hFontUi, 1)
	procSendMessageW.Call(chkStartMenu, BM_SETCHECK, BST_CHECKED, 0)
	step2Controls = append(step2Controls, chkStartMenu)

	chkProgramGroup, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Create VibeRAR &program group"))),
		WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
		222, 67, 255, 19,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(chkProgramGroup, WM_SETFONT, hFontUi, 1)
	procSendMessageW.Call(chkProgramGroup, BM_SETCHECK, 0, 0)
	step2Controls = append(step2Controls, chkProgramGroup)

	// GroupBox 3: Shell integration: Pos (210, 102) Size: 276x132
	hGrpShell, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Shell integration"))),
		WS_CHILD|BS_GROUPBOX,
		210, 102, 276, 132,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hGrpShell, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, hGrpShell)

	chkShellInteg, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("&Integrate VibeRAR into shell"))),
		WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
		222, 122, 255, 19,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(chkShellInteg, WM_SETFONT, hFontUi, 1)
	procSendMessageW.Call(chkShellInteg, BM_SETCHECK, BST_CHECKED, 0)
	step2Controls = append(step2Controls, chkShellInteg)

	chkCascaded, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Cascaded con&text menus"))),
		WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
		222, 141, 255, 19,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(chkCascaded, WM_SETFONT, hFontUi, 1)
	procSendMessageW.Call(chkCascaded, BM_SETCHECK, BST_CHECKED, 0)
	step2Controls = append(step2Controls, chkCascaded)

	chkIcons, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Ico&ns in context menus"))),
		WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
		222, 160, 255, 19,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(chkIcons, WM_SETFONT, hFontUi, 1)
	procSendMessageW.Call(chkIcons, BM_SETCHECK, BST_CHECKED, 0)
	step2Controls = append(step2Controls, chkIcons)

	chkLegacy, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("&Legacy context menus"))),
		WS_CHILD|BS_AUTOCHECKBOX|WS_TABSTOP,
		222, 178, 255, 19,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(chkLegacy, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, chkLegacy)

	hBtnContextItems, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("&Context menu items..."))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		222, 205, 210, 22,
		hwnd, uintptr(ID_CONTEXT_ITEMS), hInstance, 0,
	)
	procSendMessageW.Call(hBtnContextItems, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, hBtnContextItems)

	// GroupBox 4: Bottom notice frame: Pos (10, 236) Size: 476x117
	hGrpBottomFrame, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr(""))),
		WS_CHILD|BS_GROUPBOX,
		10, 236, 476, 117,
		hwnd, 0, hInstance, 0,
	)
	step2Controls = append(step2Controls, hGrpBottomFrame)

	step2Notice := `By default VibeRAR integrates into the top-level Windows 11 context menu using an MSIX Sparse Package and native 64-bit IExplorerCommand COM server.

If "Legacy context menus" is turned on, VibeRAR options will appear in the classic Windows "Show more options" menu.

VibeRAR automatically associates with RAR, ZIP, 7Z, and 15 popular archive formats.`

	hStep2Text, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr(step2Notice))),
		WS_CHILD|SS_LEFT,
		18, 249, 459, 97,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hStep2Text, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, hStep2Text)

	// Step 2 Bottom Buttons: OK at (136, 364), Help at (253, 364)
	hOk2, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("OK"))),
		WS_CHILD|WS_TABSTOP|BS_DEFPUSHBUTTON,
		136, 364, 105, 24,
		hwnd, uintptr(ID_STEP2_OK), hInstance, 0,
	)
	procSendMessageW.Call(hOk2, WM_SETFONT, hFontBold, 1)
	step2Controls = append(step2Controls, hOk2)

	hHelp2, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Help"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		253, 364, 105, 24,
		hwnd, uintptr(ID_STEP2_HELP), hInstance, 0,
	)
	procSendMessageW.Call(hHelp2, WM_SETFONT, hFontUi, 1)
	step2Controls = append(step2Controls, hHelp2)


	// =========================================================================
	// STEP 3 CONTROLS (Identical to WinRAR Setup Step 3 Done)
	// =========================================================================

	// Top GroupBox frame: Pos (12, 3) Size: 471x211
	hStep3TopFrame, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr(""))),
		WS_CHILD|BS_GROUPBOX,
		12, 3, 471, 211,
		hwnd, 0, hInstance, 0,
	)
	step3Controls = append(step3Controls, hStep3TopFrame)

	hS3Text1, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("VibeRAR has been successfully installed to the folder"))),
		WS_CHILD|SS_LEFT, 25, 20, 446, 14,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hS3Text1, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hS3Text1)

	hwndDestDisplay, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr(installPath))),
		WS_CHILD|SS_LEFT, 25, 36, 446, 16,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hwndDestDisplay, WM_SETFONT, hFontBold, 1)
	step3Controls = append(step3Controls, hwndDestDisplay)

	hS3Text2, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("Use \"Settings...\" command from the VibeRAR \"Options\" menu for further customization."))),
		WS_CHILD|SS_LEFT, 25, 56, 446, 26,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hS3Text2, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hS3Text2)

	hS3Text3, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("If you installed VibeRAR for the first time, please explore the retro themes and dual toolbar density options."))),
		WS_CHILD|SS_LEFT, 25, 86, 446, 26,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hS3Text3, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hS3Text3)

	hS3Text4, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("Please note that VibeRAR is 100% free software forever. You can use it indefinitely for personal and commercial purposes without purchasing any license."))),
		WS_CHILD|SS_LEFT, 25, 118, 446, 38,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hS3Text4, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hS3Text4)

	hS3Text5, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("Sponsored Web Hosting: https://www.0webhost.com/ - Buy hosting at cheap prices.\nCreated by Ahmed Samir Hany (https://ahmedsamir.com)"))),
		WS_CHILD|SS_LEFT, 25, 162, 446, 38,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hS3Text5, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hS3Text5)

	// Middle Command Buttons (2 Rows x 3 Columns)
	// Row 1: What's new at (12, 225), Run VibeRAR at (171, 225), Read help at (330, 225)
	hWhatsNew, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("&What's new"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		12, 225, 155, 24,
		hwnd, uintptr(ID_STEP3_WHATSNEW), hInstance, 0,
	)
	procSendMessageW.Call(hWhatsNew, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hWhatsNew)

	hRunBtn, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("&Run VibeRAR"))),
		WS_CHILD|WS_TABSTOP|BS_DEFPUSHBUTTON,
		171, 225, 155, 24,
		hwnd, uintptr(ID_STEP3_RUN), hInstance, 0,
	)
	procSendMessageW.Call(hRunBtn, WM_SETFONT, hFontBold, 1)
	step3Controls = append(step3Controls, hRunBtn)

	hHelp3, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Read &help"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		330, 225, 153, 24,
		hwnd, uintptr(ID_STEP3_HELP), hInstance, 0,
	)
	procSendMessageW.Call(hHelp3, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hHelp3)

	// Row 2: View license at (12, 253), Order at (171, 253), Home page at (330, 253)
	hLicenseBtn, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("View &license"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		12, 253, 155, 24,
		hwnd, uintptr(ID_STEP3_LICENSE), hInstance, 0,
	)
	procSendMessageW.Call(hLicenseBtn, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hLicenseBtn)

	hOrderBtn, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("&Order"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		171, 253, 155, 24,
		hwnd, uintptr(ID_STEP3_ORDER), hInstance, 0,
	)
	procSendMessageW.Call(hOrderBtn, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hOrderBtn)

	hHomePageBtn, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Ho&me page"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		330, 253, 153, 24,
		hwnd, uintptr(ID_STEP3_HOMEPAGE), hInstance, 0,
	)
	procSendMessageW.Call(hHomePageBtn, WM_SETFONT, hFontUi, 1)
	step3Controls = append(step3Controls, hHomePageBtn)

	// Bottom frame with "Thank you for using VibeRAR!": Pos (12, 285) Size: 471x69
	hStep3BottomFrame, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr(""))),
		WS_CHILD|BS_GROUPBOX,
		12, 285, 471, 69,
		hwnd, 0, hInstance, 0,
	)
	step3Controls = append(step3Controls, hStep3BottomFrame)

	hThankYouText, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("STATIC"))),
		uintptr(unsafe.Pointer(utf16Ptr("Thank you for using VibeRAR!"))),
		WS_CHILD|SS_LEFT, 25, 313, 446, 25,
		hwnd, 0, hInstance, 0,
	)
	procSendMessageW.Call(hThankYouText, WM_SETFONT, hFontBold, 1)
	step3Controls = append(step3Controls, hThankYouText)

	// Done button at Pos (195, 364) Size: 105x24
	hDoneBtn, _, _ := procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(utf16Ptr("BUTTON"))),
		uintptr(unsafe.Pointer(utf16Ptr("Done"))),
		WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON,
		195, 364, 105, 24,
		hwnd, uintptr(ID_STEP3_DONE), hInstance, 0,
	)
	procSendMessageW.Call(hDoneBtn, WM_SETFONT, hFontBold, 1)
	step3Controls = append(step3Controls, hDoneBtn)

	// Show initial step
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
