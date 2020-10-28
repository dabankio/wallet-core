#!/bin/bash

#确认可用的测试用例
#usage: bash ./scripts/ci-test.sh
pkg=github.com/dabankio/wallet-core
go build -o ci.out \
    && go test ${pkg}/bip39/... \
    && go test ${pkg}/bip44/... \
    && go test ${pkg}/wallet/... \
    && go test ${pkg}/core/bbc/... \
    && go test ${pkg}/core/btc/... \
    && go test ${pkg}/core/eth/... \
    && go test ${pkg}/core/trx/... \
    && go test ${pkg}/qa/bbc/... \
    && go test ${pkg}/qa/btc/... \
    && go test ${pkg}/qa/eth/... \
    && go test ${pkg}/qa/omni/... \
    && go test ${pkg}/qa/wallet/... 