# 🚀 Hackday: SvelteKit + Go gRPC-Web Integration

A "play" project exploring the seamless integration of Go gRPC and SvelteKit 5. This repository demonstrates how to bypass the traditional hurdles of gRPC-Web by using the [Connect RPC](https://connectrpc.com/) protocol to achieve end-to-end type safety.

## 🎯 The Challenge

Standard gRPC uses HTTP/2 trailing headers, which browsers cannot access. Traditionally, this required a complex Envoy proxy. The Solution: This project uses Connect RPC to allow the Go server to speak a web-friendly version of gRPC that SvelteKit can consume directly (or via SSR) with zero proxy overhead and 100% TypeScript type safety.

## Why Connect RPC

Here is a step-by-step breakdown of the advantages and disadvantages.

    ✅ Browser Native (No Proxies)
    ✅ Better Debugging (cURL & DevTools)
    ✅ Plugs into the Go Ecosystem (A Connect server is an http.Handler. You can use it with any Go router and any standard middleware you already love.)
    ✅ Three Protocols in One (gRPC, gRPC-Web, Connect)

### Disadvantages of Connect RPC

    ❌ Performance "Edge Cases" - for 99% of web apps, the difference is zero - Google’s grpc-go is slightly more optimized for raw HTTP/2 framing than a standard net/http server.
    ❌ Smaller Ecosystem
    ❌ Feature Parity (Streaming) - bidirectional streaming over HTTP/1.1 is technically impossible while Connect allows it to fail silently if the network only supports HTTP/1.1.

## 🏗 Project Architecture

```sh
├── proto/                    # Source of Truth (Protobuf definitions)
├── server/                   # Go Service (The "Source")
│   └── pkg/greeter/          # Auto-generated Go gRPC stubs
├── client/                   # SvelteKit App (The "Consumer")
│   └── src/gen/              # Auto-generated TypeScript Connect stubs
├── taskfile.yml              # The orchestrator
├── buf.gen.server.yaml       # Backend generation config
└── buf.gen.client.yaml       # Frontend generation config
```

## 🛠 Features

Type-Safe Flow: Protobuf → Go Structs → TypeScript Runes.

SvelteKit 5 Integration: Leverages $state, $derived, and $props for reactive gRPC updates.

SSR Ready: gRPC calls are made in +page.server.ts to keep backend secrets secure and improve SEO.

Dual Generation: Independent buf templates for fine-grained control over server and client code.

## 🚀 Quick Start

### 1. Requirements

Below are listed required tools:
    - Go (1.21+)
    - Node.js (v20+)
    - Buf CLI
    - Task (recommended)

### 2. The Generation (The "Secret Sauce")

The core of this project is the automated sync between the proto and the codebases:

```sh
# Sync both worlds
task gen-server
task gen-client
```

### 3. Run the "Play"

1. Start Go Backend:

    ```sh
    cd server && go run main.go
    ```

2. Start SvelteKit:

    ```sh
    cd client && npm run dev
    ```

## 📡 The Integration Blueprint

The Go Handler (Connect)

Instead of standard gRPC, we wrap our service to support the Connect protocol, making it reachable via standard HTTP/1.1 or HTTP/2:

```go
mux := http.NewServeMux()
path, handler := greeterconnect.NewGreeterHandler(&server{})
mux.Handle(path, handler)
```

The Svelte Load (SSR)
We fetch data on the server so the browser receives fully rendered HTML, with the gRPC data already populated:

```js
export const load: PageServerLoad = async () => {
  const res = await client.getAll({});
  return { todos: res.items }; // Fully typed!
};
```

## 📜 Taskfile Reference

Task,Command,Result
gen-server,buf generate ...,Updates server/pkg/greeter with Go stubs.
gen-client,buf generate ...,Updates client/src/gen with TS/Connect stubs.
