package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

// reference
// https://www.modb.pro/db/419363
// https://golangbyexample.com/
// https://dev.to/koddr/let-s-write-config-for-your-golang-web-app-on-right-way-yaml-5ggp
func main() {

	 // Generate our config based on the config supplied
    // by the user in the flags
    cfgPath, err := ParseFlags()
    if err != nil {
        logrus.Fatal(err)
    }
    cfg, err := NewConfig(cfgPath)
    if err != nil {
        logrus.Fatal(err)
    }

    logLevel := cfg.Server.LogLevel

    //TODO more case
	if logLevel == "INFO" {
		logrus.SetLevel(logrus.InfoLevel)
	} else if logLevel == "ERR" {
		logrus.SetLevel(logrus.ErrorLevel)
	}

	
	// back headers
	http.HandleFunc("/healthz", healthz)


	go func() {
		err := http.ListenAndServe(":19004", nil)

		if errors.Is(err, http.ErrServerClosed) {
			
		} else if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			  }).Warn("error starting server")
			os.Exit(1)
		}
	}()
	

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("我被优雅终止了")

}

func healthz(w http.ResponseWriter, r *http.Request) {
	logrus.Info("info accessed healthz")
    logrus.Error("error accessed healthz")

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "This is my status eq 200 page!\n")
}



type Config struct {
    Server struct {
		LogLevel    string `yaml:"log"`
    } `yaml:"server"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
    // Create config structure
    config := &Config{}

    // Open config file
    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Init new YAML decode
    d := yaml.NewDecoder(file)

    // Start YAML decoding from file
    if err := d.Decode(&config); err != nil {
        return nil, err
    }

    return config, nil
}



// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
    s, err := os.Stat(path)
    if err != nil {
        return err
    }
    if s.IsDir() {
        return fmt.Errorf("'%s' is a directory, not a normal file", path)
    }
    return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
    // String that contains the configured configuration path
    var configPath string

    // Set up a CLI flag called "-config" to allow users
    // to supply the configuration file
    flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

    // Actually parse the flags
    flag.Parse()

    // Validate the path first
    if err := ValidateConfigPath(configPath); err != nil {
        return "", err
    }

    // Return the configuration path
    return configPath, nil
}