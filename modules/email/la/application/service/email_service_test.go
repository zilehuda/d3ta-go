package service

import (
	"testing"

	appDTO "github.com/muharihar/d3ta-go/modules/email/la/application/dto"
	domSchema "github.com/muharihar/d3ta-go/modules/email/la/domain/schema"

	"github.com/muharihar/d3ta-go/system/handler"
	"github.com/muharihar/d3ta-go/system/initialize"
)

func newEmailService(t *testing.T) (*EmailService, *handler.Handler, error) {
	h, err := handler.NewHandler()
	if err != nil {
		return nil, nil, err
	}

	c, err := newConfig(t)
	if err != nil {
		return nil, nil, err
	}

	h.SetConfig(c)
	if err := initialize.LoadAllDatabase(h); err != nil {
		return nil, nil, err
	}

	s, err := NewEmailService(h)
	if err != nil {
		return nil, nil, err
	}

	return s, h, nil
}

func TestEmailService_Send(t *testing.T) {
	svc, h, err := newEmailService(t)
	if err != nil {
		t.Errorf("Error.newEmailService: %s", err.Error())
		return
	}

	req := new(appDTO.SendEmailReqDTO)
	req.TemplateCode = "test.code"
	req.From = &domSchema.MailAddress{Email: "d3tago.from@domain.tld", Name: "D3TA Golang"}
	req.To = &domSchema.MailAddress{Email: "d3tago.test@outlook.com", Name: "D3TAgo Test (Outlook)"}
	req.CC = []*domSchema.MailAddress{
		{Email: "d3tago.test@protonmail.com", Name: "D3TAgo Test CC 1 (Protonmail)"},
		{Email: "d3tago.test.cc@tutanota.com", Name: "D3TAgo Test CC 2 (Tutanota)"}}
	req.BCC = []*domSchema.MailAddress{
		{Email: "d3tago.test@tutanota.com", Name: "D3TAgo Test BCC 1 (Tutanota)"},
		{Email: "d3tago.test.bcc@outlook.com", Name: "D3TAgo Test BCC 2 (Outlook)"}}
	req.TemplateData = map[string]interface{}{
		"Header.Name":      "Name",
		"Body.URL":         "https://google.com",
		"Footer.Signature": "Customer Service",
	}

	i := newIdentity(h, t)
	i.RequestInfo.RequestObject = "/api/v1/email/send"
	i.RequestInfo.RequestAction = "POST"

	resp, err := svc.Send(req, i)
	if err != nil {
		t.Errorf("Error.EmailService.Send: %s", err.Error())
		return
	}

	if resp != nil {
		respJSON := resp.ToJSON()
		t.Logf("Resp: %s", string(respJSON))
	}
}