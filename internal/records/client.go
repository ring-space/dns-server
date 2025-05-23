package records

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// Client хранит базовый URL, токен и кэш записей.
type Client struct {
	baseURL string
	token   string
	rec     map[string][]string
	mu      sync.RWMutex
	refresh time.Duration
}

// NewClient создаёт клиента с базовым URL, интервалом обновления и токеном.
func NewClient(baseURL string, refresh time.Duration, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		rec:     make(map[string][]string),
		refresh: refresh,
	}
}

// register отправляет токен на endpoint /api/v1/dns/register.
func (c *Client) register() {
	url := c.baseURL + "/api/v1/dns/register"
	payload := map[string]string{"token": c.token}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Ошибка маршалинга JSON для регистрации: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Ошибка создания POST-запроса для регистрации: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Ошибка при регистрации: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Регистрация вернула %d: %s", resp.StatusCode, body)
		return
	}
	log.Printf("Успешная регистрация с токеном")
}

// StartAutoRefresh сначала регистрируется, затем периодически обновляет записи.
func (c *Client) StartAutoRefresh() {
	c.register()
	c.fetch()
	ticker := time.NewTicker(c.refresh)
	for range ticker.C {
		c.fetch()
	}
}

// fetch делает GET на /zones с X-Api-Key и парсит JSON-ответ.
func (c *Client) fetch() {
	url := c.baseURL + "/zones"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Ошибка создания GET-запроса к %s: %v", url, err)
		return
	}
	req.Header.Set("X-Api-Key", c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Ошибка при запросе %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка чтения тела ответа: %v", err)
		return
	}
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

// Get возвращает список IP для данного доменного имени.
func (c *Client) Get(name string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rec[name]
}
