# Download base image ubuntu 20.04
FROM ubuntu

# Disable Prompt During Packages Installation
ARG DEBIAN_FRONTEND=noninteractive

# Update Ubuntu Software repository
RUN apt update

RUN apt-get update -y

RUN apt-get install -y gccgo

WORKDIR /go/app

COPY . .

RUN gccgo -O3 -o cool main.go

CMD ["./cool"]