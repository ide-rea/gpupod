package main

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestSortPod(t *testing.T) {
	var tests = []struct {
		in       SortablePodInfos
		expected SortablePodInfos
	}{
		{
			in: SortablePodInfos{
				{namespace: "ns1", node: "n1", name: "pod3", status: corev1.PodFailed},
				{namespace: "ns1", node: "n1", name: "pod4", status: corev1.PodPending},
				{namespace: "ns1", node: "n1", name: "pod1", status: corev1.PodRunning},
				{namespace: "ns1", node: "n1", name: "pod2", status: corev1.PodSucceeded},
			},
			expected: SortablePodInfos{
				{namespace: "ns1", node: "n1", name: "pod1", status: corev1.PodRunning},
				{namespace: "ns1", node: "n1", name: "pod4", status: corev1.PodPending},
				{namespace: "ns1", node: "n1", name: "pod3", status: corev1.PodFailed},
				{namespace: "ns1", node: "n1", name: "pod2", status: corev1.PodSucceeded},
			},
		},
		{
			in: SortablePodInfos{
				{namespace: "ns4", node: "n1", name: "pod4", status: corev1.PodRunning},
				{namespace: "ns3", node: "n1", name: "pod3", status: corev1.PodRunning},
				{namespace: "ns2", node: "n1", name: "pod2", status: corev1.PodRunning},
				{namespace: "ns1", node: "n1", name: "pod1", status: corev1.PodRunning},
			},
			expected: SortablePodInfos{
				{namespace: "ns1", node: "n1", name: "pod1", status: corev1.PodRunning},
				{namespace: "ns2", node: "n1", name: "pod2", status: corev1.PodRunning},
				{namespace: "ns3", node: "n1", name: "pod3", status: corev1.PodRunning},
				{namespace: "ns4", node: "n1", name: "pod4", status: corev1.PodRunning},
			},
		},
	}
	for _, sample := range tests {
		sort.Sort(sample.in)
		sample.in.SortByFieldName("namespace")
		if !reflect.DeepEqual(sample.in, sample.expected) {
			t.Errorf("expected:\n %s \n got:\n %s", sample.expected.Marshal(), sample.in.Marshal())
		}
	}
}

func (t SortablePodInfos) Marshal() string {
	sb := strings.Builder{}
	for _, pod := range t {
		tp := reflect.TypeOf(pod)
		v := reflect.ValueOf(pod)
		for i := 0; i < tp.NumField(); i++ {
			sb.WriteString(tp.Field(i).Name)
			sb.WriteString("=")
			sb.WriteString(v.Field(i).String())
			if i+1 != tp.NumField() {
				sb.WriteString(",")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
