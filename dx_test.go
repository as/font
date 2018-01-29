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

func runTest(t *testing.T, v entry) {
	t.Helper()
	ft := NewGoMono(11)
	have, want := ft.Dx([]byte(v.text)), v.pix
	if have != want {
		t.Logf("%s: have %d, want %d", v.text, want, have)
		t.Fail()
	}
}
