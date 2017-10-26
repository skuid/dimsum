package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"strings"

	"github.com/gorilla/mux"
	"github.com/heroku/docker-registry-client/registry"
	yaml "gopkg.in/yaml.v2"
)

func errorAndExit(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}

var (
	config = flag.String("config", "config.yaml", "Configuration file location.")
)

type RegistryConfig struct {
	Name         string `yaml:"name"`
	Address      string `yaml:"address"`
	Username     string `yaml:"username,omitempty"`
	Password     string `yaml:"password,omitempty"`
	PasswordFile string `yaml:"passwordFile,omitempty"`
	Email        string `yaml:"email,omitempty"`
}

func (self *RegistryConfig) readPasswordFile() (string, error) {
	password, err := ioutil.ReadFile(self.PasswordFile)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(password)), nil
}

type AccountList struct {
	Accounts []RegistryConfig `yaml:"accounts"`
}

type Config struct {
	AccountList AccountList `yaml:"dockerRegistry"`
}

func parseConfig(cfg []byte) ([]RegistryConfig, error) {
	parsed := Config{}

	err := yaml.Unmarshal(cfg, &parsed)
	if err != nil {
		return nil, err
	}

	return parsed.AccountList.Accounts, nil
}

var RegistryCache map[string]*registry.Registry = make(map[string]*registry.Registry, 0)
var AccountCache map[string]RegistryConfig = make(map[string]RegistryConfig, 0)

func getRegistryInstance(accountName string) (*registry.Registry, error) {
	if account, ok := AccountCache[accountName]; ok {
		if account.PasswordFile != "" {
			password, err := account.readPasswordFile()

			if err != nil {
				return nil, err
			}

			return registry.New(account.Address, account.Username, password)
		} else {
			return RegistryCache[account.Name], nil
		}
	}

	return nil, fmt.Errorf("Account %s does not exist", accountName)
}

func initRegistryCache(accounts []RegistryConfig) {
	for _, account := range accounts {
		if account.PasswordFile == "" {
			hub, err := registry.New(account.Address, account.Username, account.Password)
			errorAndExit(err)
			RegistryCache[account.Name] = hub
		}
	}
}

func initAccountCache(accounts []RegistryConfig) {
	for _, account := range accounts {
		if account.PasswordFile != "" {
			if _, err := os.Stat(account.PasswordFile); os.IsNotExist(err) {
				fmt.Printf("Password file for account %s does not exist: %s\n", account.Name, account.PasswordFile)
			}
		}
		AccountCache[account.Name] = account
	}
}

func main() {

	flag.Parse()

	cfg, err := ioutil.ReadFile(*config)

	errorAndExit(err)

	accounts, err := parseConfig(cfg)
	errorAndExit(err)

	initAccountCache(accounts)
	initRegistryCache(accounts)

	r := mux.NewRouter()

	//get the manifest for an image
	r.HandleFunc("/{account}/{repository:.+\\/.+}/{tag}/metadata", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		hub, err := getRegistryInstance(vars["account"])

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		manifest, err := hub.Manifest(vars["repository"], vars["tag"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(manifest)

	})

	//get History.V1Compatibility information for an image
	r.HandleFunc("/{account}/{repository:.+\\/.+}/{tag}/history", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hub, err := getRegistryInstance(vars["account"])

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		manifest, err := hub.Manifest(vars["repository"], vars["tag"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if param := r.URL.Query().Get("level"); param != "" {
			index, _ := strconv.Atoi(param)                          //TODO: actually handle this error
			w.Write([]byte(manifest.History[index].V1Compatibility)) //TODO: handle invalid indexes
		} else {
			json.NewEncoder(w).Encode(manifest.History)
		}

	})

	if err := http.ListenAndServe(":8080", r); err != http.ErrServerClosed {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
