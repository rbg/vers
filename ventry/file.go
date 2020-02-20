package ventry

// Copyright © 2020 Robert B Gordon <rbg@h9k.io>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/gofrs/flock"
	"golang.org/x/sys/unix"

	"gopkg.in/yaml.v2"
)

// writeVersionFile updates the version file info
func writeVersionFile(path string, info *VFile) error {
	var bytes []byte

	log.Debugf("writeVersionFile: %+#v", info)
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	// get the type we can handle json or yaml
	switch ext := filepath.Ext(p); ext {
	case ".json":
		bytes, err = json.Marshal(info)
		if err != nil {
			return err
		}

	case ".yml":
		fallthrough
	case ".yaml":
		bytes, err = yaml.Marshal(info)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported file type")
	}
	err = ioutil.WriteFile(p, bytes, 0640)
	if err != nil {
		return err
	}
	return nil
}

// readVersionFile gets the version file info
func readVersionFile(path string) (*VFile, error) {
	var info VFile

	p, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	// get the type we can handle json or yaml
	switch ext := filepath.Ext(p); ext {
	case ".json":
		if err := json.Unmarshal([]byte(data), &info); err != nil {
			return nil, err
		}

	case ".yml":
		fallthrough
	case ".yaml":
		if err := yaml.Unmarshal([]byte(data), &info); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported file type")
	}
	return &info, nil
}

// Open a version entries file.
func Open(path string, creat bool) (*VEntry, error) {
	var ve VEntry

	p, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	// figure out if we should create..
	oFlg := os.O_RDWR
	if creat {
		oFlg |= os.O_CREATE
	}
	f, err := os.OpenFile(p, oFlg, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ve.path = p
	ve.lck = flock.New(p + ".lck")
	ve.ent = &VFile{
		Version: make(Entries),
		Prev:    make(Rollback),
	}
	log.Debugf("Open(): path->%s VEntry->%+#v", path, ve)
	return &ve, nil
}

// Path returns the current path
func (v *VEntry) Path() string {
	return v.path
}

// LPath returns the current lock file path
func (v *VEntry) LPath() string {
	return v.lck.Path()
}

// Add will update/add an entry
func (v *VEntry) Add(name string, ent *Vers) {
	log.Debugf("Add(): enrty->%s values->%+#v", name, ent)
	// push current values to history
	if ve, ok := v.ent.Version[name]; ok {
		v.ent.Prev[name] = *ve
	}
	v.ent.Version[name] = ent
}

// Rm will remove an entry
func (v *VEntry) Rm(name string) {
	if _, ok := v.ent.Version[name]; ok {
		delete(v.ent.Version, name)
	}
	if _, ok := v.ent.Prev[name]; ok {
		delete(v.ent.Prev, name)
	}
}

// Dump will dump entries
func (v *VEntry) Dump(format string) error {
	// get the type we can handle json or yaml
	switch format {
	case "str":
		fallthrough
	case "shell":
		for name, _ := range v.ent.Version {
			v.Print(name, format)
		}
	case "json":
		out, err := json.MarshalIndent(v.ent.Version, "", "   ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
	case "yml":
		fallthrough
	case "yaml":
		out, err := yaml.Marshal(v.ent.Version)
		if err != nil {
			return err
		}
		fmt.Println(string(out))
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

// Print will dump name(d) entries
func (v *VEntry) Print(name, format string) error {
	ve, ok := v.ent.Version[name]
	if !ok {
		return fmt.Errorf("%s; does not exist", name)
	}
	ent := make(Entries)
	ent[name] = ve
	// get the type we can handle json or yaml
	switch format {
	case "shell":
		str := fmt.Sprintf("export %s_VERS=%s%d.%d.%d%s", strings.ToUpper(name),
			ve.Prefix, ve.Major, ve.Minor, ve.Patch, ve.Suffix)
		fmt.Println(strings.ReplaceAll(str, "-", "_"))
	case "str":
		fmt.Printf("%s%d.%d.%d%s\n",
			ve.Prefix, ve.Major, ve.Minor, ve.Patch, ve.Suffix)
	case "json":
		out, err := json.MarshalIndent(ent, "", "   ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
	case "yml":
		fallthrough
	case "yaml":
		out, err := yaml.Marshal(ent)
		if err != nil {
			return err
		}
		fmt.Println(string(out))
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

// Read reads the entries file, populates the hash
func (v *VEntry) Read(retry int) error {
	var rt int
	log.Debugf("Read(): v->%#+v", v)
	for rt < retry {
		rt++
		ok, err := v.lck.TryRLock()
		if err != nil {
			return err
		}
		if ok {
			defer v.lck.Unlock()
			v.ent, err = readVersionFile(v.path)
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Debugf("Read(): did not lock file", v.Path())
	return fmt.Errorf("unable to get lock")
}

// Write will write the entries to stable store.
func (v *VEntry) Write(retry int) error {
	var rt int
	log.Debugf("Write(): v->%#+v", v)
	for rt < retry {
		rt++
		ok, err := v.lck.TryLock()
		if err != nil {
			return err
		}
		if ok {
			defer v.lck.Unlock()
			return writeVersionFile(v.path, v.ent)
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Debugf("Write(): did not lock file", v.Path())
	return fmt.Errorf("unable to get lock")
}

// Bump will inc the value of version field
func (v *VEntry) Bump(name, what string) error {
	var rt int
	for rt < 10 {
		rt++
		ok, err := v.lck.TryLock()
		if err != nil {
			return err
		}
		if ok {
			defer v.lck.Unlock()
			v.ent, err = readVersionFile(v.path)
			ve, ok := v.ent.Version[name]
			if !ok {
				return fmt.Errorf("%s; does not exist", name)
			}
			// push current values to history
			v.ent.Prev[name] = *ve
			switch what {
			case "major":
				ve.Major++
				ve.Minor = 0
				ve.Patch = 0
			case "minor":
				ve.Minor++
				ve.Patch = 0
			case "patch":
				ve.Patch++
			default:
				return errors.New("Invalid bump setting")
			}
			return writeVersionFile(v.path, v.ent)
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Debugf("Bump(): did not lock file", v.Path())
	return errors.New("Did not obtain lock")
}

// Undo restore previous value
func (v *VEntry) Undo(name string) error {
	var rt int
	for rt < 10 {
		rt++
		ok, err := v.lck.TryLock()
		if err != nil {
			return err
		}
		if ok {
			defer v.lck.Unlock()
			v.ent, err = readVersionFile(v.path)
			if _, ok := v.ent.Version[name]; !ok {
				return fmt.Errorf("%s; does not exist", name)
			}

			if ve, ok := v.ent.Prev[name]; ok {
				v.ent.Version[name] = &ve
				delete(v.ent.Prev, name)
				return writeVersionFile(v.path, v.ent)
			}
			return fmt.Errorf("%s; previous value does not exist", name)
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Debugf("Bump(): did not lock file", v.Path())
	return errors.New("Did not obtain lock")
}

// Delete will remove the entry
func (v *VEntry) Delete(name string) error {
	var rt int
	for rt < 10 {
		rt++
		ok, err := v.lck.TryLock()
		if err != nil {
			return err
		}
		if ok {
			defer v.lck.Unlock()
			v.ent, err = readVersionFile(v.path)
			_, ok := v.ent.Version[name]
			if !ok {
				return fmt.Errorf("%s; does not exist", name)
			}
			v.Rm(name)
			return writeVersionFile(v.path, v.ent)
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Debugf("Bump(): did not lock file", v.Path())
	return errors.New("Did not obtain lock")
}

// Unlock and close file
func (v *VEntry) Close() {
	unix.Unlink(v.lck.Path())
	v.lck.Close()
}
