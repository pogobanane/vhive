
make:
    #!/bin/sh
    if [[ ! -z $GOPATH ]]; then
      echo "gopath must be empty"
      exit 1
    fi
    set -x
    go mod vendor
    go install ./...

    pushd examples/invoker
    go mod vendor
    go install ./...
    popd

    pushd examples/deployer
    go mod vendor
    go install ./...
    popd

