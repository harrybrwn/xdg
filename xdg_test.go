package xdg

import (
	"os"
	"testing"
)

func TestDefaultVal(t *testing.T) {
	unsetAll()
	defer func() {
		os.Unsetenv("HOME")
		unsetAll()
	}()
	name := "go-xdg-test"
	xdg := NewXDG(NewDirFinder(name))
	os.Setenv("HOME", "/home/t")
	eq(t, "/home/t/.go-xdg-test", xdg.getDir("unknown_key"))
	arrEq(t, nil, xdg.getDirs("unknown_key"))
	eq(t, "/home/t/.config/go-xdg-test", xdg.getDir(configHomeKey))
	eq(t, "/home/t/.cache/go-xdg-test", xdg.getDir(cacheHomeKey))
	eq(t, "/home/t/.local/share/go-xdg-test", xdg.getDir(dataHomeKey))
	eq(t, "/home/t/.local/state/go-xdg-test", xdg.getDir(stateHomeKey))

	eq(t, "/home/t/.config/go-xdg-test", Config(name))
	eq(t, "/home/t/.cache/go-xdg-test", Cache(name))
	eq(t, "/home/t/.local/share/go-xdg-test", Data(name))
	eq(t, "/home/t/.local/state/go-xdg-test", State(name))
	eq(t, "", Runtime(name))
	arrEq(t, []string{"/usr/local/share/go-xdg-test", "/usr/share/go-xdg-test"}, DataDirs(name))
	arrEq(t, []string{"/etc/xdg/go-xdg-test"}, ConfigDirs(name))

	os.Setenv(configHomeKey, "/h/t/.conf")
	os.Setenv(cacheHomeKey, "/h/t/.local/cache")
	os.Setenv(dataHomeKey, "/h/t/.local/share/dat")
	os.Setenv(stateHomeKey, "/h/t/.local/share/state")
	os.Setenv(runtimeDirKey, "/h/t/.local/share/run")
	os.Setenv(dataDirsKey, "/h/t/.local/share/datas:/xdg-data")
	os.Setenv(configDirsKey, "/h/t/.conf")
	eq(t, "/h/t/.conf/go-xdg-test", xdg.getDir(configHomeKey))
	eq(t, "/h/t/.local/cache/go-xdg-test", xdg.getDir(cacheHomeKey))
	eq(t, "/h/t/.local/share/dat/go-xdg-test", xdg.getDir(dataHomeKey))
	eq(t, "/h/t/.local/share/state/go-xdg-test", xdg.getDir(stateHomeKey))
	eq(t, "/h/t/.local/share/run/go-xdg-test", xdg.getDir(runtimeDirKey))
	arrEq(t, []string{"/h/t/.local/share/datas/go-xdg-test", "/xdg-data/go-xdg-test"}, DataDirs(name))
	arrEq(t, []string{"/h/t/.conf/go-xdg-test"}, ConfigDirs(name))
}

func TestGetDir_NoHome(t *testing.T) {
	os.Unsetenv("HOME")
	name := "go-xdg-test"
	xdg := NewXDG(NewDirFinder(name))
	res := xdg.getDir(configHomeKey)
	eq(t, "", res)
}

func TestDir(t *testing.T) {
	d := Dir("/tmp/me/.local/share/run/")
	eq(t, "/tmp/me/.local/share/run/", d.String())
	arrEq(t, d.Split(), []string{"tmp", "me", ".local", "share", "run"})
	eq(t, "/tmp/me/.local/share/run/x", d.Append("x").String())
}

func TestDir_Create(t *testing.T) {
	d := Dir("/tmp/me/.local/share/run")
	eq(t, exists(string(d)), d.Exists())
	if d.Exists() {
		t.Error("dir should not exist")
	}
	if err := d.Create(); err != nil {
		t.Error(err)
	}
	eq(t, exists(string(d)), d.Exists())
	if !exists(string(d)) {
		t.Error("dir should exist")
	}
	_ = os.RemoveAll(string(d))
}

func eq[T comparable](t *testing.T, a, b T) {
	t.Helper()
	if a != b {
		t.Errorf("\"%v\" not equal to \"%v\"", a, b)
	}
}

func arrEq[T comparable](t *testing.T, a, b []T) {
	t.Helper()
	if len(a) != len(b) {
		t.Errorf("arrays are different lengths")
		return
	}
	L := len(a)
	for i := 0; i < L; i++ {
		if a[i] != b[i] {
			t.Errorf("element %d: \"%v\" not equal to \"%v\"", i, a[i], b[i])
		}
	}
}

func unsetAll() {
	for _, key := range []string{
		configHomeKey,
		cacheHomeKey,
		dataHomeKey,
		stateHomeKey,
		runtimeDirKey,
		configDirsKey,
		dataDirsKey,
	} {
		err := os.Unsetenv(key)
		if err != nil {
			panic(err)
		}
	}
}
