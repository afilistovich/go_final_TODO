package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/afilistovich/go_final_TODO/internal/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server обёртка над стандартным http.Server с настройками приложения
type Server struct {
	server *http.Server
}

// NewServer создаёт и настраивает новый экземпляр сервера
// port - порт для прослушивания (например, "7540")
// webDir - путь к директории с фронтендом (HTML, CSS, JS)
func NewServer(port string, webDir string) *Server {
	r := chi.NewRouter()

	// Подключаем стандартные middleware от chi
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Регистрируем все API маршруты
	// Маршруты определены в пакете api
	api.RegisterRoutes(r)

	// Настраиваем файловый сервер для раздачи статики (фронтенда)
	// http.FileServer автоматически раздаёт файлы из указанной директории
	// Маршрут "/*" означает, что все запросы, не совпавшие с API, пойдут сюда
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	// Создаём конфигурацию HTTP-сервера с таймаутами
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		server: srv,
	}
}

// Start запускает сервер и начинает принимать HTTP-запросы
func (s *Server) Start() error {
	slog.Info("Server is listening", "addr", s.server.Addr)
	return s.server.ListenAndServe()
}
