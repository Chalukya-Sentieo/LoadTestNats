FROM golang:1.16

WORKDIR /go/src/app
COPY . .
RUN go mod tidy
RUN go build -o GinServer
RUN cp GinServer /go/src
CMD [ "/go/src/GinServer" ]