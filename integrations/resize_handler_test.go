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

type ResizeHandleSuite struct {
	suite.Suite
	addr   string
	ctx    context.Context
	client http.Client
}

func NewResizeHandleSuite() *ResizeHandleSuite {
	return &ResizeHandleSuite{}
}

func (s *ResizeHandleSuite) SetupSuite() {
	conf := SetupSuite()
	s.addr = net.JoinHostPort(conf.HTTP.Host, strconv.Itoa(conf.HTTP.Port))
	s.ctx = context.Background()
	s.client = http.Client{
		Timeout: 30 * time.Second,
	}
}

func (s *ResizeHandleSuite) TestResizingImage() {
	for _, test := range []struct {
		width  string
		height string
		url    string
	}{
		{width: "200", height: "300", url: fmt.Sprintf("%s/gopher_333x666.jpg", nginxHost)},
		{width: "1000", height: "500", url: fmt.Sprintf("%s/gopher_2000x1000.jpg", nginxHost)},
		{width: "600", height: "500", url: fmt.Sprintf("%s/sea_632x474.jpg", nginxHost)},
		{width: "300", height: "200", url: fmt.Sprintf("%s/ubuntu_989x587.png", nginxHost)},
		{width: "600", height: "600", url: fmt.Sprintf("%s/wolf_1024x1024.jpeg", nginxHost)},
	} {
		s.Run(fmt.Sprintf("width:%s height:%s url: %s", test.width, test.height, test.url), func() {
			req, err := http.NewRequestWithContext(
				s.ctx,
				http.MethodGet,
				fmt.Sprintf("http://%s/fill/%s/%s/%s", s.addr, test.width, test.height, test.url),
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
			s.Require().Equal(test.width, strconv.Itoa(bounds.Max.X))
			s.Require().Equal(test.height, strconv.Itoa(bounds.Max.Y))
		})
	}
}

func TestResizeHandleSuite(t *testing.T) {
	suite.Run(t, NewResizeHandleSuite())
}
