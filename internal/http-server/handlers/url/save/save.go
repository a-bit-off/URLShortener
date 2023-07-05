// TODO: add: test
package save

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"golang.org/x/exp/slog"

	"URLShortener/internal/config"
	"URLShortener/internal/lib/api/response"
	"URLShortener/internal/lib/logger/sl"
	"URLShortener/internal/lib/random"
	"URLShortener/internal/storage"
)

// генерация mock
//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=URLSaver

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	GetURL(alias string) (string, error)
}

func New(cfg *config.Config, log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		alias := req.Alias
		// генерирует рандомный alias в случае отсутсвия
		err = generateUniqueAlias(cfg, urlSaver, &alias)
		if err != nil {
			log.Error("failed to generate alias", sl.Err(err))
			render.JSON(w, r, response.Error("failed to generate alias"))
			return
		}

		// сохраняем данны в бд
		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, response.Error("failed to add url"))
			return
		}

		log.Info("url added", slog.Int64("id", id))
		responseOK(w, r, alias)
	}
}

func generateUniqueAlias(cfg *config.Config, urlSaver URLSaver, alias *string) error {
	if *alias != "" {
		return nil
	}
	for {
		*alias = random.NewRandomString(cfg.AliasLength)

		//	проверка на уникальность в базе данных
		url, err := urlSaver.GetURL(*alias)
		if url == "" {
			return nil
		} else if err != nil {
			return err
		}
	}
	return errors.New("failed to generate alias")
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    alias,
	})
}
