package ventry

import "github.com/gofrs/flock"

// Vers tracks single program version.
type Vers struct {
	Prefix string
	Suffix string
	Major  int
	Minor  int
	Patch  int
}

// Entries one or more versions.
type Entries map[string]*Vers

// Vers is a file locked instance of entries
type VEntry struct {
	lck  *flock.Flock
	path string
	ent  Entries
}
