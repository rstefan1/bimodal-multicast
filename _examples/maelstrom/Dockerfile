FROM golang:1.21-bookworm

RUN mkdir -p bmmc-maelstrom/bin

RUN set -ex && \
    apt update && \
    apt install bzip2 && \
    apt install openjdk-17-jdk --yes && \
    apt install graphviz --yes && \
    apt install gnuplot --yes

RUN set -ex && \
    wget https://github.com/jepsen-io/maelstrom/releases/download/v0.2.3/maelstrom.tar.bz2 && \
    tar -xf maelstrom.tar.bz2 && \
    cp -r maelstrom/ bmmc-maelstrom/bin/

COPY go.mod bmmc-maelstrom/
COPY go.sum bmmc-maelstrom/
COPY *.go bmmc-maelstrom/

RUN set -ex && \
    cd bmmc-maelstrom && \
    sed -i 's/\tgithub\.com\/rstefan1\/bimodal-multicast\ v0\.0\.0$/\tgithub\.com\/rstefan1\/bimodal-multicast\ f73d28a6cf84c243efe44cefcf8efbeb589c6af8/g' go.mod && \
    cat go.mod && \
    go mod tidy && \
    go install .

ENTRYPOINT ./bmmc-maelstrom/bin/maelstrom/maelstrom test -w broadcast --bin $GOPATH/bin/bmmc-maelstrom --node-count 25 --time-limit 20 --rate 100 --latency 100
