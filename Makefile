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

wallet: api-gen
	go build -o wallet_server market-wallet/cmd/wallet

clean:
	rm -rf market_server wallet_server
