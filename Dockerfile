FROM golang

RUN go get github.com/vsco/tail
RUN go get github.com/mattbaird/elastigo
RUN mkdir /tail
COPY src/ /tail/src/
COPY tests/ /tail/tests/

WORKDIR /tail/src
RUN go build -o /usr/bin/log-shipper log-shipper.go

CMD log-shipper -index testindex -elastichost=elasticsearch /tail/tests/nginx.log
