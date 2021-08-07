package netutils

import (
	"strings"

	"github.com/miekg/dns"
)

const (
	DNS_TRAILING_DOT = "."
)

func GetCNameByDomainName(domainName string) (string, error) {

	//config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	config, _ := dns.ClientConfigFromReader(strings.NewReader("search domain.name\nnameserver 8.8.4.4\nnameserver 8.8.8.8"))
	c := new(dns.Client)
	m := new(dns.Msg)

	// Note the trailing dot. miekg/dns is very low-level and expects canonical names.
	m.SetQuestion(domainName+DNS_TRAILING_DOT, dns.TypeCNAME)
	m.RecursionDesired = true
	r, _, _ := c.Exchange(m, config.Servers[0]+":"+config.Port)

	//fmt.Println(r.Answer[0].(*dns.CNAME).Target) // domainName+"."
	// fmt.Println(r.Answer[0])
	if r != nil && len(r.Answer) > 0 {
		answer := r.Answer[0]
		switch v := answer.(type) {
		case *dns.CNAME:
			cname := v
			if cname != nil {
				return strings.TrimSuffix(cname.Target, DNS_TRAILING_DOT), nil
			}
			break
		}
	}

	return "", nil
}
