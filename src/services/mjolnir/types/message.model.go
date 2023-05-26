package types

type MjolnirConfig struct {
	RedisSubChannel string
	StdLogsEnable   bool
	Services        map[string]MjolnirService
}

type MjolnirService struct {
	Name    string
	LogPath string
	S3Path  string
}

type MjolnirMessage struct {
	Service    string `json:"service"`
	Level      string `json:"level"`
	TimeStamp  string `json:"timestamp,omitempty"`
	Message    string `json:"message"`
	Caller     string `json:"caller,omitempty"`
	StackTrace string `json:"stacktrace,omitempty"`
}
