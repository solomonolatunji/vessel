// Package emails embeds the email HTML templates for use by the mailer service.
package emails

import "embed"

//go:embed *.tmpl
var Templates embed.FS
