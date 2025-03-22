# debianを使用する
FROM golang:1.24-bookworm

WORKDIR /go

# sqlite3をインストール
RUN apt-get update && \
    apt-get install -y sqlite3

RUN mkdir ./project

WORKDIR ./project
