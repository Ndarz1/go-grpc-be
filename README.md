# Go gRPC Backend

[![Go Version](https://img.shields.io/badge/Go-1.19%2B-blue.svg)](https://go.dev/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](./LICENSE)
[![Build](https://img.shields.io/github/actions/workflow/status/Ndarz1/go-grpc-be/go.yml?label=Build)](https://github.com/Ndarz1/go-grpc-be/actions)
[![Repo](https://img.shields.io/badge/GitHub-Ndarz1%2Fgo--grpc--be-black?logo=github)](https://github.com/Ndarz1/go-grpc-be)

A **gRPC backend service** built with **Go (Golang)**.  
This project provides a foundation for building **high-performance microservices** using **gRPC** and **Protocol Buffers**, following **clean architecture principles**.

---

## 🚀 Features

- ⚡ **gRPC server implementation**
- 🧩 **Protocol Buffers** for efficient serialization
- 🧱 **Clean architecture** project structure
- 🔌 **Easily extendable** with new services
- 🧪 Ready for **testing and integration** with gRPC clients

---

## 🧰 Prerequisites

Make sure you have installed:

- [Go 1.19+](https://go.dev/dl/)
- [Protocol Buffer compiler (protoc)](https://grpc.io/docs/protoc-installation/)
- Go plugins for `protoc`:

  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  export PATH="$PATH:$(go env GOPATH)/bin"
  ```
  
⚙️ Installation
  ```bash
  git clone https://github.com/Ndarz1/go-grpc-be.git
  cd go-grpc-be
  go mod download
  ```

🧬 Code Generation
  ```bash
  # Option 1: Use go generate
  go generate ./...

  # Option 2: Use protoc manually
  protoc --go_out=. --go-grpc_out=. proto/*.proto
  ```

🏃 Usage
Start the gRPC Server
  ```bash
  go run cmd/server/main.go
  ```
Run a gRPC Client
  ```
  go run cmd/client/main.go
  ```
📂 Project Structure
  ```
go-grpc-be/
├── cmd/
│   ├── server/          # Entry point for gRPC server
│   └── client/          # Example gRPC client
├── pkg/
│   ├── pb/              # Generated protobuf code
│   ├── service/         # gRPC service implementations
│   └── utils/           # Utility functions
├── proto/               # Protocol Buffer (.proto) definitions
└── internal/            # Private application code (domain, repository, etc.)

  



















