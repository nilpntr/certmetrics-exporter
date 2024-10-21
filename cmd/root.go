package cmd

import (
	"fmt"
	"github.com/nilpntr/certmetrics-exporter/internal/handlers"
	"github.com/nilpntr/certmetrics-exporter/internal/k8s"
	"github.com/nilpntr/certmetrics-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "certmetrics-exporter",
	Short: "CertMetrics Exporter is a tool that exports(prometheus metrics) tls secrets metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		var clientSet *k8s.Client
		if viper.GetString("kube_env") == "dev" {
			_clientSet, err := k8s.NewLocalClient()
			if err != nil {
				return err
			}
			clientSet = _clientSet
		} else {
			_clientSet, err := k8s.NewInClusterClient()
			if err != nil {
				return err
			}
			clientSet = _clientSet
		}

		metricsClient := metrics.New(clientSet)

		mux := http.NewServeMux()

		registry := prometheus.NewRegistry()
		registry.MustRegister(metricsClient)
		mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
		mux.HandleFunc("/healthz", handlers.Healthz)

		zap.L().Sugar().Infof("Starting exporter on: %v", viper.GetInt("port"))

		if err := http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), mux); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)
}

func initConfig() {
	viper.SetDefault("log_level", "info")
	viper.SetDefault("port", 9106)
	viper.SetDefault("kube_env", "prod")
	viper.SetDefault("verify_cn", true)
	viper.SetDefault("refresh_interval", 30)

	viper.AutomaticEnv()
}

func initLogger() {
	var level zapcore.Level
	switch viper.GetString("log_level") {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}
	enc := zap.NewProductionEncoderConfig()
	enc.TimeKey = "timestamp"
	enc.EncodeTime = zapcore.ISO8601TimeEncoder

	zapCfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     enc,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}
	logger := zap.Must(zapCfg.Build())
	logger.Info("Logger initialized 🎉")
	zap.ReplaceGlobals(logger)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
