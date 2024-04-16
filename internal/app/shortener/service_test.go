package shortener

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	repository := NewRepository()
	service := NewService(repository)

	assert.Equal(t, &Service{repository}, service)
}

func TestService_ShrtByUrl(t *testing.T) {
	type fields struct {
		r map[string]string
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"success exists",
			fields{map[string]string{"a": "b", "foo": "bar", "bar": "foo", "lorem": "ipsum"}},
			args{"foo"},
			"bar",
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MemRepository{tt.fields.r}
			s := &Service{r}
			got, err := s.ShrtByUrl(tt.args.url)
			if !tt.wantErr(t, err, fmt.Sprintf("ShrtByUrl(%v)", tt.args.url)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ShrtByUrl(%v)", tt.args.url)
		})
	}
}

func TestService_ShrtByUrl_Creating(t *testing.T) {
	service := NewService(NewRepository())

	url1, err := service.ShrtByUrl("foo")
	assert.Nil(t, err)

	url2, err := service.ShrtByUrl("bar")
	assert.Nil(t, err)

	url3, err := service.ShrtByUrl("baz")
	assert.Nil(t, err)

	res := []string{url1, url2, url3}
	for _, u := range res {
		assert.Len(t, u, shortLen)
	}

	assert.NotEqual(t, url1, url2)
	assert.NotEqual(t, url3, url1)
	assert.NotEqual(t, url2, url3)

	url, _ := service.ShrtByUrl("foo")
	assert.Equal(t, url1, url)
	url, _ = service.ShrtByUrl("bar")
	assert.Equal(t, url2, url)
}

func TestService_UrlByShrt(t *testing.T) {
	type fields struct {
		r map[string]string
	}
	type args struct {
		shrt string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"suc",
			fields{map[string]string{"a": "b", "foo": "bar", "bar": "foo", "lorem": "ipsum"}},
			args{"bar"},
			"foo",
			assert.NoError,
		},
		{
			"fail",
			fields{map[string]string{"a": "b"}},
			args{"bar"},
			"foo",
			assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				&MemRepository{tt.fields.r},
			}
			got, err := s.UrlByShrt(tt.args.shrt)
			if !tt.wantErr(t, err, fmt.Sprintf("UrlByShrt(%v)", tt.args.shrt)) {
				return
			}
			if err != nil {
				return
			}
			assert.Equalf(t, tt.want, got, "UrlByShrt(%v)", tt.args.shrt)
		})
	}
}
