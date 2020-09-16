package config 

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
	"hash/fnv"
	"strconv"
	"bytes"	
	"encoding/json"
)

// Context struct that fits new configuration
type Context struct {
	ID	 			uint32
	URL 			string
	User	 		string
	Insecure	bool
}

// ContextList struct that is used for display operations
type ContextList struct {
	Entities []*Context `json:"entities,omitempty"`
}

// File contains fullpath to configfile
var File string

// activeContext contains ID of currently active context
var activeContext uint32

const (
	ntxConfigRoot string = "ntxContexts"
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil{
		panic(err)
	}

	viper.AddConfigPath(home + "/.nutactl")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
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
	activeContext = id
	return safeWriteConfig()
}

func readConfig() error {
	if File == "" {
		return fmt.Errorf("Please specify config file first through Config.File=fileName")
	}

	viper.SetConfigFile(File)

	if err := viper.ReadInConfig(); err != nil {
    if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// config file found but other error was produced
			return err
		}
		// config file not found
		os.OpenFile(File, os.O_RDONLY|os.O_CREATE, 0666)
		if err := viper.ReadInConfig(); err != nil {
			// any other error
			return err
		}
	}

	return nil
}

// makes sure that no plain pw is written to cfgfile
func safeWriteConfig() error {
	// loop through all active keys 
	configMap := viper.AllSettings()
	for _, key := range viper.AllKeys(){
		// restructured into if-else-if to make it more readable
		if strings.Contains(strings.ToLower(key), "pw") || strings.Contains(strings.ToLower(key), "password"){
			// delete any keys related to pw
			delete(configMap, key)
		} else if (! strings.Contains(key, ".")){
			// delete any keys that are not atleast one level nested
			delete(configMap, key)
		}
	}

	// new viper instance is needed, so that only needed values are written to the file
	// if no new instance is created, all keys are still there (even if used with readConfig as below)
	newViper := viper.New()
	home, err := os.UserHomeDir()
	if err != nil{
		return err
	}
	newViper.SetConfigFile(File)
	newViper.AddConfigPath(home + "/.nutactl")
	viper.SetConfigName("config")
	newViper.SetConfigType("yaml")

	encodedConfig, _ := json.MarshalIndent(configMap, "", " ")
	err = newViper.ReadConfig(bytes.NewReader(encodedConfig))
	if err != nil{
		return err
	}

	err = newViper.WriteConfig()
	if err != nil{
		return err
	}

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
			Insecure:	viper.GetBool(configPath + ".insecure"),
		}, nil
	}

	return nil, fmt.Errorf("context with specified ID was not found")
}

// GetAllContexts returns a Slice of all existing contexts
func GetAllContexts() (contexts []Context, err error){
	// currently dont know any other way than reading through all keys
	// and checking if the key contains the string 'id'
	for _, key := range viper.AllKeys(){
		if strings.Contains(strings.ToLower(key), "id"){
			context, err := GetContext(viper.GetUint32(key))
			if err != nil {
				return nil, err
			}

			contexts = append(contexts, *context)
			// contexts.Entities = append(contexts.Entities, context)
		}
	}
	
	return contexts, nil
}

// CreateContext creates new context
func CreateContext(url string, user string, insecure bool) (ID uint32, err error) {
	ID = generateID(url, user, insecure)
	configPath := getConfigPath(ID)
	viper.Set(configPath + ".id", ID)
	viper.Set(configPath + ".url", url)
	viper.Set(configPath + ".user", user)
	viper.Set(configPath + ".insecure", insecure)
	// set newely created context as active
	err = setActiveContext(ID)
	if err != nil{
		return 0, err
	}

	// write config
	err = safeWriteConfig()
	if err != nil{
		return 0, err
	}

	return ID, nil
}

// SetContext sets existing context as active
func SetContext(ID uint32) (err error) {
	// set new id of context
	err = setActiveContext(ID)
	if err != nil{
		return err
	}

	// save config before initializing context
	err = safeWriteConfig()
	if err != nil{
		return err
	}
	return nil
}

// RemoveContext removes existing context
func RemoveContext (ID string) (err error){
	// add id to contextToDelete
	// on safeWriteConfig() this wil be interpreted and deleted
	// convert ID to uint32
	idInt, err := strconv.Atoi(ID)
	id := uint32(idInt)
	if err != nil {
		return err
	}

	// remove from config
	path := getConfigPath(id)
	viper.Set(path, "nil")	

	// set activeContext to nil if it's the one to be deleted
	if id == activeContext{
		viper.Set(fmt.Sprintf("%s.active", ntxConfigRoot), "unset!")
	}

	err = safeWriteConfig()
	if err != nil{
		return err
	}

	return nil
}

// InitContext sets activeContext and defaults config options
func InitContext () {
	// initialize viper config
	err := readConfig()
	if err != nil{
		fmt.Println(err)
		panic(err)
	}

	// set active context
	activeContext = viper.GetUint32(fmt.Sprintf("%s.active", ntxConfigRoot))
	configPath := getConfigPath(activeContext)

	// only overwrite values that are not already set by env or as flags
	if ! viper.IsSet("api-url"){
		viper.Set("api-url", viper.GetString(configPath + ".url"))
	}

	if ! viper.IsSet("username"){
		viper.Set("username", viper.GetString(configPath + ".user"))
	}
	
	if ! viper.IsSet("insecure"){
		viper.Set("insecure", viper.GetBool(configPath + ".insecure"))
	}
}