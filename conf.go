package main

import (
	"github.com/tobscore/aity/mongo"
	"io/ioutil"
	"os/user"
)
import "github.com/BurntSushi/toml"

// Reads the applications' configuration from a fixed path. This is usually the user's home directoy
func ReadConf() (mongo.Conn, error) {
	conn := mongo.Conn{}

	var confBytes []byte
	userHome, err := user.Current()
	if err != nil {
		return conn, err
	}
	confBytes, err = ioutil.ReadFile(userHome.HomeDir + "/.aityconf")
	if err != nil {
		return conn, err
	}
	// Read the config file bytes as a string
	config := string(confBytes[:])

	// Decode the string representation to a conn struct
	_, err = toml.Decode(config, &conn)
	return conn, err
}
