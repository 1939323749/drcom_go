package service

import (
	"io"
	"net/http"
)

func (s *Service) checkConnectivity() (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://4.ipw.cn", nil)
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
			s.logger.Error(err.Error())
		}
	}(resp.Body)

	return resp.StatusCode == http.StatusOK, nil
}
