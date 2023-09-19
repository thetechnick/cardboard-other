package main

import (
	"fmt"
	"path/filepath"

	cardboardv1alpha1 "cardboard.package-operator.run/api/v1alpha1"
	"cardboard.package-operator.run/pkg"
	"github.com/pterm/pterm"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	project, err := readProjectFile("p.yml")
	if err != nil {
		return err
	}

	return installGoTools(project)
}

func installGoTools(project *cardboardv1alpha1.Project) error {

	depsPath, err := filepath.Abs(filepath.Join(".cache", "deps"))
	if err != nil {
		return err
	}
	deps := pkg.DependencyDirectory(depsPath)

	var goToolNames []string
	for _, goTool := range project.Spec.GoTools {
		goTool := goTool
		rebuild, err := deps.NeedsRebuild(goTool.Name, goTool.Version)
		if err != nil {
			return err
		}
		if rebuild {
			goToolNames = append(goToolNames, goTool.Name)
		}
	}

	if len(goToolNames) == 0 {
		pterm.Success.Println("Go tools up-to-date")
		return nil
	}

	p, err := pterm.DefaultProgressbar.WithTotal(len(goToolNames)).WithTitle("Installing Tools").Start()
	if err != nil {
		return err
	}

	for _, goTool := range project.Spec.GoTools {
		p.UpdateTitle(fmt.Sprintf("Installing %s@%s", goTool.Name, goTool.Version))

		if err := deps.GoInstall(
			goTool.Name, goTool.URL, goTool.Version,
		); err != nil {
			return fmt.Errorf("installing %s: %w", goTool.Name, err)
		}
		pterm.Success.Println("Installed " + goTool.Name)
		p.Increment()
	}

	return nil
}
