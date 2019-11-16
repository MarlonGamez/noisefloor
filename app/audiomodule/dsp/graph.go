package dsp

import (
	"encoding/xml"
	"io"

	"github.com/jacoblister/noisefloor/app/audiomodule/dsp/processor"
	"github.com/jacoblister/noisefloor/app/audiomodule/dsp/processor/processorbuiltin"
)

// Graph is a graph of processors and connectors, plus exported parameter map
type Graph struct {
	Name       string
	Processors []ProcessorDefinition
	Connectors []Connector
}

func (g *Graph) inputConnectorsForProcessor(processor Processor) []*Connector {
	_, procInputs, _, _ := processor.Definition()
	connectorCount := len(procInputs)
	result := make([]*Connector, connectorCount, connectorCount)
	for i := 0; i < len(result); i++ {
		result[i] = &Connector{}
	}

	for i := 0; i < len(g.Connectors); i++ {
		if g.Connectors[i].ToProcessor == processor {
			result[g.Connectors[i].ToPort] = &g.Connectors[i]
		}
	}
	return result
}

func (g *Graph) outputConnectorsForProcessor(processor Processor) [][]*Connector {
	_, _, procOutputs, _ := processor.Definition()
	connectorCount := len(procOutputs)
	result := make([][]*Connector, connectorCount, connectorCount)
	for i := 0; i < len(result); i++ {
		result[i] = make([]*Connector, 0, 0)
	}

	for i := 0; i < len(g.Connectors); i++ {
		if g.Connectors[i].FromProcessor == processor {
			result[g.Connectors[i].FromPort] = append(result[g.Connectors[i].FromPort], &g.Connectors[i])
		}
	}
	return result
}

func exampleGraph() Graph {
	graph := Graph{}

	midiInput := processorbuiltin.MIDIInput{}
	graph.Processors = append(graph.Processors,
		ProcessorDefinition{X: 16, Y: 16, Processor: &midiInput})
	osc := processor.Oscillator{}
	graph.Processors = append(graph.Processors,
		ProcessorDefinition{X: 120, Y: 16, Processor: &osc})
	env := processor.Envelope{}
	graph.Processors = append(graph.Processors,
		ProcessorDefinition{X: 120, Y: 96, Processor: &env})
	gain := processor.Gain{}
	graph.Processors = append(graph.Processors,
		ProcessorDefinition{X: 224, Y: 16, Processor: &gain})
	outputTerminal := processorbuiltin.Terminal{}
	outputTerminal.SetParameters(true, 2)
	graph.Processors = append(graph.Processors,
		ProcessorDefinition{X: 328, Y: 16, Processor: &outputTerminal})
	scope := processor.Scope{Trigger: true, Skip: 4}
	graph.Processors = append(graph.Processors,
		ProcessorDefinition{X: 328, Y: 96, Processor: &scope})
	scope2 := processor.Scope{Trigger: false, Skip: 200}
	graph.Processors = append(graph.Processors,
		ProcessorDefinition{X: 224, Y: 208, Name: "scope2", Processor: &scope2})

	graph.Connectors = append(graph.Connectors,
		Connector{FromProcessor: &midiInput, FromPort: 0, ToProcessor: &osc, ToPort: 0})
	graph.Connectors = append(graph.Connectors,
		Connector{FromProcessor: &midiInput, FromPort: 1, ToProcessor: &env, ToPort: 0})
	graph.Connectors = append(graph.Connectors,
		Connector{FromProcessor: &midiInput, FromPort: 2, ToProcessor: &env, ToPort: 1})

	graph.Connectors = append(graph.Connectors,
		Connector{FromProcessor: &osc, FromPort: 0, ToProcessor: &gain, ToPort: 0})
	graph.Connectors = append(graph.Connectors,
		Connector{FromProcessor: &env, FromPort: 0, ToProcessor: &gain, ToPort: 1})
	graph.Connectors = append(graph.Connectors,
		Connector{FromProcessor: &gain, FromPort: 0, ToProcessor: &outputTerminal, ToPort: 0})
	graph.Connectors = append(graph.Connectors,
		Connector{FromProcessor: &gain, FromPort: 0, ToProcessor: &outputTerminal, ToPort: 1})
	graph.Connectors = append(graph.Connectors,
		Connector{FromProcessor: &gain, FromPort: 0, ToProcessor: &scope, ToPort: 0})
	graph.Connectors = append(graph.Connectors,
		Connector{FromProcessor: &env, FromPort: 0, ToProcessor: &scope2, ToPort: 0})

	return graph
}

// loadProcessorGraph loads a procesor graph from file
func loadProcessorGraph(filename string) Graph {
	// just sets up a static graph for now

	return exampleGraph()
}

// saveProcessorGraph saves the graph to the provided writer
func saveProcessorGraph(graph Graph, writer io.Writer) {
	xml, _ := xml.MarshalIndent(graph, "", "   ")
	println(string(xml))

	println("save graph")
}
