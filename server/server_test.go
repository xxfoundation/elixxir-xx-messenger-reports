////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package server

import "testing"

func TestStartServer(t *testing.T) {
	err := StartServer("", "<p>%s</p>", Params{
		KeyPath:  "",
		CertPath: "",
		Port:     "5432",
	})
	if err != nil {
		t.Fatalf("Failed to start server: %+v", err)
	}
	select {}
}
