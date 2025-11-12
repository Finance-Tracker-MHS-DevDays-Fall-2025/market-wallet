COMMON_REPO=../backend-common

all: market wallet

api-gen:
	protoc --go_out=./internal/generated -I$(COMMON_REPO)/proto  $(COMMON_REPO)/proto/common/common.proto
	protoc --go-grpc_out=./internal/generated -I$(COMMON_REPO)/proto  $(COMMON_REPO)/proto/market/market.proto 
	protoc --go_out=./internal/generated -I$(COMMON_REPO)/proto  $(COMMON_REPO)/proto/market/market.proto 
	protoc --go-grpc_out=./internal/generated -I$(COMMON_REPO)/proto  $(COMMON_REPO)/proto/wallet/wallet.proto 
	protoc --go_out=./internal/generated -I$(COMMON_REPO)/proto  $(COMMON_REPO)/proto/wallet/wallet.proto 

market: api-gen
	go build -o market_server market-wallet/cmd/market

market-client: api-gen
	go build -o market-client_server market-wallet/cmd/client_market

wallet: api-gen
	go build -o wallet_server market-wallet/cmd/wallet

clean:
	rm -rf market_server wallet_server


market-img:
	docker build -f build/market/Dockerfile -t market .

wallet-img:
	docker build -f build/wallet/Dockerfile -t wallet .

upload-img:
	docker tag market:latest cr.yandex/crpkimlhn85fg9vjfj7l/market:latest
	docker tag wallet:latest cr.yandex/crpkimlhn85fg9vjfj7l/wallet:latest
	docker image push cr.yandex/crpkimlhn85fg9vjfj7l/market:latest
	docker image push cr.yandex/crpkimlhn85fg9vjfj7l/wallet:latest
