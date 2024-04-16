package shortener

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemRepository_GetByUrl(t *testing.T) {
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
		want    *Url
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"found",
			fields{map[string]string{"a": "b", "foo": "bar", "bar": "foo", "lorem": "ipsum"}},
			args{"foo"},
			&Url{"foo", "bar"},
			assert.NoError,
		},
		{
			"not found",
			fields{map[string]string{"foo": "bar"}},
			args{"baz"},
			nil,
			assert.Error,
		},
		{
			"empty rep",
			fields{make(map[string]string)},
			args{"foo"},
			nil,
			assert.Error,
		},
		{
			"empty url",
			fields{map[string]string{"": "bar"}},
			args{""},
			&Url{"", "bar"},
			assert.NoError,
		},
		{
			"empty shrt",
			fields{map[string]string{"foo": ""}},
			args{"foo"},
			&Url{"foo", ""},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MemRepository{
				r: tt.fields.r,
			}
			got, err := r.GetByUrl(tt.args.url)
			if !tt.wantErr(t, err, fmt.Sprintf("GetByUrl(%v)", tt.args.url)) {
				return
			}

			if err != nil {
				assert.ErrorIs(t, err, NotFoundError)
				return
			}
			assert.Equalf(t, tt.want, got, "GetByUrl(%v)", tt.args.url)
		})
	}
}

func TestMemRepository_GetByShortCode(t *testing.T) {
	type fields struct {
		r map[string]string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Url
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"found",
			fields{map[string]string{"a": "b", "foo": "bar", "bar": "foo", "lorem": "ipsum"}},
			args{"bar"},
			&Url{"foo", "bar"},
			assert.NoError,
		},
		{
			"not found",
			fields{map[string]string{"foo": "bar"}},
			args{"foo"},
			nil,
			assert.Error,
		},
		{
			"empty",
			fields{map[string]string{"a": "b", "foo": "", "bar": "foo", "lorem": "ipsum"}},
			args{""},
			&Url{"foo", ""},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MemRepository{
				r: tt.fields.r,
			}
			got, err := r.GetByShortCode(tt.args.code)
			if !tt.wantErr(t, err, fmt.Sprintf("GetByShortCode(%v)", tt.args.code)) {

				return
			}

			if err != nil {
				assert.ErrorIs(t, err, NotFoundError)
				return
			}
			assert.Equalf(t, tt.want, got, "GetByShortCode(%v)", tt.args.code)
		})
	}
}

func TestMemRepository_Save(t *testing.T) {
	type fields struct {
		r map[string]string
	}
	type args struct {
		url *Url
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  assert.ErrorAssertionFunc
		expError error
	}{
		{
			"success",
			fields{map[string]string{"foo": "bar"}},
			args{&Url{"new", "url"}},
			assert.NoError,
			nil,
		},
		{
			"already exists",
			fields{map[string]string{"foo": "bar"}},
			args{&Url{"foo", "baz"}},
			assert.Error,
			AlreadyExistsError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MemRepository{
				r: tt.fields.r,
			}
			argUrl := tt.args.url.Url
			err := r.Save(tt.args.url)
			url, _ := r.GetByUrl(argUrl)
			if !tt.wantErr(t, err, fmt.Sprintf("Save(%v)", tt.args.url)) {
				return
			}
			if err != nil {
				assert.ErrorIs(t, err, tt.expError)
				return
			}
			assert.Equal(t, tt.args.url, url)
		})
	}
}

func TestNewRepository(t *testing.T) {
	tests := []struct {
		name string
		want *MemRepository
	}{
		{"success", &MemRepository{r: make(map[string]string)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewRepository(), "NewRepository()")
		})
	}
}