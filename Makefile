.PHONY: orderbook

build: orderbook
orderbook:
	go build -o ./orderbookd ./cmd/main.go
clean:
	rm -i -f orderbookd
run:
	LOG_LEVEL=info ./orderbookd