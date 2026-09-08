//go:build windows

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"
	"unicode/utf16"
	"unsafe"
)

const (
	modelName      = "gemma3:1b"
	keepAlive      = "2m"
	debounceDelay  = 250 * time.Millisecond
	pollInterval   = 50 * time.Millisecond
	serverIdleStop = 3 * time.Minute
	apiBase        = "http://127.0.0.1:11434"
	createNoWindow = 0x08000000

	wmUser          = 0x0400
	wmTray          = wmUser + 1
	wmRButtonUp     = 0x0205
	wmLButtonDblClk = 0x0203
	wmDestroy       = 0x0002

	nimAdd     = 0x00000000
	nimModify  = 0x00000001
	nimDelete  = 0x00000002
	nifMessage = 0x00000001
	nifIcon    = 0x00000002
	nifTip     = 0x00000004

	mfString       = 0x00000000
	mfSeparator    = 0x00000800
	tpmRightButton = 0x0002
	tpmReturnCmd   = 0x0100

	idiApplication = 32512
	idcArrow       = 32512

	menuRelease = 1001
	menuStartup = 1002
	menuExit    = 1003

	f24VirtualKey  = 0x87
	keyeventfKeyup = 0x0002
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	shell32                 = syscall.NewLazyDLL("shell32.dll")
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	procRegisterClassExW    = user32.NewProc("RegisterClassExW")
	procCreateWindowExW     = user32.NewProc("CreateWindowExW")
	procDefWindowProcW      = user32.NewProc("DefWindowProcW")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessageW    = user32.NewProc("DispatchMessageW")
	procPostQuitMessage     = user32.NewProc("PostQuitMessage")
	procLoadIconW           = user32.NewProc("LoadIconW")
	procLoadCursorW         = user32.NewProc("LoadCursorW")
	procCreatePopupMenu     = user32.NewProc("CreatePopupMenu")
	procAppendMenuW         = user32.NewProc("AppendMenuW")
	procTrackPopupMenu      = user32.NewProc("TrackPopupMenu")
	procDestroyMenu         = user32.NewProc("DestroyMenu")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procGetCursorPos        = user32.NewProc("GetCursorPos")
	procMessageBoxW         = user32.NewProc("MessageBoxW")
	procKeybdEvent          = user32.NewProc("keybd_event")
	procShellNotifyIconW    = shell32.NewProc("Shell_NotifyIconW")
	procGetModuleHandleW    = kernel32.NewProc("GetModuleHandleW")

	globalBridge *Bridge
	globalHWND   uintptr
	trayData     notifyIconData
)

type point struct {
	X int32
	Y int32
}

type msg struct {
	HWnd     uintptr
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       point
	LPrivate uint32
}

type wndClassEx struct {
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

type notifyIconData struct {
	CbSize           uint32
	HWnd             uintptr
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            uintptr
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GuidItem         [16]byte
	HBalloonIcon     uintptr
}

type generateOptions struct {
	Temperature float64 `json:"temperature"`
	NumPredict  int     `json:"num_predict"`
}

type generateRequest struct {
	Model     string           `json:"model"`
	Prompt    string           `json:"prompt"`
	Stream    bool             `json:"stream"`
	KeepAlive interface{}      `json:"keep_alive"`
	Options   *generateOptions `json:"options,omitempty"`
}

type generateResponse struct {
	Response string `json:"response"`
}

type Bridge struct {
	dir          string
	requestFile  string
	responseFile string
	tempFile     string
	logFile      string

	client *http.Client

	mu            sync.Mutex
	ollamaCmd     *exec.Cmd
	startedByUs   bool
	lastActivity  time.Time
	lastHandled   string
	warnedMissing bool
	stopped       bool
}

func newBridge() *Bridge {
	dir := filepath.Join(os.TempDir(), "rime_ollama_bridge")
	_ = os.MkdirAll(dir, 0755)
	return &Bridge{
		dir:          dir,
		requestFile:  filepath.Join(dir, "request.txt"),
		responseFile: filepath.Join(dir, "response.txt"),
		tempFile:     filepath.Join(dir, "response.tmp"),
		logFile:      filepath.Join(dir, "bridge.log"),
		client:       &http.Client{Timeout: 25 * time.Second},
		lastActivity: time.Now(),
	}
}

func (b *Bridge) logf(format string, args ...interface{}) {
	line := fmt.Sprintf("%s %s\r\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
	f, err := os.OpenFile(b.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		_, _ = f.WriteString(line)
		_ = f.Close()
	}
}

func (b *Bridge) run() {
	b.logf("RimeTranslate bridge started")
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		b.mu.Lock()
		if b.stopped {
			b.mu.Unlock()
			return
		}
		b.mu.Unlock()

		b.stopOwnedOllamaIfIdle()

		data, err := os.ReadFile(b.requestFile)
		if err != nil {
			continue
		}
		text := strings.TrimSpace(string(data))
		if text == "" || text == b.lastHandled {
			continue
		}

		snapshot := text
		time.Sleep(debounceDelay)

		latestBytes, err := os.ReadFile(b.requestFile)
		if err != nil {
			continue
		}
		latest := strings.TrimSpace(string(latestBytes))
		if latest != snapshot {
			continue
		}

		b.lastHandled = snapshot
		if err := b.translateAndRespond(snapshot); err != nil {
			b.logf("translate error: %v", err)
			b.lastHandled = "" // allow retry
		}
	}
}

func translationPrompt(text string) (prompt, srcLabel, dstLabel string, ok bool) {
	hasHan := false
	hasLatin := false
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			hasHan = true
		}
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			hasLatin = true
		}
	}

	if hasHan {
		return "Translate the following Chinese text into natural, accurate, concise English. Output only the English translation, with no explanation: " + text, "ZH", "EN", true
	}
	if hasLatin {
		return "Translate the following English text into natural, accurate, concise Simplified Chinese. Output only the Chinese translation, with no explanation: " + text, "EN", "ZH", true
	}
	return "", "", "", false
}

func (b *Bridge) translateAndRespond(text string) error {
	prompt, srcLabel, dstLabel, ok := translationPrompt(text)
	if !ok {
		_ = os.Remove(b.responseFile)
		return nil
	}

	if err := b.ensureOllama(); err != nil {
		return err
	}
	reqBody := generateRequest{
		Model:     modelName,
		Prompt:    prompt,
		Stream:    false,
		KeepAlive: keepAlive,
		Options: &generateOptions{
			Temperature: 0.1,
			NumPredict:  32,
		},
	}
	raw, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(http.MethodPost, apiBase+"/api/generate", bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("ollama HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var out generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return err
	}
	translation := strings.TrimSpace(out.Response)
	if translation == "" {
		return errors.New("empty translation")
	}

	content := text + "\n" + translation
	if err := os.WriteFile(b.tempFile, []byte(content), 0644); err != nil {
		return err
	}
	_ = os.Remove(b.responseFile)
	if err := os.Rename(b.tempFile, b.responseFile); err != nil {
		return err
	}

	b.mu.Lock()
	b.lastActivity = time.Now()
	b.mu.Unlock()

	sendF24()
	b.logf("%s=%s | %s=%s", srcLabel, text, dstLabel, translation)
	return nil
}

func (b *Bridge) ollamaAlive() bool {
	client := &http.Client{Timeout: 650 * time.Millisecond}
	resp, err := client.Get(apiBase + "/api/tags")
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 500
}

func (b *Bridge) ensureOllama() error {
	if b.ollamaAlive() {
		return nil
	}

	b.mu.Lock()
	alreadyStarting := b.startedByUs && b.ollamaCmd != nil && b.ollamaCmd.Process != nil
	b.mu.Unlock()

	if !alreadyStarting {
		path, err := exec.LookPath("ollama.exe")
		if err != nil {
			path, err = exec.LookPath("ollama")
		}
		if err != nil {
			b.mu.Lock()
			first := !b.warnedMissing
			b.warnedMissing = true
			b.mu.Unlock()
			if first {
				showMessage("Rime Translate", "未找到 ollama.exe。请先安装 Ollama，或把 Ollama 加入 PATH。")
			}
			return errors.New("ollama.exe not found in PATH")
		}

		cmd := exec.Command(path, "serve")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: createNoWindow,
		}
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("start ollama serve: %w", err)
		}

		b.mu.Lock()
		b.ollamaCmd = cmd
		b.startedByUs = true
		b.lastActivity = time.Now()
		b.mu.Unlock()
		b.logf("started ollama serve pid=%d", cmd.Process.Pid)

		go func(c *exec.Cmd) {
			_ = c.Wait()
			b.mu.Lock()
			if b.ollamaCmd == c {
				b.ollamaCmd = nil
				b.startedByUs = false
			}
			b.mu.Unlock()
		}(cmd)
	}

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if b.ollamaAlive() {
			return nil
		}
		time.Sleep(250 * time.Millisecond)
	}
	return errors.New("Ollama did not become ready on 127.0.0.1:11434")
}

func (b *Bridge) unloadModel() {
	if !b.ollamaAlive() {
		return
	}
	raw, _ := json.Marshal(generateRequest{
		Model:     modelName,
		Prompt:    "",
		Stream:    false,
		KeepAlive: 0,
	})
	client := &http.Client{Timeout: 3 * time.Second}
	req, _ := http.NewRequest(http.MethodPost, apiBase+"/api/generate", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err == nil && resp != nil {
		_ = resp.Body.Close()
	}
}

func (b *Bridge) releaseResources() {
	b.unloadModel()

	b.mu.Lock()
	cmd := b.ollamaCmd
	owned := b.startedByUs
	b.mu.Unlock()

	if owned && cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
		b.logf("stopped owned ollama server")
	}
	_ = os.Remove(b.responseFile)
}

func (b *Bridge) stopOwnedOllamaIfIdle() {
	b.mu.Lock()
	owned := b.startedByUs
	cmd := b.ollamaCmd
	idle := time.Since(b.lastActivity)
	b.mu.Unlock()

	if owned && cmd != nil && cmd.Process != nil && idle > serverIdleStop {
		b.releaseResources()
	}
}

func (b *Bridge) stop() {
	b.mu.Lock()
	b.stopped = true
	b.mu.Unlock()
	b.releaseResources()
	b.logf("RimeTranslate bridge stopped")
}

func sendF24() {
	procKeybdEvent.Call(f24VirtualKey, 0, 0, 0)
	procKeybdEvent.Call(f24VirtualKey, 0, keyeventfKeyup, 0)
}

func utf16Ptr(s string) *uint16 {
	u := utf16.Encode([]rune(s + "\x00"))
	return &u[0]
}

func fillUTF16(dst []uint16, s string) {
	u := utf16.Encode([]rune(s))
	if len(u) >= len(dst) {
		u = u[:len(dst)-1]
	}
	copy(dst, u)
	if len(u) < len(dst) {
		dst[len(u)] = 0
	}
}

func showMessage(title, text string) {
	procMessageBoxW.Call(0, uintptr(unsafe.Pointer(utf16Ptr(text))), uintptr(unsafe.Pointer(utf16Ptr(title))), 0x40)
}

func addTrayIcon(hwnd uintptr) error {
	icon, _, _ := procLoadIconW.Call(0, idiApplication)
	trayData = notifyIconData{}
	trayData.CbSize = uint32(unsafe.Sizeof(trayData))
	trayData.HWnd = hwnd
	trayData.UID = 1
	trayData.UFlags = nifMessage | nifIcon | nifTip
	trayData.UCallbackMessage = wmTray
	trayData.HIcon = icon
	fillUTF16(trayData.SzTip[:], "Rime Translate · 中英双向")
	r, _, err := procShellNotifyIconW.Call(nimAdd, uintptr(unsafe.Pointer(&trayData)))
	if r == 0 {
		return err
	}
	return nil
}

func removeTrayIcon() {
	if trayData.HWnd != 0 {
		procShellNotifyIconW.Call(nimDelete, uintptr(unsafe.Pointer(&trayData)))
	}
}

func regCommand(args ...string) ([]byte, error) {
	cmd := exec.Command("reg.exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: createNoWindow}
	return cmd.CombinedOutput()
}

func startupEnabled() bool {
	_, err := regCommand("query", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, "/v", "RimeTranslate")
	return err == nil
}

func setStartup(enable bool) error {
	if enable {
		exe, err := os.Executable()
		if err != nil {
			return err
		}
		exe, _ = filepath.Abs(exe)
		value := fmt.Sprintf(`"%s"`, exe)
		out, err := regCommand("add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, "/v", "RimeTranslate", "/t", "REG_SZ", "/d", value, "/f")
		if err != nil {
			return fmt.Errorf("enable startup: %v (%s)", err, strings.TrimSpace(string(out)))
		}
		return nil
	}
	out, err := regCommand("delete", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, "/v", "RimeTranslate", "/f")
	if err != nil && startupEnabled() {
		return fmt.Errorf("disable startup: %v (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func showTrayMenu(hwnd uintptr) {
	menu, _, _ := procCreatePopupMenu.Call()
	if menu == 0 {
		return
	}
	defer procDestroyMenu.Call(menu)

	releaseText := utf16Ptr("释放翻译资源")
	startupLabel := "开启开机自动启动"
	if startupEnabled() {
		startupLabel = "关闭开机自动启动 ✓"
	}
	startupText := utf16Ptr(startupLabel)
	exitText := utf16Ptr("退出 Rime Translate")

	procAppendMenuW.Call(menu, mfString, menuRelease, uintptr(unsafe.Pointer(releaseText)))
	procAppendMenuW.Call(menu, mfString, menuStartup, uintptr(unsafe.Pointer(startupText)))
	procAppendMenuW.Call(menu, mfSeparator, 0, 0)
	procAppendMenuW.Call(menu, mfString, menuExit, uintptr(unsafe.Pointer(exitText)))

	var p point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	procSetForegroundWindow.Call(hwnd)
	cmd, _, _ := procTrackPopupMenu.Call(menu, tpmRightButton|tpmReturnCmd, uintptr(p.X), uintptr(p.Y), 0, hwnd, 0)

	switch cmd {
	case menuRelease:
		if globalBridge != nil {
			go globalBridge.releaseResources()
		}
	case menuStartup:
		if err := setStartup(!startupEnabled()); err != nil {
			showMessage("Rime Translate", err.Error())
		}
	case menuExit:
		if globalBridge != nil {
			globalBridge.stop()
		}
		removeTrayIcon()
		procPostQuitMessage.Call(0)
	}
}

func wndProc(hwnd uintptr, message uint32, wParam, lParam uintptr) uintptr {
	switch message {
	case wmTray:
		switch uint32(lParam) {
		case wmRButtonUp:
			showTrayMenu(hwnd)
			return 0
		case wmLButtonDblClk:
			// Double-click: immediately free model/server resources, keep bridge running.
			if globalBridge != nil {
				go globalBridge.releaseResources()
			}
			return 0
		}
	case wmDestroy:
		removeTrayIcon()
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, uintptr(message), wParam, lParam)
	return r
}

func createHiddenWindow() (uintptr, error) {
	hInstance, _, _ := procGetModuleHandleW.Call(0)
	className := utf16Ptr("RimeTranslateHiddenWindow")
	wndProcCallback := syscall.NewCallback(wndProc)
	icon, _, _ := procLoadIconW.Call(0, idiApplication)
	cursor, _, _ := procLoadCursorW.Call(0, idcArrow)

	wc := wndClassEx{
		CbSize:        uint32(unsafe.Sizeof(wndClassEx{})),
		LpfnWndProc:   wndProcCallback,
		HInstance:     hInstance,
		HIcon:         icon,
		HCursor:       cursor,
		LpszClassName: className,
		HIconSm:       icon,
	}
	atom, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	if atom == 0 {
		return 0, err
	}

	hwnd, _, err := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(utf16Ptr("Rime Translate"))),
		0,
		0, 0, 0, 0,
		0, 0, hInstance, 0,
	)
	if hwnd == 0 {
		return 0, err
	}
	return hwnd, nil
}

func runMessageLoop() {
	var m msg
	for {
		r, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func main() {
	globalBridge = newBridge()
	go globalBridge.run()

	hwnd, err := createHiddenWindow()
	if err != nil {
		globalBridge.logf("create window failed: %v", err)
		return
	}
	globalHWND = hwnd

	if err := addTrayIcon(hwnd); err != nil {
		globalBridge.logf("add tray icon failed: %v", err)
	}

	runMessageLoop()
	globalBridge.stop()
}
