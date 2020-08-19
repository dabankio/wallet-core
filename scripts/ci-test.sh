#!/bin/bash

#确认可用的测试用例
#usage: bash ./scripts/ci-test.sh


pkg=github.com/dabankio/wallet-core

export GOPROXY=https://goproxy.cn

go build \
    && go test ${pkg}/qa/bbc/... \
    && go test ${pkg}/bip39/... \
    && go test ${pkg}/bip44/... \
    && go test ${pkg}/core/bbc/... \
    && go test ${pkg}/core/btc/... \
    && go test ${pkg}/core/eth/... 