// +build windows

// The Ripple API on Windows is not officially supported and you're probably
// gonna swear a lot if you intend to use it on Windows. Caveat emptor.

package main

import (
	"log"
	"net"

	"github.com/valyala/fasthttp"
	"zxq.co/ripple/rippleapi/common"
)

func startuato(hn fasthttp.RequestHandler) {
	conf, _ := common.Load()
	var (
		l   net.Listener
		err error
	)

	// Listen on a TCP or a UNIX domain socket.
	if conf.Unix {
		l, err = net.Listen("unix", conf.ListenTo)
	} else {
		l, err = net.Listen("tcp", conf.ListenTo)
	}
	if nil != err {
		log.Fatalln(err)
	}

	fasthttp.Serve(l, hn)
}
