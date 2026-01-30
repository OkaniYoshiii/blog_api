package middleware

import (
	"net/http"
	"strings"
	"testing"
)

func TestPipe(t *testing.T) {
	buffer := strings.Builder{}
	tests := [...]struct {
		Name        string
		Middlewares []Middleware
		Output      string
	}{
		{
			Name: "Single middleware",
			Middlewares: []Middleware{
				func(next http.Handler) http.Handler {
					buffer.WriteString("0")
					return next
				},
			},
			Output: "0",
		},
		{
			Name:        "No middlewares",
			Middlewares: []Middleware{},
			Output:      "",
		},
		{
			Name: "Correct order of execution",
			Middlewares: []Middleware{
				func(next http.Handler) http.Handler {
					buffer.WriteString("0")
					return next
				},
				func(next http.Handler) http.Handler {
					buffer.WriteString("1")
					return next
				},
			},
			Output: "01",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			Pipe(test.Middlewares...)(http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {}))

			expected := test.Output
			got := buffer.String()
			if expected != got {
				t.Errorf("incorrect output : expected %s, got %s", expected, got)
			}

			buffer.Reset()
		})
	}
}
