package ui

import (
	"html/template"
	"io"
)

// const listPage_template = `
// <head>
// 	<title>List File</title>
// </head>

// <div>
// 	<table>
// 	<!--
// 		set the last iterate value as marker to send next n list
// 	-->
// 		{{range .Data}}
// 			<tr>
// 				<td>{{.Name}}</td>
// 				<td>{{.Size}} byte</td>
// 				<td><a href="{{$.Endpoint}}?filename={{.Name}}" download>Download File</a></td>
// 			</tr>
// 		{{else}}
// 			<td>No data</td>
// 		{{end}}

// 	</table>
// </div>

// <script>
//         function loadMoreWithLastName() {
//             // Get the last name from the last item in the data
//             var lastName = document.querySelector('table tr:last-of-type td:first-of-type').textContent;

//             // Construct the URL with the last name
//             var url = '/list/?max=1&start=' + encodeURIComponent(lastName);

// 			// Update the hx-get attribute of the button
//             var button = document.getElementById('loadMoreButton');
//             button.setAttribute('hx-get', url);

//             // Trigger the htmx request
//             htmx.trigger(button, 'click');
//         }
// </script>

// `

const listPage_template = `
<!DOCTYPE html>
<html>
<head>
	<title>List File</title>
	
</head>


<body>
	<div hx-get="{{.}}" hx-target="#fileTable" hx-trigger="load" >
	</div>

	<div id="fileTable">
		loading....
	</div>
</body>
`

var listPage = template.Must(template.New("listPage").Parse(listPage_template))

func RenderListPage(w io.Writer, endpoint string) error {
	return listPage.Execute(w, endpoint)
}

const listDataTable_template = `

<body>
	<style>

		table {
			text-align: center;
			border-collapse: collapse;
			width: 100%;
			margin: 0 auto;
		}

		th, td {
			padding: 12px;
			text-align: left;
			border-bottom: 1px solid #DDD;
		}

		th {
			background-color: #f2f2f2;
			font-weight: bold;
			text-transform: uppercase;
		}

		tr:hover {
			background-color: #D6EEEE;
		}

		tr:nth-child(even) {
			background-color: #f2f2f2;
		}

		a {
			text-decoration: none;
			color: #0066cc;
		}

		a:hover {
			text-decoration: underline;
		}

	</style>


	<table>
		{{range .Data}}
			<tr>
				<td>{{.Name}}</td>
				<td>{{.Size}} byte</td>
				<td><a href="{{$.DownloadEndpoint}}?filename={{.Name}}" download>Download File</a></td>
			</tr>
		{{else}} 
			<td>No data</td>
		{{end}}

	</table>
</body>

`

var listDataTable = template.Must(template.New("listResult.html").Parse(listDataTable_template))

type ListData struct {
	Name      string
	Size      int64
	SharedURL string `json:"url,omitempty"`
}

type listComponent struct {
	Data             []*ListData
	DownloadEndpoint string
}

// dlEndpoint is base endpoint to download individual data , getEndpoint to get the list of data
func RenderListResult(w io.Writer, list []*ListData, dlEndpoint string) error {
	if list == nil {
		list = []*ListData{}
	}
	lc := listComponent{
		Data:             list,
		DownloadEndpoint: dlEndpoint,
	}
	return listDataTable.Execute(w, lc)
}
