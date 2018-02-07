all:
	go build gocrawl.go Client.go UrlFinder.go CssManager.go Sets.go StringReservoir.go
run:
	go run gocrawl.go Client.go UrlFinder.go CssManager.go Sets.go StringReservoir.go
test:
	go test gocrawl.go Client.go UrlFinder.go CssManager.go Sets.go StringReservoir.go StringReservoir_test.go Sets_test.go

# go build -ldflags "-s -w"