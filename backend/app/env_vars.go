package app

import "os"

type EnvVars struct {
	ApiURL    string
	AuthToken string
	DBURL     string
}

// Global variable to hold the initialized environment variables
// use this instead of calling os.Getenv directly
var EnvVarsValues = initEnvVars()

// avoid magic strings by defining env variable names (keys) in a struct
var envVarKeys = struct {
	externalApiUrl       string
	externalApiAuthToken string
	dbUrl                string
}{
	externalApiUrl:       "EXTERNAL_API_URL",
	externalApiAuthToken: "EXTERNAL_API_AUTH_TOKEN",
	dbUrl:                "DB_URL",
}

// Initialize environment variables and validate them
func initEnvVars() EnvVars {
	envVars := map[string]string{
		envVarKeys.externalApiUrl:       os.Getenv(envVarKeys.externalApiUrl),
		envVarKeys.externalApiAuthToken: os.Getenv(envVarKeys.externalApiAuthToken),
		envVarKeys.dbUrl:                os.Getenv(envVarKeys.dbUrl),
	}

	// iterate over the environment variables
	for key, value := range envVars {
		if value == "" {
			panic("Environment variable " + key + " is not set")
		}
	}
	return EnvVars{
		ApiURL:    envVars[envVarKeys.externalApiUrl],
		AuthToken: envVars[envVarKeys.externalApiAuthToken],
		DBURL:     envVars[envVarKeys.dbUrl],
	}
}
