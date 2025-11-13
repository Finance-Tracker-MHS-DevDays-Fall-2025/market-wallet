# Market/Wallet Service

Сервисы, предоставляющие работу с внешним API инвестиций

## Быстрый старт

### 1. Установить зависимости

```bash
go mod download
```

### 2. Скачать репозиторий с общими proto

```bash
cd .. && git clone https://github.com/Finance-Tracker-MHS-DevDays-Fall-2025/backend-common.git && cd -
```

### 3. Собрать бинарники сервиса

```bash
make all
```

Сервис будет доступен на:

- gRPC: `localhost:50051`

## Команды

- `make all` - собирает сервисы market & wallet
- `make clean` - удаляет собранные артефакты
- `make api-gen` - генерирует Go код из protobuf
- `make market` - собирает market-serivce
- `make market-client` - соирает тестовый клиент для makret-service
- `make wallet` - собирает wallet-service
- `make market-img` - собирает docker образ сервиса market
- `make wallet-img` - собирает docker образ сервиса wallet
- `make upload-img` - собирает docker образ сервиса market, wallet и публикует образы во внешний registry


## Структура проекта

```
market-wallet/
├── cmd/                    # точки входа сервисов
├── build/                  # папка с докер-файлами
├── internal/
│   │── generated/          # сгенерированный из протобафов Go-код 
│   │── market/             # логика похода во внешний API для market-service
│   └── utils/              # тут лежит вспомогательный код для работы с API
├── go.mod
├── Makefile
└── Readme.md
```
