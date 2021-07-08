
make:
    #!/bin/sh
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

