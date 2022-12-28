package mid

import (
	"fmt"
	"strings"
)

const DBUG = true

func dbug(iif bool, msg ...interface{}) {
	if iif {
		fmt.Println("=== mid :", fmt.Sprint(msg...))
	}
}

// ============================================================
// pwa       | r.Host         : localhost:3030
// pwa       | r.Referer()    : http://localhost:3030/test
//
// r.Header.Get("Referer")    : http://localhost:3030/test
// r.Header.Get("Origin")     : http://localhost:3030
// ============================================================

// Match request header 'Referer' or 'Origin' value to any app-service origin.
func matchOrigin(origins []string, rO string) bool {
	if rO == "" {
		return true
	}
	var permit bool
	rO = rO + "/"
	for _, origin := range origins {
		dbug(DBUG, "matchOrigin : REQ: ", rO, " ALLOW: ", origin+"/")
		if strings.Contains(rO, origin+"/") {
			permit = true
		}
	}
	return permit
}

// Match request header 'Host' value to any app-service host.
func matchHost(hosts []string, rH string) bool {
	if rH == "" {
		return true
	}
	var permit bool
	rH = "/" + rH + "/"
	for _, host := range hosts {
		dbug(DBUG, "matchHost : REQ: ", rH, " ALLOW: ", "/"+host+"/")
		if strings.Contains("/"+host+"/", rH) {
			permit = true
		}
	}
	return permit
}

// ======================================================
// pwa       | =================
// pwa       | r.TLS          : nil
// pwa       | r.URL          : /test
// pwa       | r.URL.String() : /test
// pwa       | r.URL.Query()  : map[]
// pwa       | r.URL.Path     : /test
// pwa       | r.URL.Fragment :
// pwa       | r.Method       : GET
// pwa       | r.RemoteAddr   : 172.29.0.3:59838
// pwa       | r.Referer()    : http://localhost:3030/test
// pwa       | r.Host         : localhost:3030
// pwa       | r.RequestURI   : /test
// pwa       | r.Proto        : HTTP/1.0
// pwa       | r.UserAgent    : Mozilla/5.0 (Win...
//
// https://golang.org/pkg/net/http/#Request
// ======================================================

// ====================================================================
// Golang REMOVEs the Host header from the client-request headers map!
// ====================================================================
// All we have is r.Host as a proxy for the Host header if present,
// else r.Host is somthing else:
// "
//  ... is either the value of the "Host" header
//      or the host name given in the URL itself.
// "
// Inexplicably, Golang hides the Host header itself;
// they "promote" it to that morphodite (r.Host).
// IOW, we have NO WAY OF KNOWING if the client sent a Host header.
// https://pkg.go.dev/net/http?utm_source=godoc#Request.Header
// ====================================================================
