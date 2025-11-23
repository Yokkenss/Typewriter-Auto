//easy 6er bei Büsser

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"unicode/utf8"

	"golang.org/x/sys/windows"
)

var (
	user32      = windows.NewLazySystemDLL("user32.dll")
	procKeybdEv = user32.NewProc("keybd_event")
)

const (
	KEYEVENTF_KEYUP = 0x0002
)

func pressKeyDown(vk byte) {
	pressKeyEvent(vk, false)
}

func pressKeyUp(vk byte) {
	pressKeyEvent(vk, true)
}

func pressKeyEvent(vk byte, up bool) {
	flags := uintptr(0)
	if up {
		flags = KEYEVENTF_KEYUP
	}
	procKeybdEv.Call(uintptr(vk), 0, flags, 0)
	time.Sleep(3 * time.Millisecond)
}

func pressKey(vk byte) {
	procKeybdEv.Call(uintptr(vk), 0, 0, 0)
	time.Sleep(5 * time.Millisecond)
	procKeybdEv.Call(uintptr(vk), 0, KEYEVENTF_KEYUP, 0)
}

func typeRune(r rune) {

	if r == rune(160) {
		pressKey(0x20)
		return
	}

	if r >= 'A' && r <= 'Z' {
		pressShiftCombo(byte(r))
		return
	}

	switch r {
	case ' ':
		pressKey(0x20)
	case 'ö':
		pressKey(0xC0)
	case 'Ö':
		pressKeyDown(0x10)
		pressKey(0xC0)
		pressKeyUp(0x10)

	case 'ä':
		pressKey(0xDE)
	case 'Ä':
		pressKeyDown(0x10)
		pressKey(0xDE)
		pressKeyUp(0x10)

	case 'ü':
		pressKey(0xBA)
	case 'Ü':
		pressKeyDown(0x10)
		pressKey(0xBA)
		pressKeyUp(0x10)

	case '-':
		pressKey(0xBD)

	case 'a':
		pressKey(0x41)
	case 's':
		pressKey(0x53)
	case 'd':
		pressKey(0x44)
	case 'f':
		pressKey(0x46)
	case 'j':
		pressKey(0x4A)
	case 'k':
		pressKey(0x4B)
	case 'l':
		pressKey(0x4C)

	case '\n', '\r':
		pressKey(0x0D)

	case '.':
		pressKey(0xBE)

	case ',':
		pressKey(0xBC)

	case ':':
		pressShiftCombo(0xBE)

	case ';':
		pressShiftCombo(0xBC)

	case '?':
		pressShiftCombo(0xBF)

	case '!':
		pressShiftCombo(0x31)

	case '/':
		pressKey(0xBF)

	case '\\':
		pressAltGrCombo(0xDC)

	case '0':
		pressKey(0x30)

	case '1':
		pressKey(0x31)

	case '2':
		pressKey(0x32)

	case '3':
		pressKey(0x33)

	case '4':
		pressKey(0x34)

	case '5':
		pressKey(0x35)

	case '6':
		pressKey(0x36)

	case '7':
		pressKey(0x37)

	case '8':
		pressKey(0x38)

	case '9':
		pressKey(0x39)

	case '(':
		pressShiftCombo(0x38)

	case ')':
		pressShiftCombo(0x39)

	default:
		if r >= 'a' && r <= 'z' {
			pressKey(byte(r - 'a' + 'A'))
		} else {
			log.Printf("Unknown rune: %q (%d)\n", r, r)
		}
	}
}

func pressAltGrCombo(vk byte) {
	procKeybdEv.Call(uintptr(0x11), 0, 0, 0)
	procKeybdEv.Call(uintptr(0x12), 0, 0, 0)

	procKeybdEv.Call(uintptr(vk), 0, 0, 0)
	procKeybdEv.Call(uintptr(vk), 0, KEYEVENTF_KEYUP, 0)

	procKeybdEv.Call(uintptr(0x12), 0, KEYEVENTF_KEYUP, 0)
	procKeybdEv.Call(uintptr(0x11), 0, KEYEVENTF_KEYUP, 0)
}

func pressShiftCombo(vk byte) {

	procKeybdEv.Call(uintptr(0x10), 0, 0, 0)
	time.Sleep(2 * time.Millisecond)

	procKeybdEv.Call(uintptr(vk), 0, 0, 0)
	procKeybdEv.Call(uintptr(vk), 0, KEYEVENTF_KEYUP, 0)

	procKeybdEv.Call(uintptr(0x10), 0, KEYEVENTF_KEYUP, 0)
	time.Sleep(2 * time.Millisecond)
}

func typeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	s := string(body)
	rn, _ := utf8.DecodeRuneInString(s)
	if rn == utf8.RuneError {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Typing: %q (%d)\n", rn, rn)
	typeRune(rn)

	_, _ = w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/type", typeHandler)

	addr := ":9090"
	fmt.Println("Typewriter AutoTyper by me")
	fmt.Println("Listening on http://localhost" + addr + "/type")
	fmt.Println("insert js script in console")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
