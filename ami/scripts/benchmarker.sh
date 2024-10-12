#!/bin/bash
set -ex
sudo su - ishocon
. /home/ishocon/.bashrc

cd /tmp/
tar -zxvf /tmp/benchmarker.tar.gz
cd /tmp/benchmarker
go build -x -o /home/ishocon/benchmark *.go

# Load data into MySQL
sudo mysql -u root -pishocon1 -e 'CREATE DATABASE IF NOT EXISTS ishocon1_bench;'
sudo mysql -u root -pishocon1 -e "CREATE USER IF NOT EXISTS ishocon IDENTIFIED BY 'ishocon';"
sudo mysql -u root -pishocon1 -e 'GRANT ALL ON *.* TO ishocon;'
tar -zxvf ~/data/ishocon1.dump.tar.gz -C ~/data && sudo mysql -u root -pishocon1 ishocon1_bench < ~/data/ishocon1.dump
