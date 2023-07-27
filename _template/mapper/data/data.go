package data

import (
	"fmt"

	"k8s.io/klog/v2"

	"github.com/kubeedge/mapper-generator/mappers/Template/driver"
	"github.com/kubeedge/mapper-generator/pkg/common"
)

type PropertyData struct {
	InstanceID    string
	DeviceName    string
	Client        *driver.CustomizedClient
	Name          string
	Type          string
	VisitorConfig *driver.VisitorConfig
}

type PropertyDataModule interface {
	GetPropertyData() PropertyData
	SetPropertyData(PropertyData)

	PushData(string)
}

type PanelCtl struct {
	Property    PropertyData
	dataChan    chan string
	dataModules []PropertyDataModule
}

var (
	panelCtl *PanelCtl
)

func NewPanelCtl(propertyData PropertyData) *PanelCtl {
	panelCtl = &PanelCtl{
		Property: propertyData,
		dataChan: make(chan string),
	}

	panelCtl.RegisterModule(NewTwinData(&propertyData))
	return panelCtl
}

func (ctl *PanelCtl) RegisterModule(dataModule PropertyDataModule) {
	ctl.dataModules = append(ctl.dataModules, dataModule)
}

func (ctl *PanelCtl) InitPanel() {
	for _, module := range ctl.dataModules {
		module.SetPropertyData(ctl.Property)
	}
}

func (ctl *PanelCtl) Start() {
	sData, err := ctl.Property.GetData()
	if err != nil {
		klog.Error(err)
		return
	}
	go ctl.startModule()
	ctl.dataChan <- sData
}

func (ctl *PanelCtl) startModule() {
	for data := range ctl.dataChan {
		for _, module := range ctl.dataModules {
			go module.PushData(data)
		}
	}
}

func (pd *PropertyData) GetData() (string, error) {
	value, err := pd.Client.GetDeviceData(pd.VisitorConfig)
	if err != nil {
		return "", fmt.Errorf("get device data failed: %v", err)
	}
	sData, err := common.ConvertToString(value)
	if err != nil {
		return "", fmt.Errorf("failed to convert %s %s value as string : %v", pd.DeviceName, pd.Name, err)
	}
	if len(sData) > 30 {
		klog.V(4).Infof("Get %s : %s ,value is %s......", pd.DeviceName, pd.Name, sData[:30])
	} else {
		klog.V(4).Infof("Get %s : %s ,value is %s", pd.DeviceName, pd.Name, sData)
	}
	return sData, nil
}
