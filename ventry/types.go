package ventry

import "github.com/gofrs/flock"

// Vers tracks single program version.
type Vers struct {
	Tag    string
	Prefix string
	Suffix string
	Major  int
	Minor  int
	Patch  int
}

// Entries one or more versions.
type Entries map[string]*Vers

// Rollback is how we keep a history (for now a single item)
type Rollback map[string]Vers

// VFile represents the format we write to the version file it
// has the current version and a history/rollback hash and array
type VFile struct {
	Version Entries
	Prev    Rollback
}

// Vers is a file locked instance of entries
type VEntry struct {
	lck  *flock.Flock
	path string
	ent  *VFile
}
