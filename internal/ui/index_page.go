package ui

import (
	"html/template"
	"io"
)

const indexPage_template = `<!DOCTYPE html>
<html>

<head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<style>
		body {
			background-color: #f4f4f4;
			margin-top: 50px;
			justify-content: center;
		}
		
		h2 {
            text-align: center;
        }

		li {
			float: left;
		}
	
		li a {
			display: block;
			color: white;
			text-align: center;
			padding: 16px;
			text-decoration: none;
		}
	
		li a:hover {
			background-color: #111111;
		}
	
		ul {
			list-style-type: none;
			margin: 0;
			padding: 0;
			overflow: hidden;
			background-color: #333333;
		}

		
		.result{
			display: flex;
  			justify-content: center;
			margin-top: 20px;
		}

	</style>
</head>
<script src="https://unpkg.com/htmx.org@1.9.9"></script>
<body>


<h2>Menu</h2>
<p></p>

<ul>
	{{range .}}
		<li><a hx-get={{.Endpoint}}  hx-target=".result" hx-trigger="click" hx-swap="innerHTML" >{{.Name}}</a></li>
	{{end}}
	<li><a href="https://github.com/odit-bit/cloudfs" target="_blank" rel="noopener noreferrer" >Github</a></li>
	
</ul>

<div class="result"></div>


</body>
</html>
`

var indexPage = template.Must(template.New("index.html").Parse(indexPage_template))

func RenderIndexPage(w io.Writer, data []Menu) error {
	return indexPage.Execute(w, data)
}

type Menu struct {
	Name     string
	Endpoint string
}

type IndexCreator struct {
	menu []Menu
}

func NewIndexPage() *IndexCreator {
	idx := IndexCreator{
		menu: []Menu{},
	}
	return &idx
}

func (idx *IndexCreator) Render(w io.Writer) error {
	return RenderIndexPage(w, idx.menu)
}

func (idx *IndexCreator) AddMenu(name string, href string) {
	it := Menu{
		Name:     name,
		Endpoint: href,
	}
	idx.menu = append(idx.menu, it)
}

func (idx *IndexCreator) Value() []Menu {
	return idx.menu
}
