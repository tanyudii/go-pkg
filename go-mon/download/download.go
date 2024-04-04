package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	gologger "pkg.tanyudii.me/go-pkg/go-logger"
)

var (
	ErrDownloadFailed = fmt.Errorf("[ERROR]: Download failed")
)

func Download(url, method string) ([]byte, error) {
	return WithCtx(context.Background(), url, method)
}

func WithCtx(ctx context.Context, url, method string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			gologger.Errorf("error closing response body: %v", err)
		}
	}()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return io.ReadAll(resp.Body)
	}
	return nil, fmt.Errorf("%w: returning non 2xx http code %v", ErrDownloadFailed, resp.StatusCode)
}
