all: market wallet


market:
	go build -o market_server market-wallet/cmd/market

wallet:
	go build -o wallet_server market-wallet/cmd/wallet

clean:
	rm -rf market_server wallet_server
