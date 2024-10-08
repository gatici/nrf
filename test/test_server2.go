// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/omec-project/nrf/logger"
	"github.com/omec-project/openapi/models"
	"github.com/omec-project/util/http2_util"
	utilLogger "github.com/omec-project/util/logger"
	"github.com/omec-project/util/path_util"
)

var (
	NrfLogPath = path_util.Free5gcPath("github.com/free5gc/nrf/management/sslkeylog.log")
	NrfPemPath = path_util.Free5gcPath("free5gc/support/TLS/nrf.pem")
	NrfKeyPath = path_util.Free5gcPath("free5gc/support/TLS/nrf.key")
)

func main() {
	router := utilLogger.NewGinWithZap(logger.GinLog)

	router.POST("", func(c *gin.Context) {
		/*buf, err := c.GetRawData()
		if err != nil {
			t.Errorf(err.Error())
		}
		// Remove NL line feed, new line character
		//requestBody = string(buf[:len(b uf)-1])*/
		var ND models.NotificationData

		if err := c.ShouldBindJSON(&ND); err != nil {
			logger.UtilLog.Panicln(err.Error())
		}
		logger.UtilLog.Infoln(ND)
		c.JSON(http.StatusNoContent, gin.H{})
	})

	srv, err := http2_util.NewServer(":30678", NrfLogPath, router)
	if err != nil {
		logger.UtilLog.Panicln(err.Error())
	}

	err2 := srv.ListenAndServeTLS(NrfPemPath, NrfKeyPath)
	if err2 != nil && err2 != http.ErrServerClosed {
		logger.UtilLog.Panicln(err2.Error())
	}
}
