#!/bin/bash -ex

pushd `dirname $0` > /dev/null
base=$(pwd -P)
popd > /dev/null

# Gather some data about the repo
source $base/vars.sh

# Run the test
newman -c $base/../postman/pzsvc-lasinfo-black-box-tests.json.postman_collection -s
