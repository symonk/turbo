<img src="https://github.com/symonk/turbo/blob/main/.github/images/turbo.png" border="1" width="275" height="275"/>

[![GoDoc](https://pkg.go.dev/badge/github.com/symonk/turbo)](https://pkg.go.dev/github.com/symonk/turbo)
[![Build Status](https://github.com/symonk/turbo/actions/workflows/go_test.yml/badge.svg)](https://github.com/symonk/turbo/actions/workflows/go_test.yml)
[![codecov](https://codecov.io/gh/symonk/turbo/branch/main/graph/badge.svg)](https://codecov.io/gh/symonk/turbo)
[![Go Report Card](https://goreportcard.com/badge/github.com/symonk/turbo)](https://goreportcard.com/report/github.com/symonk/turbo)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://github.com/symonk/turbo/blob/master/LICENSE)


# ⚡ Turbo

**Turbo** is a blazing-fast, autoscaling worker pool library for Go, crafted with concurrency and simplicity in mind. It provides a minimal API for efficient goroutine scheduling and safe concurrent task execution.

> ⚠️ **Note**: Turbo is currently in early-phase development and not yet production-ready. Contributions and feedback are welcome!

## 🚀 Features

- 🔥 **Minimal API** – Designed for ease of use with just what you need  
- ⚙️ **Efficient goroutine scheduling** – Manages workers smartly for performance  
- 🔒 **Safe concurrency** – Submit and execute tasks safely in parallel  
- 📚 **Library-first design** – Turbo is a Go library, not an executable  
- 📊 **Priority queueing** *(planned)* – Execute high-priority tasks first  
- 🌊 **Draining & pause capabilities** *(planned)* – Gracefully halt or resume task execution  
- 🧠 **Autoscaling** *(planned)* – Scale worker count based on demand  

## 📦 Installation

Install the latest version using:

```go
go get github.com/symonk/turbo@v0.1.0
```

> 💡 Replace `v0.1.0` with the latest version tag from the [Releases](https://github.com/symonk/turbo/releases) page if a newer one is available.

## ⚡ Quick Start

```go
package main

import (
    "fmt"
    "github.com/symonk/turbo"
)

func main() {
    // Create a new worker pool with 4 workers
    p := turbo.NewPool(4)
    defer p.Close()

    // Submit a task to the pool
    p.Submit(func() {
        fmt.Println("Hello from Turbo!")
    })
}
```

➡️ Check out the [examples/basic](https://github.com/symonk/turbo/tree/main/examples/basic) folder for more usage patterns.

## 🤝 Contributing

Contributions are welcome! To contribute:

1. 🍴 Fork the repository  
2. 🌿 Create a branch: `git checkout -b feature/your-feature-name`  
3. 🛠️ Make your changes  
4. ✅ Run the test suite: `go test ./...`  
5. 💬 Commit your changes: `git commit -m 'Add your feature'`  
6. 🚀 Push the branch: `git push origin feature/your-feature-name`  
7. 📬 Open a pull request  
