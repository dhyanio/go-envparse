# go-envparse
This utility helps in parsing environment variables from the global environment or from a .env file. It allows the user to provide default values, validate environment variables, and customize the location of the .env file.

## Features
1. Dynamic Environment Variables: You can specify any set of environment variables to be parsed.
2. Default Value Support: If an environment variable is not found, you can provide a default value.
3. Validation: You can define custom validation functions for each environment variable to ensure it meets specific requirements (e.g., non-empty, specific format).
4. Flexible .env File Location: The parser supports custom paths for the .env file.
5. Logging: Logs each variable as it's successfully loaded, using default values or failing validation if applicable.
6. Case-Insensitive Matching: Environment variables are matched in a case-insensitive manner when read from the .env file.

```golang
package main

import (
  oktaUtils "https://github.com/dhyanio/go-envparse"
)

func main() {
	envVars := []EnvVar{
		{
			Key:          "CLIENT_ID",
			DefaultValue: "default_client_id",
			Validate: func(value string) bool {
				return len(value) > 0
			},
		},
		{
			Key:          "CLIENT_SECRET",
			DefaultValue: "default_client_secret",
			Validate: func(value string) bool {
				return len(value) >= 8
			},
		},
		{
			Key:          "ISSUER",
			DefaultValue: "https://default-issuer.com",
			Validate: func(value string) bool {
				return strings.HasPrefix(value, "https://")
			},
		},
	}

	// Call the parser with the specified variables and .env file location
	ParseEnvironment(envVars, ".env")
}
```

.env file
```env
CLIENT_ID=1234
CLIENT_SECRET=12345
ISSUER=https://dhyanio.com/oauth2/default
```
