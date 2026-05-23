FROM golang

WORKDIR /usr/src/bootdotdev-go-build-a-blog-aggregator

ENV XDG_CONFIG_HOME="/config"

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go*.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/bootdotdev-go-build-a-blog-aggregator

CMD ["bootdotdev-go-build-a-blog-aggregator"]
