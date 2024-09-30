FROM ubuntu:20.04

ENV LANG en_US.UTF-8
ENV LC_ALL=C.UTF-8
ENV TZ Asia/Tokyo
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update && \
    apt-get install -y wget tzdata mysql-server && \
    apt-get clean

# Go のインストール
ARG TARGETARCH
RUN wget -q https://dl.google.com/go/go1.13.15.linux-${TARGETARCH}.tar.gz && \
    tar -C /usr/local -xzf go1.13.15.linux-${TARGETARCH}.tar.gz && \
    rm go1.13.15.linux-${TARGETARCH}.tar.gz
ENV PATH=$PATH:/usr/local/go/bin \
    GOROOT=/usr/local/go \
    GOPATH=$HOME/.local/go

# build benchmark
COPY admin/benchmarker /root/admin/benchmarker
RUN cd /root/admin/benchmarker &&  \
    GOARCH=${TARGETARCH} go build -x -o ../../benchmark *.go

# MySQL 初期設定
COPY admin/ishocon1.dump.tar.gz /root/admin/ishocon1.dump.tar.gz

WORKDIR /root

# MySQL のインストール
RUN ["/bin/bash", "-c", "debconf-set-selections <<< 'mysql-server mysql-server/root_password password ishocon1'"]
RUN ["/bin/bash", "-c", "debconf-set-selections <<< 'mysql-service mysql-server/mysql-apt-config string 4'"]

# 初期データ投入
RUN service mysql start && \
    mysql -u root -pishocon -e 'CREATE DATABASE IF NOT EXISTS ishocon1;' && \
    mysql -u root -pishocon -e "CREATE USER IF NOT EXISTS ishocon IDENTIFIED BY 'ishocon';" && \
    mysql -u root -pishocon -e 'GRANT ALL ON *.* TO ishocon;' && \
    tar -zxvf ~/admin/ishocon1.dump.tar.gz -C ~/admin && mysql -u root -pishocon ishocon1 < ~/admin/ishocon1.dump && \
    rm -rf ~/admin/ishocon1.dump


COPY docker/benchmarker/entrypoint.sh /root/entrypoint.sh
ENTRYPOINT ["/root/entrypoint.sh"]

CMD /root/benchmark --ip $TARGET:80
