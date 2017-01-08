all:
	/usr/local/go/bin/go build gocrawl.go PageContent.go ThreadSafeStringQueue.go ThreadSafeSets.go
run:
	/usr/local/go/bin/go run gocrawl.go PageContent.go ThreadSafeStringQueue.go ThreadSafeSets.go

# go build -ldflags "-s -w"