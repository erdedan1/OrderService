package config

import "time"

type GRPCServerConfig struct {
	Address              string        `env:"GRPC_SERVER_ADDRESS" validate:"required"`
	MaxRecvMsgSize       int           `env:"GRPC_SERVER_MAX_RECV_MSG_SIZE" validate:"gte=0"`
	MaxSendMsgSize       int           `env:"GRPC_SERVER_MAX_SEND_MSG_SIZE" validate:"gte=0"`
	EnableReflection     bool          `env:"GRPC_SERVER_ENABLE_REFLECTION" validate:"-"`
	TLSCertFile          string        `env:"GRPC_SERVER_TLS_CERT_FILE" validate:"omitempty,file"`
	TLSKeyFile           string        `env:"GRPC_SERVER_TLS_KEY_FILE" validate:"omitempty,file"`
	ReadTimeout          time.Duration `env:"GRPC_SERVER_READ_TIMEOUT" validate:"gte=0"`
	WriteTimeout         time.Duration `env:"GRPC_SERVER_WRITE_TIMEOUT" validate:"gte=0"`
	EnablePrometheus     bool          `env:"GRPC_SERVER_ENABLE_PROMETHEUS" validate:"-"`
	PrometheusListenAddr string        `env:"GRPC_SERVER_PROMETHEUS_LISTEN_ADDR" validate:"required_with=EnablePrometheus,omitempty"`
}
