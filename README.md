# DNS Geo-Routing Server

Простой DNS-сервер на Go с поддержкой:

* A‑запросов с выбором ближайшего сервера по GeoIP2
* Горячей подгрузкой DNS‑записей из внешнего JSON‑сервиса
* Экспортом метрик по странам клиентов в Prometheus

---

## Требования

* Go 1.18+
* Файл GeoLite2 City `.mmdb`
* HTTP‑сервис, отдающий JSON вида `map[string][]string`

---

## Переменные окружения и флаги

Все параметры можно задавать через переменные окружения или через флаги (флаги имеют приоритет над ENV).

| Параметр       | ENV                   | Флаг             | Описание                              | Значение по умолчанию |
| -------------- | --------------------- | ---------------- | ------------------------------------- | --------------------- |
| DNS порт       | `DNS_PORT`            | `-port`          | адрес и порт DNS‑сервера              | `:53`                 |
| URL записей    | `RECORDS_SERVICE_URL` | `-records-url`   | URL JSON‑сервиса с записями           | **обязателен**        |
| GeoIP DB       | `GEOIP_DB_PATH`       | `-geoip-db`      | путь к GeoLite2‑City.mmdb             | **обязателен**        |
| Интервал фетча | `REFRESH_INTERVAL`    | `-refresh`       | как часто обновлять записи            | `30s`                 |
| Метрики порт   | `METRICS_PORT`        | `-metrics-port`  | порт HTTP‑сервера метрик (`/metrics`) | `:9090`               |
| Flush метрик   | `METRICS_FLUSH`       | `-metrics-flush` | интервал сброса и публикации метрик   | `5m`                  |

---

## Локальная сборка и запуск

```bash
# Клонируем и переходим в репозиторий
git clone https://github.com/youruser/dns-server.git
cd dns-server

# Скачиваем зависимости и собираем бинарник
go mod download
go build -o bin/dns-server ./cmd/dns-server

# Копируем GeoIP базу рядом с бинарником
cp /path/to/GeoLite2-City.mmdb ./

# Задаём ENV и запускаем сервис
export RECORDS_SERVICE_URL=http://localhost:8080/records
export GEOIP_DB_PATH=./GeoLite2-City.mmdb
export DNS_PORT=:53
export REFRESH_INTERVAL=45s
export METRICS_PORT=:9191
export METRICS_FLUSH=2m

# Запуск (флаги не обязательны, можно их передать для переопределения)
./bin/dns-server \
  -port $DNS_PORT \
  -records-url $RECORDS_SERVICE_URL \
  -geoip-db $GEOIP_DB_PATH \
  -refresh $REFRESH_INTERVAL \
  -metrics-port $METRICS_PORT \
  -metrics-flush $METRICS_FLUSH
```

Проверка:

```bash
dig @127.0.0.1 -p 53 example.com A +short
curl http://127.0.0.1:9191/metrics | grep dns_requests_by_country
```

# DNS Geo-Routing Server

Простой DNS-сервер на Go с поддержкой:

* A‑запросов с выбором ближайшего сервера по GeoIP2
* Горячей подгрузкой DNS‑записей из внешнего JSON‑сервиса
* Экспортом метрик по странам клиентов в Prometheus

---

## Требования

* Go 1.18+
* Файл GeoLite2 City `.mmdb`
* HTTP‑сервис, отдающий JSON вида `map[string][]string`

---

## Переменные окружения и флаги

Все параметры можно задавать через переменные окружения или через флаги (флаги имеют приоритет над ENV).

| Параметр       | ENV                   | Флаг             | Описание                              | Значение по умолчанию |
| -------------- | --------------------- | ---------------- | ------------------------------------- | --------------------- |
| DNS порт       | `DNS_PORT`            | `-port`          | адрес и порт DNS‑сервера              | `:53`                 |
| URL записей    | `RECORDS_SERVICE_URL` | `-records-url`   | URL JSON‑сервиса с записями           | **обязателен**        |
| GeoIP DB       | `GEOIP_DB_PATH`       | `-geoip-db`      | путь к GeoLite2‑City.mmdb             | **обязателен**        |
| Интервал фетча | `REFRESH_INTERVAL`    | `-refresh`       | как часто обновлять записи            | `30s`                 |
| Метрики порт   | `METRICS_PORT`        | `-metrics-port`  | порт HTTP‑сервера метрик (`/metrics`) | `:9090`               |
| Flush метрик   | `METRICS_FLUSH`       | `-metrics-flush` | интервал сброса и публикации метрик   | `5m`                  |

---

## Локальная сборка и запуск

```bash
# Клонируем и переходим в репозиторий
git clone https://github.com/youruser/dns-server.git
cd dns-server

# Скачиваем зависимости и собираем бинарник
go mod download
go build -o bin/dns-server ./cmd/dns-server

# Копируем GeoIP базу рядом с бинарником
cp /path/to/GeoLite2-City.mmdb ./

# Задаём ENV и запускаем сервис
export RECORDS_SERVICE_URL=http://localhost:8080/records
export GEOIP_DB_PATH=./GeoLite2-City.mmdb
export DNS_PORT=:53
export REFRESH_INTERVAL=45s
export METRICS_PORT=:9191
export METRICS_FLUSH=2m

# Запуск (флаги не обязательны, можно их передать для переопределения)
./bin/dns-server \
  -port $DNS_PORT \
  -records-url $RECORDS_SERVICE_URL \
  -geoip-db $GEOIP_DB_PATH \
  -refresh $REFRESH_INTERVAL \
  -metrics-port $METRICS_PORT \
  -metrics-flush $METRICS_FLUSH
```

Проверка:

```bash
dig @127.0.0.1 -p 53 example.com A +short
curl http://127.0.0.1:9191/metrics | grep dns_requests_by_country
```

---

**Сборка и запуск контейнера**:

```bash
docker build -t dns-server:latest .

docker run -d \
  -p 5053:53/udp \
  -p 9090:9090 \
  --add-host=host.docker.internal:host-gateway \
  -e RECORDS_SERVICE_URL=http://host.docker.internal:8080/records \
  -e GEOIP_DB_PATH=/usr/local/bin/GeoLite2-City.mmdb \
  -e DNS_PORT=:53 \
  -e REFRESH_INTERVAL=30s \
  -e METRICS_PORT=:9090 \
  -e METRICS_FLUSH=5m \
  -v $(pwd)/GeoLite2-City.mmdb:/usr/local/bin/GeoLite2-City.mmdb \
  dns-server:latest
```