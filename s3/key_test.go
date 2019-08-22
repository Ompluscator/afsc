package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs/url"
	"io/ioutil"
	"strings"
	"testing"
)

func TestAES256Key_SetHeader(t *testing.T) {
	authConfig, err := NewTestAuthConfig()
	if err != nil {
		t.Skip(err)
		return

	}
	ctx := context.Background()
	var useCases = []struct {
		description string
		URL         string
		location    string
		data        []byte
		key         string
		base64Key   string
	}{

		{
			description: "securing data with key",
			key:         strings.Repeat("xd", 16),
			location:    "folder/secret1.txt",
			URL:         fmt.Sprintf("s3://%v/", TestBucket),
			data:        []byte("this is test"),
		},

		{
			description: "securing data with base64key",
			location:    "folder/secret2.txt",
			URL:         fmt.Sprintf("s3://%v/", TestBucket),
			data:        []byte("this is test"),
			base64Key:   "eGR4ZHhkeGR4ZHhkeGR4ZHhkeGR4ZHhkeGR4ZHhkeGQ=",
		},
	}

	mgr := New(authConfig)

	defer func() {
		_ = mgr.Delete(ctx, fmt.Sprintf("s3://%v/", TestBucket))
	}()
	for _, useCase := range useCases {

		var key *CustomKey
		if useCase.key != "" {
			key = NewCustomKey([]byte(useCase.key))
		} else {
			key, err = NewBase64CustomKey(useCase.base64Key)
			assert.Nil(t, err, useCase.description)
		}

		URL := url.Join(useCase.URL, useCase.location)
		err := mgr.Upload(ctx, URL, 0644, bytes.NewReader(useCase.data), key)
		assert.Nil(t, err, useCase.description)
		_, err = mgr.DownloadWithURL(ctx, URL)
		assert.NotNil(t, err, useCase.description)
		reader, err := mgr.DownloadWithURL(ctx, URL, key)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		data, err := ioutil.ReadAll(reader)
		assert.EqualValues(t, useCase.data, data, useCase.description)
	}

}
