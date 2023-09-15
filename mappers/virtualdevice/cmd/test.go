package main

import (
	"encoding/json"
	"fmt"
	"github.com/kubeedge/mapper-generator/mappers/virtualdevice/data/dbprovider/tdengine"
	"github.com/kubeedge/mapper-generator/pkg/common"
	"time"
)

func main() {

	configdata := tdengine.ConfigData{
		Dsn: "root:taosdata@http(10.211.55.3:6041)/test",
	}

	config, _ := json.Marshal(configdata)
	dbconfig, _ := tdengine.NewDataBaseClient(config)

	err := dbconfig.InitDbClient()
	if err != nil {
		fmt.Printf("conn failed %v\n", err)
	}
	defer dbconfig.CloseSessio()
	err = tdengine.DB.Ping()
	if err != nil {
		fmt.Println("failed conn tdengine")
	} else {
		fmt.Println("connect success")
	}

	data := common.DataModel{
		DeviceName:   "device03",
		PropertyName: "wendu",
		Value:        "31.0",
		Type:         "string",
		TimeStamp:    time.Now().Unix(),
	}

	fmt.Println(data)

	//dbconfig.AddData(&data)
	//dbmodel, _ := dbconfig.GetDataByDeviceName("random_instance_01")
	//for _, datas := range dbmodel {
	//	fmt.Println(datas)
	//}

	dbmodel, _ := dbconfig.GetDataByTimeRange("random-instance-01", 1694794006, 1694794018)
	for _, datas := range dbmodel {
		fmt.Println(datas)
	}

}
