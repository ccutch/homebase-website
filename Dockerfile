FROM golang:alpine AS build-env
WORKDIR /go/src/github.com/ccutch/homebase-website
COPY . .
RUN apk add --no-cache git mercurial             \
      && go get github.com/gobuffalo/packr/...   \
    && apk del git mercurial
RUN packr build -o app


FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/ccutch/homebase-website/app /app/
COPY ./frontend /app/frontend
RUN apk add --update nodejs

ENTRYPOINT ./app
