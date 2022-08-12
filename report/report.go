package report

import (
	"fmt"
	"git.xx.network/elixxir/xx-messenger-reports/interfaces"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type handler struct {
	sendgridKey    string
	emailTemplate  string
	recipientEmail string
}

// Handler interface is the public interface of the handler package, for pushing a report to email service
type Handler interface {
	HandleReport(msg *interfaces.Report) error
}

func NewHandler(key, template, recipientEmail string) Handler {
	return &handler{
		sendgridKey:    key,
		emailTemplate:  template,
		recipientEmail: recipientEmail,
	}
}

// HandleReport accepts an interfaces.Report object and sends it as an email to the reports@xxmessenger.io account
func (h *handler) HandleReport(msg *interfaces.Report) error {
	senderAddress := mail.NewEmail("xxmessenger reports", "reports@xxmessenger.io")
	recipientAddress := mail.NewEmail("xxmessenger reports group", h.recipientEmail)
	contents := fmt.Sprintf("<h2>Incident Report</h2><table><tr><th></th><th>UserID</th><th>Username</th><th>Nickname</th><th>PubKey</th></tr><tr><th>Sender</th><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr><tr><th>Receiver</th><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr></table>", msg.Sender.UserId, msg.Sender.Username, msg.Sender.Nickname, msg.Sender.ChannelPubkey, msg.Recipient.UserId, msg.Recipient.Username, msg.Recipient.Nickname, msg.Recipient.ChannelPubkey)
	if msg.PartyName != "" || msg.PartyBlob != "" {
		groupContents := fmt.Sprintf("<h3>Party Info</h3><table><tr><th>Party Name</th><td>%s</td></tr><tr><th>Party Blob</th><td>%s</td></tr></table>", msg.PartyName, msg.PartyBlob)
		groupMembers := "<table><tr><th colspan=5>Group Members</th></tr><tr><th>#</th><th>UserID</th><th>Username</th><th>Nickname</th><th>PubKey</th></tr>"
		for i := 0; i < len(msg.PartyMembers); i++ {
			m := msg.PartyMembers[i]
			groupMembers += fmt.Sprintf("<tr><th>%d</th><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", i, m.UserId, m.Username, m.Nickname, m.ChannelPubkey)
		}
		groupMembers += "</table>"
		contents = contents + "<br>" + groupContents + "<br>" + groupMembers
	}
	htmlContents := fmt.Sprintf(h.emailTemplate, contents)

	emailMessage := mail.NewSingleEmail(senderAddress, fmt.Sprintf("xx messenger %s report", msg.Type), recipientAddress, string(contents), htmlContents)
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
	} else if 200 > response.StatusCode || response.StatusCode > 299 {
		return errors.Errorf("Received error %d when sending report email: %+v", response.StatusCode, response.Body)
	}
	return nil
}
