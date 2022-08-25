///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// The server package contains the code to run an http server, accepting post
// requests containing report data from an xx messenger user.

package server

import (
	"encoding/json"
	"fmt"
	"git.xx.network/elixxir/xx-messenger-reports/interfaces"
	"git.xx.network/elixxir/xx-messenger-reports/report"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	jww "github.com/spf13/jwalterweatherman"
	"io"
	"io/ioutil"
	"time"
)

// Params for report server
type Params struct {
	KeyPath  string
	CertPath string
	Port     string
}

type Server struct {
	comms *gin.Engine
}

// StartServer starts an http server for processing reports
func StartServer(key, template string, params Params) error {
	jww.INFO.Printf("Starting server on port %s", params.Port)

	h := report.NewHandler(key, template)

	impl := Server{}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Origin", "access-control-allow-origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.POST("/report", func(c *gin.Context) {
		err := handleReport(h, c.Request.Body)
		if err != nil {
			jww.ERROR.Printf("Failed to handle report: %+v", err)
			c.String(500, "err: %+v", err)
		} else {
			jww.DEBUG.Printf("Report accepted")
			c.String(200, "report accepted")
		}
	})

	impl.comms = r
	// Run with TLS
	if params.KeyPath == "" && params.CertPath == "" {
		//jww.WARN.Println("NO TLS CONFIGURED")
		return r.Run(fmt.Sprintf("0.0.0.0:%s", params.Port))
	} else {
		return r.RunTLS(fmt.Sprintf("0.0.0.0:%s", params.Port),
			params.CertPath, params.KeyPath)
	}
}

// handleReport helper function is called by the post endpoint for reports.
// It accepts a report.Handler interface and an io.ReadCloser body
func handleReport(handler report.Handler, body io.ReadCloser) error {
	requestData, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	msg := &interfaces.Report{}
	err = json.Unmarshal(requestData, &msg)
	if err != nil {
		return err
	}
	jww.INFO.Printf("Received report from %s", msg.Sender.Username)
	err = handler.HandleReport(msg)
	if err != nil {
		return err
	}
	return nil
}
