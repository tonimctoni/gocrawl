all:
	go build gocrawl.go PageContent.go ThreadSafeStringQueue.go ThreadSafeSets.go
run:
	go run gocrawl.go PageContent.go ThreadSafeStringQueue.go ThreadSafeSets.go