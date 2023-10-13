# Soup

Soup is a tiny library for working with HTML. It provides a simple API for extracting data.

## Quick Start

Install dependencies.

```go
go get github.com/cfichtmueller/soup
```

Extract some data.

```go
// Load a web page
res, err := http.Get("https://example.com")
if err != nil {
	return err
}

// Parse the page
p, err := soup.Parse(res)
if err != nil {
	return err
}

// Extract data from the page
products := p.AllWithClassNameR("product")
for _, product := range products {
	link := product.FirstWithTag("a")
	if link != nil {
        name := link.TextContent()
		url := link.Attr("href")
		fmt.Println(name, ":", url)
    }
}
```

The name is inspired by [jsoup](https://jsoup.org).
