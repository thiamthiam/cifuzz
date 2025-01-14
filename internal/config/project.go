package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"code-intelligence.com/cifuzz/util/fileutil"
	"code-intelligence.com/cifuzz/util/stringutil"
)

const (
	BuildSystemCMake string = "cmake"
	BuildSystemOther string = "other"
)

var buildSystemTypes = []string{BuildSystemCMake, BuildSystemOther}

type ProjectConfig struct {
	LastUpdated string
	BuildSystem string `yaml:"build-system"`
}

const projectConfigFile = "cifuzz.yaml"

//go:embed cifuzz.yaml.tmpl
var projectConfigTemplate string

// CreateProjectConfig creates a new project config in the given directory
func CreateProjectConfig(projectDir string) (string, error) {

	// try to open the target file, returns error if already exists
	configpath := filepath.Join(projectDir, projectConfigFile)
	f, err := os.OpenFile(configpath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return configpath, errors.WithStack(err)
		}
		return "", errors.WithStack(err)
	}

	// setup config struct with (default) values
	config := ProjectConfig{
		LastUpdated: time.Now().Format("2006-01-02"),
	}

	// parse the template and write it to config file
	t, err := template.New("project_config").Parse(projectConfigTemplate)
	if err != nil {
		return "", errors.WithStack(err)
	}

	err = t.Execute(f, config)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return configpath, nil
}

func ParseProjectConfig(opts interface{}) (string, error) {
	var err error

	projectDir, err := FindProjectDir()
	if err != nil {
		return "", err
	}

	configpath := filepath.Join(projectDir, projectConfigFile)
	viper.SetConfigFile(configpath)

	err = viper.ReadInConfig()
	if err != nil {
		return "", errors.WithStack(err)
	}

	// viper.Unmarshal doesn't return an error if the timeout value is
	// missing a unit, so we check that manually
	if viper.GetString("timeout") != "" {
		_, err = time.ParseDuration(viper.GetString("timeout"))
		if err != nil {
			return "", errors.WithStack(fmt.Errorf("error decoding 'timeout': %w", err))
		}
	}

	err = viper.Unmarshal(opts)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return projectDir, nil
}

func ReadProjectConfig(projectDir string) (*ProjectConfig, error) {
	var err error

	configpath := filepath.Join(projectDir, projectConfigFile)
	viper.SetConfigFile(configpath)

	// Set defaults
	useSandboxDefault := runtime.GOOS == "linux"
	viper.SetDefault("sandbox", useSandboxDefault)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	config := &ProjectConfig{
		BuildSystem: viper.GetString("build-system"),
	}

	if config.BuildSystem == "" {
		config.BuildSystem, err = DetermineBuildSystem(projectDir)
		if err != nil {
			return nil, err
		}
	} else {
		err = ValidateBuildSystem(config.BuildSystem)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func ValidateBuildSystem(buildSystem string) error {
	if !stringutil.Contains(buildSystemTypes, buildSystem) {
		return errors.Errorf("Invalid build system \"%s\"", buildSystem)
	}
	return nil
}

func DetermineBuildSystem(projectDir string) (string, error) {
	isCMakeProject, err := fileutil.Exists(filepath.Join(projectDir, "CMakeLists.txt"))
	if err != nil {
		return "", err
	}
	if isCMakeProject {
		return BuildSystemCMake, nil
	} else {
		return BuildSystemOther, nil
	}
}

func FindProjectDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", errors.WithStack(err)
	}
	configFileExists, err := fileutil.Exists(filepath.Join(dir, projectConfigFile))
	if err != nil {
		return "", err
	}
	for !configFileExists {
		if dir == filepath.Dir(dir) {
			err := fmt.Errorf("not a cifuzz project (or any of the parent directories): %s %w", projectConfigFile, os.ErrNotExist)
			return "", errors.WithStack(err)
		}
		dir = filepath.Dir(dir)
		configFileExists, err = fileutil.Exists(filepath.Join(dir, projectConfigFile))
		if err != nil {
			return "", err
		}
	}

	dir, err = fileutil.CanonicalPath(dir)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return dir, nil
}
