package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"turker.web/model"
)

func TestJsonHandler(t *testing.T) {
	assert := assert.New(t)

	var reqBody bytes.Buffer
	user := model.User{FirstName: "choi", LastName: "jaeun", Email: "abc@abc.com"}
	json.NewEncoder(&reqBody).Encode(user)

	// 목업 서버
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()
	resp, err := http.Post(ts.URL+"/json", "application/json", &reqBody)
	assert.NoError(err)
	/*
		http 메서드 중에는 DELTE가 없기 때문에,
		http.NewRequest("DELETE",ts.URL + "/path", nil)
		http.DefaultClient.DO(req) 를 사용.
	*/

	// req := httptest.NewRequest(http.MethodPost, "/user", &reqBody)
	// w := httptest.NewRecorder()
	// JsonHandler(w, req)
	// resp := w.Result()

	assert.Equal(http.StatusCreated, resp.StatusCode)
}

func TestFileHandler(t *testing.T) {
	assert := assert.New(t)
	originFilePath := "../file_test.txt"
	originFile, err := os.Open(originFilePath)
	assert.NoError(err)
	defer originFile.Close()

	_, err = originFile.Stat()
	assert.NoError(err)

	var reqBody bytes.Buffer
	writer := multipart.NewWriter(&reqBody)

	multi, err := writer.CreateFormFile("upload_file", originFile.Name())
	assert.NoError(err)

	_, err = io.Copy(multi, originFile)
	assert.NoError(err)
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/file", &reqBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	FileHandler(w, req)

	resp := w.Result()

	assert.Equal(http.StatusOK, resp.StatusCode)
}
