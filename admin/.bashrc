export PATH=$PATH:/usr/local/go/bin
export GOROOT=/usr/local/go
export GOPATH=$HOME/.local/go
export PATH="$HOME/.rbenv/bin:$PATH"
command -v rbenv >/dev/null && eval "$(rbenv init -)"
export PYENV_ROOT="$HOME/.pyenv"
export PATH="$PYENV_ROOT/bin:$PATH"
command -v pyenv >/dev/null && eval "$(pyenv init -)"

export ISHOCON1_DB_HOST="127.0.0.1"
export ISHOCON1_DB_PORT=3306
export ISHOCON1_DB_USER="ishocon"
export ISHOCON1_DB_PASSWORD="ishocon"
export ISHOCON1_DB_NAME="ishocon1"
export ISHOCON1_DB_BENCH_NAME="ishocon1_bench"
