////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// The report package handles sending received data to an email group via
// sendgrid's api.

package report

import (
	"fmt"
	"git.xx.network/elixxir/xx-messenger-reports/interfaces"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type handler struct {
	sendgridKey   string
	emailTemplate string
}

// Handler is the public interface of the report package
type Handler interface {
	// HandleReport accepts a report and pushes it to an external service
	// Currently it sends the report to an email group for review
	HandleReport(msg *interfaces.Report) error
}

// NewHandler initializes a report Handler, accepts a key & an html template
func NewHandler(key, template string) Handler {
	return &handler{
		sendgridKey:   key,
		emailTemplate: template,
	}
}

// HandleReport accepts an interfaces.Report object and sends it as an email
// to the reports@elixxir.io account
func (h *handler) HandleReport(msg *interfaces.Report) error {
	senderAddress := mail.NewEmail("xxmessenger reports",
		"reports@xxmessenger.io")
	recipientAddress := mail.NewEmail("xxmessenger reports group",
		"reports@elixxir.io")
	contents := fmt.Sprintf("<h2>Incident Report</h2><table><tr><th>"+
		"</th><th>UserID</th><th>Username</th><th>Nickname</th><th>PubKey"+
		"</th></tr><tr><th>Sender</th><td>%s</td><td>%s</td><td>%s</td><td>"+
		"%s</td></tr><tr><th>Receiver</th><td>%s</td><td>%s</td><td>%s</td>"+
		"<td>%s</td></tr></table>", msg.Sender.UserId, msg.Sender.Username,
		msg.Sender.Nickname, msg.Sender.ChannelPubkey, msg.Recipient.UserId,
		msg.Recipient.Username, msg.Recipient.Nickname,
		msg.Recipient.ChannelPubkey)
	if msg.PartyName != "" || msg.PartyBlob != "" {
		groupContents := fmt.Sprintf("<h3>Party Info</h3><table><tr>"+
			"<th>Party Name</th><td>%s</td></tr><tr><th>Party Blob</th>"+
			"<td>%s</td></tr></table>", msg.PartyName, msg.PartyBlob)
		groupMembers := "<table><tr><th colspan=5>Group Members</th>" +
			"</tr><tr><th>#</th><th>UserID</th><th>Username</th>" +
			"<th>Nickname</th><th>PubKey</th></tr>"
		for i := 0; i < len(msg.PartyMembers); i++ {
			m := msg.PartyMembers[i]
			groupMembers += fmt.Sprintf("<tr><th>%d</th><td>%s</td>"+
				"<td>%s</td><td>%s</td><td>%s</td></tr>", i, m.UserId,
				m.Username, m.Nickname, m.ChannelPubkey)
		}
		groupMembers += "</table>"
		contents = contents + "<br>" + groupContents + "<br>" + groupMembers
	}
	htmlContents := fmt.Sprintf(h.emailTemplate, contents)

	emailMessage := mail.NewSingleEmail(senderAddress,
		fmt.Sprintf("xx messenger %s report", msg.Type),
		recipientAddress, string(contents), htmlContents)
	attachment := mail.NewAttachment()
	attachment.SetContent(msg.Screenshot)
	attachment.SetType("image/png")
	attachment.SetDisposition("attachment")
	attachment.SetFilename("screenshot.png")
	emailMessage.AddAttachment(attachment)

	fmt.Println(string(contents))
	scl := sendgrid.NewSendClient(h.sendgridKey)
	response, err := scl.Send(emailMessage)
	if err != nil {
		return err
	} else if 200 > response.StatusCode ||
		response.StatusCode > 299 {
		return errors.Errorf("Received error %d when sending "+
			"report email: %+v", response.StatusCode, response.Body)
	}
	return nil
}
