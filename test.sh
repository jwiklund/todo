#!/bin/bash

set -e

test_dirs=`find . -name "*_test.go" -exec dirname {} \; | sort | uniq`
for test_dir in $test_dirs; do
    pushd $test_dir
    go test
    success=$?
    popd
    if [[ $success != 0 ]]; then
	exit 1
    fi
done
