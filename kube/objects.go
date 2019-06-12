package kube

import (
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Objects encapsulates all the objects from a Kubernetes cluster.
type Objects struct {
	Nodes                  *corev1.NodeList
	PersistentVolumes      *corev1.PersistentVolumeList
	ComponentStatuses      *corev1.ComponentStatusList
	Pods                   *corev1.PodList
	PodTemplates           *corev1.PodTemplateList
	PersistentVolumeClaims *corev1.PersistentVolumeClaimList
	ConfigMaps             *corev1.ConfigMapList
	Services               *corev1.ServiceList
	Secrets                *corev1.SecretList
	ServiceAccounts        *corev1.ServiceAccountList
	ResourceQuotas         *corev1.ResourceQuotaList
	LimitRanges            *corev1.LimitRangeList
}

// Client encapsulates a client for a Kubernetes cluster.
type Client struct {
	kubeClient kubernetes.Interface
}

// FetchObjects returns the objects from a Kubernetes cluster.
func (c *Client) FetchObjects() (*Objects, error) {
	const all = ""

	client := c.kubeClient.CoreV1()
	opts := metav1.ListOptions{}
	objects := &Objects{}

	var g errgroup.Group

	g.Go(func() (err error) {
		objects.Nodes, err = client.Nodes().List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.PersistentVolumes, err = client.PersistentVolumes().List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.ComponentStatuses, err = client.ComponentStatuses().List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.Pods, err = client.Pods(all).List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.PodTemplates, err = client.PodTemplates(all).List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.PersistentVolumeClaims, err = client.PersistentVolumeClaims(all).List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.ConfigMaps, err = client.ConfigMaps(all).List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.Secrets, err = client.Secrets(all).List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.Services, err = client.Services(all).List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.ServiceAccounts, err = client.ServiceAccounts(all).List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.ResourceQuotas, err = client.ResourceQuotas(all).List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.LimitRanges, err = client.LimitRanges(all).List(opts)
		return
	})

	err := g.Wait()
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func NewClient(configPath, configContext string) (*Client, error) {
	var config *rest.Config
	var err error

	if configContext != "" {
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: configPath},
			&clientcmd.ConfigOverrides{
				CurrentContext: configContext,
			}).ClientConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", configPath)
	}

	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		kubeClient: client,
	}, nil
}