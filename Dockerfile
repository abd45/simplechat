FROM golang

RUN apt-get update && \        
    apt-get install -y git

RUN mkdir -p $GOPATH/src/github.com/abd45 && \
    cd $GOPATH/src/github.com/abd45 && \
    git clone https://github.com/abd45/simplechat.git && \
    cd simplechat/server && \
    go build && \
    cp server /usr/bin

EXPOSE 10001
ENTRYPOINT ["/usr/bin/server", "--address", "0.0.0.0:10001"]  
