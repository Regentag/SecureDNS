package main

import (
	"io"
	"log"
	"os"

	"github.com/miekg/dns"
)

const CLOUDFLARE_DNS = "1.1.1.1:53"
const CLOUDFLARE_DOH_HOST = "cloudflare-dns.com."

const LOCAL_CACHE_FILE = "cf-host.bin"

func GetCloudflareAddress(forceQuery bool) (*dns.Msg, error) {
	cfAddrCache := GetMyDataPath(LOCAL_CACHE_FILE)

	if !forceQuery && CheckFileExists(cfAddrCache) {
		// 파일에 저장된 데이터를 읽어와 DNS 응답 메시지로 반환한다.
		return getDohHostAddrFromFile(cfAddrCache)
	} else {
		// DNS 서버에 주소를 질의하고 결과를 반환한다.
		// 결과물에 이상이 없으면 파일에 기록한다.

		r, dnsErr := getDohHostAddrFromDNS()

		if dnsErr != nil {
			WriteErrorLogMsg("DNS에서 Cloudflare DOH 서버 주소 질의 실패.", dnsErr)
			return nil, dnsErr
		}

		_ = writeDohHostAddrToFile(cfAddrCache, r)

		return r, dnsErr
	}
}

func getDohHostAddrFromDNS() (*dns.Msg, error) {
	log.Println("Cloudflare DOH 서버의 IP 주소를 질의합니다...")

	m := new(dns.Msg)
	m.SetQuestion(CLOUDFLARE_DOH_HOST, dns.TypeA)

	client := new(dns.Client)
	r, _, err := client.Exchange(m, CLOUDFLARE_DNS)
	return r, err
}

func getDohHostAddrFromFile(filename string) (*dns.Msg, error) {
	log.Println("Cloudflare DOH 서버의 IP 주소 파일에서 불러옵니다...")
	log.Println("  Cache file = " + filename)

	f, e := os.Open(filename)
	if e != nil {
		WriteErrorLogMsg("저장된 CF 데이터를 열지 못함.", e)
		return nil, e
	}
	defer f.Close()

	bin, e := io.ReadAll(f)
	if e != nil {
		WriteErrorLogMsg("저장된 CF 데이터를 읽지 못함.", e)
		return nil, e
	}

	r := new(dns.Msg)
	e = r.Unpack(bin)
	if e != nil {
		WriteErrorLogMsg("저장된 CF 데이터를 DNS 메시지로 변환하지 못함.", e)
		return nil, e
	}

	return r, e
}

func writeDohHostAddrToFile(filename string, msg *dns.Msg) error {
	bin, e := msg.Pack()

	if e != nil {
		WriteErrorLogMsg("DNS 메시지 저장을 위한 Packing 실패.", e)
		return e
	}

	// 파일에 기록한다.
	f, e := os.Create(filename)

	if e != nil {
		WriteErrorLogMsg("DNS 메시지 저장을 위한 파일 열기 실패.", e)
		return e
	}

	defer f.Close()

	_, e = f.Write(bin)

	if e != nil {
		WriteErrorLogMsg("DNS 메시지 파일에 쓰기 실패.", e)
	}

	return e
}
