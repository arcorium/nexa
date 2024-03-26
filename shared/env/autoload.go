package env

import "github.com/joho/godotenv"

func LoadEnvs(filenames ...string) error {
	return godotenv.Load(filenames...)
}
