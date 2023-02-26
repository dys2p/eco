package email

import (
	"context"
	"os/exec"
	"time"
)

// Sendmail runs /usr/bin/sendmail to queue an email for sending.
//
// Your sendmail binary probably has the setgid bit set to be able to copy the mail into the queue directory.
// If your binary runs as a sandboxed systemd service, you might have to specify one or several of these lines:
//
//	SupplementaryGroup=postdrop
//	ReadWritePaths=/var/spool/nullmailer
//	# remove anything that sets NoNewPrivileges=true, such as:
//	# PrivateDevices=true
//	# ProtectKernelTunables=true
type Sendmail struct {
	From string
}

func (mailer Sendmail) Send(to string, subject string, body []byte) error {
	mail, err := MakeEmail(mailer.From, to, subject, body)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	envelopeFrom := mailer.From
	envelopeTo := to

	sendmail := exec.CommandContext(ctx, "/usr/sbin/sendmail", "-i", "-f", envelopeFrom, "--", envelopeTo) // -i don't treat a line with only a . character as the end of input
	sendmail.Stdin = mail
	return sendmail.Run()
}
