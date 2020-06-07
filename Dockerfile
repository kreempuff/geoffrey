FROM golang:1.14 as builder
RUN mkdir /workspace
WORKDIR /workspace
COPY go.mod ./
COPY go.sum ./
COPY bin ./bin
COPY workspace ./workspace
COPY orchestrator ./orchestrator
COPY worker ./worker

FROM builder as orchestratorBuilder
RUN go build -o geoffrey-o ./bin/orchestrator/main.go

FROM ubuntu:16.04 as orchestrator
RUN adduser orchestrator --gecos ''
COPY --from=orchestratorBuilder /workspace/geoffrey-o /usr/local/bin/geoffrey-o
RUN chmod 0755 /usr/local/bin/geoffrey-o
USER orchestrator
CMD geoffrey-o

FROM builder as workerBuilder
RUN go build -o geoffrey-w ./bin/worker/main.go

FROM ubuntu:16.04 as worker
RUN adduser worker --gecos ''
COPY --from=workerBuilder /workspace/geoffrey-w /usr/local/bin/geoffrey-w
RUN chmod 0755 /usr/local/bin/geoffrey-w
USER worker
CMD geoffrey-w