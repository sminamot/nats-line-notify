package line

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
)

type Line struct {
	Message    string `json:"message"`
	ImageURL   string `json:"image_url"`
	RetryCount int    `json:"retry_count"`
}

func (s *Line) Notify(token string) error {
	buf := bytes.Buffer{}
	mw := multipart.NewWriter(&buf)
	mw.WriteField("message", s.Message)
	if s.ImageURL != "" {
		mw.WriteField("imageThumbnail", s.ImageURL)
		mw.WriteField("imageFullsize", s.ImageURL)
	}
	mw.Close()

	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", &buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("line notify failed: %d", resp.StatusCode)
	}
	return resp.Body.Close()
}
