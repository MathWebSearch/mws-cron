FROM golang:1-alpine as builder

# Build dependencies
RUN apk add --no-cache make git

# Build the cron thingy
ADD . /go/src/github.com/MathWebSearch/mws-cron
WORKDIR /go/src/github.com/MathWebSearch/mws-cron
RUN make build-local

# And add it into a 'from scratch'
FROM scratch
COPY --from=builder /go/src/github.com/MathWebSearch/mws-cron/out/mws-cron /mws-cron

ENTRYPOINT [ "/mws-cron" ]