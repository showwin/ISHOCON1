ARG BASE_IMAGE=ishocon1-app-base:latest
FROM ${BASE_IMAGE}

# Python のインストール
RUN sudo apt-get update && \
    sudo apt-get install -y zlib1g-dev libbz2-dev libffi-dev libsqlite3-dev liblzma-dev libmariadb-dev pkgconf && \
    sudo apt-get clean
RUN git clone https://github.com/pyenv/pyenv.git ~/.pyenv && \
    PYENV_ROOT="$HOME/.pyenv" && PATH="$PYENV_ROOT/bin:$PATH" && \
    eval "$(pyenv init -)" && \
    pyenv install 3.8.5 && pyenv global 3.8.5

# アプリケーション
COPY --chown=ishocon:ishocon webapp/ /home/ishocon/webapp

# ライブラリのインストール
RUN sudo apt-get install -y libmysqlclient-dev
RUN LC_ALL=C.UTF-8 && LANG=C.UTF-8 && cd /home/ishocon/webapp/python && \
    /home/ishocon/.pyenv/shims/pip install -r requirements.txt

WORKDIR /home/ishocon

COPY run.sh /home/ishocon/run.sh
