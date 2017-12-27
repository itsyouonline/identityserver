FROM golang:1.7
MAINTAINER itsyou.online

ARG GOBINDATAVERSION=a0ff2567cfb70903282db057e799fd826784d41d

RUN git clone https://github.com/jteeuwen/go-bindata.git $GOPATH/src/github.com/jteeuwen/go-bindata
WORKDIR $GOPATH/src/github.com/jteeuwen/go-bindata
RUN git checkout $GOBINDATAVERSION
RUN go get github.com/jteeuwen/go-bindata/...


ARG GORAMLVERSION=a5aed17caa97fd822abd86c59abf4465b07172a6

RUN git clone https://github.com/Jumpscale/go-raml.git $GOPATH/src/github.com/Jumpscale/go-raml
WORKDIR $GOPATH/src/github.com/Jumpscale/go-raml
RUN git checkout $GORAMLVERSION
RUN ./build.sh

# autopep8 is required by go-raml 
RUN apt-get update
RUN apt-get install -y python-pip
RUN pip install autopep8

ENV CGO_ENABLED 0
WORKDIR /go/src/github.com/itsyouonline/identityserver

EXPOSE 8080

ENTRYPOINT go generate && go build && ./identityserver -d
