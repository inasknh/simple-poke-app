package config

type DatabaseConfiguration struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type AppConfiguration struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
}

type Cache struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Api struct {
	Client   string `yaml:"client"`
	Host     string `yaml:"host"`
	BerryURL string `yaml:"berry_url"`
}
type Configurations struct {
	App      AppConfiguration      `yaml:"app"`
	Database DatabaseConfiguration `yaml:"database"`
	Cache    Cache                 `yaml:"cache"`
	Api      Api                   `yaml:"api"`
}
