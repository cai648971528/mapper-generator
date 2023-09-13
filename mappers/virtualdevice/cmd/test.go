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
		Dsn: "root:taosdata@http(10.211.55.3:6041)/device",
	}
	standard := tdengine.DataStandard{
		SuperTable: "device",
		TagLabel:   "shumeipao002",
	}

	config, _ := json.Marshal(configdata)
	stand, _ := json.Marshal(standard)
	dbconfig, _ := tdengine.NewDataBaseClient(config, stand)

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
		DeviceName:   "device02",
		PropertyName: "humidity",
		Value:        "35.0",
		Type:         "string",
		TimeStamp:    time.Now().Unix(),
	}

	fmt.Println(data)

	datatime := time.Unix(data.TimeStamp, 0).Format("2006-01-02 15:04:05")
	fmt.Println(datatime)

	insertSQL := fmt.Sprintf("INSERT INTO %s USING %s TAGS ('%s') VALUES('%v','%s', '%s', '%s', '%s');",
		data.DeviceName, dbconfig.Standard.SuperTable, dbconfig.Standard.TagLabel, datatime, data.DeviceName, data.PropertyName, data.Value, data.Type)

	fmt.Println(insertSQL)
	//_, err = tdengine.DB.Exec(insertSQL)
	//
	////err = dbconfig.AddData(&data, dbclient)
	//if err != nil {
	//	fmt.Printf("add failed :%v", err.Error())
	//} else {
	//	fmt.Println("add success")
	//}

	//dbconfig.AddData(&data)
	//dbmodel, _ := dbconfig.GetDataByDeviceName("device02")
	//for _, datas := range dbmodel {
	//	fmt.Println(datas)
	//}

}
