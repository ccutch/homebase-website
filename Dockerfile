#FROM golang:1.9 as builder


#RUN go get github.com/gobuffalo/packr/...
#COPY . .
#RUN CGO_ENABLED=0 GOOS=linux packr build -installsuffix cgo 


##FROM alpine:latest
##WORKDIR /root/
##COPY --from=builder /go/src/github.com/ccutch/homebase-website/app .
#CMD ["./app"]

FROM golang:alpine AS build-env

WORKDIR /go/src/github.com/ccutch/homebase-website
COPY . .
RUN apk add --no-cache git mercurial             \
      && go get github.com/gobuffalo/packr/...   \
    && apk del git mercurial

RUN packr build -o app

# final stage
FROM alpine

WORKDIR /app
COPY --from=build-env /go/src/github.com/ccutch/homebase-website/app /app/
ENTRYPOINT ./app
