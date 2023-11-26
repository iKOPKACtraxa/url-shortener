package delete

import (
	"log/slog"
	"net/http"

	resp "url-shortener/internal/lib/api/responce"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLDeleter

type Response struct {
	resp.Response
}

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")
		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}
		log.Info("deleted", slog.String("alias", alias))
		render.JSON(w, r, resp.OK())
	}
}
