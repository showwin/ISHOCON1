FROM ubuntu:20.04

ENV LANG=en_US.UTF-8 \
    LC_ALL=C.UTF-8 \
    TZ=Asia/Tokyo
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update && \
    apt-get install -y sudo wget less vim tzdata nginx \
    # 各言語のインストールに必要なもの下準備 \
    curl git gcc make libssl-dev libreadline-dev && \
    apt-get clean

# ishocon ユーザ作成
RUN groupadd -g 1001 ishocon && \
    useradd  -g ishocon -G sudo -m -s /bin/bash ishocon && \
    echo 'ishocon:ishocon' | chpasswd
RUN echo 'ishocon ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# *env 系の設定が入った .bashrc のコピー
COPY docker/app/base/.bashrc /home/ishocon/.bashrc

# Nginx の設定
COPY docker/app/base/nginx.conf /etc/nginx/nginx.conf

# MySQL のインストールと dump, cnf の配置
RUN ["/bin/bash", "-c", "debconf-set-selections <<< 'mysql-server mysql-server/root_password password ishocon1'"]
RUN ["/bin/bash", "-c", "debconf-set-selections <<< 'mysql-server mysql-server/root_password_again password ishocon1'"]
RUN ["/bin/bash", "-c", "debconf-set-selections <<< 'mysql-service mysql-server/mysql-apt-config string 4'"]
RUN apt-get install -y mysql-server
RUN mkdir -m 777 /home/ishocon/data
COPY admin/ishocon1.dump.tar.gz /home/ishocon/data/ishocon1.dump.tar.gz
# COPY admin/config/my.cnf /etc/mysql/my.cnf

USER ishocon

EXPOSE 3306 80
