package app

import (
	"bytes"
	"net/http"

	usTm "github.com/kyorohiro/k07me/user/template"
)

var userTemplate = usTm.NewUserTemplate(userConfig)

func init() {
	var buffer *bytes.Buffer = bytes.NewBufferString("")
	buffer.WriteString("<html><title>K07ME</title><body>")
	buffer.WriteString("<div>")
	buffer.WriteString("<a href=\"/api/v1/twitter/tokenurl/redirect\">redirect</a>")
	buffer.WriteString("</div>")
	buffer.WriteString("</body></html>")

	userTemplate.InitUserApi()
	http.Handle("/", http.FileServer(http.Dir("web")))

}
