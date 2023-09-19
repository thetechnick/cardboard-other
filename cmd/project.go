package main

import (
	"os"

	"sigs.k8s.io/yaml"

	cardboardv1alpha1 "cardboard.package-operator.run/api/v1alpha1"
)

func readProjectFile(file string) (
	*cardboardv1alpha1.Project, error,
) {
	c, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	p := &cardboardv1alpha1.Project{}
	if err := yaml.Unmarshal(c, p); err != nil {
		return nil, err
	}
	return p, nil
}
