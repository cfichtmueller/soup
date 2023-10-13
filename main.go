// Copyright 2023 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package soup

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Selector struct {
	// Selects an element with a given id. Takes precedence over ClassName
	Id string
	// Selects an element with a given class. Takes precedence over Tag
	ClassName string
	// Selects an element with a given tag
	Tag string
	// Perform a recursive search. That is, include the node's children in the search.
	Recursive bool
}

type Node struct {
	backing *html.Node
}

func newNode(b *html.Node) *Node {
	return &Node{backing: b}
}

func newNodes(b []*html.Node) []*Node {
	r := make([]*Node, 0, len(b))
	for _, n := range b {
		r = append(r, newNode(n))
	}
	return r
}

// AllWithClassName returns all child nodes that have the given class name.
func (n *Node) AllWithClassName(className string) []*Node {
	return newNodes(AllWithClassName(n.backing, className))
}

// AllWithClassNameR is the recursive variant of AllWithClassName.
func (n *Node) AllWithClassNameR(className string) []*Node {
	return newNodes(AllWithClassNameR(n.backing, className))
}

// AllWithTag returns all child nodes with the given tag.
func (n *Node) AllWithTag(tagName string) []*Node {
	return newNodes(AllWithTag(n.backing, tagName))
}

// AllWithTagR is the recursive variant of AllWithTag
func (n *Node) AllWithTagR(tagName string) []*Node {
	return newNodes(AllWithTagR(n.backing, tagName))
}

// Attr returns the attribute value or an empty string if the attribute isn't found.
func (n *Node) Attr(attr string) string {
	return Attr(n.backing, attr)
}

// FirstWithClassName returns the first child with the given class.
func (n *Node) FirstWithClassName(className string) *Node {
	res := FirstWithClassName(n.backing, className)
	if res != nil {
		return newNode(res)
	}
	return nil
}

// FirstWithClassNameR is the recursive variant of FirstWithClassName.
func (n *Node) FirstWithClassNameR(className string) *Node {
	res := FirstWithClassNameR(n.backing, className)
	if res != nil {
		return newNode(res)
	}
	return nil
}

// FirstWithId returns the first child with the given id.
func (n *Node) FirstWithId(id string) *Node {
	res := FirstWithId(n.backing, id)
	if res != nil {
		return newNode(res)
	}
	return nil
}

// FirstWithIdR is the recursive variant of FirstWithId.
func (n *Node) FirstWithIdR(id string) *Node {
	res := FirstWithIdR(n.backing, id)
	if res != nil {
		return newNode(res)
	}
	return nil
}

// FirstWithTag returns the first child node with the given tag.
func (n *Node) FirstWithTag(tag string) *Node {
	res := FirstWithTag(n.backing, tag)
	if res != nil {
		return newNode(res)
	}
	return nil
}

// FirstWithTagR is the recursive variant of FirstWithTag.
func (n *Node) FirstWithTagR(tag string) *Node {
	res := FirstWithTagR(n.backing, tag)
	if res != nil {
		return newNode(res)
	}
	return nil
}

// HasClass returns true if the node has the given class
func (n *Node) HasClass(className string) bool {
	return HasClass(n.backing, className)
}

// SelectAll selects all child node that match the given Selector.
func (n *Node) SelectAll(selector Selector) []*Node {
	res := SelectAll(n.backing, selector)
	if res != nil {
		return newNodes(res)
	}
	return nil
}

// SelectFirst selects the first child node that matches the given selector.
func (n *Node) SelectFirst(selector Selector) *Node {
	res := SelectFirst(n.backing, selector)
	if res != nil {
		return newNode(res)
	}
	return nil
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.backing.Data)
}

// TextContent returns the text content of the node
func (n *Node) TextContent() string {
	return TextContent(n.backing)
}

// Parse parses a node from a reader
func Parse(r io.Reader) (*Node, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return newNode(root), nil
}

// SelectAll selects all child node that match the given Selector
func SelectAll(node *html.Node, selector Selector) []*html.Node {
	if len(selector.Id) > 0 {
		node := firstWithId(node, selector.Id, selector.Recursive)
		if node == nil {
			return []*html.Node{node}
		}
		return []*html.Node{}
	}
	if len(selector.ClassName) > 0 {
		return allWithClassName(node, selector.ClassName, selector.Recursive)
	}
	if len(selector.Tag) > 0 {
		return allWithTag(node, selector.Tag, selector.Recursive)
	}
	return nil
}

// SelectFirst selects the first child node that matches the given selector
func SelectFirst(node *html.Node, selector Selector) *html.Node {
	if len(selector.Id) > 0 {
		return firstWithId(node, selector.Id, selector.Recursive)
	}
	if len(selector.ClassName) > 0 {
		return firstWithClassName(node, selector.ClassName, selector.Recursive)
	}
	if len(selector.Tag) > 0 {
		return firstWithTag(node, selector.Tag, selector.Recursive)
	}
	return nil
}

// FirstWithId returns the first child with the given id.
func FirstWithId(node *html.Node, id string) *html.Node {
	return firstWithId(node, id, false)
}

// FirstWithIdR is the recursive variant of FirstWithId
func FirstWithIdR(node *html.Node, id string) *html.Node {
	return firstWithId(node, id, true)
}

func firstWithId(node *html.Node, id string, recursive bool) *html.Node {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && Attr(c, "id") == id {
			return c
		}
	}
	if !recursive {
		return nil
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if el := firstWithId(c, id, true); el != nil {
				return el
			}
		}
	}
	return nil
}

// FirstWithClassName returns the first child with the given class.
func FirstWithClassName(node *html.Node, className string) *html.Node {
	return firstWithClassName(node, className, false)
}

// FirstWithClassNameR is the recursive variant of FirstWithClassName.
func FirstWithClassNameR(node *html.Node, className string) *html.Node {
	return firstWithClassName(node, className, true)
}

func firstWithClassName(node *html.Node, className string, recursive bool) *html.Node {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if HasClass(c, className) {
			return c
		}
	}
	if !recursive {
		return nil
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if n := firstWithClassName(c, className, true); n != nil {
			return n
		}
	}
	return nil
}

// AllWithClassName returns all children that have the given class name.
func AllWithClassName(node *html.Node, className string) []*html.Node {
	return allWithClassName(node, className, false)
}

// AllWithClassNameR is the recursive variant of AllWithClassName
func AllWithClassNameR(node *html.Node, className string) []*html.Node {
	return allWithClassName(node, className, true)
}

func allWithClassName(node *html.Node, className string, recursive bool) []*html.Node {
	res := make([]*html.Node, 0)
	if HasClass(node, className) {
		res = append(res, node)
	}
	if !recursive {
		return res
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		res = append(res, allWithClassName(c, className, true)...)
	}
	return res
}

// FirstWithTag returns the first child with the given tag name
func FirstWithTag(node *html.Node, tagName string) *html.Node {
	return firstWithTag(node, tagName, false)
}

// FirstWithTagR is the recursive variant of FirstWithTag
func FirstWithTagR(node *html.Node, tagName string) *html.Node {
	return firstWithTag(node, tagName, true)
}

func firstWithTag(node *html.Node, tagName string, recursive bool) *html.Node {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tagName {
			return c
		}
	}
	if !recursive {
		return nil
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if n := firstWithTag(c, tagName, true); n != nil {
			return n
		}
	}
	return nil
}

// AllWithTag returns all children with the given tag.
func AllWithTag(node *html.Node, tagName string) []*html.Node {
	return allWithTag(node, tagName, false)
}

// AllWithTagR is the recursive variant of AllWithTag
func AllWithTagR(node *html.Node, tagName string) []*html.Node {
	return allWithTag(node, tagName, true)
}

func allWithTag(node *html.Node, tagName string, recursive bool) []*html.Node {
	res := make([]*html.Node, 0)
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tagName {
			res = append(res, c)
		}
		if recursive {
			res = append(res, allWithTag(c, tagName, true)...)
		}
	}
	return res
}

func HasClass(node *html.Node, className string) bool {
	for _, a := range node.Attr {
		if a.Key == "class" {
			for _, class := range strings.Split(a.Val, " ") {
				if class == className {
					return true
				}
			}
			return false
		}
	}
	return false
}

func TextContent(node *html.Node) string {
	if node.Type == html.TextNode {
		return strings.Trim(node.Data, " \t\n\r")
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			return strings.Trim(c.Data, " \t\n\n")
		}
	}
	return ""
}

// Attr returns the attribute value or an empty string if the attribute isn't found.
func Attr(node *html.Node, attr string) string {
	for _, a := range node.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return ""
}
