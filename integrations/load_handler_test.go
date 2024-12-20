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

type LoadHandleSuite struct {
	suite.Suite
	addr   string
	ctx    context.Context
	client http.Client
	width  string
	height string
	URL    string
}

func NewLoadHandleSuite() *LoadHandleSuite {
	return &LoadHandleSuite{}
}

func (s *LoadHandleSuite) SetupSuite() {
	conf := SetupSuite()
	s.addr = net.JoinHostPort(conf.HTTP.Host, strconv.Itoa(conf.HTTP.Port))
	s.ctx = context.Background()
	s.client = http.Client{
		Timeout: 30 * time.Second,
	}
}

func (s *LoadHandleSuite) SetupTest() {
	s.width = "250"
	s.height = "250"
	s.URL = fmt.Sprintf("%s/gopher_333x666.jpg", nginxHost)
}

const (
	timeLimit = 1 * time.Millisecond
)

func (s *LoadHandleSuite) TestLoadingImage() {
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/fill/%s/%s/%s", s.addr, s.width, s.height, s.URL),
		bytes.NewReader(nil),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()

	contentType := response.Header.Get("Content-Type")
	is := strings.Contains(contentType, "image/")
	s.Require().True(is)

	s.Require().Equal(http.StatusOK, response.StatusCode)
	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)
	s.Require().Len(respBody, 24426)
}

func (s *LoadHandleSuite) TestLoadingImageCached() {
	start := time.Now()
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/fill/%s/%s/%s", s.addr, s.width, s.height, s.URL),
		bytes.NewReader(nil),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Require().Equal(http.StatusOK, response.StatusCode)
	_, err = io.ReadAll(response.Body)
	s.Require().NoError(err)

	elapsed := time.Since(start)
	s.Require().Less(elapsed, int64(timeLimit), "the program is too slow")
}

func TestLoadHandleSuite(t *testing.T) {
	suite.Run(t, NewLoadHandleSuite())
}
