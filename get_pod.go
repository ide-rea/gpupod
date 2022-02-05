package main

import (
	"context"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/liggitt/tabwriter"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgwatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	tabwriterMinWidth = 6
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
	tabwriterFlags    = tabwriter.RememberWidths
)

var PodPhaseToRank = map[corev1.PodPhase]uint{
	corev1.PodRunning:   1,
	corev1.PodPending:   2,
	corev1.PodFailed:    3,
	corev1.PodSucceeded: 4,
	corev1.PodUnknown:   5,
}

type PodInfo struct {
	namespace   string
	node        string
	name        string
	image       string
	gpuUsage    string
	status      corev1.PodPhase
	createdTime time.Time
	reason      string
}

type SortablePodInfos []*PodInfo

func listWatchGpuPod() {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	config.Timeout = time.Second * time.Duration(timeout)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	podList, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	var podInfos SortablePodInfos
	for _, pod := range podList.Items {
		if !successPod && pod.Status.Phase == corev1.PodSucceeded {
			continue
		}
		if podInfo := extractPodInfo(&pod); podInfo != nil {
			podInfos = append(podInfos, podInfo)
		}
	}
	// First sort by pod status, pod rank sequence is Running < Pending < Failed < Succeeded < Unknown.
	sort.Sort(podInfos)
	// Second Sort by pod namespaces, in each pod group divided by pod status.
	// (ns1, Running) < (ns2, Running) < (ns1, Pending).
	podInfos.SortByFieldName("namespace")
	writer := NewWriter()
	printer := newPrinter(createdTime, reason, writer)
	printer.PrintHeader()
	for _, podInfo := range podInfos {
		printer.PrintRow(podInfo)
	}
	writer.Flush()
	accessor := meta.NewAccessor()
	if watch {
		rv := ""
		if rv, err = accessor.ResourceVersion(podList); err != nil {
			panic(err)
		}
		w, err := clientset.CoreV1().Pods("").Watch(context.TODO(), metav1.ListOptions{ResourceVersion: rv})
		if err != nil {
			panic(err)
		}
		for {
			event, ok := <-w.ResultChan()
			if !ok {
				retry := 3
				for {
					if retry == 0 {
						panic(err)
					}
					w, err = clientset.CoreV1().Pods("").Watch(context.TODO(), metav1.ListOptions{ResourceVersion: rv})
					retry--
					if apierrors.IsGone(err) || apierrors.IsResourceExpired(err) {
						rv = ""
						continue
					} else if err != nil {
						continue
					} else {
						break
					}
				}
				continue
			}
			if event.Type == pkgwatch.Error {
				continue
			}
			if event.Type == pkgwatch.Bookmark {
				rv, err = accessor.ResourceVersion(event.Object)
				continue
			}
			pod := event.Object.(*corev1.Pod)
			rv, err = accessor.ResourceVersion(pod)
			if err != nil {
				continue
			}
			if podInfo := extractPodInfo(pod); podInfo != nil {
				printer.PrintRow(podInfo)
				writer.Flush()
			}
		}
	}
}

func extractPodInfo(pod *corev1.Pod) *PodInfo {
	var gpuCards int64
	var images []string
	for _, c := range pod.Spec.Containers {
		if quantity, ok := c.Resources.Requests["nvidia.com/gpu"]; ok {
			images = append(images, c.Image)
			gpuCards += quantity.Value()
		}
	}
	if gpuCards == 0 {
		return nil
	}
	return &PodInfo{
		namespace:   pod.Namespace,
		node:        pod.Spec.NodeName,
		name:        pod.Name,
		image:       strings.Join(images, "|"),
		gpuUsage:    strconv.Itoa(int(gpuCards)),
		reason:      pod.Status.Reason,
		status:      pod.Status.Phase,
		createdTime: pod.GetObjectMeta().GetCreationTimestamp().Time,
	}
}

func (t SortablePodInfos) Less(i, j int) bool {
	if GetPodPhaseRank(t[i].status) <= GetPodPhaseRank(t[j].status) {
		return true
	}
	return false
}

func (t SortablePodInfos) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t SortablePodInfos) Len() int {
	return len(t)
}

func (t SortablePodInfos) SortByFieldName(name string) {
	if len(t) <= 1 {
		return
	}
	var s int
	status := t[0].status
	for k := 0; k < len(t); k++ {
		if t[k].status == status {
			j := k
			v := reflect.ValueOf(t[j]).Elem().FieldByName(name).String()
			for j > s {
				if v < reflect.ValueOf(t[j-1]).Elem().FieldByName(name).String() {
					t[j], t[j-1] = t[j-1], t[j]
					j--
				} else {
					break
				}
			}
		} else {
			status = t[k].status
			s = k
		}
	}
}

func newPrinter(withCreatedTime, withReason bool, writer *tabwriter.Writer) *Printer {
	p := &Printer{
		withCreatedTime: withCreatedTime,
		withReason:      withReason,
		w:               writer,
	}
	return p
}

type Printer struct {
	withCreatedTime bool
	withReason      bool
	w               *tabwriter.Writer
}

func (p *Printer) PrintRow(info *PodInfo) {
	imageFields := strings.Split(info.image, "/")
	cells := []string{
		info.namespace,
		info.node,
		imageFields[len(imageFields)-1],
		info.name,
		info.gpuUsage,
		string(info.status),
	}
	if p.withCreatedTime {
		cells = append(cells, info.createdTime.Local().String())
	}
	if p.withReason {
		cells = append(cells, info.reason)
	}
	data := strings.Join(cells, "\t") + "\n"
	p.w.Write([]byte(data))
}

func (p *Printer) PrintHeader() {
	header := []string{"namespace", "node", "image", "podName", "gpu", "status"}
	if p.withCreatedTime {
		header = append(header, "createdTime")
	}
	if p.withReason {
		header = append(header, "reason")
	}
	data := strings.Join(header, "\t") + "\n"
	p.w.Write([]byte(data))
}

func NewWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)
}

func GetPodPhaseRank(phase corev1.PodPhase) uint {
	if rank, ok := PodPhaseToRank[phase]; ok {
		return rank
	}
	return PodPhaseToRank[corev1.PodUnknown]
}
