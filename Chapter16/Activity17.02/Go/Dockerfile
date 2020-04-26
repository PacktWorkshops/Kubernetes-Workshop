FROM golang:alpine as builder

LABEL version="1.0" \
      maintainer="fmasood@redhat.com"

RUN mkdir /go/src/app && apk update && apk add --no-cache git && go get -u github.com/golang/dep/cmd/dep
COPY ./*.go ./Gopkg.toml /go/src/app/

WORKDIR /go/src/app

RUN dep ensure -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM registry.fedoraproject.org/fedora-minimal:32
RUN groupadd appgroup && useradd appuser -G appgroup
COPY --from=builder /go/src/app/main /app/
WORKDIR /app
USER appuser
CMD ["./main"]





