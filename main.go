package main

import (
	"flag"
	"log"

	syslog "gopkg.in/mcuadros/go-syslog.v2"
)

var (
	debug      bool
	serverName string
)

func receiveSyslog(ch syslog.LogPartsChannel) {
	var (
		l   *logEntry
		err error
	)

	for msg := range ch {
		if debug {
			log.Printf("%s: %s", msg["hostname"], msg["content"])
		}

		if l, err = parseSyslogMessage(msg); err != nil {
			log.Printf("Unable to parse message: %s", err)
			continue
		}

		if l.jrpc_method != "" {
			l.uri = l.jrpc_method
		} else {
			l.jrpc_method = ""
			l.uri = "/"
		}

		prometheusMetricsRegister(l)
	}
}

func main() {
	var (
		listenSyslog  string
		listenHTTP    string
		usePrometheus bool
	)

	flag.StringVar(&listenSyslog, "listenSyslog", "0.0.0.0:514", "ip:port to listen for syslog messages")
	flag.StringVar(&listenHTTP, "listenHTTP", "0.0.0.0:9100", "ip:port to listen for http requests")
	flag.StringVar(&serverName, "serverName", "", "Server name to use in Prometheus metrics (leave empty for dynamic hostname)")
	flag.BoolVar(&usePrometheus, "usePrometheus", true, "Enable posting metrics to Prometheus")
	flag.BoolVar(&debug, "debug", false, "Enable debug")
	flag.Parse()

	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	log.Printf("Starting Syslog UDP listener: %s", listenSyslog)
	server := syslog.NewServer()
	server.SetFormat(syslog.RFC3164)
	server.SetHandler(handler)
	server.ListenUDP(listenSyslog)
	server.Boot()

	go receiveSyslog(channel)

	if usePrometheus {
		log.Printf("Starting Prometheus HTTP listener: %s", listenHTTP)
		go prometheusListener(listenHTTP)
	}
	server.Wait()
}
