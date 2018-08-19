# go-cdds
cyclone DDS go interface


This is currently just wrapper (rough implementation).
Will be organized for go-like.

## Quic start
1. build and install cyclone DDS
2. in `/examples/helloworld` directory, generate `HelloWorld.o` (copy from cyclone DDS's example directory should be easy)
3. Run bellow in separate terminals
    - `go run publisherpublish.go`
    - `go run subscriber/subscribe.go`
