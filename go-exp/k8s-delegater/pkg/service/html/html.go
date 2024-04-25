package html

import (
	"bytes"
	"context"
	"text/template"

	"k8s-delegater/pkg/common"
	"k8s-delegater/pkg/dao/keyvault"
)

var emailURL string = ""

func GetEmailURL() (string, error) {
	if emailURL != "" {
		return emailURL, nil
	}

	kv, err := keyvault.NewKeyVaultClient()
	if err != nil {
		common.GetLog().Errorf("keyvault NewKeyVaultClient error: %v", err)
		return "", err
	}

	url, err := kv.GetSecret(context.Background(), "EmailURL")
	if err != nil {
		common.GetLog().Errorf("keyvault GetSecret error: %v", err)
		return "", err
	}
	emailURL = url
	return emailURL, nil
}

func GetTableHTML(table *HTMLTable) (string, error) {
	t := template.Must(template.New("").Parse(`
	<table style="border-collapse:collapse;word-break:keep-all;white-space:nowrap;font-size:14px;">
	{{with .Title}}<caption style="padding:15px;color:#fff;background-color:#48a6fb;font-size:18px;">{{.}}</caption>{{end}}
	<tr>
	{{range .Header}}
	<th style="border-color:#000;border-width:1px;border-style:solid;padding:15px;">{{.}}</th>
	{{end}}
	</tr>
	{{range .Row}}
	<tr>
	{{range .}}
	<td style="border-color:#000;border-width:1px;border-style:solid;padding:5px;">{{.}}</td>
	{{end}}
	</tr>
	{{end}}
	</table>`))

	var body bytes.Buffer
	if err := t.Execute(&body, table); err != nil {
		common.GetLog().Errorf("template Execute error: %v", err)
		return "", nil
	}
	return body.String(), nil
}
