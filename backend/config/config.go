package config

type Config struct {
    Port   string
    DBPath string
}

func Load() *Config {
    return &Config{
        Port:   ":8080",
        DBPath: "invoices.db",
    }
}