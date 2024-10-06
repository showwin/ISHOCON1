#!/bin/bash
set -ex
sudo su - ishocon
. /home/ishocon/.bashrc

cd /tmp/
tar -zxvf /tmp/benchmarker.tar.gz
cd /tmp/benchmarker
go build -x -o /home/ishocon/benchmark *.go
