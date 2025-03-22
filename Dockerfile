# debianを使用する
FROM golang:1.24-bookworm

WORKDIR /go/src

# sqlite3をインストール
RUN apt-get update && \
    apt-get install -y sqlite3

RUN mkdir ./code

RUN mkdir ./data

