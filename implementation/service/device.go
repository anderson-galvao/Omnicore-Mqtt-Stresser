package service

import (
	"encoding/json"

	"github.com/RacoWireless/iot-gw-stresser/model"
	"github.com/rs/zerolog/log"
)

const DEVICEPATH = "device/"

func (d *deviceIotPubsub) CreateDevicePublish(dev model.Device) error {

	PubStruct := model.Publish{Operation: model.POST, Entity: "Device", Data: dev, Path: DEVICEPATH + dev.Parent}

	msg, err := json.Marshal(PubStruct)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	err = Publish(dev.Tenant, d.pubTopic, msg)

	return err
}
func (d *deviceIotPubsub) UpdateDevicePublish(dev model.Device, updateMask string) error {

	PubStruct := model.Publish{Operation: model.PATCH, Entity: "Device", Data: dev, Path: DEVICEPATH + dev.Parent + "?updateMask=" + updateMask}
	msg, err := json.Marshal(PubStruct)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	err = Publish(dev.Tenant, d.pubTopic, msg)
	return err
}
func (d *deviceIotPubsub) DeleteDevicePublish(dev model.Device) error {

	PubStruct := model.Publish{Operation: model.DELETE, Entity: "Device", Data: dev, Path: DEVICEPATH + dev.Parent}

	msg, err := json.Marshal(PubStruct)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	err = Publish(dev.Tenant, d.pubTopic, msg)

	return err
}
func (d *deviceIotPubsub) AddDevCertificatePublish(dev model.DeviceCert) error {

	PubStruct := model.Publish{Operation: model.POST, Entity: "Device", Data: dev, Path: DEVICEPATH + dev.Parent}
	msg, err := json.Marshal(PubStruct)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	err = Publish(dev.Tenant, d.pubTopic, msg)
	return err
}
func (d *deviceIotPubsub) DeleteDevCertificatePublish(dev model.DeviceCert) error {

	PubStruct := model.Publish{Operation: model.DELETE, Entity: "Device", Data: dev, Path: DEVICEPATH + dev.Parent}
	msg, err := json.Marshal(PubStruct)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	err = Publish(dev.Tenant, d.pubTopic, msg)
	return err
}
