#!/bin/bash
set -ex

export LANG=en_US.UTF-8
export LC_ALL=C.UTF-8
export TZ=Asia/Tokyo
sudo ln -snf /usr/share/zoneinfo/$TZ /etc/localtime
echo $TZ | sudo tee /etc/timezone

sudo apt-get update
sudo apt-get install -y sudo wget less vim tzdata nginx \
  curl git gcc make libssl-dev libreadline-dev
sudo apt-get clean

# ishocon ユーザ作成
sudo groupadd -g 1001 ishocon
sudo useradd  -g ishocon -G sudo -m -s /bin/bash ishocon
# sudo echo 'ishocon:ishocon' | chpasswd
echo 'ishocon ALL=(ALL) NOPASSWD:ALL' | sudo tee /etc/sudoers.d/ishocon
sudo mkdir -m 775 /home/ishocon/webapp
sudo mkdir -m 777 /home/ishocon/data
sudo chown -R ishocon:ishocon /home/ishocon
sudo chmod 777 /home/ishocon
