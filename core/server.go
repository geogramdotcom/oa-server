package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/geogramdotcom/oa-server/core/api"
	"github.com/geogramdotcom/oa-server/core/auth"
	"github.com/geogramdotcom/oa-server/core/model"
	"github.com/geogramdotcom/oa-server/core/model/db"
	"github.com/geogramdotcom/oa-server/core/model/types"
	"github.com/geogramdotcom/oa-server/core/util"
)

func main() {
	//filename is the path to the json config file
	var config types.Config
	file, err := os.Open("./config.json")

	if err != nil {
		log.Fatal(fmt.Errorf("failed to open ./config.json with: %s", err.Error()))
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		log.Fatal(fmt.Errorf("failed to decode ./config.json with: %s", err.Error()))
	}

	// connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s",
	// 	config.User,
	// 	config.Password,
	// 	config.DatabaseAddress,
	// 	config.Database)
	fmt.Printf("%s:%s@tcp([%s])/%s", config.User, config.Password, config.DatabaseAddress, config.Database)
	fmt.Println("\n...")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.User, config.Password, config.DatabaseAddress, config.Database)

	// config.User + ":" + config.Password + "@" + config.DatabaseAddress + "/" + config.Database

	db, err := db.NewDB(connectionString)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to database with: %s", err.Error()))
	}

	bc := &util.StandardBcrypt{}

	model.NewModel(db, bc, config)
	auth.NewAuthService(db, bc)

	app, err := api.Init(config.ApiPrefix)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create api instance with: %s", err.Error()))
	}

	if config.CertFile == "" || config.KeyFile == "" {
		err = http.ListenAndServe(config.Address+":"+strconv.Itoa(config.Port), app.MakeHandler())
	} else {
		err = http.ListenAndServeTLS(config.Address+":"+strconv.Itoa(config.Port), config.CertFile, config.KeyFile, app.MakeHandler())
	}
	log.Fatal(fmt.Errorf("failed to start server with: %s", err.Error()))
}
