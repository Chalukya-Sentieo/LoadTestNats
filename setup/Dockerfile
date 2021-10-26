FROM golang:1.17.2

WORKDIR /app/
COPY ../go.* .
COPY ../src/ .
RUN go mod edit --replace github.com/Chalukya-Sentieo/LoadTestNats/src/modules=/app/modules
RUN go mod tidy
#RUN go build -o /GinServer
#RUN rm -rf *
# CMD [ "/GinServer" ]
CMD [ "tail", "-f", "/dev/null" ]
