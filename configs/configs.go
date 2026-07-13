package configs

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type APIConfig struct {
	ServerPort   int    `envconfig:"SERVER_PORT" default:"8080"`
	DatabasePath string `envconfig:"DATABASE_PATH" default:"/data/app.db"`
	DatabaseSeed string `envconfig:"DATABASE_SEED" default:"true"`

	ServerJWTSecret string `envconfig:"SERVER_JWT_SECRET" required:"true"`

	Bootstrap BootstrapConfig
}

type BootstrapConfig struct {
	Mode string `envconfig:"BOOTSTRAP_MODE" default:"auto"`
	// auto  = seed/bootstrap only if DB is empty
	// none  = never seed/bootstrap
	// force = seed/bootstrap every startup, probably dev/test only

	OrgName       string `envconfig:"BOOTSTRAP_ORG_NAME" default:"Default Organization"`
	OrgPublicID   string `envconfig:"BOOTSTRAP_ORG_PUBLIC_ID" default:"org_default"`
	ProjectName   string `envconfig:"BOOTSTRAP_PROJECT_NAME" default:"Default Project"`
	AdminEmail    string `envconfig:"BOOTSTRAP_ADMIN_EMAIL" default:"admin@example.com"`
	AdminPassword string `envconfig:"BOOTSTRAP_ADMIN_PASSWORD" default:"change-me"`
}

func LoadConfig[T any](prefix string) (*T, error) {
	var conf T

	if err := loadEnvFile(".env", godotenv.Load); err != nil {
		return nil, fmt.Errorf("load base environment: %w", err)
	}

	envFiles := collectEnvFiles(conf)
	for _, file := range envFiles {
		if err := loadEnvFile(file, godotenv.Overload); err != nil {
			return nil, fmt.Errorf("load environment from %s: %w", file, err)
		}
	}

	if err := envconfig.Process(prefix, &conf); err != nil {
		return nil, fmt.Errorf("cannot create environment params: %w", err)
	}
	return &conf, nil
}

type envFileProvider interface {
	EnvFiles() []string
}

type envLoader func(filenames ...string) error

func loadEnvFile(filename string, loader envLoader) error {
	if strings.TrimSpace(filename) == "" {
		return nil
	}

	if err := loader(filename); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && errors.Is(pathErr.Err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return nil
}

func collectEnvFiles[T any](conf T) []string {
	seen := map[string]struct{}{}
	files := make([]string, 0)

	if provider, ok := any(&conf).(envFileProvider); ok {
		for _, file := range provider.EnvFiles() {
			trimmed := strings.TrimSpace(file)
			if trimmed == "" {
				continue
			}
			if _, exists := seen[trimmed]; exists {
				continue
			}
			files = append(files, trimmed)
			seen[trimmed] = struct{}{}
		}
	}

	if file := defaultEnvFileFor(conf); file != "" {
		if _, exists := seen[file]; !exists {
			files = append(files, file)
		}
	}

	return files
}

func defaultEnvFileFor[T any](conf T) string {
	t := reflect.TypeOf(conf)
	if t == nil {
		return ""
	}
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	name := strings.TrimSpace(t.Name())
	if name == "" {
		return ""
	}
	name = strings.TrimSuffix(name, "Config")
	if name == "" {
		return ""
	}
	return fmt.Sprintf(".env.%s", strings.ToLower(name))
}
