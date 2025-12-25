package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"time"
)

// --- COLORES ---
const (
	ColorReset   = "\033[0m"
	ColorCyan    = "\033[36m"
	ColorYellow  = "\033[33m"
	ColorGreen   = "\033[32m"
	ColorMagenta = "\033[35m"
)

var nodes = []string{
	"http://localhost:8080",
	"http://localhost:8081",
	"http://localhost:8082",
}

func getShard(key string) (string, int) {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	index := int(h.Sum32()) % len(nodes)
	return nodes[index], index
}

func main() {
	printProxyBanner()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		key := r.URL.Query().Get("key")

		if key == "" {
			http.Error(w, "MISSING_KEY", 400)
			return
		}

		targetNode, shardID := getShard(key)

		// Visualización de la ruta
		fmt.Printf("%s[INCOMING]%s %s %s ", ColorCyan, ColorReset, r.Method, key)
		fmt.Printf("%s----(hash)---->%s [Shard %d] %s\n", ColorYellow, ColorReset, shardID, targetNode)

		// Redirección
		redirectUrl := targetNode + r.URL.Path + "?" + r.URL.RawQuery
		proxyReq, _ := http.NewRequest(r.Method, redirectUrl, r.Body)

		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			fmt.Printf("%s[ERROR]%s Node %s is DOWN!\n", "\033[31m", ColorReset, targetNode)
			http.Error(w, "Node Unavailable", 502)
			return
		}
		defer resp.Body.Close()

		io.Copy(w, resp.Body)

		fmt.Printf("%s[DONE]%s Request completed in %v\n", ColorGreen, ColorReset, time.Since(start))
	})

	log.Fatal(http.ListenAndServe(":9000", nil))
}

func printProxyBanner() {
	fmt.Print(ColorMagenta)
	fmt.Println(`
  ____  ____   _____  ____   __
 |  _ \|  _ \ / _ \ \/ /\ \ / /
 | |_) | |_) | | | \  /  \ V / 
 |  __/|  _ <| |_| /  \   | |  
 |_|   |_| \_\\___/_/\_\  |_|  `)
	fmt.Println("  :: Load Balancer & Sharding ::")
	fmt.Println(ColorReset)
	fmt.Println("  +--- Cluster Topology ---+")
	for i, node := range nodes {
		fmt.Printf("  | Shard %d -> %s |\n", i, node)
	}
	fmt.Println("  +------------------------+")
	fmt.Println("  Listening on port :9000...")
	fmt.Println()
}
