package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// records — статические DNS-записи (домены → список IP)
var records = map[string][]string{
	"example.com.":  {"93.184.216.34", "198.51.100.1"},
	"test.local.":   {"192.168.0.100", "192.168.1.100"},
	"devgitops.ru.": {"217.25.226.183"},
}

// recordsHandler отдаёт JSON со списком записей
func recordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(records); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	// Регистрация обработчика
	http.HandleFunc("/records", recordsHandler)

	addr := ":8080"
	log.Printf("Starting Records Service on %s…", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Records Service failed: %v", err)
	}
}
