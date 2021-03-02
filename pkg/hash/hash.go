package hash

import (
	"fmt"
	"hash/fnv"
)

// FNV64a return 64bit hash
func FNV64a(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum64())
}
