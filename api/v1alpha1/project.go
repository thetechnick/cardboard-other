package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

type ProjectGoTool struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Version string `json:"version"`
}

func init() { register(&Project{}) }
