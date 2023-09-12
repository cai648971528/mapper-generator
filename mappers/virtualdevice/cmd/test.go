package main

import (
	"encoding/json"
	"fmt"
	"github.com/kubeedge/mapper-generator/mappers/virtualdevice/data/dbprovider/tdengine"
)

func main() {

	configdata := tdengine.ConfigData{
		Dsn: "root:taosdata@tcp(10.211.55.3:6030)/",
	}
	standard := tdengine.DataStandard{
		SuperTable: "device",
		TagLabel:   "local",
		TagGroupId: "1",
	}

	config, _ := json.Marshal(configdata)
	stand, _ := json.Marshal(standard)
	dbconfig, _ := tdengine.NewDataBaseClient(config, stand)
	dbclient, _ := dbconfig.InitDbClient()
	err := dbclient.Ping()
	if err != nil {
		fmt.Println("failed conn tdengine")
	}

}
