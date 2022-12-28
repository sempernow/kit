package web_test

import (
	"net"
	"os"
	"strings"
	"testing"

	"github.com/sempernow/kit/testkit"
	"github.com/sempernow/kit/web"
)

func TestSRI(t *testing.T) {
	t.Log("@ CSP/SRI ...")
	bb, err := os.ReadFile("testSRI.txt")
	testkit.Log(t, "Read resource file", err)
	got := web.SRI(bb)
	exp := "sha384-AZcucad51j6X9oVsNdkqjE23BWuEY1t5PdGmVP6EYYEz3ukrI8xIzdtYkc9lLnot"
	testkit.LogDiff(t, "Sign it", got, exp)
}
func TestGetOutboundIP(t *testing.T) {

	t.Skip()

	t.Log(web.GetOutboundIP())
	t.Log(web.GetOutboundAddr())
	t.Log(web.GetHostname())
	host, _ := net.LookupHost("localhost")
	t.Log(`LookupHost("localhost"):`, host)
	iplocal, _ := net.LookupIP("localhost")
	t.Log(`LookupIP("localhost")  :`, iplocal)
	ip := strings.Split(web.GetOutboundAddr(), ":")
	names, _ := net.LookupAddr(ip[0])
	t.Log(`LookupAddr("`+ip[0]+`"):`, names[0])
	names, _ = net.LookupAddr("0.0.0.0")
	t.Log(`LookupAddr("0.0.0.0"):`, names[0])
	names, _ = net.LookupAddr("127.0.0.1")
	t.Log(`LookupAddr("127.0.0.1"):`, names[0])

	// kit_test.go:39: 192.168.1.101
	// kit_test.go:40: 192.168.1.101
	// kit_test.go:41: XPC
	// kit_test.go:43: LookupHost("localhost"): [127.0.0.1]
	// kit_test.go:45: LookupIP("localhost")  : [127.0.0.1]
	// kit_test.go:48: LookupAddr("192.168.1.101"): host.docker.internal.
	// kit_test.go:50: LookupAddr("0.0.0.0"): ___id___.c.mystat-in.net.
	// kit_test.go:52: LookupAddr("127.0.0.1"): localhost
	// kit_test.go:56: LookupCNAME("0.0.0.0"):
	// kit_test.go:58: LookupCNAME("192.168.1.101"):
	// kit_test.go:60: LookupTXT("localhost"): []
	// kit_test.go:62: LookupTXT("192.168.1.101"): []
	// kit_test.go:64: LookupTXT("192.168.1.101"): []

	// Below return NOTHING:
	dns, _ := net.LookupCNAME("0.0.0.0")
	t.Log(`LookupCNAME("0.0.0.0"):`, dns)
	dns, _ = net.LookupCNAME(ip[0])
	t.Log(`LookupCNAME("`+ip[0]+`"):`, dns)
	txt, _ := net.LookupTXT(names[0])
	t.Log(`LookupTXT("`+names[0]+`"):`, txt)
	txt, _ = net.LookupTXT(ip[0])
	t.Log(`LookupTXT("`+ip[0]+`"):`, txt)
	txt, _ = net.LookupTXT(web.GetOutboundAddr())
	t.Log(`LookupTXT("`+web.GetOutboundAddr()+`"):`, txt)
	//t.Fatal("=== end")
}
