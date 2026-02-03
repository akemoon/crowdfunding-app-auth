package metrics

import "github.com/prometheus/client_golang/prometheus"

type AuthMetrics struct {
	AuthSignInTotal *prometheus.CounterVec
}

func NewAuthMetrics(reg prometheus.Registerer) *AuthMetrics {
	authSignInTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_signin_total",
			Help: "Total number of signin attempts",
		},
		[]string{"result"},
	)

	reg.MustRegister(authSignInTotal)

	return &AuthMetrics{
		AuthSignInTotal: authSignInTotal,
	}
}
