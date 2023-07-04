package delete

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	"URLShortener/internal/lib/api/response"
	"URLShortener/internal/lib/logger/sl"
	"URLShortener/internal/storage"
)

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias in empty")
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, response.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		log.Info("delete alias", slog.String("url", alias))
		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    alias,
	})
}
