FROM golang:alpine3.22
LABEL maintainer="Adeboye"
RUN apk update && mkdir /go/src/app && mkdir /root/.aws
WORKDIR /go/src/app
COPY . .
EXPOSE 8080
CMD ["go", "run", "golang_api.go"]