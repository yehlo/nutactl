package config 

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"os"
	"hash/fnv"
	"strconv"
)

// Context struct that fits new configuration
type Context struct {
	ID	 			uint32
	URL 			string
	User	 		string
	PW				string
	Insecure	bool
}

// File contains fullpath to configfile
var File string

// Prefix contains viper prefix to use
var Prefix string

const (
	ntxConfigRoot string = "ntxContexts"
)

func init() { 
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	// initialize viper config
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(Prefix)
	viper.AutomaticEnv()
}

func generateID(url string, user string, insecure bool) (id uint32) {
	s := url + user + strconv.FormatBool(insecure)
	h := fnv.New32a()
	h.Write([]byte(s))
	id = h.Sum32()
	return id
}

func getConfigPath(ID uint32) string {
	configPath := fmt.Sprintf("%s.%d", ntxConfigRoot, ID) 
	return configPath
}

func deactivateContext(ID uint32) {
	
}

func getActiveContext() {

}

// GetContext returns context of specified ID as struct, returned value is nil if none is active
func GetContext(ID uint32) (*Context, error){
	if File == "" {
		return nil, fmt.Errorf("Please specify config file first through Config.File=fileName")
	}

	viper.SetConfigFile(File)

	if err := viper.ReadInConfig(); err != nil {
    if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, creating file
			os.OpenFile(File, os.O_RDONLY|os.O_CREATE, 0666)
			//since no context exists, double nil is returned
			return nil, nil
		}
		// config file found but other error was produced
		return nil, err
	}

	// set configpath and check if key exists. Return expected context
	configPath := getConfigPath(ID)
	if viper.IsSet(configPath) {
		return &Context{
			ID: 			viper.GetUint32("id"),
			URL: 			viper.GetString("api-url"),
			User:			viper.GetString("username"),
			PW:				viper.GetString("password"),
			Insecure:	viper.GetBool("insecure"),
		}, nil
	} 
	return nil, nil

}

// CreateContext creates new context
func CreateContext(url string, user string, pw string, insecure bool) (ID uint32) {
	ID = generateID(url, user, insecure)
	configPath := getConfigPath(ID)
	viper.Set(configPath + ".id", ID)
	viper.Set(configPath + ".url", url)
	viper.Set(configPath + ".user", user)
	viper.Set(configPath + ".pw", pw)
	viper.Set(configPath + ".insecure", insecure)
	viper.WriteConfig

	return ID
}

