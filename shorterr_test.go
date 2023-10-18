package shorterr_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	se "github.com/ansiwen/shorterr"
)

func Example() {
	f1 := func() (string, error) {
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

	f2 := func() (name string, err error) {
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

	_, err1 := f1()
	_, err2 := f2()
	fmt.Println(err1)
	fmt.Println(err2)

	// Output:
	// open data.json: no such file or directory
	// open data.json: no such file or directory
}

func TestCheck(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer se.PassTo(&err)
		se.Check(errFunc(x))
		return
	})
}

func TestTry(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer se.PassTo(&err)
		a = argsToSlice(se.Try(errFunc1(x)))
		return
	})
}

func TestTry2(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer se.PassTo(&err)
		a = argsToSlice(se.Try2(errFunc2(x)))
		return
	})
}

func TestTry3(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer se.PassTo(&err)
		a = argsToSlice(se.Try3(errFunc3(x)))
		return
	})
}

func TestTry4(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer se.PassTo(&err)
		a = argsToSlice(se.Try4(errFunc4(x)))
		return
	})
}

func TestTry5(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer se.PassTo(&err)
		a = argsToSlice(se.Try5(errFunc5(x)))
		return
	})
}

func TestDo(t *testing.T) {
	assert(t, "failed2: failed", func(x bool) (a []int, err error) {
		defer se.PassTo(&err)
		a = argsToSlice(se.Do(errFunc1(x)).Or("failed2"))
		return
	})
}

func TestDo2(t *testing.T) {
	assert(t, "failed2: failed", func(x bool) (a []int, err error) {
		defer se.PassTo(&err)
		a = argsToSlice(se.Do2(errFunc2(x)).Or("failed2"))
		return
	})
}

func TestDo3(t *testing.T) {
	assert(t, "failed2: failed", func(x bool) (a []int, err error) {
		defer se.PassTo(&err)
		a = argsToSlice(se.Do3(errFunc3(x)).Or("failed2"))
		return
	})
}

func TestDo4(t *testing.T) {
	assert(t, "failed2: failed", func(x bool) (a []int, err error) {
		defer se.PassTo(&err)
		a = argsToSlice(se.Do4(errFunc4(x)).Or("failed2"))
		return
	})
}

func TestDo5(t *testing.T) {
	assert(t, "failed2: failed", func(x bool) (a []int, err error) {
		defer se.PassTo(&err)
		a = argsToSlice(se.Do5(errFunc5(x)).Or("failed2"))
		return
	})
}

func TestAssert(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer se.PassTo(&err)
		se.Assert(x, "failed")
		return
	})
}

func TestOtherPanic(t *testing.T) {
	f := func() (err error) {
		defer se.PassTo(&err)
		panic("bla")
	}
	panicked := false
	func() {
		defer func() {
			s, ok := recover().(string)
			panicked = ok && s == "bla"
		}()
		f()
	}()
	if !panicked {
		t.Fatal("Expected panic")
	}
}

type testError error

func errFunc(b bool) error {
	if !b {
		return testError(fmt.Errorf("failed"))
	}
	return nil
}

func errFunc1(b bool) (int, error) {
	if !b {
		return 0, errFunc(false)
	}
	return 1, nil
}

func errFunc2(b bool) (int, int, error) {
	if !b {
		return 0, 0, errFunc(false)
	}
	return 1, 1, nil
}

func errFunc3(b bool) (int, int, int, error) {
	if !b {
		return 0, 0, 0, errFunc(false)
	}
	return 1, 1, 1, nil
}

func errFunc4(b bool) (int, int, int, int, error) {
	if !b {
		return 0, 0, 0, 0, errFunc(false)
	}
	return 1, 1, 1, 1, nil
}

func errFunc5(b bool) (int, int, int, int, int, error) {
	if !b {
		return 0, 0, 0, 0, 0, errFunc(false)
	}
	return 1, 1, 1, 1, 1, nil
}

func argsToSlice(i ...int) []int {
	return i
}

func all(a []int, x int) bool {
	for _, i := range a {
		if i != x {
			return false
		}
	}
	return true
}

func assert(t *testing.T, msg string, f func(x bool) (a []int, err error)) {
	t.Helper()
	a, err := f(true)
	if !all(a, 1) {
		t.Fatal("Expected non-zero return values")
	}
	if err != nil {
		t.Fatal("Expected no error")
	}
	a, err = f(false)
	if !all(a, 0) {
		t.Fatal("Expected zero return values")
	}
	if _, ok := err.(testError); err == nil || !ok {
		t.Fatal("Expected testError type")
	}
	if err.Error() != msg {
		t.Fatalf("expected: %s got: %s", msg, err.Error())
	}
}
