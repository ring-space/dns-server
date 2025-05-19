package records

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// Client хранит URL и кэш записей
type Client struct {
	url     string
	rec     map[string][]string
	mu      sync.RWMutex
	refresh time.Duration
}

// NewClient создаёт клиента
func NewClient(url string, refresh time.Duration) *Client {
	return &Client{url: url, rec: make(map[string][]string), refresh: refresh}
}

// StartAutoRefresh запускает периодический фетч
func (c *Client) StartAutoRefresh() {
	c.fetch()
	for range time.Tick(c.refresh) {
		c.fetch()
	}
}

// fetch загружает записи
func (c *Client) fetch() {
	resp, err := http.Get(c.url)
	if err != nil {
		log.Printf("Ошибка запроса к %s: %v", c.url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка чтения тела ответа: %v", err)
		return
	}

	// вот тут выводим чистый ответ (может быть JSON или что-то ещё)
	log.Printf("[DEBUG] raw response: %s", body)

	var data map[string][]string
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		return
	}

	c.mu.Lock()
	c.rec = data
	c.mu.Unlock()
	log.Printf("[REFRESH] загружено %d доменов", len(data))
}

// Get возвращает записи для домена
func (c *Client) Get(name string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rec[name]
}
