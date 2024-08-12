FROM  golang:1.22.4 AS build
COPY *.go /app/
COPY main.go /app/main.go
COPY go.mod /app/
COPY ./proto /app/proto
WORKDIR /app

## i need to figure out how to properly build the protobuffers from scratch. My workaround is to just import them from my host machine
#RUN apt-get update
#RUN apt-get install -y protobuf-compiler golang-goprotobuf-dev
#RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
#RUN go install google.golang.org/protobuf/cmd/protoc-gen-go

#RUN protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative ./proto/verifiable_processing.proto

RUN go get
COPY client_fix.txt /go/pkg/mod/github.com/camunda/zeebe/clients/go/v8@v8.5.3/pkg/zbc/client.go
RUN go build -o /bin/server

FROM debian AS final
#RUN apt-get update -y
#RUN apt-get install --reinstall ca-certificates -y

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

# Copy the executable from the "build" stage.
COPY --from=build /bin/server /bin/server

# Expose the port that the application listens on.
#EXPOSE 2650

# What the container should run when it is started.
ENTRYPOINT [ "/bin/server" ]
