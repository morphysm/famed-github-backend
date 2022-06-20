package buildinfo

import (
	"runtime"
	"runtime/debug"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/rotisserie/eris"
)

const (
	// ProjectName is the stylized project name.
	ProjectName = "Famed"
	// ProgramName is the program name (should be in ASCII character for wide portability).
	ProgramName = "famed"
	// ProjectWebsite is the project's web page.
	ProjectWebsite = "https://www.famed.morphysm.com/"
)

type BuildInfo struct {
	Version         *semver.Version
	Date            time.Time
	Revision        string
	Target          string
	CompilerVersion string
}

func (i BuildInfo) UserAgent() string {
	return ProgramName + "/" + i.Version.String() + " (" + runtime.Version() + "; " + i.Target + ")"
}

func NewBuildInfo() (*BuildInfo, error) {
	buildInfo := &BuildInfo{
		Target:          runtime.GOOS + "/" + runtime.GOARCH,
		CompilerVersion: runtime.Version(),
	}

	// TODO: Retrieve version from vcs tag
	version, err := semver.NewVersion("0.0.0")
	if err != nil {
		return nil, eris.Wrap(err, "failed to instantiate new semver")
	}

	buildInfo.Version = version

	debugBuildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, eris.Wrap(err, "failed to read build info")
	}

	for _, kv := range debugBuildInfo.Settings {
		switch kv.Key {
		case "vcs.revision":
			buildInfo.Revision = kv.Value
		case "vcs.time":
			buildInfo.Date, _ = time.Parse(time.RFC3339, kv.Value)
		}
	}

	return buildInfo, nil
}
