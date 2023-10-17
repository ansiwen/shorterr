# shorterr

## Introduction

This is a simple library that implements short-circuit style error handling,
similar to the `?` operator in Rust. That is, in case of errors or other
failures, the current function is interrupted and the error is returned.

## Usage

The function that wants to use `shorterr` must install the `PassTo` function
with a `defer`, and then make use of the various check and wrapper functions,
provided for different function signatures.


Code that looks like that:

```go
func myFunc() (string, error) {
	file, err := os.Open("data.json")
	if err != nil {
		return "", fmt.Errorf("open data.json: %w", err)
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	var myData map[string]any
	json.Unmarshal(jsonData, &myData)
	if err != nil {
		return "", fmt.Errorf("unmarshalling: %w", err)
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

can be written like that:

```go
...

import se "github.com/ansiwen/shorterr"

...

func myFunc() (name string, err error) {
	defer se.PassTo(&err)

	file := se.Do(os.Open("data.json")).Or("open data.json")

	jsonData := se.Must(io.ReadAll(file))

	var myData map[string]any
	se.Check(json.Unmarshal(jsonData, &myData), "unmarshalling")

	val, ok := myData["name"]
	se.Assert(ok, "missing name property")

	name, ok = val.(string)
	se.Assert(ok, "invalid name property")

	return
}
```
