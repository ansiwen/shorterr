# shorterr

[![Go Report Card](https://goreportcard.com/badge/github.com/ansiwen/shorterr)](https://goreportcard.com/report/github.com/ansiwen/shorterr)
[![GoDoc](https://pkg.go.dev/badge/github.com/ansiwen/shorterr?status.svg)](https://pkg.go.dev/github.com/ansiwen/shorterr?tab=doc)

## Introduction

This is a simple library that implements short-circuit style error handling,
similar to the `?` operator in Rust. That is, in case of errors or other
failures, the current function is interrupted and the error is returned.

## Usage

A function that wants to use `shorterr` must install the `PassTo` function
with a `defer`, and then make use of the various check and wrapper functions,
provided for different function signatures.

Code that originally looks like this

```go
func myFunc() (string, error) {
	file, err := os.Open("data.json")
	if err != nil {
		return "", err
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("can't read data.json: %w", err)
	}

	var myData map[string]any
	json.Unmarshal(jsonData, &myData)
	if err != nil {
		return "", fmt.Errorf("unmarshalling failed: %w", err)
	}

	val, ok := myData["name"]
	if !ok {
		return "", fmt.Errorf("missing name property")
	}

	name, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("invalid name property")
	}

	return name, nil
}
```

can instead be written like this

```go
...

import se "github.com/ansiwen/shorterr"

...

func myFunc() (name string, err error) {
	defer se.PassTo(&err)

	file := se.Try(os.Open("data.json"))

	jsonData := se.Do(io.ReadAll(file)).Or("can't read data.json")

	var myData map[string]any
	se.Check(json.Unmarshal(jsonData, &myData), "unmarshalling failed")

	val, ok := myData["name"]
	se.Assert(ok, "missing name property")

	name, ok = val.(string)
	se.Assert(ok, "invalid name property")

	return
}
```
