package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestValidateWidth(t *testing.T) {
	req := &Request{}
	t.Run("validate width: correct width", func(t *testing.T) {
		width := "300"
		err := req.validateWidth(width)
		require.NoError(t, err)
	})

	t.Run("error validate width: wrong width", func(t *testing.T) {
		width := "wrong width"
		err := req.validateWidth(width)
		require.Error(t, err)
		require.EqualError(t, err, "wrong width: wrong width")
	})

	t.Run("error validate width: width too large", func(t *testing.T) {
		width := "10150"
		err := req.validateWidth(width)
		require.Error(t, err)
		require.EqualError(t, err, "wrong width: 10150")
	})

	t.Run("error validate width: width less or equal zero", func(t *testing.T) {
		width := "-10"
		err := req.validateWidth(width)
		require.Error(t, err)
		require.EqualError(t, err, "wrong width: -10")
	})
}

func TestRequestValidateHeight(t *testing.T) {
	req := &Request{}
	t.Run("validate height: correct height", func(t *testing.T) {
		err := req.validateHeight("500")
		require.NoError(t, err)
	})

	t.Run("error validate height: wrong height", func(t *testing.T) {
		height := "wrong height"
		err := req.validateHeight(height)
		require.Error(t, err)
		require.EqualError(t, err, "wrong height: wrong height")
	})

	t.Run("error validate height: height too large", func(t *testing.T) {
		height := "10500"
		err := req.validateHeight(height)
		require.Error(t, err)
		require.EqualError(t, err, "wrong height: 10500")
	})

	t.Run("error validate height: height less or equal zero", func(t *testing.T) {
		height := "-5"
		err := req.validateHeight(height)
		require.Error(t, err)
		require.EqualError(t, err, "wrong height: -5", height)
	})
}

func TestRequestValidateURL(t *testing.T) {
	req := &Request{}
	t.Run("validate url: correct url", func(t *testing.T) {
		req.Protocol = "https"
		err := req.validateURL("localhost/gopher_500x500.jpg")
		require.NoError(t, err)
	})

	t.Run("error validate url: empty", func(t *testing.T) {
		err := req.validateURL("")
		require.Error(t, err)
		require.EqualError(t, err, "loading url is empty")
	})

	t.Run("validate url: correct url", func(t *testing.T) {
		req.Protocol = ""
		err := req.validateURL("example.com")
		require.Error(t, err)
		require.EqualError(t, err, "wrong url")
	})
}

func TestRequestValidateExt(t *testing.T) {
	req := &Request{}
	t.Run("error validate url: empty extension", func(t *testing.T) {
		err := req.validateExt("github.com/stretchr/testify/require")
		require.Error(t, err)
		require.EqualError(t, err, "loading image extension is empty")
	})

	t.Run("error validate url: wrong extension", func(t *testing.T) {
		err := req.validateExt("github.com/stretchr/testify/require.ttf")
		require.Error(t, err)
		require.EqualError(t, err, "loading image has wrong extension: ttf")
	})
}
