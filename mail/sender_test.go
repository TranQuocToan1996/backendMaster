package mail

import (
	"testing"

	"github.com/TranQuocToan1996/backendMaster/util"
	"github.com/stretchr/testify/require"
)

type SpySender struct {
	expectErr error
}

func NewSpySender(expectErr error) EmailSender {
	return &SpySender{expectErr}
}

func (sender *SpySender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {

	return sender.expectErr
}

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	var expectSuccess error = nil
	sender := NewSpySender(expectSuccess)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>Special thanks to <a href="http://techschool.guru">Tech School</a></p>
	`
	to := []string{config.EmailSenderAddress}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, []string{}, []string{}, attachFiles)
	require.NoError(t, err)
}
