package cmd

import (
	"ksniff/pkg/config"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"testing"
)

func TestComplete_NotEnoughArguments(t *testing.T) {
	// given
	settings := config.NewKsniffSettings(genericclioptions.IOStreams{})
	sniff := NewKsniff(settings)
	cmd := &cobra.Command{}
	var commands []string

	// when
	err := sniff.Complete(cmd, commands)

	// then
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "not enough arguments"))
}

func TestComplete_EmptyPodName(t *testing.T) {
	// given
	settings := config.NewKsniffSettings(genericclioptions.IOStreams{})
	sniff := NewKsniff(settings)
	cmd := &cobra.Command{}
	var commands []string

	// when
	err := sniff.Complete(cmd, append(commands, ""))

	// then
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "pod name is empty"))
}

// TODO: This needs to be replaced with a function that tests a valid Podname (or move API checks to outside of Complete)
// func TestComplete_PodNameSpecified(t *testing.T) {
// 	// given
// 	settings := config.NewKsniffSettings(genericclioptions.IOStreams{})
// 	sniff := NewKsniff(settings)
// 	cmd := NewCmdSniff(genericclioptions.IOStreams{})
// 	var commands []string

// 	// when
// 	err := sniff.Complete(cmd, append(commands, "pod-name"))

// 	// then
// 	assert.Nil(t, err)
// 	assert.Equal(t, "pod-name", settings.UserSpecifiedPodName)
// }

// TODO Write a test for Pod / DaemonSet / Deployment / Node and ensure resources are parsed correctly
// TODO Mock API and test that queries are fired from Validation correctly
// TODO Mock the API for Run() to ensure that Pods creation requests are pushed out
