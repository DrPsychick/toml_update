package main

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"os"
	"strconv"
	"strings"
)

func getEnv(pfx string) map[string]string {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, pfx) {
			continue
		}
		// PFX_VARIABLE=section.sub.name="value"
		parts := strings.SplitN(e, "=", 3)
		if len(parts) < 3 {
			continue
		}

		// trim quotes from the value
		envs[parts[1]] = strings.Trim(parts[2], "\"")
	}
	return envs
}

func parseValue(in string) interface{} {
	if i, err := strconv.ParseInt(in, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(in, 64); err == nil {
		return f
	}
	if b, err := strconv.ParseBool(in); err == nil {
		return b
	}
	if strings.HasPrefix(in, "[") {
		// ["foo", "bar"] - we have to remove the quotes or they will be escaped and included
		values := strings.Trim(in, "[]")
		if values == "" {
			return []string{}
		}

		list := strings.Split(values, ",")
		for i, v := range list {
			// trim quotes and spaces from each element
			list[i] = strings.Trim(v, "\" ")
		}
		return list
	}

	return string(in)
}

func updateConfigFromEnv(conf []byte, pfx string) ([]byte, error) {
	config, err := toml.Load(string(conf))
	if err != nil {
		return nil, err
	}

	envs := getEnv(pfx)
	if len(envs) == 0 {
		return nil, fmt.Errorf("No environment variables with prefix %s", pfx)
	}

	for k, v := range envs {
		withComment := false
		if strings.HasPrefix(v, "#") {
			withComment = true
			v = strings.TrimLeft(v, "#")
		}

		val := parseValue(v)

		if withComment {
			config.SetWithComment(k, "", true, val)
		} else {
			config.Set(k, val)
		}
	}

	out, err := toml.Marshal(config)
	return out, err
}

func main() {
	conf_file := os.Getenv("CONF_UPDATE")
	prefix := os.Getenv("CONF_PREFIX")
	if conf_file == "" || prefix == "" {
		fmt.Println("# No CONF_UPDATE or CONF_PREFIX defined - exiting.")
		os.Exit(0)
	}

	var buf []byte
	var err error
	if buf, err = os.ReadFile(conf_file); err != nil {
		fmt.Printf("Failed to read config file: %s\n", conf_file)
		os.Exit(1)
	}
	if buf, err = updateConfigFromEnv(buf, prefix); err != nil {
		fmt.Printf("Failed to update config from ENV: %s\n", err)
		os.Exit(1)
	}
	if err = os.WriteFile(conf_file, buf, 0644); err != nil {
		fmt.Printf("Failed to write back config to file '%s': %s\n", conf_file, err)
		os.Exit(1)
	}
}
