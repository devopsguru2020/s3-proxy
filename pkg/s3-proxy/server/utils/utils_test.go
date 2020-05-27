// +build unit

package utils

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/config"
	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/log"
)

func TestHandleInternalServerError(t *testing.T) {
	headers := http.Header{}
	headers.Add("Content-Type", "text/html; charset=utf-8")
	type args struct {
		rw          http.ResponseWriter
		err         error
		requestPath string
		tplCfg      *config.TemplateConfig
	}
	tests := []struct {
		name               string
		args               args
		expectedHTTPWriter *respWriterTest
	}{
		{
			name: "Template should be ok",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				err:         errors.New("fake"),
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "../../../../templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "../../../../templates/unauthorized.tpl",
					Forbidden:           "../../../../templates/forbidden.tpl",
					BadRequest:          "../../../../templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Internal Server Error</h1>
    <p>fake</p>
  </body>
</html>
`),
			},
		},
		{
			name: "Template not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				err:         errors.New("fake"),
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "templates/internal-server-error.tpl",
					Unauthorized:        "templates/unauthorized.tpl",
					Forbidden:           "templates/forbidden.tpl",
					BadRequest:          "templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`
<!DOCTYPE html>
<html>
  <body>
	<h1>Internal Server Error</h1>
	<p>open templates/internal-server-error.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleInternalServerError(log.NewLogger(), tt.args.rw, tt.args.tplCfg, tt.args.requestPath, tt.args.err)
			if !reflect.DeepEqual(tt.expectedHTTPWriter, tt.args.rw) {
				t.Errorf("HandleInternalServerError() => httpWriter = %+v, want %+v", tt.args.rw, tt.expectedHTTPWriter)
			}
		})
	}
}

func TestHandleInternalServerErrorWithTemplate(t *testing.T) {
	headers := http.Header{}
	headers.Add("Content-Type", "text/html; charset=utf-8")
	type args struct {
		tplString   string
		rw          http.ResponseWriter
		err         error
		requestPath string
		tplCfg      *config.TemplateConfig
	}
	tests := []struct {
		name               string
		args               args
		expectedHTTPWriter *respWriterTest
	}{
		{
			name: "Without template should be ok",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				err:         errors.New("fake"),
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "../../../../templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "../../../../templates/unauthorized.tpl",
					Forbidden:           "../../../../templates/forbidden.tpl",
					BadRequest:          "../../../../templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Internal Server Error</h1>
    <p>fake</p>
  </body>
</html>
`),
			},
		},
		{
			name: "Template string should be ok",
			args: args{
				tplString: `Fake template`,
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				err:         errors.New("fake"),
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "../../../../templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "../../../../templates/unauthorized.tpl",
					Forbidden:           "../../../../templates/forbidden.tpl",
					BadRequest:          "../../../../templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp:    []byte(`Fake template`),
			},
		},
		{
			name: "Template not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				err:         errors.New("fake"),
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "templates/internal-server-error.tpl",
					Unauthorized:        "templates/unauthorized.tpl",
					Forbidden:           "templates/forbidden.tpl",
					BadRequest:          "templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`
<!DOCTYPE html>
<html>
  <body>
	<h1>Internal Server Error</h1>
	<p>open templates/internal-server-error.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleInternalServerErrorWithTemplate(log.NewLogger(), tt.args.rw, tt.args.tplCfg, tt.args.tplString, tt.args.requestPath, tt.args.err)
			if !reflect.DeepEqual(tt.expectedHTTPWriter, tt.args.rw) {
				t.Errorf("HandleInternalServerError() => httpWriter = %+v, want %+v", tt.args.rw, tt.expectedHTTPWriter)
			}
		})
	}
}

func TestHandleNotFound(t *testing.T) {
	headers := http.Header{}
	headers.Add("Content-Type", "text/html; charset=utf-8")
	type args struct {
		rw          http.ResponseWriter
		requestPath string
		tplCfg      *config.TemplateConfig
	}
	tests := []struct {
		name               string
		args               args
		expectedHTTPWriter *respWriterTest
	}{
		{
			name: "Template should be ok",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "../../../../templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "../../../../templates/unauthorized.tpl",
					Forbidden:           "../../../../templates/forbidden.tpl",
					BadRequest:          "../../../../templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  404,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Not Found /request1</h1>
  </body>
</html>
`),
			},
		},
		{
			name: "Template not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "../../../../templates/unauthorized.tpl",
					Forbidden:           "../../../../templates/forbidden.tpl",
					BadRequest:          "../../../../templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Internal Server Error</h1>
    <p>open templates/not-found.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
		{
			name: "All templates not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "templates/internal-server-error.tpl",
					Unauthorized:        "templates/unauthorized.tpl",
					Forbidden:           "templates/forbidden.tpl",
					BadRequest:          "templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`
<!DOCTYPE html>
<html>
  <body>
	<h1>Internal Server Error</h1>
	<p>open templates/internal-server-error.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleNotFound(log.NewLogger(), tt.args.rw, tt.args.tplCfg, tt.args.requestPath)
			if !reflect.DeepEqual(tt.expectedHTTPWriter, tt.args.rw) {
				t.Errorf("HandleNotFound() => httpWriter = %+v, want %+v", tt.args.rw, tt.expectedHTTPWriter)
			}
		})
	}
}

func TestHandleNotFoundWithTemplate(t *testing.T) {
	headers := http.Header{}
	headers.Add("Content-Type", "text/html; charset=utf-8")
	type args struct {
		tplString   string
		rw          http.ResponseWriter
		requestPath string
		tplCfg      *config.TemplateConfig
	}
	tests := []struct {
		name               string
		args               args
		expectedHTTPWriter *respWriterTest
	}{
		{
			name: "Without template should be ok",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "../../../../templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "../../../../templates/unauthorized.tpl",
					Forbidden:           "../../../../templates/forbidden.tpl",
					BadRequest:          "../../../../templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  404,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Not Found /request1</h1>
  </body>
</html>
`),
			},
		},
		{
			name: "Template string should be ok",
			args: args{
				tplString: "Fake template",
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "../../../../templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "../../../../templates/unauthorized.tpl",
					Forbidden:           "../../../../templates/forbidden.tpl",
					BadRequest:          "../../../../templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  404,
				Resp:    []byte("Fake template"),
			},
		},
		{
			name: "Template not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "../../../../templates/unauthorized.tpl",
					Forbidden:           "../../../../templates/forbidden.tpl",
					BadRequest:          "../../../../templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Internal Server Error</h1>
    <p>open templates/not-found.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
		{
			name: "All templates not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "templates/internal-server-error.tpl",
					Unauthorized:        "templates/unauthorized.tpl",
					Forbidden:           "templates/forbidden.tpl",
					BadRequest:          "templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`
<!DOCTYPE html>
<html>
  <body>
	<h1>Internal Server Error</h1>
	<p>open templates/internal-server-error.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleNotFoundWithTemplate(log.NewLogger(), tt.args.rw, tt.args.tplCfg, tt.args.tplString, tt.args.requestPath)
			if !reflect.DeepEqual(tt.expectedHTTPWriter, tt.args.rw) {
				t.Errorf("HandleNotFound() => httpWriter = %+v, want %+v", tt.args.rw, tt.expectedHTTPWriter)
			}
		})
	}
}

func TestHandleUnauthorized(t *testing.T) {
	headers := http.Header{}
	headers.Add("Content-Type", "text/html; charset=utf-8")
	type args struct {
		rw          http.ResponseWriter
		requestPath string
		tplCfg      *config.TemplateConfig
	}
	tests := []struct {
		name               string
		args               args
		expectedHTTPWriter *respWriterTest
	}{
		{
			name: "Template should be ok",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "../../../../templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "../../../../templates/unauthorized.tpl",
					Forbidden:           "../../../../templates/forbidden.tpl",
					BadRequest:          "../../../../templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  401,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Unauthorized</h1>
  </body>
</html>
`),
			},
		},
		{
			name: "Template not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "templates/unauthorized.tpl",
					Forbidden:           "templates/forbidden.tpl",
					BadRequest:          "templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Internal Server Error</h1>
    <p>open templates/unauthorized.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
		{
			name: "All templates not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "templates/internal-server-error.tpl",
					Unauthorized:        "templates/unauthorized.tpl",
					Forbidden:           "templates/forbidden.tpl",
					BadRequest:          "templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`
<!DOCTYPE html>
<html>
  <body>
	<h1>Internal Server Error</h1>
	<p>open templates/internal-server-error.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleUnauthorized(log.NewLogger(), tt.args.rw, tt.args.tplCfg, tt.args.requestPath)
			if !reflect.DeepEqual(tt.expectedHTTPWriter, tt.args.rw) {
				t.Errorf("HandleUnauthorized() => httpWriter = %+v, want %+v", tt.args.rw, tt.expectedHTTPWriter)
			}
		})
	}
}

func TestHandleBadRequest(t *testing.T) {
	headers := http.Header{}
	headers.Add("Content-Type", "text/html; charset=utf-8")
	type args struct {
		rw          http.ResponseWriter
		requestPath string
		err         error
		tplCfg      *config.TemplateConfig
	}
	tests := []struct {
		name               string
		args               args
		expectedHTTPWriter *respWriterTest
	}{
		{
			name: "Template should be ok",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				err:         errors.New("fake"),
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "../../../../templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "../../../../templates/unauthorized.tpl",
					Forbidden:           "../../../../templates/forbidden.tpl",
					BadRequest:          "../../../../templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  400,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Bad Request</h1>
    <p>fake</p>
  </body>
</html>
`),
			},
		},
		{
			name: "Template not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "templates/unauthorized.tpl",
					Forbidden:           "templates/forbidden.tpl",
					BadRequest:          "templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Internal Server Error</h1>
    <p>open templates/bad-request.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
		{
			name: "All templates not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "templates/internal-server-error.tpl",
					Unauthorized:        "templates/unauthorized.tpl",
					Forbidden:           "templates/forbidden.tpl",
					BadRequest:          "templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`
<!DOCTYPE html>
<html>
  <body>
	<h1>Internal Server Error</h1>
	<p>open templates/internal-server-error.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleBadRequest(log.NewLogger(), tt.args.rw, tt.args.tplCfg, tt.args.requestPath, tt.args.err)
			if !reflect.DeepEqual(tt.expectedHTTPWriter, tt.args.rw) {
				t.Errorf("HandleBadRequest() => httpWriter = %+v, want %+v", tt.args.rw, tt.expectedHTTPWriter)
			}
		})
	}
}

func TestHandleForbiddenWithTemplate(t *testing.T) {
	headers := http.Header{}
	headers.Add("Content-Type", "text/html; charset=utf-8")
	type args struct {
		tplString   string
		rw          http.ResponseWriter
		requestPath string
		tplCfg      *config.TemplateConfig
	}
	tests := []struct {
		name               string
		args               args
		expectedHTTPWriter *respWriterTest
	}{
		{
			name: "Template should be ok",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "../../../../templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "../../../../templates/unauthorized.tpl",
					Forbidden:           "../../../../templates/forbidden.tpl",
					BadRequest:          "../../../../templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  403,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Forbidden</h1>
  </body>
</html>
`),
			},
		},
		{
			name: "Template string should be ok",
			args: args{
				tplString: "Fake template",
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "../../../../templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "../../../../templates/unauthorized.tpl",
					Forbidden:           "../../../../templates/forbidden.tpl",
					BadRequest:          "../../../../templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  403,
				Resp:    []byte(`Fake template`),
			},
		},
		{
			name: "Template not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "templates/unauthorized.tpl",
					Forbidden:           "templates/forbidden.tpl",
					BadRequest:          "templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Internal Server Error</h1>
    <p>open templates/forbidden.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
		{
			name: "All templates not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "templates/internal-server-error.tpl",
					Unauthorized:        "templates/unauthorized.tpl",
					Forbidden:           "templates/forbidden.tpl",
					BadRequest:          "templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`
<!DOCTYPE html>
<html>
  <body>
	<h1>Internal Server Error</h1>
	<p>open templates/internal-server-error.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleForbiddenWithTemplate(log.NewLogger(), tt.args.rw, tt.args.tplCfg, tt.args.tplString, tt.args.requestPath)
			if !reflect.DeepEqual(tt.expectedHTTPWriter, tt.args.rw) {
				t.Errorf("HandleForbidden() => httpWriter = %+v, want %+v", tt.args.rw, tt.expectedHTTPWriter)
			}
		})
	}
}

func TestHandleForbidden(t *testing.T) {
	headers := http.Header{}
	headers.Add("Content-Type", "text/html; charset=utf-8")
	type args struct {
		rw          http.ResponseWriter
		requestPath string
		tplCfg      *config.TemplateConfig
	}
	tests := []struct {
		name               string
		args               args
		expectedHTTPWriter *respWriterTest
	}{
		{
			name: "Template should be ok",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "../../../../templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "../../../../templates/unauthorized.tpl",
					Forbidden:           "../../../../templates/forbidden.tpl",
					BadRequest:          "../../../../templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  403,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Forbidden</h1>
  </body>
</html>
`),
			},
		},
		{
			name: "Template not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "../../../../templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "../../../../templates/internal-server-error.tpl",
					Unauthorized:        "templates/unauthorized.tpl",
					Forbidden:           "templates/forbidden.tpl",
					BadRequest:          "templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`<!DOCTYPE html>
<html>
  <body>
    <h1>Internal Server Error</h1>
    <p>open templates/forbidden.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
		{
			name: "All templates not found",
			args: args{
				rw: &respWriterTest{
					Headers: http.Header{},
				},
				requestPath: "/request1",
				tplCfg: &config.TemplateConfig{
					TargetList:          "templates/target-list.tpl",
					NotFound:            "templates/not-found.tpl",
					InternalServerError: "templates/internal-server-error.tpl",
					Unauthorized:        "templates/unauthorized.tpl",
					Forbidden:           "templates/forbidden.tpl",
					BadRequest:          "templates/bad-request.tpl",
				},
			},
			expectedHTTPWriter: &respWriterTest{
				Headers: headers,
				Status:  500,
				Resp: []byte(`
<!DOCTYPE html>
<html>
  <body>
	<h1>Internal Server Error</h1>
	<p>open templates/internal-server-error.tpl: no such file or directory</p>
  </body>
</html>
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleForbidden(log.NewLogger(), tt.args.rw, tt.args.tplCfg, tt.args.requestPath)
			if !reflect.DeepEqual(tt.expectedHTTPWriter, tt.args.rw) {
				t.Errorf("HandleForbidden() => httpWriter = %+v, want %+v", tt.args.rw, tt.expectedHTTPWriter)
			}
		})
	}
}
