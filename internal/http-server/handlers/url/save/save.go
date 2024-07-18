package save

import (
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	URL   string `json:"url" validate: "required, url"`
	Alias string `json: "alias" validate: "required"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias, omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSAve, alias string) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.save.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode reqest"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid reqest", sl.Err(err))
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		alias := req.Alias

		if alias == "" {
			log.Info("failed to add url, empty alias")
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		err = urlSaver.SaveURL(req.URL, alias)

		if err != nil {
			log.Info("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		log.Info("url added")

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
