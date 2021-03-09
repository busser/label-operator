/*
Copyright 2021.

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

package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Define utility constants for object names and testing timeouts/durations and intervals.
const (
	PodNamePrefix = "test-pod-"
	PodNamespace  = "default"

	timeout  = time.Second * 3
	duration = time.Second * 3
	interval = time.Millisecond * 250
)

var _ = Describe("Pod controller", func() {

	ctx := context.Background()

	Context("When a Pod has the annotation but not the label", func() {
		It("Should add the label", func() {
			podName := PodNamePrefix + "with-annotation-without-label"
			createPod(ctx, podName, true, false)
			validatePod(ctx, podName, true)
		})
	})

	Context("When a Pod has the annotation and the label", func() {
		It("Should keep the label", func() {
			podName := PodNamePrefix + "with-annotation-with-label"
			createPod(ctx, podName, true, true)
			validatePod(ctx, podName, true)
		})
	})

	Context("When a Pod does not have the annotation but has the label", func() {
		It("Should remove the label", func() {
			podName := PodNamePrefix + "without-annotation-with-label"
			createPod(ctx, podName, false, true)
			validatePod(ctx, podName, false)
		})
	})

	Context("When a Pod has neither the annotation nor the label", func() {
		It("Should not add the label", func() {
			podName := PodNamePrefix + "without-annotation-without-label"
			createPod(ctx, podName, false, false)
			validatePod(ctx, podName, false)
		})
	})
})

func createPod(ctx context.Context, name string, withAnnotation, withLabel bool) {
	/*
		Create a new Pod.
	*/

	By("Creating a new Pod")

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: PodNamespace,
		},
		Spec: corev1.PodSpec{
			// For simplicity, we only fill out the required fields.
			Containers: []corev1.Container{
				{
					Name:  "test-container",
					Image: "test-image",
				},
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
		},
	}

	if withAnnotation {
		pod.Annotations = map[string]string{
			addPodNameLabelAnnotation: "true",
		}
	}

	if withLabel {
		pod.Labels = map[string]string{
			podNameLabel: name,
		}
	}

	Expect(k8sClient.Create(ctx, pod)).Should(Succeed())
}

func validatePod(ctx context.Context, name string, shouldHaveLabel bool) {
	/*
		If the Pod has the annotation, make sure it eventually and
		consistently has the label.
		If the Pod does not have the annotation, make sure it
		eventually and consistently does not have the label.
	*/

	By("Checking the Pod has or does not have the label")

	podIsValid := func() bool {
		podLookupKey := types.NamespacedName{Name: name, Namespace: PodNamespace}
		var createdPod corev1.Pod

		err := k8sClient.Get(ctx, podLookupKey, &createdPod)
		if err != nil {
			return false
		}

		if shouldHaveLabel {
			return (createdPod.Labels[podNameLabel] == name) == shouldHaveLabel
		}

		_, hasLabel := createdPod.Labels[podNameLabel]
		return !hasLabel
	}

	Eventually(podIsValid, timeout, interval).Should(BeTrue())
	Consistently(podIsValid, duration, interval).Should(BeTrue())
}
