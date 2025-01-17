// +build js

package vdom

import (
	"syscall/js"
)

func addEventHandler(svgNamespace bool, element *Element, domNode js.Value, handler *EventHandler) {
	domNode.Call("addEventListener", handler.Type, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		jsEvent := args[0]
		eventData := map[string]interface{}{}
		jsEvent.Call("preventDefault")

		switch handler.Type {
		case "change":
			eventData["Value"] = jsEvent.Get("target").Get("value").String()
		case "click",
			"mousedown",
			"mouseup",
			"mouseenter",
			"mouseleave",
			"mousemove",
			"contextmenu":
			eventData["Buttons"] = jsEvent.Get("buttons").Int()
			eventData["ClientX"] = jsEvent.Get("clientX").Int()
			eventData["ClientY"] = jsEvent.Get("clientY").Int()
			if svgNamespace {
				bbox := domNode.Call("getBoundingClientRect")
				eventData["OffsetX"] = jsEvent.Get("clientX").Int() - bbox.Get("x").Int()
				eventData["OffsetY"] = jsEvent.Get("clientY").Int() - bbox.Get("y").Int()
			}
		case "keydown",
			"keyup":
			eventData["key"] = jsEvent.Get("key").String()
			eventData["KeyCode"] = jsEvent.Get("keyCode").Int()
		}

		event := Event{Type: handler.Type, Data: eventData}

		updateDomBegin()
		handler.handlerFunc(element, &event)
		patch := updateDomEnd()
		applyPatchToDom(patch)

		return nil
	}))
}

func createElementRecursive(svgNamespace bool, element *Element) js.Value {
	document := js.Global().Get("document")

	var node js.Value
	if svgNamespace == true {
		node = document.Call("createElementNS", "http://www.w3.org/2000/svg", element.Name)
	} else {
		node = document.Call("createElement", element.Name)
	}

	for name, value := range element.Attrs {
		node.Call("setAttribute", name, value)
	}

	for _, handler := range element.EventHandlers {
		addEventHandler(svgNamespace, element, node, &handler)
	}

	for _, child := range element.Children {
		switch child.Type {
		case Normal:
			childNode := createElementRecursive(svgNamespace, &child)
			node.Call("appendChild", childNode)
		case Text:
			if svgNamespace {
				node.Set("textContent", child.Attrs["Text"])
			} else {
				node.Set("innerText", child.Attrs["Text"])
			}
		}
	}

	return node
}

func getElementByPath(path []int) js.Value {
	element := js.Global().Get("document").Get("body").Get("firstElementChild")
	if !element.IsNull() {
		for i := 0; i < len(path); i++ {
			element = element.Get("children").Index(path[i])
		}
	}
	return element
}

func removeHeaderElements() {
	head := js.Global().Get("document").Get("head")
	children := head.Get("children")
	for i := children.Get("length").Int() - 1; i > 0; i-- {
		if !(children.Index(i).InstanceOf(js.Global().Get("HTMLScriptElement"))) {
			head.Call("removeChild", head, children.Index(i))
		}
	}
}

func applyPatchToDom(patchList PatchList) {
	for i := 0; i < len(patchList.Patch); i++ {
		patch := patchList.Patch[i]

		parent := js.Global().Get("document").Get("body")
		target := getElementByPath(patch.Path)
		if !target.IsNull() {
			parent = target.Get("parentElement")
		}

		switch patch.Type {
		case Header:
			removeHeaderElements()

			head := js.Global().Get("document").Get("head")
			for j := 0; j < len(patch.Element.Children); j++ {
				element := createElementRecursive(false, &patch.Element.Children[j])
				head.Call("appendChild", element)
			}
		case Replace:
			element := createElementRecursive(patchList.SVGNamespace, &patch.Element)

			if !target.IsNull() {
				parent.Call("replaceChild", element, target)
			} else {
				parent.Call("appendChild", element)
			}
		case AttrSet:
			target.Call("setAttribute", patch.Attr.Name, patch.Attr.Value)
		case AttrRemove:
			target.Call("removeAttribute", patch.Attr.Name)
		case ValueSet:
			target.Set("value", patch.Attr.Value)
		case TextSet:
			target.Set("innerText", patch.Attr.Value)
		}
	}
}

//componentUpdateListen reads and applies background component state changes
func componentUpdateListen(c chan Component) {
	for {
		component := <-c

		updateDomBegin()
		UpdateComponent(component)
		patch := updateDomEnd()
		applyPatchToDom(patch)
	}
}

//ListenAndServe starts the javascript target
func ListenAndServe(args ...interface{}) {
	componentUpdate = make(chan Component, 10)
	go componentUpdateListen(componentUpdate)

	applyPatchToDom(fullDomPatch())

	done := make(chan struct{})
	<-done
}
