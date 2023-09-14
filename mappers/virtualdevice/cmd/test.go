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
		PropertyName: "shidu",
		Value:        "21.0",
		Type:         "string",
		TimeStamp:    time.Now().Unix(),
	}

	fmt.Println(data)

	datatime := time.Unix(data.TimeStamp, 0).Format("2006-01-02 15:04:05")
	fmt.Println(datatime)

	//insertSQL := fmt.Sprintf("INSERT INTO %s USING %s TAGS ('%s') VALUES('%v','%s', '%s', '%s', '%s');",
	//	data.PropertyName, dbconfig.Standard.SuperTable, dbconfig.Standard.TagLabel, datatime, data.DeviceName, data.PropertyName, data.Value, data.Type)
	//
	//fmt.Println(insertSQL)
	//_, err = tdengine.DB.Exec(insertSQL)
	//

	dbconfig.AddData(&data)
	//dbmodel, _ := dbconfig.GetDataByDeviceName("device02")
	//for _, datas := range dbmodel {
	//	fmt.Println(datas)
	//}

	//stabel := fmt.Sprintf("CREATE STABLE %s (ts timestamp, devicename binary(64), propertyname binary(64), data binary(64),type binary(64)) TAGS (%s binary(64));", data.DeviceName, data.PropertyName)
	//fmt.Println(stabel)

}
