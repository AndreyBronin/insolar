#!/bin/#!/usr/bin/env bash

WORKDIR=~/go/src/github.com/insolar
REPOS=( "insolar-api" "insolar-internal-api" "insolar-observer-api" )
PORTS=( "8080" "8081" "8082" )
CURENT_DIR=$(pwd)

for i in {0..2} ; do
  repo_name=${REPOS[i]}
  port=${PORTS[i]}
  lsof -ti:${port} | xargs kill
  if [[ -d ${WORKDIR}/${repo_name}/ ]]; then
    cd ${WORKDIR}/${repo_name}/
    npm run export -- --target=api-watch.yaml
    nohup npm run watch >ignore.log 2>&1 </dev/null &
    echo ${repo_name} started on http://localhost:${port}/
  fi
done

cd ${CURENT_DIR}
