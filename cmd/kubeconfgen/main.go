package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/certificates"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
	"github.com/spf13/pflag"

	apiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/yaml"
)

func writeConfig(clusterName string, kubeConfig []byte) (string, error) {
	confFilename := clusterName + ".config"
	if err := ioutil.WriteFile(confFilename, kubeConfig, 0644); err != nil {
		return "", errors.New("failed to write file")
	}

	cwd := getCWD()

	log.Println("Wrote config file to" + cwd + confFilename)
	return confFilename, nil
}

func getCWD() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	cwd := filepath.Dir(exe)

	return cwd
}

func generateConfig(url, domain, clusterName, k8sAPIAddress, CACertificatePEM string) (kj []byte) {
	user := "user"
	contextName := user + "@" + clusterName

	namedCluster := apiv1.NamedCluster{
		Name: clusterName,
		Cluster: apiv1.Cluster{
			Server:                   k8sAPIAddress,
			CertificateAuthorityData: []byte(CACertificatePEM),
		},
	}

	namedContext := apiv1.NamedContext{
		Name: contextName,
		Context: apiv1.Context{
			Cluster:  clusterName,
			AuthInfo: user,
		},
	}

	namedAuthInfo := apiv1.NamedAuthInfo{
		Name: user,
		AuthInfo: apiv1.AuthInfo{
			Exec: &apiv1.ExecConfig{
				Command: "client-keystone-auth",
				Args: []string{
					"--domain-name=" + domain,
					"--keystone-url=" + url,
				},
				APIVersion: "client.authentication.k8s.io/v1beta1",
			},
		},
	}

	kubeConfig := &apiv1.Config{
		Kind:        "Config",
		APIVersion:  "v1",
		Preferences: apiv1.Preferences{},
		Clusters: []apiv1.NamedCluster{
			namedCluster,
		},
		AuthInfos: []apiv1.NamedAuthInfo{
			namedAuthInfo,
		},
		Contexts: []apiv1.NamedContext{
			namedContext,
		},
		CurrentContext: contextName,
	}

	kubeConfigYAML, err := yaml.Marshal(kubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	return kubeConfigYAML
}

// Version information
var (
	GitCommit string
	BuildDate string
	GoVersion string
	OperatingSystem string
	Architecture string
)

func main() {
	var (
		url         string
		domain      string
		user        string
		projectName string
		projectID   string
		password    string
		clusterName string
		version 	bool
	)

	pflag.StringVar(&url, "keystone-url", os.Getenv("OS_AUTH_URL"), "URL for the OpenStack Keystone API")
	pflag.StringVar(&domain, "domain-name", os.Getenv("OS_DOMAIN_NAME"), "Keystone domain name")
	pflag.StringVar(&user, "user-name", os.Getenv("OS_USERNAME"), "User name")
	pflag.StringVar(&projectName, "project-name", os.Getenv("OS_PROJECT_NAME"), "Keystone project name")
	pflag.StringVar(&projectID, "project-id", os.Getenv("OS_PROJECT_ID"), "Keystone project ID")
	pflag.StringVar(&password, "password", "*****", "Password")
	pflag.StringVar(&clusterName, "cluster-name", "", "Magnum cluster name")
	pflag.BoolVar(&version, "version", false, "Display version information")
	pflag.Parse()

	if version {
		fmt.Print("kubegenconf:")
		fmt.Printf("\n\tVersion:\t%s\n", "N/A")
		fmt.Printf("\tGo version:\t%s\n", GoVersion)
		fmt.Printf("\tGit commit:\t%s\n", GitCommit)
		fmt.Printf("\tBuilt:\t\t%s\n", BuildDate)
		fmt.Printf("\tOS/Arch:\t%s/%s\n", OperatingSystem, Architecture)
		os.Exit(0)
	}

	if password == "*****" {
		password = os.Getenv("OS_PASSWORD")
	}

	if clusterName == "" {
		log.Println("Please provide a cluster name!")
		os.Exit(1)
	}

	if domain == "" {
		domain = "default"
	}

	options := gophercloud.AuthOptions{
		IdentityEndpoint: url,
		Username:         user,
		TenantName:       projectName,
		Password:         password,
		DomainName:       domain,
	}

	provider, err := openstack.AuthenticatedClient(options)
	if err != nil {
		log.Fatal(err)
	}

	client, err := openstack.NewContainerInfraV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatal(err)
	}

	cluster, err := clusters.Get(client, clusterName).Extract()
	if err != nil {
		log.Fatal(err)
	}

	CACertificate, err := certificates.Get(client, cluster.UUID).Extract()
	if err != nil {
		log.Fatal(err)
	}
	CACertificatePEM := CACertificate.PEM

	kubeConfig := generateConfig(url, domain, clusterName, cluster.APIAddress, CACertificatePEM)

	filename, err := writeConfig(clusterName, kubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	cwd := getCWD()
	fmt.Println("export KUBECONFIG=" + cwd + "/" + filename)
}
