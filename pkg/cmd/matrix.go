/*
Copyright 2022 Arda Güçlü.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"

	_ "k8s.io/client-go/plugin/pkg/client/auth" // combined authprovider import
)

const (
	example = `
	# show pods statuses in default namespace per each node
	kubectl pod-node-matrix

	# show pods statuses in given namespace per each node
	kubectl pod-node-matrix -n test-namespace

	# show pods statuses in all namespaces per each node
	kubectl pod-node-matrix -A

	# show pods statuses in verbose format
	# it is suggested only there are few nodes.
	kubectl pod-node-matrix --verbose
`
	longDesc = `
	Shows pods statuses and nodes in table view to
	mitigate troubleshooting effort.
	
	When there is a problem at pods, detecting whether the problem
	stems from nodes or deployments is arduous task.
	
	pod-node-matrix shows pods and nodes in table view.
	If one of columns is red, thereby it is reasonable to interfere
	that node has a problem. If one of rows is red,
	node might not be the cause instead it is better to check
	pods, deployments, etc.
`
)

// PodNodeMatrixOptions provides information required
// to show the pod node matrix.
type PodNodeMatrixOptions struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams

	namespace     string
	allNamespaces bool
	verbose       bool
	kubeClient    kubernetes.Interface
}

// NewPodNodeMatrixOptions provides an instance of PodNodeMatrixOptions with default values.
func NewPodNodeMatrixOptions(streams genericclioptions.IOStreams) *PodNodeMatrixOptions {
	return &PodNodeMatrixOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

// NewCmdPodNodeMatrix provides a cobra command wrapping PodNodeMatrixOptions
func NewCmdPodNodeMatrix(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewPodNodeMatrixOptions(streams)

	cmd := &cobra.Command{
		Use:          "pod-node-matrix [flags]",
		Short:        "show [pod status x node] matrix in table view",
		Long:         longDesc,
		SilenceUsage: true,
		Example:      example,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&o.verbose, "verbose", o.verbose, "If present, show pods statuses and node names in verbose format instead shortcuts. It is suggested only when you have few nodes.")
	cmd.Flags().BoolVarP(&o.allNamespaces, "all-namespaces", "A", o.allNamespaces, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	o.configFlags.AddFlags(cmd.Flags())
	return cmd
}

// Complete sets all information required to show pod node matrix
func (o *PodNodeMatrixOptions) Complete() error {
	config, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	o.kubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	if o.allNamespaces {
		o.namespace = ""
	} else {
		namespace, _, err := o.configFlags.ToRawKubeConfigLoader().Namespace()
		if err != nil || len(namespace) == 0 {
			namespace = "default"
		}

		o.namespace = namespace
	}

	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *PodNodeMatrixOptions) Validate(args []string) error {
	if len(args) > 0 {
		return errors.New("unknown argument")
	}

	return nil
}

// nodeWrap is Node object with index value in addition to that.
type nodeWrap struct {
	v1.Node
	index int
}

// Run shows pod-node matrix according to the given parameters
func (o *PodNodeMatrixOptions) Run() error {
	pods, err := o.kubeClient.CoreV1().Pods(o.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(pods.Items) == 0 {
		return fmt.Errorf("no pod is found in namespace %s", o.namespace)
	}

	nodes, err := o.kubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(nodes.Items) == 0 {
		return errors.New("no node is found which is weird")
	}

	indx := 0
	nodeMap := make(map[string]nodeWrap)
	var nodeNames []string
	for _, n := range nodes.Items {
		nodeMap[n.Name] = nodeWrap{
			Node:  n,
			index: indx,
		}
		nodeNames = append(nodeNames, n.Name)
		indx += 1
	}

	err = o.Print(pods.Items, nodeMap, nodeNames)
	if err != nil {
		return err
	}
	return nil
}

// Print prints table view with colors to emphasize results
func (o *PodNodeMatrixOptions) Print(pods []v1.Pod, nm map[string]nodeWrap, nodeNames []string) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, pod := range pods {
		data := make([]string, len(nm)+1)
		colors := make([]tablewriter.Colors, len(nm)+1)
		data[0] = pod.Name
		colors[0] = tablewriter.Colors{}

		if val, ok := nm[pod.Spec.NodeName]; ok {
			data[val.index+1] = string(pod.Status.Phase)

			switch pod.Status.Phase {
			case v1.PodPending:
				if !o.verbose {
					data[val.index+1] = "-"
				}
				colors[val.index+1] = tablewriter.Colors{tablewriter.FgYellowColor}
			case v1.PodRunning:
				fallthrough
			case v1.PodSucceeded:
				if !o.verbose {
					data[val.index+1] = "✓"
				}
				colors[val.index+1] = tablewriter.Colors{tablewriter.FgGreenColor}
			case v1.PodFailed:
				if !o.verbose {
					data[val.index+1] = "x"
				}
				colors[val.index+1] = tablewriter.Colors{tablewriter.FgRedColor}
			}
		} else {
			colors[0] = tablewriter.Colors{tablewriter.BgRedColor}
		}

		table.Rich(data, colors)
	}

	nodeNames, footers := o.getHeadersAndFooters(nodeNames)
	header := []string{"Pods"}
	header = append(header, nodeNames...)
	table.SetHeader(header)
	table.Render()

	if !o.verbose && footers != nil {
		nodeMapping := tablewriter.NewWriter(os.Stdout)
		nodeMapping.SetAlignment(tablewriter.ALIGN_LEFT)
		nodeMapping.AppendBulk(footers)
		nodeMapping.Render()
	}
	return nil
}

func (o *PodNodeMatrixOptions) getHeadersAndFooters(nodeNames []string) (headers []string, footers [][]string) {
	if o.verbose {
		return nodeNames, nil
	}

	for i, val := range nodeNames {
		ind := strconv.Itoa(i)
		headers = append(headers, ind)
		footers = append(footers, []string{ind, val})
	}

	return
}
