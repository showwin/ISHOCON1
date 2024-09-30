ARG BASE_IMAGE=ishocon2-app-base:latest
FROM ${BASE_IMAGE}

ENV RUBY_VERSION=2.7.1

# Ruby のインストール
RUN sudo apt-get update && \
    sudo apt-get install -y ruby-dev libmysqlclient-dev libffi-dev libyaml-dev bzip2 && \
    sudo apt-get clean && \
    git clone https://github.com/sstephenson/rbenv.git ~/.rbenv
RUN PATH="$HOME/.rbenv/bin:$PATH" && \
    eval "$(rbenv init -)" && \
    git clone https://github.com/sstephenson/ruby-build.git ~/.rbenv/plugins/ruby-build && \
    rbenv install $RUBY_VERSION && rbenv rehash && rbenv global $RUBY_VERSION

# アプリケーション
COPY --chown=ishocon:ishocon webapp/ /home/ishocon/webapp

WORKDIR /home/ishocon

RUN sudo gem install bundler -v "1.16.1" && \
    cd webapp/ruby && \
    . /home/ishocon/.bashrc && \
    bundle install

COPY run.sh /home/ishocon/run.sh
