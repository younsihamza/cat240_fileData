FROM golang:1.24.1-alpine3.21
RUN apk add --no-cache libc6-compat && mkdir /app
WORKDIR /app
COPY . /app/
RUN go mod tidy && \
	go build -o app && \
	chmod +x /app/app && \
	find . -maxdepth 1 -not -name 'app' -not -name '.' -not -name '..' -exec rm -rf {} \; && \
	mkdir -p /app/log
CMD ["./app"]