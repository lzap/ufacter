package common

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type netmaskTPair struct {
	in  uint8
	out string
	err error
}

type netmaskBPair struct {
	in  []byte
	out string
	err error
}

type byteTPair struct {
	in      uint64
	inUnit  string
	out     float64
	outUnit string
}

var EPSILON float64 = 1

func floatEquals(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func TestConvertBytes(t *testing.T) {
	testpairs := []byteTPair{
		{0, "", 0, "bytes"},
		{5238784, "", 4.996094, "MiB"},
		{10485760, "", 10.000000, "MiB"},
		{100910080, "", 96.235352, "MiB"},
		{3267260416, "", 3.042873, "GiB"},
		{306816327680, "", 285.744972, "GiB"},
	}
	for _, pair := range testpairs {
		out, outUnit, err := ConvertBytes(pair.in)
		if err != nil {
			t.Logf("input: %v; %v != nil", pair, err)
			t.Fail()
		}
		if strings.Compare(outUnit, pair.outUnit) != 0 {
			t.Logf("input: %v; '%v' != '%v'", pair, outUnit, pair.outUnit)
			t.Fail()
		}
		equality := floatEquals(out, pair.out)
		if equality == false {
			t.Logf("input: %v; %f != %f; diff: %f, %v", pair, out, pair.out,
				(out - pair.out), equality)
			t.Fail()
		}
	}
}

func TestConvertBytesTo(t *testing.T) {
	testpairs := []byteTPair{
		{0, "MiB", 0, ""},
		{5238784, "kB", 5116, ""},
		{10485760, "kB", 10240, ""},
		{100910080, "MiB", 96.235352, ""},
		{3267260416, "MiB", 3115.902344, ""},
		{306816327680, "MiB", 292602.851562, ""},
	}
	for _, pair := range testpairs {
		out, outUnit, err := ConvertBytesTo(pair.in, pair.inUnit)
		if err != nil {
			t.Logf("input: %v; %v != nil", pair, err)
			t.Fail()
		}
		if strings.Compare(outUnit, pair.inUnit) != 0 {
			t.Logf("input: %v; '%v' != '%v'", pair, outUnit, pair.inUnit)
			t.Fail()
		}
		equality := floatEquals(out, pair.out)
		if equality == false {
			t.Logf("input: %v; %f != %f; diff: %f, %v", pair, out, pair.out,
				(out - pair.out), equality)
			t.Fail()
		}
	}
}

func TestConvertNetmask(t *testing.T) {
	testpairs := []netmaskTPair{
		{8, "255.0.0.0", nil},
		{17, "255.255.128.0", nil},
		{24, "255.255.255.0", nil},
		{28, "255.255.255.240", nil},
		{33, "", fmt.Errorf("Invalid Netmask given.")},
	}
	for _, pair := range testpairs {
		out, err := ConvertNetmask(pair.in)
		if err == nil && pair.err != nil {
			t.Fatalf("%v != %v", err, pair.err)
		}
		if pair.err != nil {
			// probably not safe to continue
			continue
		}
		if out != pair.out {
			t.Fatalf("%v != %v", out, pair.out)
		}
	}
}
func TestIPMaskToString4(t *testing.T) {
	testpairs := []netmaskBPair{
		{[]byte{255, 0, 0, 0}, "255.0.0.0", nil},
		{[]byte{255, 255, 0, 0}, "255.255.0.0", nil},
		{[]byte{255, 255, 255, 0}, "255.255.255.0", nil},
		{[]byte{255, 255, 42, 0}, "255.255.42.0", nil},
	}
	for _, pair := range testpairs {
		out := IPMaskToString4(pair.in)
		if out != pair.out {
			t.Fatalf("%v != %v", out, pair.out)
		}
	}
}

func TestIPMaskToString6(t *testing.T) {
	testpairs := []netmaskBPair{
		{[]byte{255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, "ff0::", nil},
		{[]byte{255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, "ffff::", nil},
		{[]byte{255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, "ffff:ff0::", nil},
	}
	for _, pair := range testpairs {
		out := IPMaskToString6(pair.in)
		if out != pair.out {
			t.Fatalf("%v != %v", out, pair.out)
		}
	}
}

func TestGetHostEtc(t *testing.T) {
	testValue := "test_value"
	err := os.Setenv("HOST_ETC", testValue)
	if err != nil {
		t.Fatalf("%v", err)
	}
	value := GetHostEtc()
	if strings.Compare(value, testValue) != 0 {
		t.Fatalf("%v != %v", value, testValue)
	}
}

func TestGetHostEtcNotSet(t *testing.T) {
	expectedVal := "/etc"
	err := os.Unsetenv("HOST_ETC")
	if err != nil {
		t.Fatalf("%v", err)
	}
	value := GetHostEtc()
	if strings.Compare(value, expectedVal) != 0 {
		t.Fatalf("%v != %v", value, expectedVal)
	}
}

func TestGetHostProc(t *testing.T) {
	testValue := "test_value"
	err := os.Setenv("HOST_PROC", testValue)
	if err != nil {
		t.Fatalf("%v", err)
	}
	value := GetHostProc()
	if strings.Compare(value, testValue) != 0 {
		t.Fatalf("%v != %v", value, testValue)
	}
}

func TestGetHostProcNotSet(t *testing.T) {
	expectedVal := "/proc"
	err := os.Unsetenv("HOST_PROC")
	if err != nil {
		t.Fatalf("%v", err)
	}
	value := GetHostProc()
	if strings.Compare(value, expectedVal) != 0 {
		t.Fatalf("%v != %v", value, expectedVal)
	}
}

func TestGetHostSys(t *testing.T) {
	testValue := "test_value"
	err := os.Setenv("HOST_SYS", testValue)
	if err != nil {
		t.Fatalf("%v", err)
	}
	value := GetHostSys()
	if strings.Compare(value, testValue) != 0 {
		t.Fatalf("%v != %v", value, testValue)
	}
}

func TestGetHostSysNotSet(t *testing.T) {
	expectedVal := "/sys"
	err := os.Unsetenv("HOST_SYS")
	if err != nil {
		t.Fatalf("%v", err)
	}
	value := GetHostSys()
	if strings.Compare(value, expectedVal) != 0 {
		t.Fatalf("%v != %v", value, expectedVal)
	}
}
