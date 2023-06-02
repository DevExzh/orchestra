package utils

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shirou/gopsutil/v3/disk"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type (
	RuntimeVersion struct {
		Name string
		Path string
	}

	DetectedRuntime struct {
		name         string
		arguments    []string
		versions     []RuntimeVersion
		defaultIndex int
	}

	Searcher interface {
		GetFiles(string) []string
	}

	RuntimeSearcher struct {
		runtimes map[string]*DetectedRuntime
		database *RuntimeDatabase
		opt      []bool // 0: SaveToDatabase, 1: LoadFromDatabase
	}

	RuntimeDatabase struct {
		stat bool
		dbp  string  // Path of the database file
		db   *sql.DB // pointer representing a local database
	}

	RuntimeDeclaration struct {
		Type     string `json:"type"`
		Deployer string `json:"deployer,omitempty"`
		Version  string `json:"version,omitempty"`
	}

	RuntimeConfig struct {
		Type        string `json:"type"`
		Description string `json:"description,omitempty"`
		Versions    []struct {
			Version            string `json:"version"`
			Path               string `json:"path"`
			Platform           string `json:"platform,omitempty"`
			PlatformMinVersion string `json:"platform-min-version,omitempty"`
		} `json:"versions,omitempty"`
		DefaultVersion string `json:"default-version,omitempty"`
		Permissions    []struct {
			Name  string `json:"name,omitempty"`
			Allow bool   `json:"allow,omitempty"`
		} `json:"permissions,omitempty"`
		Environments []struct {
			Name    string `json:"name,omitempty"`
			Version string `json:"version,omitempty"`
			Value   string `json:"value,omitempty"`
		} `json:"environments,omitempty"`
		ApplyEnvironments bool `json:"apply-environments,omitempty"`
	}
)

func RuntimePath() ([]string, error) {
	var paths []string
	if pts, err := disk.Partitions(false); err == nil { // Get all partitions' information
		for _, pt := range pts {
			// Specific partition
			mp := pt.Mountpoint
			if !bytes.HasSuffix(ConvertStringToByteSlice(mp), []byte{os.PathSeparator}) {
				mp += ConvertByteSliceToString([]byte{os.PathSeparator})
			}
			rt := filepath.Join(mp, ".runtime") // Runtime Directory
			if _, st := os.Stat(rt); os.IsNotExist(st) {
				continue
			} else {
				j := filepath.Join(rt, "runtimes.json")
				if _, jt := os.Stat(j); os.IsNotExist(jt) {
					continue
				} else {
					d := &RuntimeDeclaration{}
					jf, er := os.ReadFile(j)
					if er != nil {
						continue
					}
					if m := json.Unmarshal(jf, &d); m != nil {
						continue
					}
					if d.Type == "application/x-devexzh-runtimes" {
						paths = append(paths, rt)
					} else {
						continue
					}
				}
			}
		}
		if len(paths) != 0 {
			return paths, nil
		} else {
			return nil, errors.New("cannot find the drive")
		}
	} else {
		return nil, err
	}
}

/**
DetectedRuntime SECTION START
*/

func (receiver *DetectedRuntime) AppendVersion(version RuntimeVersion) {
	receiver.versions = append(receiver.versions, version)
}

func (receiver *DetectedRuntime) DefaultVersion() *RuntimeVersion {
	return &receiver.versions[receiver.defaultIndex]
}

func (receiver *DetectedRuntime) Versions() []RuntimeVersion {
	return receiver.versions
}

func (receiver *DetectedRuntime) VersionCount() int {
	return len(receiver.versions)
}

func (receiver *DetectedRuntime) HasVersionAt(path string) bool {
	for _, v := range receiver.versions {
		a, _ := filepath.Abs(path)
		b, _ := filepath.Abs(v.Path)
		if a == b {
			return true
		}
	}
	return false
}

func (receiver *DetectedRuntime) Name() string {
	return receiver.name
}

/**
  DetectedRuntime SECTION END
*/

/**
  RuntimeSearcher SECTION START
*/

func (s *RuntimeSearcher) SetLoadFromDatabase(flag bool) {
	s.opt[1] = flag
}

func (s *RuntimeSearcher) SetSaveToDatabase(flag bool) {
	s.opt[0] = flag
}

func (s *RuntimeSearcher) GetFiles(path string) []string {
	f, _ := os.Open(path)
	abs, _ := filepath.Abs(path)
	files, _ := f.ReadDir(-1)
	if f.Close() != nil {
		fmt.Fprintln(os.Stderr)
	}
	retStr := make([]string, 16)
	for _, file := range files {
		if !file.IsDir() { // Append files only
			retStr = append(retStr, filepath.Join(abs, file.Name()))
		}
	}
	return retStr
}

func (s *RuntimeSearcher) findRuntime(displayName string, runtimeName string, processFunc func(string) string) *DetectedRuntime {
	var run = DetectedRuntime{
		name: displayName,
	}
	r := runtimeName
	if runtime.GOOS == "windows" { // DO NOT REMOVE THIS LINE
		r += ".exe"
	}
	for _, loc := range strings.Split(os.Getenv("Path"), string(os.PathListSeparator)) {
		_, err := os.Stat(loc)
		if err == nil { // The Path that environment variable refers to truly exists
			for _, file := range s.GetFiles(loc) {
				var f string = file
				if runtime.GOOS == "windows" && filepath.Ext(f) != ".exe" {
					continue
				}
				if fi, _ := os.Lstat(f); fi.Mode()&os.ModeSymlink != 0 {
					f, _ = os.Readlink(f)
				}
				if filepath.Base(f) == r && !run.HasVersionAt(f) {
					run.versions = append(run.versions, RuntimeVersion{Name: func(str string) string {
						if str == "" {
							return "unknown"
						} else {
							return str
						}
					}(processFunc(f)), Path: file})
				}
			}
		}
	}
	return &run
}

func (s *RuntimeSearcher) beautifyName(name string) string {
	switch name {
	case "python", "py":
		return "Python"
	case "node":
		return "Node.js"
	case "java":
		return "Java"
	default:
		return "Unknown language"
	}
}

func (s *RuntimeSearcher) FindRuntimesInPath() {
	p, err := RuntimePath()
	if err != nil {
		return
	}
	for _, pt := range p {
		f, _ := os.Open(pt)
		files, _ := f.ReadDir(-1)
		er := f.Close()
		if er != nil {
			continue
		}
		process := func(root, resolved string) string {
			res, cur := ConvertStringToByteSlice(resolved), ConvertStringToByteSlice(root)
			bytes.ReplaceAll(res, []byte{'/'}, []byte{os.PathSeparator})
			if bytes.HasPrefix(res, []byte{'.', '.'}) {
				bytes.Replace(res, []byte{'.', '.'}, ConvertStringToByteSlice(filepath.Base(root)), 1)
			}
			if bytes.HasPrefix(res, []byte{'.'}) {
				bytes.Replace(res, []byte{'.'}, cur, 1)
			}
			return ConvertByteSliceToString(res)
		}
		for _, file := range files {
			if file.IsDir() {
				rtc := filepath.Join(pt, file.Name(), "runtime.config.json")
				if _, e := os.Stat(rtc); e == nil {
					if ctn, r := os.ReadFile(rtc); r == nil {
						config := &RuntimeConfig{}
						if json.Unmarshal(ctn, &config) != nil {
							continue
						}
						d := s.runtimes[strings.ToLower(config.Type)]
						for _, ver := range config.Versions {
							d.AppendVersion(RuntimeVersion{
								Name: ver.Version,
								Path: process(filepath.Join(pt, file.Name()), ver.Path),
							})
						}
					} else {
						continue
					}
				} else {
					continue
				}
			}
		}
	}
}

func (s *RuntimeSearcher) FetchRuntimes() {
	// Fetch all runtimes in the current system via running with --version flag
	s.runtimes["python"] = s.findRuntime("Python", "python", func(loc string) string {
		o, _ := exec.Command(loc, "-V").CombinedOutput()
		return strings.TrimPrefix(strings.TrimSpace(string(o)), "Python ")
	})
	s.runtimes["java"] = s.findRuntime("Java", "java", func(loc string) string {
		o, _ := exec.Command(loc, "-version").CombinedOutput()
		ls := strings.Split(string(o), "\n")
		for _, l := range ls {
			i := strings.Index(l, "version")
			if i != -1 {
				var indices []int
				for index, character := range l {
					if character == '"' {
						indices = append(indices, index)
					}
				}
				if len(indices) == 2 {
					return l[indices[0]+1 : indices[1]]
				}
			}
		}
		/* Fetch information via reading the release file */
		j, _ := filepath.Abs(filepath.Join(filepath.Dir(loc), "../", "release"))
		_, e := os.Stat(j)
		if e == nil {
			if f, err := os.Open(j); err == nil {
				defer func(f *os.File) {
					_ = f.Close()
				}(f)
				r := bufio.NewReader(f)
				for {
					lb, er := r.ReadString('\n')
					switch {
					case er != nil:
						panic(er)
					case er == io.EOF:
						break
					}
					if !strings.Contains(lb, "JAVA_VERSION") {
						continue
					}
					v := strings.TrimSpace(strings.Split(lb, "=")[1])
					v = v[1 : len(v)-1]
					return v
				}
			} else {
				return ""
			}
		} else {
			return ""
		}
	})
	s.runtimes["node"] = s.findRuntime("Node.js", "node", func(loc string) string {
		o, _ := exec.Command(loc, "-v").Output()
		return strings.TrimPrefix(strings.TrimSpace(string(o)), "v")
	})
}

func (s *RuntimeSearcher) FindRuntimes() error {
	if len(s.runtimes) > 0 {
		return nil
	} // there's no need to load data again if objects already exist
	if s.database.IsReady() && s.opt[1] {
		r, _ := s.database.Query("select lang, ver, path from runtime;")
		defer r.Close()
		save := func() error {
			var lang, ver, vp string
			if err := r.Scan(&lang, &ver, &vp); err != nil {
				return err
			} // assign values to pointers
			value, ok := s.runtimes[lang]
			if ok {
				value.name = s.beautifyName(lang)
				value.versions = append(value.versions, RuntimeVersion{Name: ver, Path: vp})
			} else {
				rt := DetectedRuntime{
					name: s.beautifyName(lang),
				}
				rt.versions = append(rt.versions, RuntimeVersion{Name: ver, Path: vp})
				s.runtimes[lang] = &rt
			}
			return nil
		}
		if !r.Next() {
			s.FetchRuntimes()
			return s.StoreRuntimes()
		} else {
			if err := save(); err != nil {
				return err
			}
		}
		for r.Next() {
			if err := save(); err != nil {
				return err
			}
		}
	} else {
		s.FetchRuntimes()
		s.FindRuntimesInPath()
	}
	return nil
}

func (s *RuntimeSearcher) StoreRuntimes() error {
	if !s.database.stat {
		return errors.New("database not ready")
	}
	for k, v := range s.runtimes {
		for _, ver := range v.versions {
			var ve = ver.Name
			if ver.Name == "" {
				ve = "null"
			}
			_, er := s.database.db.Exec("insert into runtime(id,lang,ver,path) values (null,?,?,?);", k, ve, ver.Path)
			if er != nil {
				return er
			}
		}
	}
	return nil
}

func (s *RuntimeSearcher) SetDatabasePath(path string) error {
	return s.database.SetDatabasePath(path)
}

func (s *RuntimeSearcher) GetRuntime(lang string) *DetectedRuntime {
	return s.runtimes[lang]
}

func (s *RuntimeSearcher) Close() {
	if d := s.database.db; d != nil {
		d.Close()
	}
}

func (s *RuntimeSearcher) Runtimes() []DetectedRuntime {
	rts := make([]DetectedRuntime, len(s.runtimes))
	var i int = 0
	for _, rt := range s.runtimes {
		rts[i] = *rt
		i++
	}
	return rts
}

func (s *RuntimeSearcher) Runtime(language string) *DetectedRuntime {
	return s.runtimes[language]
}

/**
  RuntimeSearcher SECTION END
*/

func (db *RuntimeDatabase) Close() error {
	if db.db != nil {
		if err := db.db.Close(); err != nil {
			return err
		}
		db.db = nil
	}
	return nil
}

func (db *RuntimeDatabase) Renew() error {
	if err := db.prepare(); err != nil {
		return err
	}
	if _, err := db.db.Exec("drop table if exists runtime;"); err != nil {
		return err
	}
	if _, err := db.db.Exec("create table if not exists runtime(id integer primary key, lang varchar(16) not null, ver varchar(16), path text not null);"); err != nil {
		return err
	}
	return nil
}

func (db *RuntimeDatabase) IsReady() bool {
	return db.stat
}

func (db *RuntimeDatabase) SetDatabasePath(path string) error {
	stat, err := os.Stat(path)
	var er error
	if stat.IsDir() {
		db.dbp, er = filepath.Abs(filepath.Join(path, "runtimes.dat"))
		_, err = os.Stat(db.dbp)
	} else {
		db.dbp, er = filepath.Abs(path)
	}
	db.db, er = sql.Open("sqlite3", ProcessSQLDataSource(db.dbp))
	if er != nil {
		return er
	}
	if os.IsNotExist(err) {
		_, e := db.db.Exec("create table if not exists runtime(id integer primary key, lang varchar(16) not null, ver varchar(16), path text not null);")
		if e != nil {
			return e
		}
		db.stat = true
	}
	if er == nil {
		db.stat = true
	}
	return er
}

func (db *RuntimeDatabase) prepare() error {
	if !db.stat {
		if db.db == nil { // Initialize the database
			var err error
			db.db, err = sql.Open("sqlite3", ProcessSQLDataSource(db.dbp))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (db *RuntimeDatabase) Exec(sql string, args ...any) error {
	_, er := db.db.Exec(sql, args...)
	return er
}

func (db *RuntimeDatabase) Query(sql string, args ...any) (*sql.Rows, error) {
	return db.db.Query(sql, args...)
}

func NewRuntimeSearcher() *RuntimeSearcher {
	return &RuntimeSearcher{
		runtimes: make(map[string]*DetectedRuntime),
		database: &RuntimeDatabase{stat: false, dbp: "", db: nil},
		opt:      make([]bool, 2),
	}
}
