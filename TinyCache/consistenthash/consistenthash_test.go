package consistenthash

import (
	"strconv"
	"testing"
)

func TestConsistentHash(t *testing.T) {
	m := NewMap(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	m.Add("2", "4", "6")

	testCase := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCase {
		if m.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	m.Add("8")
	testCase["27"] = "8"

	for k, v := range testCase {
		if m.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
