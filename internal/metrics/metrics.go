package metrics

import (
	"context"
	"github.com/nilpntr/certmetrics-exporter/internal/k8s"
	"github.com/nilpntr/certmetrics-exporter/internal/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

const (
	collectTimeout = 10 * time.Second
	namespace      = "cert_metrics"
)

type Metrics struct {
	certMetric *prometheus.Desc
	kubeClient *k8s.Client
}

func New(kubeClient *k8s.Client) *Metrics {
	var (
		certMetric = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "cert_expiration_timestamp_seconds"),
			"The date after which the certificate expires. Expressed as a Unix Epoch Time.",
			[]string{"name", "domain", "namespace"},
			prometheus.Labels{},
		)
	)

	m := &Metrics{
		certMetric: certMetric,
		kubeClient: kubeClient,
	}

	return m
}

func (m *Metrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.certMetric
}

func (m *Metrics) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), collectTimeout)
	defer cancel()

	secrets, err := m.kubeClient.ListSecrets(ctx)
	if err != nil {
		zap.L().Sugar().Errorf("Error collecting certificates: %v", err)
		return
	}

	for _, secret := range secrets {
		if tlsCrt, ok := secret.Data["tls.crt"]; ok {
			cert, err := utils.DecodeCert(tlsCrt)
			if err != nil {
				zap.L().Sugar().Errorf("Error calculating certificate cert: %v", err)
				continue
			}

			// Make sure it also doesn't end with .svc
			if viper.GetBool("verify_cn") && !utils.IsValidDomain(cert.CommonName) {
				zap.L().Sugar().Debugf("skipping domain: %s", cert.CommonName)
				continue
			}

			ch <- prometheus.MustNewConstMetric(m.certMetric, prometheus.GaugeValue, cert.Expiry, secret.Name, cert.CommonName, secret.Namespace)
		}
	}
}
