package interfaces

// Report is the Incoming report structure, marshallable to json
type Report struct {
	Screenshot string `json:"screenshot,omitempty"`
	ReportInfo
}

// ReportInfo encapsulates all data of a report except the screenshot
type ReportInfo struct {
	Sender       User   `json:"sender"`
	Recipient    User   `json:"recipient"`
	Type         string `json:"type"` // “dm”, “group”, “channel”
	PartyName    string `json:"partyName,omitempty"`
	PartyBlob    string `json:"partyBlob,omitempty"` // include the id for the group, the entire channel data will be included when channels are available
	PartyMembers []User `json:"partyMembers,omitempty"`
}

// User encapsulates the data on a user in a report
type User struct {
	UserId        string `json:"userId"`                  //base64
	ChannelPubkey string `json:"channelPubkey,omitempty"` //optional
	Username      string `json:"username"`
	Nickname      string `json:"nickname,omitempty"` //optional
}
