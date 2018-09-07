package main

import (
	"github.com/tobscore/aity/mongo"
	"os/user"
)
import "github.com/BurntSushi/toml"

// Reads the applications' configuration from a fixed path. This is usually the user's home directoy
func ReadConf() (mongo.Conn, error) {
	conn := mongo.Conn{}

	userHome, err := user.Current()
	if err != nil {
		return conn, err
	}
	confFile := userHome.HomeDir + "/.aityconf"

	// Decode the string representation to a conn struct
	_, err = toml.DecodeFile(confFile, &conn)
	return conn, err
}
