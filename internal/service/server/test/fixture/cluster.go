package fixture

import (
	"time"

	"github.com/spacechunks/control-plane/internal/service/server/kube"
)

func Cluster(mod ...func(c *kube.Cluster)) kube.Cluster {
	c := kube.Cluster{
		ID:        "f54991c4-8441-4232-a320-3d5e36f52b2d",
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
