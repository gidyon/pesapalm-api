package sms

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gidyon/gomicro/utils/errs"
	"github.com/gidyon/pesapalm/pkg/api/sms"
	"github.com/gidyon/pesapalm/pkg/utils/httputils"
	"github.com/rs/zerolog/log"
)

var httpClient = &http.Client{
	Transport: &http.Transport{
		ForceAttemptHTTP2: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
	Timeout: time.Second * 10,
}

func ValidateAuth(pb *sms.SMSAuth) error {
	var err error
	switch {
	case pb == nil:
		err = errs.MissingField("sms auth")
	case pb.ClientId == "":
		err = errs.MissingField("client id")
	// case pb.AccessKey == "":
	// 	err = errs.MissingField("access key")
	case pb.ApiKey == "":
		err = errs.MissingField("api key")
	case pb.SenderId == "":
		err = errs.MissingField("sender id")
	}
	return err
}

func ValidateSms(pb *sms.SMS) error {
	var err error
	switch {
	case pb == nil:
		err = errs.MissingField("sms data")
	case len(pb.DestinationPhones) == 0:
		err = errs.MissingField("destination phones")
	case pb.Message == "":
		err = errs.MissingField("message")
	}
	return err
}

func SendSMS(ctx context.Context, req *sms.SendSMSRequest, env string) error {
	// Validation
	switch {
	case req == nil:
		return errs.MissingField("request")
	default:
		err := ValidateSms(req.Sms)
		if err != nil {
			return err
		}
		err = ValidateAuth(req.Auth)
		if err != nil {
			return err
		}
	}

	go sendSmsOnfon(req, env)

	return nil
}

func firstVal(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

type senderError struct {
	phone string
	err   error
}

func (s *senderError) Error() string {
	return fmt.Sprintf("%s: %v", s.phone, s.err)
}

func sendSmsOnfon(sendRequest *sms.SendSMSRequest, env string) {
	var (
		url     = firstVal(sendRequest.GetAuth().GetApiUrl())
		method  = "POST"
		errChan = make(chan error, len(sendRequest.GetSms().GetDestinationPhones()))
		errors  = make([]error, 0)
		sem     = make(chan struct{}, 5)
		msg     = fmt.Sprintf("%s%s", env, sendRequest.GetSms().GetMessage())
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for _, phone := range sendRequest.GetSms().GetDestinationPhones() {
		go func(phone string) {
			sem <- struct{}{}
			defer func() {
				<-sem
			}()

			payload := strings.NewReader(
				fmt.Sprintf(
					"{\"SenderId\": \"%s\",\"IsUnicode\": true,\"IsFlash\": true,\"MessageParameters\": [{\"Number\": \"%s\",\"Text\": \"%s\"}],\"ApiKey\": \"%s\",\"ClientId\": \"%s\"}",
					sendRequest.GetAuth().GetSenderId(),
					phone,
					msg,
					sendRequest.GetAuth().GetApiKey(),
					sendRequest.GetAuth().GetClientId(),
				),
			)

			req, err := http.NewRequest(method, url, payload)

			if err != nil {
				errChan <- &senderError{phone: phone, err: err}
				return
			}

			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("AccessKey", sendRequest.GetAuth().GetAccessKey())

			httputils.DumpRequest(req, "ONFON SMS GATEWAY REQUEST")

			res, err := httpClient.Do(req)
			if err != nil {
				errChan <- &senderError{phone: phone, err: err}
				return
			}
			defer res.Body.Close()

			httputils.DumpResponse(res, "ONFON SMS GATEWAY RESPONSE")

			resMap := map[string]interface{}{}
			err = json.NewDecoder(res.Body).Decode(&resMap)
			if err != nil {
				errChan <- &senderError{phone: phone, err: err}
				return
			}

			if val, ok := resMap["ErrorCode"]; !ok || (fmt.Sprint(val) != "0") {
				errChan <- &senderError{phone: phone, err: err}
				return
			}

			errChan <- nil
		}(phone)
	}

	for range sendRequest.GetSms().GetDestinationPhones() {
		select {
		case <-ctx.Done():
			return
		case err := <-errChan:
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		for _, err := range errors {
			log.Err(err)
		}
	}
}
