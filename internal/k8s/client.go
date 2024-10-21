package k8s

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

type Client struct {
	Clientset *kubernetes.Clientset
}

func NewInClusterClient() (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &Client{Clientset: clientset}, nil
}

func NewLocalClient() (*Client, error) {
	kubeConfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from flags: %w", err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &Client{Clientset: clientSet}, nil
}

// ListSecrets gets all the secrets from the cluster and checks if the type is in allowedSecretTypes
func (c *Client) ListSecrets(ctx context.Context) ([]v1.Secret, error) {
	var secrets []v1.Secret

	namespaces, err := c.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	for _, ns := range namespaces.Items {
		secretItems, err := c.Clientset.CoreV1().Secrets(ns.Name).List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("type=%s", v1.SecretTypeTLS),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list secrets in namespace %s: %w", ns.Name, err)
		}

		for _, secretItem := range secretItems.Items {
			secrets = append(secrets, secretItem)
		}
	}

	return secrets, nil
}
