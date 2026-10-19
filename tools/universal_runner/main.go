package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

//go:embed dist/*
var distFS embed.FS

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	shell32                      = syscall.NewLazyDLL("shell32.dll")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procSetWindowTextW           = user32.NewProc("SetWindowTextW")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procGetWindowTextW           = user32.NewProc("GetWindowTextW")
	procShellExecuteW            = shell32.NewProc("ShellExecuteW")
)

func utf16Ptr(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}

func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func enforceWindowTitle(title string) {
	titlePtr := utf16Ptr(title)
	cb := syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
		if ret, _, _ := procIsWindowVisible.Call(hwnd); ret != 0 {
			var buf [512]uint16
			procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), 512)
			current := syscall.UTF16ToString(buf[:])
			if strings.Contains(current, "127.0.0.1") || strings.Contains(current, "localhost") || strings.Contains(current, "VibeRAR") {
				procSetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(titlePtr)))
			}
		}
		return 1
	})

	for i := 0; i < 20; i++ {
		procEnumWindows.Call(cb, 0)
		time.Sleep(200 * time.Millisecond)
	}
}

func launchAppWindow(targetURL string) *exec.Cmd {
	edgePaths := []string{
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Microsoft", "Edge", "Application", "msedge.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "Microsoft", "Edge", "Application", "msedge.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Edge", "Application", "msedge.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "Google", "Chrome", "Application", "chrome.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Google", "Chrome", "Application", "chrome.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "Application", "chrome.exe"),
	}

	args := []string{
		fmt.Sprintf("--app=%s", targetURL),
		"--window-size=1040,720",
		"--window-position=100,60",
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-features=Translate,OptimizationHints,MediaRouter,DialMediaRouteProvider",
		"--disable-extensions",
		"--disable-sync",
		"--disable-background-networking",
		"--disable-component-update",
		"--app-id=VibeRAR.Desktop",
		"--class=VibeRAR",
	}

	for _, p := range edgePaths {
		if _, err := os.Stat(p); err == nil {
			cmd := exec.Command(p, args...)
			return cmd
		}
	}

	return nil
}

func main() {
	runtime.LockOSThread()

	subFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		return
	}

	port, err := findAvailablePort()
	if err != nil {
		port = 59124
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	fileServer := http.FileServer(http.FS(subFS))
	mux.Handle("/", fileServer)

	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: mux,
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return
	}

	go func() {
		_ = server.Serve(listener)
	}()

	time.Sleep(80 * time.Millisecond)

	targetURL := fmt.Sprintf("http://127.0.0.1:%d/", port)

	cmd := launchAppWindow(targetURL)
	if cmd != nil {
		if err := cmd.Start(); err == nil {
			go enforceWindowTitle("VibeRAR - High Performance Archive Manager")
		}
	} else {
		procShellExecuteW.Call(
			0,
			uintptr(unsafe.Pointer(utf16Ptr("open"))),
			uintptr(unsafe.Pointer(utf16Ptr(targetURL))),
			0, 0, 1,
		)
	}

	select {}
}
