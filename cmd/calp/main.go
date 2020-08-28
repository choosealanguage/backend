package main

import (
	"log"
	"strings"

	"github.com/choosealanguage/backend/internal/webserver"
	"github.com/choosealanguage/backend/pkg/filewatcher"
	"github.com/choosealanguage/backend/pkg/provider"
	"github.com/spf13/viper"
	"gopkg.in/fsnotify.v1"
)

// viperSetup initializes viper and specified
// config values from defined configuration
// providers.
func viperSetup() error {
	viper.SetDefault("providers", []string{})
	viper.SetDefault("webserver.debug", false)
	viper.SetDefault("webserver.cors.origin", "https://127.0.0.1:3000")
	viper.SetDefault("webserver.address", "0.0.0.0:8080")

	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.calp")
	viper.AddConfigPath("/etc/calp")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.SetEnvPrefix("CALP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return nil
}

func main() {
	// setup config via viper
	if err := viperSetup(); err != nil {
		log.Fatal("failed parsing config: ", err)
	}

	// setup provider host
	prov := provider.New()

	// setup file watcher
	watcher, err := filewatcher.New(true)
	if err != nil {
		log.Fatal("failed creating file watcher: ", err)
	}

	// setup provider watch paths
	providers := viper.GetStringSlice("providers")
	for _, provider := range providers {
		if err = watcher.AddPath(provider); err != nil {
			log.Fatal("failed adding file watcher path: ", err)
		}
	}

	// register file watch handlers
	watcher.Handle(fsnotify.Create|fsnotify.Write, func(e fsnotify.Event) {
		if err = prov.UpdateFromFile(e.Name); err != nil {
			log.Print("error : failed updating providers: ", err)
		}
		log.Print("providers updated: ", e.Name)
	})

	// register file watcher error handler
	watcher.HandleError(func(err error) {
		log.Print("error : file watcher: ", err)
	})

	// start file watcher event loop
	watcher.Run()
	log.Println("file watcher started")
	defer watcher.Close()

	// initialize and run web server
	ws := webserver.New(prov, webserver.Config{
		Debug:      viper.GetBool("webserver.debug"),
		Address:    viper.GetString("webserver.address"),
		CorsOrigin: viper.GetString("webserver.cors.origin"),
	})
	if err = ws.Run(); err != nil {
		log.Fatal("failed starting web server: ", err)
	}
	log.Println("web server started")

	// block until file watcher loop closes or
	// prcess exists
	<-watcher.Done()
}
