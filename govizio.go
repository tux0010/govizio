package govizio

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var (
	VolumeDown  = Key{5, 0, "KEYPRESS"}
	VolumeUp    = Key{5, 1, "KEYPRESS"}
	MuteOff     = Key{5, 2, "KEYPRESS"}
	MuteOn      = Key{5, 3, "KEYPRESS"}
	MuteToggle  = Key{5, 4, "KEYPRESS"}
	CycleInput  = Key{7, 1, "KEYPRESS"}
	ChannelDown = Key{8, 0, "KEYPRESS"}
	ChannelUp   = Key{8, 1, "KEYPRESS"}
	PreviousCh  = Key{8, 2, "KEYPRESS"}
	PowerOff    = Key{11, 0, "KEYPRESS"}
	PowerOn     = Key{11, 1, "KEYPRESS"}
	PowerToggle = Key{11, 2, "KEYPRESS"}
)

func (s *SmartCast) StartPairing() (*StartPairingResp, error) {
	data := struct {
		DeviceName string `json:"DEVICE_NAME"`
		DeviceID   string `json:"DEVICE_ID"`
	}{s.Name, s.ID}

	js, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	resp, err := s.apiCall(http.MethodPut, "pairing/start", bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}

	return &StartPairingResp{
		resp.Item.PairingReqToken,
		resp.Item.ChallengeType,
	}, nil
}

func (s *SmartCast) SubmitChallenge(pr *StartPairingResp, respValue string) (*PairResp, error) {
	data := struct {
		DeviceID      string `json:"DEVICE_ID"`
		ChallengeType int    `json:"CHALLENGE_TYPE"`
		ResponseValue string `json:"RESPONSE_VALUE"`
		ReqToken      int    `json:"PAIRING_REQ_TOKEN"`
	}{
		DeviceID:      s.ID,
		ChallengeType: pr.ChallengeType,
		ResponseValue: respValue,
		ReqToken:      pr.PairingReqToken,
	}

	js, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	resp, err := s.apiCall(http.MethodPut, "pairing/pair", bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}

	return &PairResp{resp.Item.AuthToken}, nil
}

func (s *SmartCast) CancelPairing() error {
	data := struct {
		DeviceName string `json:"DEVICE_NAME"`
		DeviceID   string `json:"DEVICE_ID"`
	}{s.Name, s.ID}

	js, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	_, err = s.apiCall("PUT", "pairing/cancel", bytes.NewBuffer(js))
	return err
}

func (s *SmartCast) KeyCommand(k Key) error {
	data := struct {
		Keylist []Key `json:"KEYLIST"`
	}{
		[]Key{k},
	}

	js, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	_, err = s.apiCall(http.MethodPut, "key_command/", bytes.NewBuffer(js))
	return err
}
