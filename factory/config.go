// SPDX-FileCopyrightText: 2025 Canonical Ltd
// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0

/*
 * NRF Configuration Factory
 */

package factory

import (
	"os"
	"strconv"

	"github.com/omec-project/nrf/logger"
	"github.com/omec-project/openapi/models"
	utilLogger "github.com/omec-project/util/logger"
)

const (
	NRF_EXPECTED_CONFIG_VERSION = "1.0.0"
	NRF_DEFAULT_IPV4            = "127.0.0.10"
	NRF_DEFAULT_PORT            = "8000"
	NRF_DEFAULT_PORT_INT        = 8000
	NRF_DEFAULT_SCHEME          = "https"
	NRF_NFM_RES_URI_PREFIX      = "/nnrf-nfm/v1"
	NRF_DISC_RES_URI_PREFIX     = "/nnrf-disc/v1"
)

type Config struct {
	Info          *Info              `yaml:"info"`
	Configuration *Configuration     `yaml:"configuration"`
	Logger        *utilLogger.Logger `yaml:"logger"`
	CfgLocation   string
}

type Info struct {
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type Configuration struct {
	Sbi                   *Sbi     `yaml:"sbi,omitempty"`
	MongoDBName           string   `yaml:"MongoDBName"`
	MongoDBUrl            string   `yaml:"MongoDBUrl"`
	WebuiUri              string   `yaml:"webuiUri"`
	ServiceNameList       []string `yaml:"serviceNameList,omitempty"`
	NfKeepAliveTime       int32    `yaml:"nfKeepAliveTime,omitempty"`
	MongoDBStreamEnable   bool     `yaml:"mongoDBStreamEnable"`
	NfProfileExpiryEnable bool     `yaml:"nfProfileExpiryEnable"`
}

type PlmnSupportItem struct {
	PlmnId     models.PlmnId   `yaml:"plmnId"`
	SNssaiList []models.Snssai `yaml:"snssaiList,omitempty"`
}

type Sbi struct {
	Scheme       string `yaml:"scheme"`
	TLS          *TLS   `yaml:"tls"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty"` // IP that is serviced or registered at another NRF.
	BindingIPv4  string `yaml:"bindingIPv4,omitempty"`  // IP used to run the server in the node.
	Port         int    `yaml:"port,omitempty"`
}

type TLS struct {
	PEM string `yaml:"pem,omitempty"`
	Key string `yaml:"key,omitempty"`
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}

func (c *Config) GetSbiScheme() string {
	if c.Configuration != nil && c.Configuration.Sbi != nil && c.Configuration.Sbi.Scheme != "" {
		return c.Configuration.Sbi.Scheme
	}
	return NRF_DEFAULT_SCHEME
}

func (c *Config) GetSbiPort() int {
	if c.Configuration != nil && c.Configuration.Sbi != nil && c.Configuration.Sbi.Port != 0 {
		return c.Configuration.Sbi.Port
	}
	return NRF_DEFAULT_PORT_INT
}

func (c *Config) GetSbiBindingAddr() string {
	var bindAddr string
	if c.Configuration == nil || c.Configuration.Sbi == nil {
		return "0.0.0.0:" + NRF_DEFAULT_PORT
	}
	if c.Configuration.Sbi.BindingIPv4 != "" {
		if bindIPv4 := os.Getenv(c.Configuration.Sbi.BindingIPv4); bindIPv4 != "" {
			logger.CfgLog.Infof("parsing ServerIPv4 [%s] from ENV Variable", bindIPv4)
			bindAddr = bindIPv4 + ":"
		} else {
			bindAddr = c.Configuration.Sbi.BindingIPv4 + ":"
		}
	} else {
		bindAddr = "0.0.0.0:"
	}
	if c.Configuration.Sbi.Port != 0 {
		bindAddr = bindAddr + strconv.Itoa(c.Configuration.Sbi.Port)
	} else {
		bindAddr = bindAddr + NRF_DEFAULT_PORT
	}
	return bindAddr
}

func (c *Config) GetSbiRegisterIP() string {
	if c.Configuration != nil && c.Configuration.Sbi != nil && c.Configuration.Sbi.RegisterIPv4 != "" {
		return c.Configuration.Sbi.RegisterIPv4
	}
	return NRF_DEFAULT_IPV4
}

func (c *Config) GetSbiRegisterAddr() string {
	regAddr := c.GetSbiRegisterIP() + ":"
	if c.Configuration.Sbi.Port != 0 {
		regAddr = regAddr + strconv.Itoa(c.Configuration.Sbi.Port)
	} else {
		regAddr = regAddr + NRF_DEFAULT_PORT
	}
	return regAddr
}

func (c *Config) GetSbiUri() string {
	return c.GetSbiScheme() + "://" + c.GetSbiRegisterAddr()
}
