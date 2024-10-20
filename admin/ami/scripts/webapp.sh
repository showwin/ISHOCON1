#!/bin/bash
set -ex

# Move files
cp /tmp/.bashrc /home/ishocon/.bashrc
sudo cp /tmp/nginx.conf /etc/nginx/nginx.conf
cp /tmp/ishocon1.dump.tar.gz /home/ishocon/data/ishocon1.dump.tar.gz
cd /home/ishocon
tar -zxvf /tmp/webapp.tar.gz
chown -R ishocon:ishocon /home/ishocon/webapp

# Load .bashrc
. /home/ishocon/.bashrc

# # Install Ruby libraries
cd /home/ishocon/webapp/ruby
gem install bundler -v "2.5.22"
bundle install

# # Install Python libraries
cd /home/ishocon/webapp/python
sudo apt-get install -y libmysqlclient-dev
pip install -r requirements.txt

# Install Go libraries
ls /home/ishocon
ls /home/ishocon/webapp
cd /home/ishocon/webapp/go
go build -o webapp *.go

# Load data into MySQL
sudo chown -R mysql:mysql /var/lib/mysql
sudo service mysql start
sudo mysql -u root -pishocon1 -e 'CREATE DATABASE IF NOT EXISTS ishocon1;'
sudo mysql -u root -pishocon1 -e "CREATE USER IF NOT EXISTS ishocon IDENTIFIED BY 'ishocon';"
sudo mysql -u root -pishocon1 -e 'GRANT ALL ON *.* TO ishocon;'
tar -zxvf ~/data/ishocon1.dump.tar.gz -C ~/data && sudo mysql -u root -pishocon1 ishocon1 < ~/data/ishocon1.dump

# Nginx
sudo nginx -t
sudo service nginx start
