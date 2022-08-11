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
