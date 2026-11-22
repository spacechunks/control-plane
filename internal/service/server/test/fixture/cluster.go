package fixture

import (
	"time"

	"github.com/spacechunks/control-plane/internal/service/server/kube"
)

func Cluster(mod ...func(c *kube.Cluster)) kube.Cluster {
	c := kube.Cluster{
		ID:        "",
		Name:      "test-cluster",
		Slots:     10,
		CreatedAt: time.Date(2026, 2, 28, 10, 26, 0, 0, time.UTC),
		UpdateAt:  time.Date(2027, 2, 28, 10, 26, 0, 0, time.UTC),
	}

	for _, fn := range mod {
		fn(&c)
	}

	return c
}
