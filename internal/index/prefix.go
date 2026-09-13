package index

//import "fmt"

// T can be a string or a slice of bytes
func AppendByte[T ~string | ~[]byte](prefix T, b byte) T {
	// We convert to []byte to perform the append, then back to T
	// This works because both ~string and ~[]byte can be cast to []byte
	res := append([]byte(prefix), b)
	return T(res)
}

func LongestCommonPrefix[T custom](s1, s2 T) T {
	var prefix T = s1[:0] // initialize prefix with the same type as s1 but with length 0
	for i := 0; i < Min_len(s1, s2); i++ {
		if s1[i] == s2[i] {
			prefix = AppendByte(prefix, s1[i])
		} else {
			break
		}
	}
	return prefix
}

type custom interface {
	string | []byte
}

func Min_len[T custom](s1, s2 T) int {
	if len(s1) < len(s2) {
		return len(s1)
	} else {
		return len(s2)
	}
}
