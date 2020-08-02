package config 

import (
	"fmt"
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
// var Prefix string

// activeContext contains ID of currently active context
var activeContext uint32

const (
	ntxConfigRoot string = "ntxContexts"
)

func init() { 
	// initialize viper config
	viper.SetConfigType("yaml")
	err := readConfig()
	fmt.Println(err)
	fmt.Println("config was initialized")
	fmt.Println(viper.AllKeys())

	// set active context
	activeContext = viper.GetUint32(fmt.Sprintf("%s.active", ntxConfigRoot))
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

func setActiveContext(id uint32) error{
	viper.Set(fmt.Sprintf("%s.active", ntxConfigRoot), id)
	return viper.WriteConfig()
}

func readConfig() error {
	if File == "" {
		return fmt.Errorf("Please specify config file first through Config.File=fileName")
	}

	viper.SetConfigFile(File)

	if err := viper.ReadInConfig(); err != nil {
    if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, creating file
			os.OpenFile(File, os.O_RDONLY|os.O_CREATE, 0666)
		}
		// config file found but other error was produced
		return err
	}

	fmt.Println("myconfig")
	fmt.Println(viper.ConfigFileUsed())

	return nil
}

// GetActiveContext returns currently active context
func GetActiveContext() (*Context, error){
	return GetContext(activeContext)
}

// GetContext returns context of specified ID as struct, returned value is nil if none is active
func GetContext(ID uint32) (*Context, error){
	// set configpath and check if key exists. Return expected context
	configPath := getConfigPath(ID)
	if viper.IsSet(configPath) {
		return &Context{
			ID: 			viper.GetUint32(configPath + ".id"),
			URL: 			viper.GetString(configPath + ".api-url"),
			User:			viper.GetString(configPath + ".username"),
			PW:				viper.GetString(configPath + ".password"),
			Insecure:	viper.GetBool(configPath + ".insecure"),
		}, nil
	}

	return nil, fmt.Errorf("context with specified ID was not found")
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
	viper.WriteConfig()

	return ID
}

// SetContext sets existing context as active
func SetContext(ID uint32, pass string) (err error) {
	// nil password of previously active context
	configPath := getConfigPath(activeContext)
	viper.Set(configPath + ".pw", nil)

	// set new id of context
	err = setActiveContext(ID)
	if err != nil{
		return err
	}

	// set password to newly active context
	configPath = getConfigPath(activeContext)
	viper.Set(configPath + ".pw", pass)

	// save config before initializing context
	err = viper.SafeWriteConfig()
	if err != nil{
		return err
	}

	// initialize context
	InitContext()

	return nil
}

// InitContext creates viper alias for currently active config
// needed because the config file is organized different
func InitContext () {
	configPath := getConfigPath(activeContext)
	viper.RegisterAlias("id", configPath + ".id")
	viper.RegisterAlias("url", configPath + ".url")
	viper.RegisterAlias("user", configPath + ".user")
	viper.RegisterAlias("pw", configPath + ".pw")
	viper.RegisterAlias("insecure", configPath + ".insecure")
}