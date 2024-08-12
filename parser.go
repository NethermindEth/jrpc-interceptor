package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type logEntry struct {
	// Tags
	server        		string
	scheme        		string
	method        		string
	hostname      		string
	status        		string
	protocol      		string
	uri           		string
	jrpc_method   		string

	// Fields
	clientIP      		net.IP
	duration    		float64
	bytesSent   		uint64
	bytesReceived      	uint64
	response_duration   float64
}


func parseSyslogMessage(msg format.LogParts) (l *logEntry, err error) {
	content := msg["content"].(string)

	chunks := strings.Split(content, "|")
	if len(chunks) != 12 {
		return nil, fmt.Errorf("wrong number of fields in message: %s", content)
	}

	l = &logEntry{
		server:        msg["hostname"].(string),
		scheme:        chunks[1],
		hostname:      chunks[2],
		method:        chunks[3],
		protocol:      chunks[4],
		uri:           strings.Split(chunks[5], "?")[0],
		status:        chunks[6],
		jrpc_method:   chunks[11],
	}

	if l.clientIP = net.ParseIP(chunks[0]); l.clientIP == nil {
		return nil, fmt.Errorf("unable to parse clientIP")
	}

	if l.duration, err = strconv.ParseFloat(chunks[7], 64); err != nil {
		return nil, fmt.Errorf("unable to parse request duration as float: %s", err)
	}

	if l.bytesReceived, err = strconv.ParseUint(chunks[8], 10, 64); err != nil {
		return nil, fmt.Errorf("unable to parse bytesReceived as uint: %s", err)
	}

	if l.bytesSent, err = strconv.ParseUint(chunks[9], 10, 64); err != nil {
		return nil, fmt.Errorf("unable to parse bytesSent as uint: %s", err)
	}

	if l.response_duration, err = strconv.ParseFloat(chunks[10], 64); err != nil {
		return nil, fmt.Errorf("unable to parse response duration as float: %s", err)
	}
	return
}