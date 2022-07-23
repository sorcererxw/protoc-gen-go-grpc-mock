# protoc-gen-go-grpc-mock

`protoc-gen-go-grpc-mock` generates gRPC service mocks
compatible with `golang/mock`(aka `gomock`).

## What different from `gomock`

`gomock` extracts the interfaces from Go source code and
generates the mocks. The overhead of parsing Go source code
is much higher than that of Protobuf. Using `gomock` to
generate a mock for a bunch of `*_grpc.pb.go` files is very
slow.

As mentioned above, that's why I wrote this plugin to speed
up Protobuf Go generations.

## Benchmark

[Here](./example) is the `petstore.proto` and pre
generated `petstore_grpc.pb.go`.

#### gomock

```shell
hyperfine "mockgen -source=petstore_grpc.pb.go -destination=mock/petstore_grpc_mock.pb.go" 
```

```
Benchmark #1: mockgen -source=petstore_grpc.pb.go -destination=petstore_grpc_mock.pb.go
  Time (mean ± σ):     622.5 ms ± 117.0 ms    [User: 527.2 ms, System: 768.7 ms]
  Range (min … max):   503.9 ms … 918.9 ms    10 runs
```

#### protoc-gen-go-grpc-mock

```shell
hyperfine "protoc --go-grpc-mock_out=. petstore.proto" 
```

```
Benchmark #1: protoc --go-grpc-mock_out=. petstore.proto
  Time (mean ± σ):      26.4 ms ±   7.5 ms    [User: 13.3 ms, System: 12.7 ms]
  Range (min … max):    17.8 ms …  62.8 ms    44 runs
```

## Installation

```
go install github.com/sorcererxw/protoc-gen-go-grpc-mock@latest
```

Also required:

- [protoc](https://github.com/google/protobuf)
- [protoc-gen-go](https://github.com/golang/protobuf)

## Usage

### with protoc

```shell
protoc --go_out=. --go-grpc_out=. --go-grpc-mock_out=. petstore.proto 
```

### with buf

```yaml
version: v1

plugins:
  - name: go
    out: .
    opt: paths=source_relative

  - name: go-grpc
    out: .
    opt: paths=source_relative

  - name: go-grpc-mock
    out: .
    opt: paths=source_relative
```