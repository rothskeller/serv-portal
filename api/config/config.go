// Package config reads config.json and provides site-specific and/or private
// data to the rest of the application.
package config

import (
        "encoding/json"
        "os"
)

var config map[string]string

// Get returns the named configuration variable.
func Get(key string) string {
        if config == nil {
                var (
                        cf  *os.File
                        err error
                )
                if cf, err = os.Open("config.json"); err != nil {
                        panic("can't read config.json: " + err.Error())
                }
                defer cf.Close()
                if err = json.NewDecoder(cf).Decode(&config); err != nil {
                        panic("can't parse config.json: " + err.Error())
                }
        }
        return config[key]
}
