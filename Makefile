build:
	go build -o ./bin/kkrh main.go

test:
	go test -v ./...

run:
	export MODE=dev && export DEV_PORT=8080 && export DEV_DATABASE_URL=postgres://viveknathani:root@localhost:5432/kkrhdb && export DEV_REDIS_URL=127.0.0.1:6379 && export DEV_JWT_SECRET=hey && ./bin/kkrh