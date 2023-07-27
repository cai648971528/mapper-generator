package data

import (
	"encoding/json"
	"fmt"
	"strings"

	"k8s.io/klog/v2"

	dmiapi "github.com/kubeedge/kubeedge/pkg/apis/dmi/v1alpha1"
	"github.com/kubeedge/mapper-generator/pkg/common"
	"github.com/kubeedge/mapper-generator/pkg/util/grpcclient"
	"github.com/kubeedge/mapper-generator/pkg/util/parse"
)

type TwinData struct {
	Topic string
	PropertyData
}

func NewTwinData(propertyData *PropertyData) *TwinData {
	return &TwinData{
		Topic:        fmt.Sprintf(common.TopicTwinUpdate, propertyData.InstanceID),
		PropertyData: *propertyData,
	}
}

func (td *TwinData) SetPropertyData(propertyData PropertyData) {
	td.PropertyData = propertyData
}

func (td *TwinData) GetPropertyData() PropertyData {
	return td.PropertyData
}

func (td *TwinData) PushData(data string) {
	payload, err := td.GetPayLoad(data)
	if err != nil {
		klog.Errorf("twindata %s unmarshal failed, err: %s", td.Name, err)
		return
	}

	var msg common.DeviceTwinUpdate
	if err = json.Unmarshal(payload, &msg); err != nil {
		klog.Errorf("twindata %s unmarshal failed, err: %s", td.Name, err)
		return
	}

	twins := parse.ConvMsgTwinToGrpc(msg.Twin)

	var rdsr = &dmiapi.ReportDeviceStatusRequest{
		DeviceName: td.DeviceName,
		ReportedDevice: &dmiapi.DeviceStatus{
			Twins: twins,
			State: "OK",
		},
	}

	if err := grpcclient.ReportDeviceStatus(rdsr); err != nil {
		klog.Errorf("fail to report device status of %s with err: %+v", rdsr.DeviceName, err)
	}
}

func (td *TwinData) GetPayLoad(data string) ([]byte, error) {
	var payload []byte
	var err error
	if strings.Contains(td.Topic, "$hw") {
		if payload, err = common.CreateMessageTwinUpdate(td.Name, td.Type, data); err != nil {
			return nil, fmt.Errorf("create message twin update failed: %v", err)
		}
	} else {
		if payload, err = common.CreateMessageData(td.Name, td.Type, data); err != nil {
			return nil, fmt.Errorf("create message data failed: %v", err)
		}
	}
	return payload, nil
}
