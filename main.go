package main

import (
	"TinyCache"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *TinyCache.Group {
	return TinyCache.NewGroup("scores", 2<<20, TinyCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startApiServer(apiAddr string, group *TinyCache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := group.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("frontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func startCacheServer(addr string, addrs []string, group *TinyCache.Group) {
	peers := TinyCache.NewHTTPPool(addr)
	peers.Set(addrs...)
	group.RegisterPeers(peers)
	log.Println("tinycache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Tinycache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	group := createGroup()
	if api {
		go startApiServer(apiAddr, group)
	}

	startCacheServer(addrMap[port], addrs, group)
}
