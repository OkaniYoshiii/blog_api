package routes

import (
	"net/http"
)

func LoginHandler(deps Dependencies) http.HandlerFunc {
	logger := deps.Logger
	db := deps.
	queries := deps.Queries

	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
			case http.MethodPost:
				PostLoginHandler(deps, writer, request)
			default:
				writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func PostLoginHandler(deps Dependencies, writer http.ResponseWriter, request *http.Request) {

}
