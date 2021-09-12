package morsica

import "testing"

func performEncodeTest(input string) string {
	alphabet := NewAlphabet()
	return alphabet.Encode(input)
}

func performDecodeTest(input string) string {
	alphabet := NewAlphabet()
	return alphabet.Decode(input)
}

func printTestResult(name string, want string, got string, t *testing.T) {
	if want != got {
		t.Fatalf(`%q: expected %q, got %q`, name, want, got)
	}
}

func TestEncode(t *testing.T) {
	want := "...   ---   ...       .----"
	morse := performEncodeTest("SOS 1")
	printTestResult(`Encode("SOS 1")`, want, morse, t)
}

func TestEncodeEmpty(t *testing.T) {
	want := ""
	morse := performEncodeTest("")
	printTestResult(`Encode("")`, want, morse, t)
}

func TestDecode(t *testing.T) {
	want := "sos 1"
	txt := performDecodeTest("...   ---   ...       .----")
	printTestResult(`Decode("...   ---   ...       .----")`, want, txt, t)
}
