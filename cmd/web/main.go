package main

import (
	"flag"
	"github/len4ernova/lets_go/internal/models"
	"net/http"
	"os"
	"text/template"

	"github.com/go-playground/form/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "modernc.org/sqlite"
)

type application struct {
	logger *zap.Logger
	//snippets      *models.SnippetModel
	works         *models.WorkModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	ip := flag.String("ip", "localhost", "HTTP network ip")
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "access.db", "Sqlite data source name")
	flag.Parse()

	// Настройка логгера: вывода логов в консоль в формате JSON
	configZap := zap.Config{
		Encoding:         "json", // формат вывода
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout"}, // вывод в консоль
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}
	logger, _ := configZap.Build()
	defer logger.Sync()

	//	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := models.OpenDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Sugar().Error(err.Error())
		os.Exit(1)
	}
	//Initialize a decoder instance...
	formDecoder := form.NewDecoder()
	// add it to the application dependencies.
	app := &application{
		logger: logger,
		//snippets:      &models.SnippetModel{DB: db},
		works:         &models.WorkModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}
	logger.Sugar().Info("starting server addr: ", *ip+*addr)
	err = http.ListenAndServe(*ip+*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
