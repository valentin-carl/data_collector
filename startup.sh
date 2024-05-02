#!/bin/sh

sudo apt-get update
sudo apt-get install -y unzip
sudo apt-get install -y build-essential

wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

mkdir /app
cd /app
wget https://github.com/valentin-carl/data_collector/archive/refs/heads/main.zip
unzip main.zip

cd data_collector-main/ && sudo /usr/local/go/bin/go build
sudo ./data_collector > /var/log/data_collector.log 2>&1 &
