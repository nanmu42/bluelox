package version

import (
	"fmt"
	"time"
)

const (
	// Name Program set name
	Name = "Bluelox"
)

// The following value is injected by build script
var (
	// Version git version
	Version = "delve"
	// BuildDate build date
	BuildDate = "unknown-build-date"
	// subName name specific to a command
	subName string
)

// The following value is initialized when program starts
//noinspection GoUnusedGlobalVariable
var (
	// FullName full name
	FullName string
	// FullNameWithBuildDate full name with build date
	FullNameWithBuildDate string
	// StartedAt time when this instance starts
	StartedAt time.Time
)

// SubName command name
func SubName() string {
	return subName
}

func SetSubName(newSubName string) {
	subName = newSubName
	updateFullNames()
}

func updateFullNames() {
	if subName != "" {
		FullName = fmt.Sprintf("%s-%s %s", Name, subName, Version)
		FullNameWithBuildDate = fmt.Sprintf("%s-%s %s (%s)", Name, subName, Version, BuildDate)
		return
	}

	FullName = fmt.Sprintf("%s %s", Name, Version)
	FullNameWithBuildDate = fmt.Sprintf("%s %s (%s)", Name, Version, BuildDate)
}

func init() {
	StartedAt = time.Now()
	updateFullNames()
}
