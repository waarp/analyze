#!/usr/bin/env bash

set -e

VERSION=$(grep "const VERSION_NUM" src/cmd/waarp-analyze/main.go | tr -d -c '0-9.')
TMP_DIR=$(mktemp -d)
DIST_DIR=dist
CURDIR=$(cd $(dirname $0) && pwd)

function cleanup {
    if [[ -d $TMP_DIR ]]; then
        rm  -fr $TMP_DIR
    fi
}
trap cleanup EXIT

architecture() {
    case $1 in
    amd64)
        echo -n "x86_64"
        ;;
    386)
        echo -n "x86"
        ;;
    *)
        echo "unsupported"
        return 1
        ;;
    esac
}

make_package() {
  if [[ $# != 2 ]]; then
      echo "Usage: make_package GOOS GOARCH"
      echo "ERR: missing GOOS and GOARCH variables"
      exit 1
  fi

  local binfile="waarp-analyze-$1-$2"

  GOOS=$1 GOARCH=$2 CGO_ENABLED=0 GOPATH="$CURDIR:$CURDIR/vendor" go build -v \
        -o bin/$binfile -ldflags '-s -w' ./src/cmd/waarp-analyze

  mkdir -p $TMP_DIR/usr/bin/
  cp bin/$binfile $TMP_DIR/usr/bin/waarp-analyze

  fpm -s dir -t rpm -C $TMP_DIR                             \
      --package dist/                                       \
      --name waarp-analyze                                  \
      --version $VERSION                                    \
      --iteration 1                                         \
      --architecture $(architecture $2)                     \
      --vendor "Waarp SAS"                                  \
      --maintainer "Bruno CARLIN <bruno.carlin@waarp.fr>"   \
      --url "http://www.waarp.fr"                           \
      --rpm-auto-add-directories                            \
      --rpm-dist el6                                        \
      .
  fpm -s dir -t rpm -C $TMP_DIR                             \
      --package dist/                                       \
      --name waarp-analyze                                  \
      --version $VERSION                                    \
      --iteration 1                                         \
      --architecture $(architecture $2)                     \
      --vendor "Waarp SAS"                                  \
      --maintainer "Bruno CARLIN <bruno.carlin@waarp.fr>"   \
      --url "http://www.waarp.fr"                           \
      --rpm-auto-add-directories                            \
      --rpm-dist el7                                         \
      .
  rm -rf $TMP_DIR/*

  mkdir -p $TMP_DIR/waarp-analyze-$VERSION-$1-$2
  cp bin/$binfile $TMP_DIR/waarp-analyze-$VERSION-$1-$2
  cd $TMP_DIR
  tar cfz waarp-analyze-$VERSION-$1-$2.tar.gz waarp-analyze-$VERSION-$1-$2
  cd -
  mv $TMP_DIR/waarp-analyze-$VERSION-$1-$2.tar.gz dist
  rm -rf $TMP_DIR/*
}

build() {
  mkdir -p dist
  make_package "linux" "amd64"
  make_package "linux" "386"
}

bump() {
  if [[ $1 == "" ]]; then
    echo "ERROR: Provide a version as argument"
    exit 3
  fi

  sed -i -e "s|version = '[0-9\.]\+'|version = '"$1"'|" doc/source/conf.py
  sed -i -e "s|release = '[0-9\.]\+'|release = '"$1"'|" doc/source/conf.py
  sed -i -e 's|const VERSION_NUM = "[0-9\.]\+"|const VERSION_NUM = "'$1'"|' src/cmd/waarp-analyze/main.go

  git add doc/source/conf.py src/cmd/waarp-analyze/main.go
  git commit -m "version bump to $1"
  git tag $1
}

usage() {
  echo "Usage: $0 task"
  echo "TASKS"
  echo "  build         Builds packages"
  echo "  bump VERSION  Bumps version everywhere it is used"
  echo "  help          Show this help"
}

case $1 in
  build)
    build
    ;;
  bump)
    bump $2
    ;;
  *)
    usage
    ;;
esac