package sniffer

import (
	"io"

	"ksniff/kube"
	"ksniff/pkg/config"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

var defaultInterface = "ens192"

type NodeSnifferService struct {
	settings                *config.KsniffSettings
	privilegedPod           *v1.Pod
	privilegedContainerName string
	targetInterface         string
	nodeName                string
	kubernetesApiService    kube.KubernetesApiService
}

func NewNodeSnifferService(options *config.KsniffSettings, service kube.KubernetesApiService) SnifferService {
	return &NodeSnifferService{settings: options, privilegedContainerName: "node-sniff", kubernetesApiService: service, nodeName: options.DetectedPodNodeName, targetInterface: defaultInterface}
}

func (nss *NodeSnifferService) Setup() error {
	var err error

	log.Infof("creating privileged pod on node: '%s'", nss.settings.DetectedPodNodeName)
	log.Infof("creating pod with options: '%v'", nss.settings)
	log.Debug("initiating sniff on node with option: '%v'", nss)

	// TODO: Allow overload
	if nss.settings.UseDefaultImage {
		nss.settings.Image = "maintained/tcpdump"
	}

	if nss.settings.UseDefaultTCPDumpImage {
		nss.settings.TCPDumpImage = ""
	}

	nss.privilegedPod, err = nss.kubernetesApiService.CreatePrivilegedPod(
		nss.settings.DetectedPodNodeName,
		nss.privilegedContainerName,
		nss.settings.Image,
		"",
		nss.settings.UserSpecifiedPodCreateTimeout,
	)
	if err != nil {
		log.WithError(err).Errorf("failed to create privileged pod on node: '%s'", nss.nodeName)
		return err
	}

	log.Infof("pod: '%s' created successfully on node: '%s'", nss.privilegedPod.Name, nss.settings.DetectedPodNodeName)

	return nil
}

func (nss *NodeSnifferService) Cleanup() error {
	log.Infof("removing pod: '%s'", nss.privilegedPod.Name)

	err := nss.kubernetesApiService.DeletePod(nss.privilegedPod.Name)
	if err != nil {
		log.WithError(err).Errorf("failed to remove pod: '%s", nss.privilegedPod.Name)
		return err
	}

	log.Infof("pod: '%s' removed successfully", nss.privilegedPod.Name)

	return nil
}

func buildTcpdumpCommand(netInterface string, filter string, tcpdumpImage string) []string {
	return []string{"tcpdump", "-i", netInterface, "-U", "-w", "-", filter}
}

func (nss *NodeSnifferService) Start(stdOut io.Writer) error {
	log.Info("starting remote sniffing using privileged pod")

	command := buildTcpdumpCommand(nss.targetInterface, nss.settings.UserSpecifiedFilter, nss.settings.TCPDumpImage)

	exitCode, err := nss.kubernetesApiService.ExecuteCommand(nss.privilegedPod.Name, nss.privilegedContainerName, command, stdOut)
	if err != nil {
		log.WithError(err).Errorf("failed to start sniffing using privileged pod, exit code: '%d'", exitCode)
		return err
	}

	log.Info("remote sniffing using privileged pod completed")

	return nil
}
