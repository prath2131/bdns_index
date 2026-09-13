package index

import (
	"bytes"
)

func reverseDomain(dns []byte) []byte {
	dns = bytes.TrimSpace(dns)
	prev_dot_index := len(dns)
	reversedomain := make([]byte, 0, len(dns))
	for i := len(dns) - 1; i > -1; i-- {
		if 'A' <= dns[i] && dns[i] <= 'Z' {
			dns[i] += 32
		}
		if dns[i] == '.' {
			reversedomain = append(reversedomain, dns[i+1:prev_dot_index]...)
			prev_dot_index = i
		}
	}
	reversedomain = append(reversedomain, dns[:prev_dot_index]...)
	return reversedomain
}
