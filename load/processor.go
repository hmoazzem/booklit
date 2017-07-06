package load

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/stages"
)

type Processor struct {
	PluginFactories []booklit.PluginFactory
}

func (processor *Processor) LoadFile(path string) (*booklit.Section, error) {
	result, err := ast.ParseFile(path)
	if err != nil {
		return nil, err
	}

	return processor.loadNode(path, result.(ast.Node))
}

func (processor *Processor) LoadSource(path string, source []byte) (*booklit.Section, error) {
	result, err := ast.Parse(path, source)
	if err != nil {
		return nil, err
	}

	return processor.loadNode(path, result.(ast.Node))
}

func (processor *Processor) EvaluateSection(path string, node ast.Node) (*booklit.Section, error) {
	section := &booklit.Section{
		Path: path,

		Title: booklit.Empty,
		Body:  booklit.Empty,
	}

	plugins := []booklit.Plugin{}
	for _, pf := range processor.PluginFactories {
		plugins = append(plugins, pf.NewPlugin(section))
	}

	evaluator := &stages.Evaluate{
		Plugins: plugins,
	}

	err := node.Visit(evaluator)
	if err != nil {
		return nil, err
	}

	if evaluator.Result != nil {
		section.Body = evaluator.Result
	}

	return section, nil
}

func (processor *Processor) loadNode(path string, node ast.Node) (*booklit.Section, error) {
	section, err := processor.EvaluateSection(path, node)
	if err != nil {
		return nil, err
	}

	resolver := &stages.Resolve{
		Section: section,
	}

	err = section.Visit(resolver)
	if err != nil {
		return nil, err
	}

	return section, nil
}
