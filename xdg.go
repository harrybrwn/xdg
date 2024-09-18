// Package xdg is a helper package to get configuration directories according
// to the XDG Base Directory Specification
//
// See docs:
//
//	https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
package xdg

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	listSeparator     = string(filepath.ListSeparator)
	defaultHomeBase   = ".config"
	defaultCacheBase  = ".cache"
	defaultDataBase   = ".local/share"
	defaultStateBase  = ".local/state"
	defaultDataDirs   = "/usr/local/share/:/usr/share/"
	defaultConfigDirs = "/etc/xdg"
)

const (
	configHomeKey = "XDG_CONFIG_HOME"
	cacheHomeKey  = "XDG_CACHE_HOME"
	dataHomeKey   = "XDG_DATA_HOME"
	stateHomeKey  = "XDG_STATE_HOME"
	runtimeDirKey = "XDG_RUNTIME_DIR"

	dataDirsKey   = "XDG_DATA_DIRS"
	configDirsKey = "XDG_CONFIG_DIRS"
)

func Config(name string) string       { return newXdg(name).Config() }
func State(name string) string        { return newXdg(name).State() }
func Data(name string) string         { return newXdg(name).Data() }
func Cache(name string) string        { return newXdg(name).Cache() }
func Runtime(name string) string      { return newXdg(name).Runtime() }
func ConfigDirs(name string) []string { return newXdg(name).ConfigDirs() }
func DataDirs(name string) []string   { return newXdg(name).DataDirs() }

func newXdg(name string) *XDG { return NewXDG(NewDirFinder(name)) }

type Dir string

func (d Dir) Exists() bool           { return exists(string(d)) }
func (d Dir) Create() error          { return os.MkdirAll(string(d), 0755) }
func (d Dir) String() string         { return string(d) }
func (d Dir) Append(name string) Dir { return Dir(filepath.Join(string(d), name)) }

func (d Dir) Split() []string {
	p := strings.Split(string(d), string(filepath.Separator))
	if len(p) > 0 {
		if len(p[0]) == 0 {
			p = p[1:]
		}
		if len(p[len(p)-1]) == 0 {
			p = p[0 : len(p)-1]
		}
	}
	return p
}

type DirFinder interface {
	Name() string
}

type XDG struct {
	finder DirFinder
}

func NewXDG(finder DirFinder) *XDG { return &XDG{finder: finder} }

func (xdg *XDG) Config() string       { return xdg.getDir(configHomeKey) }
func (xdg *XDG) Cache() string        { return xdg.getDir(cacheHomeKey) }
func (xdg *XDG) Data() string         { return xdg.getDir(dataHomeKey) }
func (xdg *XDG) State() string        { return xdg.getDir(stateHomeKey) }
func (xdg *XDG) Runtime() string      { return xdg.getDir(runtimeDirKey) }
func (xdg *XDG) ConfigDirs() []string { return xdg.getDirs(configDirsKey) }
func (xdg *XDG) DataDirs() []string   { return xdg.getDirs(dataDirsKey) }

func (xdg *XDG) getDir(key string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return filepath.Join(val, xdg.finder.Name())
	}
	switch key {
	case runtimeDirKey:
		return ""
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	def := xdg.defaultVal(home, key)
	if len(def) > 0 {
		return def
	}
	return filepath.Join(home, "."+xdg.finder.Name())
}

func (xdg *XDG) getDirs(key string) []string {
	var p string
	v, ok := os.LookupEnv(key)
	if ok {
		p = v
	} else {
		switch key {
		case dataDirsKey:
			p = defaultDataDirs
		case configDirsKey:
			p = defaultConfigDirs
		default:
			p = ""
		}
	}
	if len(p) > 0 {
		paths := filepath.SplitList(p)
		name := xdg.finder.Name()
		for i := range paths {
			paths[i] = filepath.Join(paths[i], name)
		}
		return paths
	}
	return nil
}

func (xdg *XDG) defaultVal(home, key string) string {
	var base string
	switch strings.ToUpper(key) {
	case configHomeKey:
		base = filepath.Join(home, defaultHomeBase)
	case cacheHomeKey:
		base = filepath.Join(home, defaultCacheBase)
	case dataHomeKey:
		base = filepath.Join(home, defaultDataBase)
	case stateHomeKey:
		base = filepath.Join(home, defaultStateBase)
	default:
		return ""
	}
	return filepath.Join(base, xdg.finder.Name())
}

func NewDirFinder(name string) *dirFinder { return &dirFinder{name} }

type dirFinder struct{ name string }

func (df *dirFinder) Name() string { return df.name }

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
