package kube

import (
	"bufio"
	"context"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type Service interface {
	// TailLogs returns a stream of log messages read from a container. The channel is closed once
	// the context is done or the container exists.
	TailLogs(ctx context.Context, podName string, containerName string) (chan LogMessage, error)
	PodByLabels(ctx context.Context, l map[string]string) (*corev1.PodList, error)
	CreateJob(ctx context.Context, name string, image string, labels map[string]string) error
}

type Client struct {
	clientset *kubernetes.Clientset
	namespace string
}

type LogMessage struct {
	Msg string
	Err error
}

func NewClient(c *kubernetes.Clientset, namespace string) *Client {
	return &Client{
		clientset: c,
		namespace: namespace,
	}
}

func (c Client) TailLogs(ctx context.Context, podName string, containerName string) (chan LogMessage, error) {
	req := c.clientset.CoreV1().Pods(c.namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: containerName,
		Follow:    true,
	})

	logs, err := req.Stream(ctx)
	if err != nil {
		return nil, err
	}

	var (
		stream = make(chan LogMessage)
		sc     = bufio.NewScanner(logs)
	)

	go func() {
		for sc.Scan() {
			select {
			case stream <- LogMessage{Msg: sc.Text()}:
			case <-ctx.Done():
				close(stream)
				logs.Close()
				return
			}
		}

		if sc.Err() != nil {
			stream <- LogMessage{
				Err: sc.Err(),
			}
			return
		}

		close(stream)
		logs.Close()
	}()

	return stream, nil
}

func (c Client) PodByLabels(ctx context.Context, l map[string]string) (*corev1.PodList, error) {
	list, err := c.clientset.CoreV1().Pods(c.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(l).String(),
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c Client) CreateJob(
	ctx context.Context,
	name string,
	image string,
	labels map[string]string,
) error {
	spec := &batchv1.Job{
		Name:      name,
		Namespace: c.namespace,
		Labels:    labels,
		Spec: batchv1.JobSpec{
			BackoffLimit: new(int32(0)),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Name:   name,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mc",
							Image: image,
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			TTLSecondsAfterFinished: new(int32(60)),
		},
	}

	_, err := c.clientset.BatchV1().Jobs(c.namespace).Create(ctx, spec, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
