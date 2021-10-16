#!/usr/bin/env bash

set -e

export LD_PRELOAD=/usr/lib/x86_64-linux-gnu/libnss_wrapper.so
export NSS_WRAPPER_PASSWD=/tmp/passwd
export NSS_WRAPPER_GROUP=/etc/group

if ! whoami &> /dev/null; then
  echo "git:x:$(id -u):$(id -g):Git user:${HOME}:/sbin/nologin" > "${NSS_WRAPPER_PASSWD}"
fi

exec "$@"
