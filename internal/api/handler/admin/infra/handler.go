package infra

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewConfigHandler,
	NewFileConfigHandler,
	NewFileHandler,
	NewApiAccessLogHandler,
	NewApiErrorLogHandler,
	NewJobHandler,
	NewJobLogHandler,
	NewWebSocketHandler,
	NewHandlers,
)

type Handlers struct {
	Config       *ConfigHandler
	FileConfig   *FileConfigHandler
	File         *FileHandler
	ApiAccessLog *ApiAccessLogHandler
	ApiErrorLog  *ApiErrorLogHandler
	Job          *JobHandler
	JobLog       *JobLogHandler
	WebSocket    *WebSocketHandler
}

func NewHandlers(
	config *ConfigHandler,
	fileConfig *FileConfigHandler,
	file *FileHandler,
	apiAccessLog *ApiAccessLogHandler,
	apiErrorLog *ApiErrorLogHandler,
	job *JobHandler,
	jobLog *JobLogHandler,
	websocket *WebSocketHandler,
) *Handlers {
	return &Handlers{
		Config:       config,
		FileConfig:   fileConfig,
		File:         file,
		ApiAccessLog: apiAccessLog,
		ApiErrorLog:  apiErrorLog,
		Job:          job,
		JobLog:       jobLog,
		WebSocket:    websocket,
	}
}
