package save

import (
	"errors"
	"log/slog"
	"net/http"

	resp "url-shortener/internal/lib/api/responce"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver

type Request struct {
	URL string `json:"url" validate:"required,url"`
	// URL   string `json:"url" validate:"url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config
const aliasLenght = 6

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
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
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			validateErr := err.(validator.ValidationErrors)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLenght) // TODO: random alias can be duplicated
		}
		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("URL already exist", slog.String("URL", req.URL))
				render.JSON(w, r, resp.Error("URL already exist")) // тут эрор, вероятно нужно возвращать не ошибку? (пока ошибку так как нет апдейта линков и прочего)
				return
			}
			log.Error("failed to add URL", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add URL"))
			return
		}
		log.Info("URL added", slog.Int64("id", id))
		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
