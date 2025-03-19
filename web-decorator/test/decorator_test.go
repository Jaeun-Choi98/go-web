package test

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"root/handler"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(handler.NewHandler())
	defer ts.Close()

	var buf bytes.Buffer
	log.SetOutput(&buf)

	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	reader := bufio.NewReader(&buf)
	line, err := reader.ReadString('\n')
	assert.NoError(err)
	assert.Contains(line, "[LOGGER] Started")
}
