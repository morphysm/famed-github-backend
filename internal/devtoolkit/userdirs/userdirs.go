// Package userdirs implements a structure and methods to handle user-specific directories.
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
// https://wiki.archlinux.org/title/XDG_Base_Directory
package userdirs

import (
	"os"
	"path/filepath"

	"github.com/rotisserie/eris"
)

// UserDirs structure contains the different user-specific paths.
type UserDirs struct {
	// ConfigHome (XDG_CONFIG_HOME) defines the base directory relative to which user-specific configuration files should be stored.
	ConfigHome string
	// CacheHome (XDG_CACHE_HOME) defines the base directory relative to which user-specific non-essential data files should be stored.
	CacheHome string
}

// NewUserDirs create a new userDirs object with correct paths.
func NewUserDirs(programName string) (*UserDirs, error) {
	userDirs := &UserDirs{}

	if programName == "" {
		return nil, eris.New("program name cannot be empty")
	}

	// ConfigHome
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, eris.Wrapf(err, "failed to find ConfigHome path for %q", programName)
	}

	userDirs.ConfigHome = filepath.Join(userConfigDir, programName)

	// CacheHome
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, eris.Wrapf(err, "failed to find CacheHome path for %q", programName)
	}

	userDirs.CacheHome = filepath.Join(userCacheDir, programName)

	err = userDirs.makePaths()
	if err != nil {
		return nil, eris.Wrapf(err, "failed to make paths for %q", programName)
	}

	return userDirs, nil
}

// MakePaths makes sure every paths exists, with the defaultFileMode permissions.
// This function also check if directories exists.
// You should only call this method once.
func (u *UserDirs) makePaths() error {
	const fileMode = os.FileMode(0o755)

	// ConfigHome
	err := os.MkdirAll(u.ConfigHome, fileMode)
	if err != nil {
		return eris.Wrap(err, "failed to make ConfigHome directory")
	}

	// CacheHome
	err = os.MkdirAll(u.CacheHome, fileMode)
	if err != nil {
		return eris.Wrap(err, "failed to make CacheHome directory")
	}

	return nil
}
