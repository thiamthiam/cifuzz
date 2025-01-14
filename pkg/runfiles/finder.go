package runfiles

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

type RunfilesFinderImpl struct {
	InstallDir string
}

func (f RunfilesFinderImpl) CIFuzzIncludePath() (string, error) {
	return f.findFollowSymlinks("share/cifuzz/include/cifuzz")
}

func (f RunfilesFinderImpl) ClangPath() (string, error) {
	path, err := exec.LookPath("clang")
	return path, errors.WithStack(err)
}

func (f RunfilesFinderImpl) CMakePath() (string, error) {
	path, err := exec.LookPath("cmake")
	return path, errors.WithStack(err)
}

func (f RunfilesFinderImpl) CMakePresetsPath() (string, error) {
	return f.findFollowSymlinks("share/cifuzz/share/CMakePresets.json")
}

func (f RunfilesFinderImpl) JazzerAgentDeployJarPath() (string, error) {
	return f.findFollowSymlinks("bin/jazzer_driver")
}

func (f RunfilesFinderImpl) JazzerDriverPath() (string, error) {
	return f.findFollowSymlinks("bin/jazzer_driver")
}

func (f RunfilesFinderImpl) LLVMCovPath() (string, error) {
	path, err := exec.LookPath("llvm-cov")
	return path, errors.WithStack(err)
}

func (f RunfilesFinderImpl) LLVMProfDataPath() (string, error) {
	path, err := exec.LookPath("llvm-profdata")
	return path, errors.WithStack(err)
}

func (f RunfilesFinderImpl) LLVMSymbolizerPath() (string, error) {
	path, err := exec.LookPath("llvm-symbolizer")
	return path, errors.WithStack(err)
}

func (f RunfilesFinderImpl) Minijail0Path() (string, error) {
	return f.findFollowSymlinks("bin/minijail0")
}

func (f RunfilesFinderImpl) ProcessWrapperPath() (string, error) {
	return f.findFollowSymlinks("lib/process_wrapper")
}

func (f RunfilesFinderImpl) ReplayerSourcePath() (string, error) {
	return f.findFollowSymlinks("share/cifuzz/src/replayer.c")
}

func (f RunfilesFinderImpl) VSCodeTasksPath() (string, error) {
	return f.findFollowSymlinks("share/cifuzz/share/tasks.json")
}

func (f RunfilesFinderImpl) findFollowSymlinks(relativePath string) (string, error) {
	absolutePath := filepath.Join(f.InstallDir, relativePath)

	resolvedPath, err := filepath.EvalSymlinks(absolutePath)
	if err != nil {
		return "", errors.Wrapf(err, "path: %s", absolutePath)
	}
	_, err = os.Stat(resolvedPath)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return resolvedPath, nil
}
