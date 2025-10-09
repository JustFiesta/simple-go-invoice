package config

type Config struct {
    Port   string
    DBPath string
}

func Load() *Config {
    return &Config{
        Port:   ":8000",
        DBPath: "invoices.db",
    }
}