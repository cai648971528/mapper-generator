# Mapper Generator
Golang mapper-generators used to implement [mappers](https://github.com/kubeedge/mappers-go).

# How to create your own mappers

## 1. Design the device model and device instance CRDs
If you don't know how to define CRD, you can get more details in this [page](https://kubeedge.io/docs/developer/device_crd/).

> <font color= red>**Warning**</font>  
> After introducing DMI, the following points should be noted when defining CRD:
> 1. `instance.yaml`: `spec.protocol` must contain `opcua/modbus/bluetooth/customizedProtocol`.   
     > If you don't want to introduce redundant information, you can define like this
>```yaml
> protocol:
>   customizedProtocol:
>     protocolName: foo
>```
> 2. `model.yaml`:  `spec.protocol` is a required field. For example
> ```yaml
>apiVersion: devices.kubeedge.io/v1alpha2
>kind: DeviceModel
>metadata:
>  name: foo-model
>  namespace: default
>spec:
>  protocol: foo
>```

## 2. Generate the code template
The mapper template is to generate a framework for the customized mapper. Run the command and input your mapper's name:
```shell
make template
Please input the mapper name (like 'Bluetooth', 'BLE'): foo
```
A floder named as you input will be generated. The file tree is as below:
```
mapper
├── cmd  ----------------------- Main process.
│   └── main.go ---------------- Almost need not change.
├── config.yaml ---------------- Configuration file including DMI's grpc settting
├── device --------------------- Implementation layer, almost need not change
│   ├── device.go -------------- Device control layer, almost need not change
│   └── twindata.go ------------ Report twin data to edgecore's DMI server, almost need not change
├── Dockerfile
├── driver --------------------- Device driver layer, complete TODO item in this 
│   ├── devicetype.go ---------- Refine the struct as your CRD.
│   └── driver.go -------------- Fill in the functions like getting data/setting register.
├── hack
│   └── make-rules
│       └── mapper.sh
├── Makefile
└── README.md

```

# Where does it come from?
mapper-generator is synced from https://github.com/kubeedge/kubeedge/tree/master/staging/src/github.com/kubeedge/mapper-generator. Code changes are made in that location, merged into kubeedge and later synced here.