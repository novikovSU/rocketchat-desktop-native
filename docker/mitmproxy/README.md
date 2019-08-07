# Use of mitmproxy

**NEED TO CORRECT AND APPEND**

Steps:

1. Start mitmproxy (one of its part)
2. Launch application with proxy. Ex. `http_proxy=localhost:8080 go run *.go`

Start dumping to console (non-interactive)

    docker-compose up dump

Start interactive console application

    docker-compose run proxy

Start webapp (UI: http://localhost:8081)

    docker-compose up web
