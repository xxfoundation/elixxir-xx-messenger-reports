package interfaces

import (
	"encoding/base64"
	"encoding/json"
	"gitlab.com/xx_network/primitives/id"
	"math/rand"
	"testing"
)

// Test marshalling & unmarshalling report message struct, with and without optional fields
func TestReportMessage_MarshalUnmarshal(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	uid1 := id.ID{}
	rng.Read(uid1[:])
	uid1.SetType(id.User)

	uid2 := id.ID{}
	rng.Read(uid2[:])
	uid1.SetType(id.User)

	uid3 := id.ID{}
	rng.Read(uid3[:])
	uid1.SetType(id.User)

	msg := Report{
		ReportInfo: ReportInfo{
			Sender: User{
				UserId:   base64.StdEncoding.EncodeToString(uid1.Bytes()),
				Username: "zezima",
			},
			Recipient: User{
				UserId:   base64.StdEncoding.EncodeToString(uid2.Bytes()),
				Username: "jakexx360",
			},
		},
		Screenshot: "",
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))

	loaded := Report{}
	err = json.Unmarshal(jsonData, &loaded)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(loaded)

	msg = Report{
		ReportInfo: ReportInfo{
			Sender: User{
				UserId:   base64.StdEncoding.EncodeToString(uid1.Bytes()),
				Username: "zezima",
			},
			Recipient: User{
				UserId:   base64.StdEncoding.EncodeToString(uid2.Bytes()),
				Username: "jakexx360",
			},
			PartyName: "test group",
			PartyMembers: []User{{
				UserId:   base64.StdEncoding.EncodeToString(uid1.Bytes()),
				Username: "zezima",
			}, {
				UserId:   base64.StdEncoding.EncodeToString(uid2.Bytes()),
				Username: "jakexx360",
			}, {
				UserId:   base64.StdEncoding.EncodeToString(uid3.Bytes()),
				Username: "iceburglettu",
			}},
		},

		Screenshot: "",
	}

	jsonData, err = json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))
	loaded = Report{}
	err = json.Unmarshal(jsonData, &loaded)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(loaded)
}
