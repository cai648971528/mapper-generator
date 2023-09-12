package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kubeedge/mapper-generator/pkg/common"
	_ "github.com/kubeedge/mapper-generator/pkg/global"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

type DataBaseConfig struct {
	Config *ConfigData
}
type ConfigData struct {
	Addr         string `json:"host,omitempty"`
	Password     string `json:"password,omitempty"`
	DB           int    `json:"DB,omitempty"`
	PoolSize     int    `json:"poolSize,omitempty"`
	MinIdleConns int    `json:"minIdleConns,omitempty"`
}

func NewDataBaseClient(config json.RawMessage) (*DataBaseConfig, error) {
	configdata := new(ConfigData)
	err := json.Unmarshal(config, configdata)
	if err != nil {
		return nil, err
	}
	return &DataBaseConfig{Config: configdata}, nil
}
func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("data/dbprovider/config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config init success")
}

func (d *DataBaseConfig) InitDbClient() (*redis.Client, error) {
	////load config file
	//initConfig()
	////init redis client
	//d.Addr = viper.GetString("redis.addr")
	//d.Password = viper.GetString("redis.password")
	//d.DB = viper.GetInt("redis.DB")
	//d.PoolSize = viper.GetInt("redis.poolSize")
	//d.MinIdleConns = viper.GetInt("redis.minIdleConn")
	redisClient := redis.NewClient(&redis.Options{
		Addr:         d.Config.Addr,
		Password:     d.Config.Password,
		DB:           d.Config.DB,
		PoolSize:     d.Config.PoolSize,
		MinIdleConns: d.Config.MinIdleConns,
	})
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("init redis ...", err)
		return nil, err
	} else {
		fmt.Println("redis inited..", pong)
	}
	return redisClient, nil
	//TODO implement me
	//panic("implement me")
}
func (d *DataBaseConfig) CloseSession(client *redis.Client) {
	//TODO implement me
	err := client.Close()
	if err != nil {
		klog.V(4).Info("close database failed")
	}
	//panic("implement me")
}
func (d *DataBaseConfig) AddData(data *common.DataModel, client redis.Client) error {
	ctx := context.Background()
	// 构建有序集合的键，这里使用 DeviceName 作为键
	klog.Fatal("deviceName:%s", data.DeviceName)
	// 检查是否存在该有序集合
	exists, err := client.Exists(ctx, data.DeviceName).Result()
	if err != nil {
		klog.V(4).Info("")
	}
	// 数据转为 JSON 格式
	dataJSON, err := json.Marshal(data)
	if err != nil {
		klog.V(4).Info("")
	}
	if exists == 0 {
		// 有序集合不存在，创建新的有序集合并添加数据
		_, err = client.ZAdd(ctx, data.DeviceName, &redis.Z{
			Score:  float64(data.TimeStamp),
			Member: string(dataJSON),
		}).Result()
		if err != nil {
			klog.V(4).Infof("failed add device %v:\n", err)
		}
	} else {
		// 有序集合已存在，直接添加数据
		_, err = client.ZAdd(ctx, data.DeviceName, &redis.Z{
			Score:  float64(data.TimeStamp),
			Member: string(dataJSON),
		}).Result()
		if err != nil {
			klog.V(4).Infof("failed add device %v:\n", err)
		}
	}
	return nil
	//TODO implement me
	//panic("implement me")
}

//func (d *DataBaseConfig) GetDataByDeviceName(deviceName string) ([]*common.DataModel, error) {
//	ctx := context.Background()
//	// 构建有序集合的键，这里使用 DeviceName 作为键
//	klog.Infof("deviceName:%s", deviceName)
//	// 使用 ZREVRANGE 命令获取有序集合的所有成员，按分数从高到低排序
//	dataJSON, err := d.Client.ZRevRange(ctx, deviceName, 0, -1).Result()
//	if err != nil {
//		fmt.Printf("failed  query data : %v ", err)
//	}
//	var dataModels []*common.DataModel
//	for _, jsonStr := range dataJSON {
//		var data common.DataModel
//		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
//			klog.V(4).Infof("Error unmarshaling data: %v\n", err)
//			continue
//		}
//		dataModels = append(dataModels, &data)
//	}
//	return dataModels, nil
//	//TODO implement me
//	//panic("implement me")
//}
//func (d *DataBaseConfig) GetPropertyDataByDeviceName(deviceName string, propertyData string) ([]*common.DataModel, error) {
//	ctx := context.Background()
//	keys, err := d.Client.Keys(ctx, deviceName).Result()
//	if err != nil {
//		klog.Infof("query device %s failed", deviceName)
//		return nil, err
//	}
//	var datamodels []*common.DataModel
//	for _, key := range keys {
//		dataJson, err := d.Client.Get(ctx, key).Result()
//		if err != nil {
//			klog.V(4).Info("Error getting device data:", err)
//			return nil, err
//		}
//		var dm common.DataModel
//		if err := json.Unmarshal([]byte(dataJson), &dm); err != nil {
//			klog.Infof("Error unmarshaling data: %v\n", err)
//			continue
//		}
//		if dm.Value == propertyData {
//			datamodels = append(datamodels, &dm)
//		}
//	}
//	return datamodels, nil
//	//TODO implement me
//	//panic("implement me")
//}
//func (d *DataBaseConfig) GetDataByTimeRange(start int64, end int64) ([]*common.DataModel, error) {
//	ctx := context.Background()
//	// 使用 ZRANGEBYSCORE 命令获取在时间范围内的数据的所有成员，按分数从低到高排序,
//	//需要传key，待优化
//	dataJSON, err := d.Client.ZRangeByScore(ctx, "device1", &redis.ZRangeBy{
//		Min: fmt.Sprintf("%d", start),
//		Max: fmt.Sprintf("%d", end),
//	}).Result()
//	if err != nil {
//		klog.Infof("Select device data failed %v\n", err)
//		return nil, err
//	}
//	var dataModels []*common.DataModel
//	for _, jsonStr := range dataJSON {
//		var data common.DataModel
//		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
//			klog.Infof("Error unmarshaling data: %v\n", err)
//			continue
//		}
//		dataModels = append(dataModels, &data)
//	}
//	return dataModels, nil
//	//TODO implement me
//	//panic("implement me")
//}
//func (d *DataBaseConfig) GetKeys() ([]string, error) {
//	keys, err := d.Client.Keys(context.Background(), "*").Result()
//	if err != nil {
//		klog.V(4).Info("select keys failed")
//	}
//	return keys, err
//}
//func (d *DataBaseConfig) DeleteDataByTimeRange(start int64, end int64) ([]*common.DataModel, error) {
//	ctx := context.Background()
//	// 构建有序集合的键，这里假设使用时间戳作为分数
//	// 使用 ZREMRANGEBYSCORE 命令来删除在时间范围内的数据
//	//需要传key，待优化
//	res, err := d.Client.ZRemRangeByScore(ctx, "device1", fmt.Sprintf("(%d", start), fmt.Sprintf("(%d", end)).Result()
//	fmt.Println(res)
//	if err != nil {
//		klog.Infof("delete %d--%d data failed: %v\n", start, end, err)
//		return nil, err
//	}
//	klog.Infof("delete %d--%d data success\n", start, end)
//	return nil, nil
//	//TODO implement me
//	//panic("implement me")
//}
