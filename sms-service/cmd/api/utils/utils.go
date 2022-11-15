package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

type XMLNode struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:"-"`
	Content []byte     `xml:",innerxml"`
	Nodes   []XMLNode  `xml:",any"`
}

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

var logEvent LogPayload

// readJSON tries to read the body of a request and converts it into JSON
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// writeJSON takes a response status code and arbitrary data and writes a json response to the client
func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// Log the events and send to the logger service
func LogEvent(l LogPayload) {

	jsonData, _ := json.MarshalIndent(l, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("We have trouble in submit log event.")
		log.Println(err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Println("We have trouble in submit log event.")
		log.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		log.Println("We have trouble in submit log event.")
		log.Println(response.StatusCode)
		log.Println(response.Body)
		return
	}
}

// errorJSON takes an error, and optionally a response status code, and generates and sends
// a json error response
func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return WriteJSON(w, statusCode, payload)
}

// trimlastchars return last characters from string by size parameter value
func TrimLastChars(s string, stringSize int) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == stringSize) {
		size = 0
	} else {
		size = stringSize
	}
	return s[len(s)-size:]
}

// sendPostRequest , take  url , data , token parameter and make http post request

func SendPostRequest(url string, data string, token string) (string, int) {
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		log.Println(err)
		return "", 502
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		return string(body), resp.StatusCode
	}
	return string(body), resp.StatusCode
}

// SendPostFormRequest , take  url , data , token parameter and make http post request

func SendPostFormRequest(apiurl string, data map[string]string, token string) (string, int) {
	form := url.Values{}
	for key, value := range data {
		form.Add(key, value)
	}
	request, _ := http.NewRequest("POST", apiurl, strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		log.Println(err)
		return "", 502
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		return string(body), resp.StatusCode
	}
	return string(body), resp.StatusCode
}

// sendGetRequest , take  url , data , token parameter and make http post request

func SendGetRequest(url string, data map[string]string, token string) (string, int) {
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Authorization", "Bearer "+token)
	q := request.URL.Query()
	for key, value := range data {
		q.Add(key, value)
	}
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), resp.StatusCode
}

// send sms single soap request

func SendSMSSoap(url string, method string, username string, password string, to []string, from string, msg string, flash bool, udh string, recid []int64) ([]byte, error) {

	_to := "<string>" + strings.Join(to, "</string><string>") + "</string>"
	_recid := "<long>" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(recid)), "</long><long>"), "[]") + "</long>"

	wsReq := "<?xml version=\"1.0\" encoding=\"utf-8\"?><soap:Envelope xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\"><soap:Body><" + method + " xmlns=\"http://tempuri.org/\"><username>" + username + "</username><password>" + password + "</password><to>" + _to + "</to><from>" + from + "</from><text>" + msg + "</text><isflash>" + strconv.FormatBool(flash) + "</isflash><udh>" + udh + "</udh><recId>" + _recid + "</recId></" + method + "></soap:Body></soap:Envelope>"

	dataAsByte := []byte(wsReq)
	response, err := http.Post(url, "text/xml; charset=utf-8", bytes.NewBuffer(dataAsByte))

	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		log.Printf("The HTTP request failed with error %s\n", err)
		return []byte{0}, err
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		return data, nil
	}
}

// send sms array soap request

func SendSMSArraySoap(url string, method string, username string, password string, to []string, from string, msg []string, flash bool, udh string, recid []int64) ([]byte, error) {

	_to := "<string>" + strings.Join(to, "</string><string>") + "</string>"
	_message := "<string>" + strings.Join(to, "</string><string>") + "</string>"
	_recid := "<long>" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(recid)), "</long><long>"), "[]") + "</long>"

	wsReq := "<?xml version=\"1.0\" encoding=\"utf-8\"?><soap:Envelope xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\"><soap:Body><" + method + " xmlns=\"http://tempuri.org/\"><username>" + username + "</username><password>" + password + "</password><to>" + _to + "</to><from>" + from + "</from><text>" + _message + "</text><isflash>" + strconv.FormatBool(flash) + "</isflash><udh>" + udh + "</udh><recId>" + _recid + "</recId></" + method + "></soap:Body></soap:Envelope>"

	dataAsByte := []byte(wsReq)
	response, err := http.Post(url, "text/xml; charset=utf-8", bytes.NewBuffer(dataAsByte))

	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		log.Printf("The HTTP request failed with error %s\n", err)
		return []byte{0}, err
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		return data, nil
	}
}

// Parse XML
func UnmarshalXML(data []byte, xmlname string) ([]string, error) {

	var xmlContent []string
	buf := bytes.NewBuffer(data)
	dec := xml.NewDecoder(buf)

	var n XMLNode
	err := dec.Decode(&n)
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		log.Println(err)
		return xmlContent, err
	}
	UnmarshalXMLWalk([]XMLNode{n}, func(n XMLNode) bool {
		if n.XMLName.Local == xmlname {
			xmlContent = append(xmlContent, string(n.Content))
		}
		return true
	})

	return xmlContent, nil
}

func UnmarshalXMLWalk(nodes []XMLNode, f func(XMLNode) bool) {
	for _, n := range nodes {
		if f(n) {
			UnmarshalXMLWalk(n.Nodes, f)
		}
	}
}
