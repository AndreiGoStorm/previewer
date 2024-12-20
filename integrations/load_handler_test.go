package integrations

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/disintegration/imaging"
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
	s.URL = "raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/gopher_500x500.jpg"
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

	s.Require().Equal(http.StatusOK, response.StatusCode)
	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)

	contentType := response.Header.Get("Content-Type")
	is := strings.Contains(contentType, "image/")
	s.Require().True(is)

	toFile, _ := os.CreateTemp("/tmp/", "*.jpg")
	defer os.Remove(toFile.Name())

	_, err = toFile.Write(respBody)
	s.Require().NoError(err)

	file, err := os.Open(toFile.Name())
	s.Require().NoError(err)
	defer file.Close()

	img, err := imaging.Decode(file)
	s.Require().NoError(err)

	bounds := img.Bounds()
	s.Require().Equal(s.width, strconv.Itoa(bounds.Max.X))
	s.Require().Equal(s.height, strconv.Itoa(bounds.Max.Y))
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
	fmt.Println(elapsed)
}

func TestLoadHandleSuite(t *testing.T) {
	suite.Run(t, NewLoadHandleSuite())
}
