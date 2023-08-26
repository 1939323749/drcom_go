package service

import (
	"io"
	"net/http"
)

func (s *Service) checkConnectivity() (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://ip.jlu.edu.cn", nil)
	if err != nil {
		return false, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			close(s.logoutCh)
		}
	}(resp.Body)

	return resp.StatusCode == http.StatusOK, nil
}
