FROM golang:1.11-alpine3.8 as builder
ADD . /go/src/github.com/twonegatives/coinsph_challenge
WORKDIR /go/src/github.com/twonegatives/coinsph_challenge
RUN apk add --update git bash && rm -rf /var/cache/apk/*
ENV GO111MODULE=on
ENV CGO_ENABLED=0
RUN GOGC=off go build -i -v \
    -installsuffix 'static' \
    -o service ./cmd/service
RUN go get -v github.com/rubenv/sql-migrate/...
RUN echo -e '' > ~/.netrc

FROM alpine:3.8
RUN apk add --update tzdata ca-certificates bash
WORKDIR /app/
COPY --from=builder /go/src/github.com/twonegatives/coinsph_challenge/service .
COPY --from=builder /go/src/github.com/twonegatives/coinsph_challenge/migrations ./migrations/
COPY --from=builder /go/src/github.com/twonegatives/coinsph_challenge/dbconfig.yml .
COPY --from=builder /go/src/github.com/twonegatives/coinsph_challenge/bin /bin/
COPY --from=builder /go/bin/sql-migrate /usr/local/bin/

CMD ["/app/service"]
EXPOSE 80
