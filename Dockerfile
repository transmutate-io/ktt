FROM golang:1.13.4-alpine3.10 AS builder

ADD . /tmp/src/code

RUN apk add make git && \
    cd /tmp/src/code && \
    make clean && \
    make

FROM builder

COPY --from=builder /tmp/src/code/build/ktt /ktt

CMD /ktt