package sqlitedb

import "github.com/squeakycheese75/open-incentives/configs"

type Config struct {
	Path      string
	Bootstrap configs.BootstrapConfig
}
