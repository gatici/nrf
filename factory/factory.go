// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0

/*
 * NRF Configuration Factory
 */

package factory

import (
	"fmt"
	"os"
	"time"

	protos "github.com/omec-project/config5g/proto/sdcoreConfig"
	"google.golang.org/grpc/connectivity"
	"gopkg.in/yaml.v2"

	"github.com/omec-project/nrf/logger"
	"github.com/sirupsen/logrus"
)

var ManagedByConfigPod bool

var NrfConfig Config

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) error {
	if content, err := os.ReadFile(f); err != nil {
		return err
	} else {
		NrfConfig = Config{}

		if yamlErr := yaml.Unmarshal(content, &NrfConfig); yamlErr != nil {
			return yamlErr
		}
		if NrfConfig.Configuration.WebuiUri == "" {
			NrfConfig.Configuration.WebuiUri = "webui:9876"
		}
		initLog.Infof("DefaultPlmnId Mnc %v , Mcc %v \n", NrfConfig.Configuration.DefaultPlmnId.Mnc, NrfConfig.Configuration.DefaultPlmnId.Mcc)
		roc := os.Getenv("MANAGED_BY_CONFIG_POD")
		if roc == "true" {
			initLog.Infoln("MANAGED_BY_CONFIG_POD is true")
			var client ConfClient
			client = ConnectToConfigServer(NrfConfig.Configuration.WebuiUri)
			for {
				if client != nil {
					initLog.Infoln("GRPC client existed.")
					UpdateConfig(client)
					time.Sleep(time.Second * 15)
					if client.GetConfigClientConn().GetState() != connectivity.Ready {
						client.GetConfigClientConn().Close()
						client = nil
						continue
					} else {
						client = ConnectToConfigServer(NrfConfig.Configuration.WebuiUri)
						continue
					}
				}
			}
		}
	}
	return nil
}

func UpdateConfig(client ConfClient) {
	var stream protos.ConfigService_NetworkSliceSubscribeClient
	for {
		stream = client.ConnectToGrpcServer()
		if stream == nil {
			time.Sleep(time.Second * 10)
			continue
		}
		break
	}
	configChannel := client.PublishOnConfigChange(true, stream)
	ManagedByConfigPod = true
	go NrfConfig.updateConfig(configChannel)

}

func CheckConfigVersion() error {
	currentVersion := NrfConfig.GetVersion()

	if currentVersion != NRF_EXPECTED_CONFIG_VERSION {
		return fmt.Errorf("config version is [%s], but expected is [%s]",
			currentVersion, NRF_EXPECTED_CONFIG_VERSION)
	}

	logger.CfgLog.Infof("config version [%s]", currentVersion)

	return nil
}
