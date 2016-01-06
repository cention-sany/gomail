package gomail

import (
	"github.com/cention-sany/mime"
	"github.com/cention-sany/mime/quotedprintable"
)

var newQPWriter = quotedprintable.NewWriter

type mimeEncoder struct {
	mime.WordEncoder
}

var (
	bEncoding = mimeEncoder{mime.BEncoding}
	qEncoding = mimeEncoder{mime.QEncoding}
)
