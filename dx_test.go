package font

import "testing"

type entry struct {
	text string
	pix  int
}

func TestDx(t *testing.T) {
	var tab = []entry{
		{"@", 7},
		{"@@", 14},
		{"hello", 35},
	}
	for _, v := range tab {
		runTest(t, v)
	}
}

func TestGoMonoCacheReplacer(t *testing.T) {
	var tab = []entry{
		{"@", 7},
		{"@@", 14},
		{"hello", 35},
		{"hello\xfdc", 57},
		{"hello\xfec", 57},
		{"hello\xffc", 57},
		{"hell\x01\xffc", 65},
	}
	for _, v := range tab {
		runReplacerTest(NewGoMono, t, v)
	}
}

func TestGoRegularCacheReplacer(t *testing.T) {
	var tab = []entry{
		{"@", 12},
		{"@@", 24},
		{"hello", 27},
		{"hello\xfdc", 48},
		{"hello\xfec", 48},
		{"hello\xffc", 48},
		{"hell\x01\xffc", 56},
	}
	for _, v := range tab {
		runReplacerTest(NewGoRegular, t, v)
	}
}

func TestSizeOfRuneFF(t *testing.T) {
	t.Skip("not a test")
	ft := NewCache(Replacer(NewGoRegular(11), NewHex(11), nil))
	t.Logf("%d\n", ft.Dx([]byte("\x00")))
	t.Logf("%d\n", ft.Dx([]byte("\xff")))
	t.Logf("%d\n", ft.Dx([]byte("\x00\x00\x00\xFF")))
}

func runReplacerTest(fn func(int) Face, t *testing.T, v entry) {
	t.Helper()
	ft := NewCache(Replacer(fn(11), NewHex(11), nil))
	have, want := ft.Dx([]byte(v.text)), v.pix
	if have != want {
		t.Logf("%s: have %d, want %d", v.text, want, have)
		t.Fail()
	}
}
func runTest(t *testing.T, v entry) {
	t.Helper()
	ft := NewGoMono(11)
	have, want := ft.Dx([]byte(v.text)), v.pix
	if have != want {
		t.Logf("%s: have %d, want %d", v.text, want, have)
		t.Fail()
	}
}
