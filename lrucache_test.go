package TinyCache

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	g := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slow DB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[v]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		if view, err := g.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of cache")
		}
		if _, err := g.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if view, err := g.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}

func TestHTTPPool(t *testing.T) {
	NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slow DB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	addr := "localhost:9999"
	peers := newHTTPPool(addr)
	log.Println("tinycache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
