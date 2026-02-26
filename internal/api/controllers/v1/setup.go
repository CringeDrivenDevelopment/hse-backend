package v1

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// setup - добавить маршрут до эндпоинта
func (h *Auth) setup(router huma.API) {
	huma.Register(router, huma.Operation{
		OperationID: "auth",
		Path:        "/api/v1/auth",
		Method:      http.MethodPost,
		Errors: []int{
			401,
			422,
			500,
		},
		Tags: []string{
			"auth",
		},
		Summary:     "Login",
		Description: "Получить токен для взаимодействия. Нуждается в Raw строке из Telegram Mini App. Действует 1 час",
	}, h.login)
}

// setup - добавить маршрут до эндпоинтов
func (h *Playlist) setup(router huma.API) {
	huma.Register(router, huma.Operation{
		OperationID: "playlist-by-id",
		Path:        "/api/v1/playlists/{id}",
		Method:      http.MethodGet,
		Errors: []int{
			401,
			404,
			422,
			500,
		},
		Tags: []string{
			"playlist",
		},
		Summary:     "Get by ID",
		Description: "Получить плейлист по ID. Для получения требуется, чтобы у юзера были права на просмотр плейлиста. При получении вернёт массив треков",
		Security: []map[string][]string{
			{
				"jwt": []string{},
			},
		},
	}, h.getById)

	huma.Register(router, huma.Operation{
		OperationID: "playlist-all",
		Path:        "/api/v1/playlists",
		Method:      http.MethodGet,
		Errors: []int{
			401,
			500,
		},
		Tags: []string{
			"playlist",
		},
		Summary:     "All",
		Description: "Получить весь список плейлистов. Вернёт только те плейлисты, к которым у пользователя есть доступ. При получении не вернёт массив треков",
		Security: []map[string][]string{
			{
				"jwt": []string{},
			},
		},
	}, h.getAll)

	huma.Register(router, huma.Operation{
		OperationID: "playlist-export",
		Path:        "/api/v1/playlists/{id}/export",
		Method:      http.MethodPost,
		Errors: []int{
			401,
			500,
		},
		Tags: []string{
			"playlist",
		},
		Summary:     "Export",
		Description: "Экспортировать плейлист (можно в m4a из youtube и создать плейлист в spotify)",
		Security: []map[string][]string{
			{
				"jwt": []string{},
			},
		},
	}, h.export)
}

func (h *Track) setup(router huma.API) {
	huma.Register(router, huma.Operation{
		OperationID: "track-search",
		Path:        "/api/v1/search",
		Method:      http.MethodGet,
		Errors: []int{
			400,
			401,
			422,
			500,
		},
		Tags: []string{
			"tracks",
		},
		Summary:     "Search",
		Description: "Найти трек по запросу. Поиск по Youtube Music / Spotify",
		Security: []map[string][]string{
			{
				"jwt": []string{},
			},
		},
	}, h.search)

	huma.Register(router, huma.Operation{
		OperationID: "tracks-submit",
		Path:        "/api/v1/playlists/{playlist_id}/{track_id}/submit",
		Method:      http.MethodPost,
		Errors: []int{
			401,
			404,
			422,
			500,
		},
		Tags: []string{
			"tracks",
		},
		Summary:     "Submit",
		Description: "Добавить трек в плейлист, если юзер есть в плейлисте. Если у юзера права админа, то трек добавляется в разрешённые, иначе на модерацию",
		Security: []map[string][]string{
			{
				"jwt": []string{},
			},
		},
	}, h.submit)

	huma.Register(router, huma.Operation{
		OperationID: "tracks-unapprove",
		Path:        "/api/v1/playlists/{playlist_id}/{track_id}/unapprove",
		Method:      http.MethodDelete,
		Errors: []int{
			401,
			404,
			422,
			500,
		},
		Tags: []string{
			"tracks",
		},
		Summary:     "Unapprove",
		Description: "Убрать трек из разрешённых. У юзера должны быть права админа",
		Security: []map[string][]string{
			{
				"jwt": []string{},
			},
		},
	}, h.unapprove)

	huma.Register(router, huma.Operation{
		OperationID: "tracks-approve",
		Path:        "/api/v1/playlists/{playlist_id}/{track_id}/approve",
		Method:      http.MethodPatch,
		Errors: []int{
			401,
			404,
			500,
		},
		Tags: []string{
			"tracks",
		},
		Summary:     "Approve",
		Description: "Добавить трек в разрешённые. У юзера должны быть права админа",
		Security: []map[string][]string{
			{
				"jwt": []string{},
			},
		},
	}, h.approve)

	huma.Register(router, huma.Operation{
		OperationID: "tracks-decline",
		Path:        "/api/v1/playlists/{playlist_id}/{track_id}/decline",
		Method:      http.MethodDelete,
		Errors: []int{
			401,
			403,
			404,
			422,
			500,
		},
		Tags: []string{
			"tracks",
		},
		Summary:     "Decline",
		Description: "Удалить трек из кандидатов в плейлист. У юзера должны быть права админа",
		Security: []map[string][]string{
			{
				"jwt": []string{},
			},
		},
	}, h.decline)
}
