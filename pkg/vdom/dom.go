package vdom

var dom = Element{Type: Root}           // dom is the current vdom as written to the dom
var svgNamespace bool                   // svgNamespace indicates an SVG vs HTML document
var headerElements []Element            // headerElements is elements list for the head section of the document
var rootComponent Component             // rootComponent is the root of the component tree
var componentMap map[Component]*Element // elementMap links active components to elements
var componentUpdate chan Component      // componentUpdate is the background component update channel

//SetSVGNamespace set the DOM namespace to SVG (default is HTML)
func SetSVGNamespace() {
	svgNamespace = true
}

//SetHeaderElements sets up links for the <head> section of the HTML document
func SetHeaderElements(elements []Element) {
	headerElements = elements
}

// SetDomRootElement sets the root DOM element
func SetDomRootElement(element *Element) {
	dom = *element
}

// RenderComponentToDom renders a VDOM component
func RenderComponentToDom(component Component) {
	rootComponent = component
	rootElement := component.Render()

	// componentMap = map[Component]*Element{}
	updateDomTreeRecursive(&rootElement, []int{})
	updateDomTreeRecursive(&rootElement, []int{})
	SetDomRootElement(&rootElement)
}

// UpdateComponent is called when a state change in a component occurs
func UpdateComponent(component Component) {
}

// UpdateComponentBackground allows a background process to
// notify a state change in a component
func UpdateComponentBackground(component Component) {
	componentUpdate <- component
}

// updateDomTreeRecursive updates the dom element path and componenet map for the whole tree
func updateDomTreeRecursive(element *Element, path []int) {
	element.Path = make([]int, len(path))
	copy(element.Path, path)

	// TODO - handle componentmap update
	// if element.Component != nil {
	// 	componentMap[element.Component] = element
	// }

	for i := 0; i < len(element.Children); i++ {
		childPath := append(path, i)
		updateDomTreeRecursive(&element.Children[i], childPath)
	}
}

// updateDomBegin notifies a DOM update cycle is starting
func updateDomBegin() {
}

// updateDomEnd notifies a DOM update cycle has ended,
// and returns a patch of DOM changes for the update cycle
func updateDomEnd() PatchList {
	newDom := rootComponent.Render()
	updateDomTreeRecursive(&newDom, []int{})

	// patchList := PatchList{SVGNamespace: svgNamespace, Patch: []Patch{Patch{Type: Replace, Path: []int{}, Element: newDom}}}
	patchList := diffElementTrees(&dom, &newDom)
	dom = newDom
	return patchList
}

// fullDomPatch returns a patch to fully populate the DOM
func fullDomPatch() PatchList {
	patchElements := []Patch{}

	if headerElements != nil {
		patchElements = append(patchElements, Patch{Type: Header, Path: []int{}, Element: Element{Children: headerElements}})
	}
	patchElements = append(patchElements, Patch{Type: Replace, Path: []int{}, Element: dom})

	patchList := PatchList{SVGNamespace: svgNamespace, Patch: patchElements}

	return patchList
}
