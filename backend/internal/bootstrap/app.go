package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"agendamento-automotivo/backend/internal/application"
	httpapi "agendamento-automotivo/backend/internal/httpapi"
	mysqladapter "agendamento-automotivo/backend/internal/infrastructure/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct{ server *http.Server }

func New() (*App, error) {
	db, err := connectDatabase()
	if err != nil {
		return nil, fmt.Errorf("conectar ao banco: %w", err)
	}

	store := mysqladapter.NewStore(db)
	if err := store.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("executar migrações: %w", err)
	}
	if err := store.Seed(context.Background()); err != nil {
		return nil, fmt.Errorf("popular dados iniciais: %w", err)
	}

	clientes := application.NewClienteService(mysqladapter.NewClienteRepository(db))
	servicos := application.NewServicoService(mysqladapter.NewServicoRepository(db))
	agendamentos := application.NewAgendamentoService(mysqladapter.NewAgendamentoRepository(db))
	handler := httpapi.NewHandler(clientes, servicos, agendamentos)
	port := env("PORT", "8080")
	return &App{server: &http.Server{Addr: ":" + port, Handler: handler.Router(), ReadHeaderTimeout: 5 * time.Second}}, nil
}

func (a *App) Run() error {
	log.Printf("API disponível na porta %s", a.server.Addr)
	return a.server.ListenAndServe()
}

func connectDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", env("DB_USER", "appuser"), env("DB_PASSWORD", "apppassword"), env("DB_HOST", "127.0.0.1"), env("DB_PORT", "3306"), env("DB_NAME", "agendamento_db"))
	var lastErr error
	for attempt := 1; attempt <= 30; attempt++ {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			var sqlDBErr error
			sqlDB, sqlDBErr := db.DB()
			if sqlDBErr == nil {
				sqlDBErr = sqlDB.Ping()
			}
			if sqlDBErr == nil {
				return db, nil
			}
			err = sqlDBErr
		}
		lastErr = err
		log.Printf("MySQL indisponível (tentativa %d/30): %v", attempt, err)
		time.Sleep(2 * time.Second)
	}
	return nil, lastErr
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
