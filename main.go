package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// --- COLORES PARA LA TERMINAL ---
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
)

type RequestData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var (
	store  *Store
	logger *Logger
	nodeID string
)

func main() {
	port := flag.String("port", "8080", "Puerto del servidor")
	dbFile := flag.String("db", "data.log", "Archivo de persistencia")
	flag.Parse()
	nodeID = "NODE-" + *port

	printBanner(*port, *dbFile)

	// Simulación de carga de módulos (Puro teatro para que se vea pro)
	logInfo("Initializing Memory Manager...")
	time.Sleep(200 * time.Millisecond)
	logInfo("Mounting Write-Ahead-Log (WAL) subsystem...")
	time.Sleep(200 * time.Millisecond)

	store = NewStore()
	var err error
	logger, err = NewLogger(*dbFile)
	if err != nil {
		logError("Failed to mount disk partition: " + err.Error())
		os.Exit(1)
	}
	defer logger.Close()

	logInfo("Replaying transaction logs for consistency...")
	logger.Recover(store)
	logSuccess(fmt.Sprintf("System State Restored. Key Count: %d", len(store.data)))

	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/set", handleSet)
	http.HandleFunc("/status", handleStatus) // Nuevo endpoint de estado

	logSuccess(fmt.Sprintf("NanoKV Shard is ONLINE on port :%s", *port))
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	key := r.URL.Query().Get("key")

	if key == "" {
		http.Error(w, "MISSING_KEY", http.StatusBadRequest)
		return
	}

	val, ok := store.Get(key)
	duration := time.Since(start)

	if !ok {
		logWarn(fmt.Sprintf("[READ] Key='%s' -> NOT FOUND (%s)", key, duration))
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	logInfo(fmt.Sprintf("[READ] Key='%s' -> FOUND (%d bytes) (%s)", key, len(val), duration))
	fmt.Fprintf(w, val)
}

func handleSet(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data RequestData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := logger.Write("SET", data.Key, data.Value); err != nil {
		logError("Disk I/O Failure: " + err.Error())
		http.Error(w, "Persistence Error", http.StatusInternalServerError)
		return
	}

	store.Set(data.Key, data.Value)
	duration := time.Since(start)

	logSuccess(fmt.Sprintf("[WRITE] Key='%s' -> PERSISTED (WAL+MEM) (%s)", data.Key, duration))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	// Endpoint para ver el estado del nodo
	status := fmt.Sprintf("Node: %s | Keys: %d | Uptime: OK", nodeID, len(store.data))
	fmt.Fprintf(w, status)
}

// --- UTILIDADES VISUALES ---

func printBanner(port, db string) {
	fmt.Print(ColorCyan)
	fmt.Println(`
  _   _                   _  ____     __
 | \ | | __ _ _ __   ___ | |/ /\ \   / /
 |  \| |/ _' | '_ \ / _ \| ' /  \ \ / / 
 | |\  | (_| | | | | (_) | . \   \ V /  
 |_| \_|\__,_|_| |_|\___/|_|\_\   \_/   `)
	fmt.Println("  :: Distributed Key-Value Store :: v1.2.0")
	fmt.Println(ColorReset)
	fmt.Printf("  %s[%s]%s Config: Port=%s | DB=%s\n", ColorPurple, time.Now().Format("15:04:05"), ColorReset, port, db)
	fmt.Println("  ------------------------------------------------")
}

func logInfo(msg string) {
	fmt.Printf("%s[INFO]%s  %s\n", ColorBlue, ColorReset, msg)
}

func logSuccess(msg string) {
	fmt.Printf("%s[OK]%s    %s\n", ColorGreen, ColorReset, msg)
}

func logWarn(msg string) {
	fmt.Printf("%s[WARN]%s  %s\n", ColorYellow, ColorReset, msg)
}

func logError(msg string) {
	fmt.Printf("%s[ERR]%s   %s\n", ColorRed, ColorReset, msg)
}
