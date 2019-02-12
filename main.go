package main

import (
	"golang.org/x/sys/windows/svc"
	//"golang.org/x/sys/windows/svc/debug"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path/filepath"
	"time"
)

type ServContext struct {
	dnsSvcStop SvrStopFunc
}

// svc.Handler 인터페이스 구현
func (srv *ServContext) Execute(args []string, req <-chan svc.ChangeRequest, stat chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	stat <- svc.Status{State: svc.StartPending}

	// 실제 서비스 내용
	stopChan := make(chan bool, 1)
	srv.runBody()

	stat <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

LOOP:
	for {
		// 서비스 변경 요청에 대해 핸들링
		switch r := <-req; r.Cmd {
		case svc.Stop, svc.Shutdown:
			stopChan <- true
			break LOOP

		case svc.Interrogate:
			stat <- r.CurrentStatus
			time.Sleep(100 * time.Millisecond)
			stat <- r.CurrentStatus

			//case svc.Pause:
			//case svc.Continue:
		}
	}

	stat <- svc.Status{State: svc.StopPending}

	// Stop DNS server
	log.Println("Shutting down...")

	de := srv.dnsSvcStop()
	if de != nil {
		WriteErrorLogMsg("DNS service shutdown error: ", de)
	} else {
		log.Println("DNS service stopped.")
	}

	log.Println("SecDNS was stopped.")
	return
}

func (srv *ServContext) runBody() {
	// DNS 서버를 go routine으로 시작하고
	// 서버 종료를 위한 함수를 얻어 저장한다.
	stopFunc, err := RunDNS(53, func(err error) {
		WriteErrorLogMsg("DNS service error: ", err)
	})

	if err != nil {
		WriteErrorLogMsgF("Can't start DNS service. ", err)
	} else {
		srv.dnsSvcStop = stopFunc
		log.Println("SecDNS service started.")
	}
}

func logPath(filename string) string {
	ex, err := os.Executable()
	if err != nil {
		ex = "." // current working directory
	}

	exPath := filepath.Dir(ex)
	return filepath.Join(exPath, filename)
}

func main() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   logPath("sec-dns.log"),
		MaxSize:    10, // megabytes
		MaxBackups: 1,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	})

	log.Println("Initializing...")

	err := svc.Run("SecDNS", &ServContext{})
	//err = debug.Run("DummyService", &dummyService{}) //콘솔출력 디버깅시
	if err != nil {
		WriteErrorLogMsgF("Fatal service error: ", err)
		panic(err)
	}
}
