// Package inbound implements dns server for inbound connection.
package main

import (
	"os"
	"sync"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/miekg/dns"
)

type Server struct {
	BindAddress string
	Protocols []string
	Forwarders []*Forwarder
	Cache *Cache
}

func (server *Server) Run() {

	mux := dns.NewServeMux()
	mux.Handle(".", server)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	log.Info("Start octodns on " + server.BindAddress)

	for _, protocol := range server.Protocols {
		go func(protocol string) {
			err := dns.ListenAndServe(server.BindAddress, protocol, mux)
			if err != nil {
				log.Fatal("Listen "+protocol+" failed: ", err)
				os.Exit(1)
			}
		}(protocol)
	}

	wg.Wait()
}

func (server *Server) ServeDNS(responseWriter dns.ResponseWriter, question *dns.Msg) {

	inboundIP, _, _ := net.SplitHostPort(responseWriter.RemoteAddr().String())
	proto := responseWriter.RemoteAddr().Network()

	// Check log level before so we don't evaluate the string literal.
	if isDebug() {
		log.Debug("Question from " + inboundIP + "(" + proto + ")" + ": " + question.Question[0].String())
	}

	found := false
	lastResponse := server.Cache.Hit(Key(question.Question[0]), question.Id)
	fromCache := false

	if lastResponse != nil {
		if isDebug() {
			log.Debug("Cache Hit: " + Key(question.Question[0]))
		}
		found = true
		fromCache = true
	}



	for _, forwarder := range server.Forwarders {
		isUdp := false
		isTcp := false
		for _, protocol := range forwarder.Protocols {
			if protocol == "udp" {
				isUdp = true
			}else if protocol == "tcp" {
				isTcp = true
			}
		}

		// Always try UDP first, unless it's disabled
		if isUdp && !found {
			client := &dns.Client{
				Net:     "udp",
				Timeout: time.Duration(forwarder.Timeout) * time.Second,
			}

			res, _, err := client.Exchange(question, forwarder.Address)
			lastResponse = res

			if err != nil{
				if err == dns.ErrTruncated {
					// Truncated, so just try TCP instead.

					log.Debug("Reply truncated, trying with TCP")
					isTcp = true
				}
				log.Warn("An error occurred connecting to " + forwarder.Name + "(UDP" + forwarder.Address + "): ", err)
			}

			if lastResponse == nil {
				lastResponse = question.Copy()
				lastResponse.Rcode = dns.RcodeServerFailure
			} else if lastResponse.Rcode == dns.RcodeNameError {
				// Don't continue to try TCP if the connection succeeded.
				isTcp = false
			}

			if isDebug() {
				for _, answer := range lastResponse.Answer {
					log.Debug("Answer from " + forwarder.Name + "(UDP" + forwarder.Address + "): " + answer.String())
				}
			}

			if lastResponse.Rcode == dns.RcodeSuccess {
				found = true
			}

		}

		if isTcp && !found{
			client := &dns.Client{
				Net:     "tcp",
				Timeout: time.Duration(forwarder.Timeout) * time.Second,
			}

			res, _, err := client.Exchange(question, forwarder.Address)
			lastResponse = res

			if err != nil{
				log.Warn("An error occurred connecting to " + forwarder.Name + "(TCP:" + forwarder.Address + "): ", err)
			}

			if lastResponse == nil {
				lastResponse = question.Copy()
				lastResponse.Rcode = dns.RcodeServerFailure
			}

			if isDebug() {
				for _, answer := range lastResponse.Answer {
					log.Debug("Answer from " + forwarder.Name + "(UDP" + forwarder.Address + "): " + answer.String())
				}
			}

			if lastResponse.Rcode == dns.RcodeSuccess {
				found = true
			}
		}
	}

	if !fromCache {
		server.Cache.InsertMessage(Key(question.Question[0]), lastResponse)
	}

	if isDebug() {
		log.Debug("FINAL Answer for " + inboundIP + "(" + proto + ")" + ": " + dns.RcodeToString[lastResponse.Rcode])
	}


	responseWriter.WriteMsg(lastResponse)

}


func isDebug() bool{
	return log.GetLevel() >= log.DebugLevel
}
