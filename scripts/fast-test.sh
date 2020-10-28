#!/bin/bash

#快速测试

pkg=github.com/dabankio/wallet-core

go build -o ci.out \
    && go test --short ${pkg}/bip39/... \
    && go test --short ${pkg}/bip44/... \
    && go test --short ${pkg}/wallet/... \
    && go test --short ${pkg}/core/bbc/... \
    && go test --short ${pkg}/core/btc/... \
    && go test --short ${pkg}/core/eth/... 