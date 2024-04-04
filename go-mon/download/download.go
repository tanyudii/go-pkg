package download

import (
	"context"
	"io"
	"net/http"
	gologger "pkg.tanyudii.me/go-pkg/go-logger"
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
	return io.ReadAll(resp.Body)
}
