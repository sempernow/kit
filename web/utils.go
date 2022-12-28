package web

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"html"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"gd9/prj3/kit/convert"
	"gd9/prj3/kit/id"
)

// Healthcheck to satisfy Docker healthcheck of any service endpoint
// without requiring any externals; neither (GNU) utilities nor shell.
func Healthcheck(url string) error {
	//fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		//fmt.Println("FAIL @ resp:", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println("FAIL @ body:", err)
		return err
	}
	if len(body) == 0 {
		//fmt.Println("FAIL @ len(body):", 0)
		return errors.New("1")
	}
	//fmt.Println("Success : body:", string(body))
	return nil
}

// Nonce of size bytes as a URL-encoded string.
func Nonce(size int) (string, error) {
	bb, err := id.Nonce(size)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bb), nil
}

// SRI generates the Sub-Resource Integrity (CSP/SRI) string of a resource per SHA384.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy
// https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity
func SRI(bb []byte) string {
	bb48 := sha512.Sum384(bb)
	return "sha384-" + base64.StdEncoding.EncodeToString(bb48[:])
}

// GetHostname returns os.Hostname(); the Docker container ID.
func GetHostname() string {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}
	return host
}

// GetForwardedHostIP from X-Forwarded-Host request header. Also see GetInboundIP(..) .
func GetForwardedHostIP(r *http.Request) (string, error) {
	ips, _ := net.LookupIP(r.Header.Get("X-Forwarded-Host"))
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return fmt.Sprint(ipv4), nil
		}
	}
	return "", errors.New("No valid ip found")
}

// GetInboundIP from request headers; X-REAL-IP > X-FORWARDED-FOR > r.RemoteAddr
// (Reverse Proxy Servers may set X-REAL-IP and X-FORWARDED-FOR headers.)
// https://golangbyexample.com/golang-ip-address-http-request/
//
//	PROBLEM:
//	@ Docker, service ctnrs have no access to client-IP address!
//	IP is of overlay network LB @ "Containers:".
//
// Client IP Address is MASKED by Docker Swarm
// "... to service containers, all incoming traffic appears to arrive from
// the same set of private network ingress network node IPs ...
// service containers cannot distinguish individual clients by IP."
//
//	SOLUTION 1: Set Nginx "ports:" (@ YAML) of "80:${CTNR_PORT}" to "host" mode.
//	SOLUTION 2: "https://github.com/newsnowlabs/docker-ingress-routing-daemon"
//	SOLUTION 3: Prefetch client side (async); set as header: X-Client-Ip
//		Prefetch @ "https://api.ipify.org?format=json"
func GetInboundIP(r *http.Request) (string, error) {

	// Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-Real-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	// Else get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	// Else get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	return "", errors.New("No valid ip found")
}

// GetOutboundIP returns the preferred outbound IP address of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "UNKNOWN"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

// GetOutboundAddr returns the preferred outbound IP address of this machine
func GetOutboundAddr() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "UNKNOWN"
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

// RuneToHTMLEntity https://gist.github.com/brandonrachal/10605780
//
//	'𝓢' //=> &#1d4e2;
func RuneToHTMLEntity(entity rune) string {
	if entity < 128 {
		return string(entity)
	}
	return "&#" + strconv.FormatInt(int64(entity), 16) + ";"
}

// StringToHTMLEntities ...
//
//	"𝓢 foo" //=> &#1d4e2; foo
func StringToHTMLEntities(s string) string {
	var encoded string
	for _, x := range s {
		encoded += RuneToHTMLEntity(x)
	}
	return encoded
}

// HTMLEntitiesToString ...
//
//	"&ldquo;Joe&#39;s Diner&rdquo; &#x2627; &lt;bogus@email.addr&gt;"
//	//=> “Joe's Diner” ☧ <bogus@email.addr>
func HTMLEntitiesToString(s string) string {
	return html.UnescapeString(s)
}

// ValidateJS requires Node.js and returns error on invalid script,
// invalid path, or absence of Node.js, which acts as the validator.
func ValidateJS(file string) error {
	cmd := "node"
	//target := "cfg.js"
	//file, _ := filepath.Abs(filepath.Join(h.meta.PathWebRoot, "sa", "scripts", target))
	_, err := ioutil.ReadFile(file)
	if err != nil {
		panic(errors.Wrap(err, "unreadable file @ "+file))
	}
	out, err := exec.Command(cmd, file).CombinedOutput()
	if err != nil {
		f := strings.Split(file, "/")
		msg := fmt.Sprintf("\nFile: %v\n ERR: %v\n", f[len(f)-1], convert.BytesToString(out))
		return errors.Wrap(err, "invalid script or Node.js not installed @ "+msg)
		// Response body: {"error":"Internal Server Error"}
	}
	return nil
}
