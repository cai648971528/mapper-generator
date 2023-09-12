package tdengine

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/kubeedge/mapper-generator/pkg/common"
	_ "github.com/taosdata/driver-go/v3/taosSql"
	"k8s.io/klog/v2"
)

type DataBaseConfig struct {
	Config   *ConfigData   `json:"config,omitempty"`
	Standard *DataStandard `json:"dataStandard,omitempty"`
}
type ConfigData struct {
	Dsn string `json:"dsn,omitempty"`
}
type DataStandard struct {
	SuperTable string `json:"superTable,omitempty"`
	TagLabel   string `json:"tagLabel,omitempty"`
	TagGroupId string `json:"tagGroupId,omitempty"`
}

func NewDataBaseClient(config json.RawMessage, standard json.RawMessage) (*DataBaseConfig, error) {
	configdata := new(ConfigData)
	datastandard := new(DataStandard)
	err := json.Unmarshal(config, configdata)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(standard, datastandard)
	if err != nil {
		return nil, err
	}
	return &DataBaseConfig{
		Config:   configdata,
		Standard: datastandard,
	}, nil
}
func (d *DataBaseConfig) InitDbClient() (*sql.DB, error) {
	taos, err := sql.Open("taosSql", d.Config.Dsn)
	if err != nil {
		klog.Infof("failed to connect TDengine, err:", err)
	}
	return taos, err
	//TODO implement me
	//panic("implement me")
}
func (d *DataBaseConfig) CloseSession(db *sql.DB) {
	err := db.Close()
	if err != nil {
		klog.Infoln("failded disconnect taosDB")
	}
	//TODO implement me
	//panic("implement me")
}
func (d *DataBaseConfig) AddData(data *common.DataModel, db *sql.DB) error {
	tablename := data.DeviceName
	_, err := db.Exec("INSERT INTO "+tablename+" USING "+d.Standard.SuperTable+" TAGS ("+d.Standard.TagLabel+", "+d.Standard.TagGroupId+") VALUES (?, ?, ?, ?,?);", data.DeviceName, data.PropertyName, data.Value, data.Type, data.TimeStamp)
	if err != nil {
		klog.Infoln("failed add data to tdengine")
	}
	return nil
	//TODO implement me
	//panic("implement me")
}
func (d *DataBaseConfig) GetDataByDeviceName(deviceName string, db *sql.DB) ([]*common.DataModel, error) {
	query := fmt.Sprintf("SELECT DeviceName, PropertyName, Value, Type, TimeStamp FROM %s", deviceName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []*common.DataModel
	for rows.Next() {
		var result *common.DataModel
		err := rows.Scan(&result.DeviceName, &result.PropertyName, &result.Value, &result.Type, &result.TimeStamp)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
	//TODO implement me
	//panic("implement me")
}
func (d *DataBaseConfig) GetPropertyDataByDeviceName(deviceName string, propertyData string) ([]*common.DataModel, error) {
	//TODO implement me
	panic("implement me")
}
func (d *DataBaseConfig) GetDataByTimeRange(start int64, end int64) ([]*common.DataModel, error) {
	//TODO implement me
	panic("implement me")
}
func (d *DataBaseConfig) DeleteDataByTimeRange(start int64, end int64) ([]*common.DataModel, error) {
	//TODO implement me
	panic("implement me")
}
