#!/bin/bash
set -ex

cd /home/ishocon

# Install MySQL
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password password ishocon1'
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password_again password ishocon1'
sudo debconf-set-selections <<< 'mysql-service mysql-server/mysql-apt-config string 4'
sudo apt-get install -y mysql-server

# Install Ruby
export RUBY_VERSION=3.3.5
sudo apt-get install -y ruby-dev libmysqlclient-dev libffi-dev libyaml-dev bzip2
sudo apt-get clean
git clone https://github.com/sstephenson/rbenv.git ~/.rbenv
export PATH="$HOME/.rbenv/bin:$PATH"
eval "$(rbenv init -)"
git clone https://github.com/sstephenson/ruby-build.git ~/.rbenv/plugins/ruby-build
rbenv install $RUBY_VERSION && rbenv rehash && rbenv global $RUBY_VERSION

# Install Python
export PYTHON_VERSION=3.8.5
sudo apt-get install -y zlib1g-dev libbz2-dev libffi-dev libsqlite3-dev liblzma-dev libmariadb-dev pkgconf
sudo apt-get clean
git clone https://github.com/pyenv/pyenv.git ~/.pyenv
export PYENV_ROOT="$HOME/.pyenv"
export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init -)"
pyenv install $PYTHON_VERSION && pyenv global $PYTHON_VERSION

# Install Go
export GO_VERSION=1.23.1
export TARGETARCH=amd64
sudo wget -q https://dl.google.com/go/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz
sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-${TARGETARCH}.tar.gz
sudo rm go${GO_VERSION}.linux-${TARGETARCH}.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GOROOT=/usr/local/go
export GOPATH=/home/ishocon/.local/go
export PATH=$PATH:$GOROOT/bin
