package kube

import (
	"encoding/json"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterSet map[string]*kubernetes.Clientset

func (s *ClusterSet) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("expected an object of cluster name -> kubeconfig string: %w", err)
	}

	out := make(ClusterSet, len(raw))
	for name, kubeconfig := range raw {
		if name == "" {
			return fmt.Errorf("cluster name must not be empty")
		}

		if kubeconfig == "" {
			return fmt.Errorf("ClusterSet: cluster %q: kubeconfig must not be empty", name)
		}

		cfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
		if err != nil {
			return fmt.Errorf("cluster %s: parsing kubeconfig: %w", name, err)
		}

		cs, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			return fmt.Errorf("cluster %s: creating clientset: %w", name, err)
		}
		out[name] = cs
	}

	*s = out
	return nil
}
