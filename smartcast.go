package govizio

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type SmartCast struct {
	ID        string
	Name      string
	AuthToken string
	IP        string
	client    *http.Client
}

type APIRespItem struct {
	AuthToken       string `json:"AUTH_TOKEN,omitempty"`
	PairingReqToken int    `json:"PAIRING_REQ_TOKEN,omitempty"`
	ChallengeType   int    `json:"CHALLENGE_TYPE,omitempty"`
}

type APIResp struct {
	Status struct {
		Result string `json:"RESULT"`
		Detail string `json:"DETAIL"`
	} `json:"STATUS"`
	Item APIRespItem `json:"item",omitempty`
}

type StartPairingResp struct {
	PairingReqToken int
	ChallengeType   int
}

type PairResp struct {
	AuthToken string
}

type Key struct {
	Codeset int    `json:"CODESET"`
	Code    int    `json:"CODE"`
	Action  string `json:"ACTION"`
}

func NewSmartCast(ip, id, name string) *SmartCast {
	// Ignore SSL cerificate for Vizio
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &SmartCast{
		IP:     ip,
		ID:     id,
		Name:   name,
		client: &http.Client{Transport: tr},
	}
}

func (s *SmartCast) SetAuthToken(token string) {
	s.AuthToken = token
}

func (s *SmartCast) apiCall(method, endpoint string, data io.Reader) (*APIResp, error) {
	uri := fmt.Sprintf("https://%s:9000/%s", s.IP, endpoint)

	req, err := http.NewRequest(method, uri, data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if s.AuthToken != "" {
		req.Header.Set("Auth", s.AuthToken)
	}

	log.WithFields(log.Fields{
		"header":   req.Header,
		"method":   method,
		"endpoint": endpoint,
		"body":     data,
	}).Println("Calling API")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	txt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{"data": string(txt)}).Println("API response")

	var r APIResp
	err = json.Unmarshal(txt, &r)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to decode JSON: %s", string(txt))
	}

	if r.Status.Result != "SUCCESS" {
		return nil, errors.New(r.Status.Detail)
	}

	return &r, nil
}
