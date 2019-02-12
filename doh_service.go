package main

// https://developers.cloudflare.com/1.1.1.1/dns-over-https/wireformat/

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/miekg/dns"
)

/// Declare error clase
type DohError struct {
	msg string
}

func (e *DohError) Error() string {
	return e.msg
}

func newErr(msg string) error {
	return &DohError{msg}
}

///
const CLOUDFLARE_DNS = "1.1.1.1:53"
const CLOUDFLARE_DOH_HOST = "cloudflare-dns.com."
const CLOUDFLARE_DOH_URL = "https://cloudflare-dns.com/dns-query"

// Create HTTPS request and POST.
func makeHttpsRequest(wire []byte) (respWire []byte, err error) {
	// disable security check for client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	buff := bytes.NewBuffer(wire)

	resp, err := client.Post(CLOUDFLARE_DOH_URL,
		"application/dns-udpwireformat", buff)

	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil, newErr("HTTP error code " + resp.Status)
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return respBody, nil
		} else {
			// io: read error
			return nil, err
		}
	} else {
		// http error
		return nil, err
	}
}

type SecHandler struct {
	ServiceType string
	Host        *dns.Msg
}

func (s SecHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	// Check for "cloudflare-dns.com" host.
	if len(r.Question) > 0 {
		if r.Question[0].Name == CLOUDFLARE_DOH_HOST &&
			r.Question[0].Qtype == dns.TypeA {
			s.Host.SetReply(r)
			w.WriteMsg(s.Host)

			// End func.
			return
		}
	}

	wire, err := r.Pack()

	if err == nil {
		resp, err := makeHttpsRequest(wire)
		if err == nil {
			// Good response then
			m := new(dns.Msg)
			err := m.Unpack(resp)
			if err == nil {
				m.SetReply(r)
				w.WriteMsg(m)
			} else {
				// Can't unpack message from wireformat.
				dns.HandleFailed(w, r)
			}
		} else {
			// HTTPS Request failed.
			dns.HandleFailed(w, r)
		}
	} else {
		// Can't pack message to wire format
		dns.HandleFailed(w, r)
	}
}

func getDohHostAddr() (*dns.Msg, error) {
	m := new(dns.Msg)
	m.SetQuestion(CLOUDFLARE_DOH_HOST, dns.TypeA)

	client := new(dns.Client)
	r, _, err := client.Exchange(m, CLOUDFLARE_DNS)
	return r, err
}

type SvrStopFunc func() error
type SvrErrorHandlerFunc func(err error)

func RunDNS(port int, errHandler SvrErrorHandlerFunc) (SvrStopFunc, error) {
	// get DOH host address
	h, e := getDohHostAddr()
	if e != nil {
		return nil, newErr("Failed to lookup Cloudflare DOH server address. " +
			e.Error())
	}

	handler := SecHandler{"UDP", h}
	srv := new(dns.Server)
	srv.Addr = ":" + strconv.Itoa(port)
	srv.Net = "udp"
	srv.Handler = handler

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errHandler(err)
		}
	}()

	return func() error {
		if srv != nil {
			return srv.Shutdown()
		} else {
			return newErr("No DNS server instance.")
		}
	}, nil
}
