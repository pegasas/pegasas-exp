package icm

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/satori/go.uuid"

	"encoding/base64"
	"encoding/xml"

	"k8s-delegater/pkg/common"
	"k8s-delegater/pkg/common/response"
	"k8s-delegater/pkg/dao/keyvault"
)

var createIncidentUrl = "https://prod.microsofticm.com/Connector3/ConnectorIncidentManager.svc?wsdl"
var mitigateIncidentUrl = "https://prod.microsofticm.com/api/cert/incidents({incidentId})/MitigateIncident"

var client *http.Client

// var client = InitClient()

func CreateIncident(incident map[string]interface{}) (int, string) {
	data, _ := ioutil.ReadFile("config/envelope.xml")

	replacer := strings.NewReplacer(
		"{message_id}", uuid.NewV1().String(),
		"{connector_id}", "c7f5a896-934f-43e2-a8f5-88320267bfc3",
		"{incident_content}", dictionaryToXml(incident, 12),
		"{routing_options}", "None",
	)
	request := replacer.Replace(string(data))

	result := post(createIncidentUrl, request, "application/soap+xml; charset=utf-8")

	response := CreateIncidentResponse{}
	xml.Unmarshal([]byte(result), &response)

	incidentId := response.Body.AddOrUpdateIncident2Response.AddOrUpdateIncident2Result.IncidentId
	status := response.Body.AddOrUpdateIncident2Response.AddOrUpdateIncident2Result.Status

	return incidentId, status
}

func MitigateIncident(incidentId int, data string) string {
	url := strings.Replace(mitigateIncidentUrl, "{incidentId}", strconv.Itoa(incidentId), 1)
	result := post(url, data, "application/json")
	return result
}

func post(url string, data string, contentType string) string {
	return invoke("POST", url, data, contentType)
}

func invoke(method, url string, data string, contentType string) string {
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		common.GetLog().Errorf("http NewRequest error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}
	req.Header.Add("Content-Type", contentType)

	res, err := client.Do(req)
	if err != nil {
		common.GetLog().Errorf("http NewRequest error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		common.GetLog().Errorf("http NewRequest error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}
	return string(body)
}

func InitClient() *http.Client {

	kv, err := keyvault.NewKeyVaultClient()
	if err != nil {
		return nil
	}

	certificate, _ := kv.GetSecret(context.Background(), "XPayICMCertificate")
	certificateDec, _ := base64.StdEncoding.DecodeString(certificate)
	ioutil.WriteFile("./CERTIFICATE.pem", certificateDec, 0666)

	privateKey, _ := kv.GetSecret(context.Background(), "XPayICMPrivateKey")
	privateKeyDec, _ := base64.StdEncoding.DecodeString(privateKey)
	ioutil.WriteFile("./PRIVATE_KEY.pem", privateKeyDec, 0666)

	certs, _ := tls.LoadX509KeyPair("./CERTIFICATE.pem", "./PRIVATE_KEY.pem")
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Renegotiation: tls.RenegotiateOnceAsClient,
				Certificates:  []tls.Certificate{certs},
			},
		},
	}
	return client
}

func dictionaryToXml(dictionary map[string]interface{}, indent int) string {
	ns := "b:"
	ws := "\n"

	out_xml := ""
	_indent := ""
	for i := 0; i < indent; i++ {
		_indent += " "
	}

	key_o := strings.Replace("<{ns}{key}>", "{ns}", ns, 1)
	key_c := strings.Replace("</{ns}{key}>", "{ns}", ns, 1) + ws

	var keys []string
	for k := range dictionary {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := dictionary[k]
		switch value := v.(type) {
		case map[string]interface{}:
			out_xml += _indent + strings.Replace(key_o, "{key}", k, 1) + ws
			out_xml += dictionaryToXml(value, indent+3) + ws
			out_xml += _indent + strings.Replace(key_c, "{key}", k, 1)
		case []interface{}:
			out_xml += _indent + strings.Replace(key_o, "{key}", k, 1) + ws
			for _, item := range value {
				out_xml += dictionaryToXml(item.(map[string]interface{}), indent+3) + ws
			}
			out_xml += _indent + strings.Replace(key_c, "{key}", k, 1)
		case bool:
			if value {
				out_xml += _indent + strings.Replace(key_o, "{key}", k, 1) + "true" + strings.Replace(key_c, "{key}", k, 1)
			} else {
				out_xml += _indent + strings.Replace(key_o, "{key}", k, 1) + "true" + strings.Replace(key_c, "{key}", k, 1)
			}
		case string:
			if value == "$now$" {
				_now := time.Now().Format(time.RFC3339)
				out_xml += _indent + strings.Replace(key_o, "{key}", k, 1) + _now + strings.Replace(key_c, "{key}", k, 1)
			} else if value == "$uuid$" {
				_uuid := uuid.NewV1().String()
				out_xml += _indent + strings.Replace(key_o, "{key}", k, 1) + _uuid + strings.Replace(key_c, "{key}", k, 1)
			} else {
				out_xml += _indent + strings.Replace(key_o, "{key}", k, 1) + escape(value) + strings.Replace(key_c, "{key}", k, 1)
			}
		case int:
			out_xml += _indent + strings.Replace(key_o, "{key}", k, 1) + strconv.Itoa(value) + strings.Replace(key_c, "{key}", k, 1)
		case nil:
			out_xml += _indent + "<" + ns + k + " i:nil=\"true\" />" + ws
		default:
			common.GetLog().Errorf("unknown type error: %v", value)
		}
	}
	return out_xml[:len(out_xml)-1]
}

func escape(data string) string {
	data = strings.ReplaceAll(data, "&", "&amp;")
	data = strings.ReplaceAll(data, ">", "&gt;")
	data = strings.ReplaceAll(data, "<", "&lt;")
	return data
}
