// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package component

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Index(uploadURL, listViewURL string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html><head><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><!-- UIkit CSS --><link rel=\"stylesheet\" href=\"https://cdn.jsdelivr.net/npm/uikit@3.21.7/dist/css/uikit.min.css\"><!-- UIkit JS --><script src=\"https://cdn.jsdelivr.net/npm/uikit@3.21.7/dist/js/uikit.min.js\"></script><script src=\"https://cdn.jsdelivr.net/npm/uikit@3.21.7/dist/js/uikit-icons.min.js\"></script><!--HTMX --><script src=\"https://unpkg.com/htmx.org@1.9.9\"></script><!--hyperscript --><script src=\"https://unpkg.com/hyperscript.org@0.9.12\"></script><style>\n\t\t\t\tbody {\n\t\t\t\t\tbackground-color: #f4f4f4;\n\t\t\t\t\tmargin-top: 50px;\n\t\t\t\t\tjustify-content: center;\n\t\t\t\t}\n\t\t\t\th2 {\n\t\t\t\t\ttext-align: center;\n\t\t\t\t}\n\t\t\t\tli {\n\t\t\t\t\tfloat: left;\n\t\t\t\t}\n\t\t\t\tli a {\n\t\t\t\t\tdisplay: block;\n\t\t\t\t\tcolor: white;\n\t\t\t\t\ttext-align: center;\n\t\t\t\t\tpadding: 16px;\n\t\t\t\t\ttext-decoration: none;\n\t\t\t\t}\n\t\t\t\tli a:hover {\n\t\t\t\t\tbackground-color: #111111;\n\t\t\t\t}\n\t\t\t\tul {\n\t\t\t\t\tlist-style-type: none;\n\t\t\t\t\tmargin: 0;\n\t\t\t\t\tpadding: 0;\n\t\t\t\t\toverflow: hidden;\n\t\t\t\t\tbackground-color: #333333;\n\t\t\t\t}\n\t\t\t\t.upload{\n\t\t\t\t\ttext-align: center;\n\t\t\t\t\tpadding: 10px;\n\t\t\t\t\tborder: 1px solid #ccc;\n\t\t\t\t\tmargin-top: 20px;\n\t\t\t\t}\n\t\t\t\t.result{\n\t\t\t\t\tdisplay: flex;\n\t\t\t\t\tjustify-content: center;\n\t\t\t\t\tmargin-top: 20px;\n\t\t\t\t}\n\t\t\t\tform {\n\t\t\t\t\ttext-align: center;\n\t\t\t\t\tmargin-top: 20px;\n\t\t\t\t}\n\t\t\t\t/* Modal styles */\n\t\t\t\t.modal {\n\t\t\t\t\tdisplay: none;\n\t\t\t\t\tposition: fixed;\n\t\t\t\t\ttop: 50%;\n\t\t\t\t\tleft: 50%;\n\t\t\t\t\ttransform: translate(-50%, -50%);\n\t\t\t\t\tpadding: 20px;\n\t\t\t\t\tbackground-color: #fff;\n\t\t\t\t\tborder: 1px solid #ccc;\n\t\t\t\t\tbox-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);\n\t\t\t\t\tz-index: 1000;\n\t\t\t\t}\n\t\t\t\ttr.htmx-swapping td {\n\t\t\t\topacity: 0;\n\t\t\t\ttransition: opacity 1s ease-out;\n\t\t\t\t}\n\t\t\t</style><style>\n\t\t\t \t.htmx-indicator{\n\t\t\t\t\tdisplay:none;\n\t\t\t\t}\n\t\t\t\t.htmx-request .htmx-indicator{\n\t\t\t\t\tdisplay:inline;\n\t\t\t\t}\n\t\t\t\t.htmx-request.htmx-indicator{\n\t\t\t\t\tdisplay:inline;\n\t\t\t\t}\n\t\t\t</style></head><body><div class=\"uk-container\"><div><div class=\"uk-text-center uk-position-z-index\" uk-sticky=\"end: !.uk-height-large; offset: 0\"><form id=\"uploadForm\" hx-encoding=\"multipart/form-data\" hx-post=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(uploadURL)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `component/index.templ`, Line: 105, Col: 81}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#responseModal\" hx-on::after-request=\" if(event.detail.successful) this.reset()\"><input type=\"file\" name=\"file\"> <input type=\"submit\" value=\"Upload\"> <progress id=\"progress\" class=\"htmx-indicator\" value=\"0\" max=\"100\"></progress></form><!-- Modal for response --><div class=\"modal\" id=\"responseModal\"><p id=\"responseText\"></p></div></div><div><div hx-get=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(listViewURL)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `component/index.templ`, Line: 120, Col: 31}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-trigger=\"load,newObject from:body\" hx-target=\"#listBody\"></div><div id=\"listBody\"></div></div></div></div><script>\n\t\t\thtmx.on('#uploadForm', 'htmx:xhr:progress', function(evt) {\n\t\t\t\thtmx.find('#progress').setAttribute('value', evt.detail.loaded/evt.detail.total * 100)\n\t\t\t\t});\n\n\t\t\t\t// // Handle response \n\t\t\t\t// document.body.addEventListener('htmx:response', function(event) {\n\t\t\t\t// \tconst responseText = event.detail.xhr.responseText.trim();\n\t\t\t\t// \tconst responseModal = document.getElementById('responseModal');\n\t\t\t\t// \tconst responseTextElement = document.getElementById('responseText');\n\n\n\t\t\t\t// \tresponseModal.style.display = 'block';\n\n\t\t\t\t// \t// Close modal after 3 seconds (adjust as needed)\n\t\t\t\t// \tsetTimeout(function() {\n\t\t\t\t// \t\tresponseModal.style.display = 'none';\n\t\t\t\t// \t}, 3000);\n\t\t\t\t// });\t\t\n  \t\t  </script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
