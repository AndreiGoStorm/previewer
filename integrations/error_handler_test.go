package integrations

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ErrorHandleSuite struct {
	suite.Suite
	addr   string
	ctx    context.Context
	client http.Client
}

func NewErrorHandleSuite() *ErrorHandleSuite {
	return &ErrorHandleSuite{}
}

func (s *ErrorHandleSuite) SetupSuite() {
	conf := SetupSuite()
	s.addr = net.JoinHostPort(conf.HTTP.Host, strconv.Itoa(conf.HTTP.Port))
	s.ctx = context.Background()
	s.client = http.Client{
		Timeout: 30 * time.Second,
	}
}

func (s *ErrorHandleSuite) TestMethodNotFound() {
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/wrong/300/300", s.addr),
		bytes.NewReader(nil),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()
	_, err = io.ReadAll(response.Body)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, response.StatusCode)
}

func (s *ErrorHandleSuite) TestPostMethodNotAllowed() {
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s/fill/300/300", s.addr),
		bytes.NewReader(nil),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusMethodNotAllowed, response.StatusCode)
	is := strings.Contains(string(respBody), fmt.Sprintf("method %s not not supported on uri", http.MethodPost))
	s.Require().True(is)
}

func (s *ErrorHandleSuite) TestWrongWidth() {
	width := "wrong"
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/fill/%s/300", s.addr, width),
		bytes.NewReader(nil),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnprocessableEntity, response.StatusCode)
	is := strings.Contains(string(respBody), fmt.Sprintf("wrong width: %s", width))
	s.Require().True(is)
}

func (s *ErrorHandleSuite) TestWrongHeight() {
	height := "0"
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/fill/500/%s", s.addr, height),
		bytes.NewReader(nil),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnprocessableEntity, response.StatusCode)
	is := strings.Contains(string(respBody), fmt.Sprintf("wrong height: %s", height))
	s.Require().True(is)
}

func (s *ErrorHandleSuite) TestEmptyURL() {
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/fill/500/500", s.addr),
		bytes.NewReader(nil),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnprocessableEntity, response.StatusCode)
	is := strings.Contains(string(respBody), "loading url is empty")
	s.Require().True(is)
}

func (s *ErrorHandleSuite) TestEmptyExt() {
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/fill/500/500/example/url", s.addr),
		bytes.NewReader(nil),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnprocessableEntity, response.StatusCode)
	is := strings.Contains(string(respBody), "loading image extension is empty")
	s.Require().True(is)
}

func (s *ErrorHandleSuite) TestWrongExt() {
	ext := "com"
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/fill/500/500/example.%s", s.addr, ext),
		bytes.NewReader(nil),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnprocessableEntity, response.StatusCode)
	is := strings.Contains(string(respBody), fmt.Sprintf("loading image has wrong extension: %s", ext))
	s.Require().True(is)
}

func (s *ErrorHandleSuite) TestLoadingError() {
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/fill/500/500/localhost.jpeg", s.addr),
		bytes.NewReader(nil),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadGateway, response.StatusCode)
	is := strings.Contains(string(respBody), "no such host")
	s.Require().True(is)
}

func TestErrorHandleSuite(t *testing.T) {
	suite.Run(t, NewErrorHandleSuite())
}
