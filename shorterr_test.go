package shorterr_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ansiwen/shorterr"
)

func Example() {
	f1 := func() (string, error) {
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

	f2 := func() (name string, err error) {
		defer shorterr.PassTo(&err)

		file := shorterr.Do(os.Open("data.json")).Or("open data.json")

		jsonData := shorterr.Must(io.ReadAll(file))

		var myData map[string]any
		shorterr.Check(json.Unmarshal(jsonData, &myData), "unmarshalling")

		val, ok := myData["name"]
		shorterr.Assert(ok, "missing name property")

		name, ok = val.(string)
		shorterr.Assert(ok, "invalid name property")

		return
	}

	f1()
	f2()
}

func TestCheck(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer shorterr.PassTo(&err)
		shorterr.Check(errFunc(x))
		return
	})
}

func TestMust(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer shorterr.PassTo(&err)
		a = argsToSlice(shorterr.Must(errFunc1(x)))
		return
	})
}

func TestMust2(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer shorterr.PassTo(&err)
		a = argsToSlice(shorterr.Must2(errFunc2(x)))
		return
	})
}

func TestMust3(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer shorterr.PassTo(&err)
		a = argsToSlice(shorterr.Must3(errFunc3(x)))
		return
	})
}

func TestMust4(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer shorterr.PassTo(&err)
		a = argsToSlice(shorterr.Must4(errFunc4(x)))
		return
	})
}

func TestMust5(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer shorterr.PassTo(&err)
		a = argsToSlice(shorterr.Must5(errFunc5(x)))
		return
	})
}

func TestDo(t *testing.T) {
	assert(t, "failed2: failed", func(x bool) (a []int, err error) {
		defer shorterr.PassTo(&err)
		a = argsToSlice(shorterr.Do(errFunc1(x)).Or("failed2"))
		return
	})
}

func TestDo2(t *testing.T) {
	assert(t, "failed2: failed", func(x bool) (a []int, err error) {
		defer shorterr.PassTo(&err)
		a = argsToSlice(shorterr.Do2(errFunc2(x)).Or("failed2"))
		return
	})
}

func TestDo3(t *testing.T) {
	assert(t, "failed2: failed", func(x bool) (a []int, err error) {
		defer shorterr.PassTo(&err)
		a = argsToSlice(shorterr.Do3(errFunc3(x)).Or("failed2"))
		return
	})
}

func TestDo4(t *testing.T) {
	assert(t, "failed2: failed", func(x bool) (a []int, err error) {
		defer shorterr.PassTo(&err)
		a = argsToSlice(shorterr.Do4(errFunc4(x)).Or("failed2"))
		return
	})
}

func TestDo5(t *testing.T) {
	assert(t, "failed2: failed", func(x bool) (a []int, err error) {
		defer shorterr.PassTo(&err)
		a = argsToSlice(shorterr.Do5(errFunc5(x)).Or("failed2"))
		return
	})
}

func TestAssert(t *testing.T) {
	assert(t, "failed", func(x bool) (a []int, err error) {
		defer shorterr.PassTo(&err)
		shorterr.Assert(x, "failed")
		return
	})
}

func TestOtherPanic(t *testing.T) {
	f := func() (err error) {
		defer shorterr.PassTo(&err)
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
