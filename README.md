<img src="https://github.com/symonk/turbo/blob/main/.github/images/turbo.png" border="1" width="275" height="275"/>

[![GoDoc](https://pkg.go.dev/badge/github.com/symonk/turbo)](https://pkg.go.dev/github.com/symonk/turbo)
[![Build Status](https://github.com/symonk/turbo/actions/workflows/go_test.yml/badge.svg)](https://github.com/symonk/turbo/actions/workflows/go_test.yml)
[![codecov](https://codecov.io/gh/symonk/turbo/branch/main/graph/badge.svg)](https://codecov.io/gh/symonk/turbo)
[![Go Report Card](https://goreportcard.com/badge/github.com/symonk/turbo)](https://goreportcard.com/report/github.com/symonk/turbo)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://github.com/symonk/turbo/blob/master/LICENSE)


> [!CAUTION]
> turbo is currently in early phase development and not fit for use

# Turbo 🧬

⚡ A blazing-fast worker pool library for Go, built with concurrency and simplicity in mind.

## Features

- 🌀 Minimal API
- ⚙️ Efficient goroutine scheduling
- 🧵 Safe concurrent task execution
- 📦 Designed as a pure library – no executables

## Example

```go
p := turbo.NewPool(4)
defer p.Close()

p.Submit(func() {
    fmt.Println("Hello from a worker!")
})

p.Wait()