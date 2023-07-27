package templates

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"text/template"
)

var (
	//go:embed assets/deployment.yaml
	deploymentTpl string
)

type Renderoptions struct {
	Group     string
	Version   string
	Resource  string
	Namespace string
	Name      string
	Tag       string
}

func Values(opts Renderoptions) map[string]string {
	if len(opts.Name) == 0 {
		opts.Name = fmt.Sprintf("%s-controller", opts.Resource)
	}

	if len(opts.Namespace) == 0 {
		opts.Namespace = "default"
	}

	if len(opts.Tag) == 0 {
		opts.Tag = os.Getenv("COMPOSITION_CONTROLLER_IMAGE_TAG")
	}

	return map[string]string{
		"apiGroup":   opts.Group,
		"apiVersion": opts.Version,
		"resource":   opts.Resource,
		"name":       opts.Name,
		"namespace":  opts.Namespace,
		"tag":        opts.Tag,
	}
}

func RenderDeployment(values map[string]string) ([]byte, error) {
	tpl, err := template.New("deployment").Funcs(TxtFuncMap()).Parse(deploymentTpl)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	if err := tpl.Execute(&buf, values); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
