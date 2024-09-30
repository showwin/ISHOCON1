#!/bin/bash -eux

service mysql restart

echo 'setup completed.'

exec "$@"
