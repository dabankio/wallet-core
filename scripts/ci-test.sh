#!/bin/bash

#确认可用的测试用例
#usage: bash ./scripts/ci-test.sh


pkg=github.com/dabankio/wallet-core

export GOPROXY=https://goproxy.cn

go test -v ${pkg}/qa/bbc \
    && go test -v -run ^TestSimplemultisigGanache$ ${pkg}/core/eth/internalized \
    && go test -v -run ^TestSimplemultisigGanacheERC20$ ${pkg}/core/eth/internalized \