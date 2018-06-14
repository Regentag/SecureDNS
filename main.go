package main

import (
	"net"
	"strings"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const LOCALHOST = "127.0.0.1"

type SDMainWindow struct {
	*walk.MainWindow
	ifaces *walk.ComboBox
	btn    *walk.PushButton
	text   *walk.TextEdit

	dnsSvcStop SvrStopFunc
}

func (w *SDMainWindow) onWindowLoad() {
	// Attach MainWindow Closing Event Handler.
	handle := w.Closing().Attach(w.onCloseWindow)

	// startup DNS service
	stopFunc, err := RunDNS(53, func(err error) {
		// DNS server error

		// Detach closing event handler.
		// If DNS server error has occured, doesn't ask to user before closing.
		w.Closing().Detach(handle)

		// display error message
		walk.MsgBox(w, "SecureDNS 오류", "DNS 오류: "+err.Error(), walk.MsgBoxOK)

		// and close main window.
		w.Close()
	})

	if err != nil {
		walk.MsgBox(w, "SecureDNS 오류", "DNS 오류: "+err.Error(), walk.MsgBoxOK)
	} else {
		w.dnsSvcStop = stopFunc
	}

	// refresh dns info
	w.updateDNSInfo()
}

func (w *SDMainWindow) updateDNSInfo() {
	var sb strings.Builder
	sb.WriteString("Cloudflare DNS Over HTTPS가 작동중입니다.\r\n")
	sb.WriteString("안전한 DNS 사용을 위하여 네트워크 설정에서\r\n")
	sb.WriteString("시스템의 DNS 주소를 127.0.0.1로 변경하여 적용하십시오.\r\n")

	w.text.SetText(sb.String())
}

func (w *SDMainWindow) toggleSecDNS() {

}

func (w *SDMainWindow) onMenuAbout() {
	var sb strings.Builder
	sb.WriteString("SecureDNS 0.1\r\n\r\n")
	sb.WriteString("Cloudflare DNS Over HTTPS를 사용한 안전한 DNS 연결 제공\r\n")
	sb.WriteString("https://github.com/Regentag/SecureDNS")

	walk.MsgBox(w, "정보...", sb.String(), walk.MsgBoxOK)
}

func (w *SDMainWindow) onCloseWindow(canceled *bool, reason walk.CloseReason) {
	var sb strings.Builder
	sb.WriteString("Cloudflare DNS Over HTTPS가 작동중입니다.\r\n")
	sb.WriteString("SecureDNS를 종료한 후 인터넷을 사용하시려면\r\n")
	sb.WriteString("네트워크 설정에서 시스템의 DNS 주소를 본래의 주소로 되돌리십시오.\r\n")
	sb.WriteString("\r\nSecureDNS를 종료할까요?")

	yn := walk.MsgBox(w, "종료", sb.String(), walk.MsgBoxYesNo)
	if yn != 6 { // not IDYES(6)
		*canceled = true
	} else {
		if w.dnsSvcStop != nil {
			w.dnsSvcStop()
		}
	}
}

func getNetworkIfaces() []string {
	ifaces, _ := net.Interfaces()
	var sl []string

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		sl = append(sl, iface.Name)
	}

	return sl
}

func main() {
	ifaces := getNetworkIfaces() // Network interfaces

	mw := &SDMainWindow{dnsSvcStop: nil}

	MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "SecureDNS 0.1",
		Size:     Size{Width: 500, Height: 500},
		Layout:   VBox{MarginsZero: true},
		Font: Font{
			Family:    "Segoe UI",
			PointSize: 13,
		},
		MenuItems: []MenuItem{
			Menu{
				Text: "파일(&F)",
				Items: []MenuItem{
					Action{
						Text:        "Exit",
						OnTriggered: func() { mw.Close() },
					},
				},
			},
			Menu{
				Text: "도움말(&H)",
				Items: []MenuItem{
					Action{
						Text:        "SecureDNS 정보...",
						OnTriggered: mw.onMenuAbout,
					},
				},
			},
		},
		Children: []Widget{
			Composite{
				Layout: HBox{MarginsZero: false, SpacingZero: false},
				Children: []Widget{
					Label{Text: "네트워크 연결: "},
					ComboBox{
						AssignTo: &mw.ifaces,
						Model:    ifaces,
						OnCurrentIndexChanged: mw.updateDNSInfo,
					},
					HSpacer{},
					/*
						PushButton{
							Text:      "DNS 보안 켜기",
							AssignTo:  &mw.btn,
							Enabled:   false,
							OnClicked: mw.toggleSecDNS,
						},
					*/
				},
			},

			TextEdit{
				AssignTo:   &mw.text,
				HScroll:    true,
				VScroll:    true,
				ReadOnly:   true,
				Text:       "Ready.",
				Background: SolidColorBrush{Color: walk.RGB(255, 255, 240)},
			},
		},
	}.Create()
	mw.Starting().Attach(mw.onWindowLoad)
	mw.ifaces.SetCurrentIndex(0)
	mw.Run()
}
