// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package component

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Upload(uploadURL string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<body><style>\n\t\t.upload{\n\t\t\ttext-align: center;\n\t\t\tpadding: 10px;\n            border: 1px solid #ccc;\n            margin-top: 20px;\n\t\t}\n\n\t\tform {\n\t\t\ttext-align: center;\n\t\t\tmargin-top: 20px;\n\t\t}\n\n\t\t/* Modal styles */\n        .modal {\n            display: none;\n            position: fixed;\n            top: 50%;\n            left: 50%;\n            transform: translate(-50%, -50%);\n            padding: 20px;\n            background-color: #fff;\n            border: 1px solid #ccc;\n            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);\n            z-index: 1000;\n        }\n\n\t\t.htmx-indicator{\n\t\t\tdisplay:none;\n\t\t}\n\t\t.htmx-request .htmx-indicator{\n\t\t\tdisplay:inline;\n\t\t}\n\t\t.htmx-request.htmx-indicator{\n\t\t\tdisplay:inline;\n\t\t}\n\n\t\t</style><div class=\"upload\"><form id=\"uploadForm\" hx-encoding=\"multipart/form-data\" hx-post=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(uploadURL)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `component/upload.templ`, Line: 44, Col: 78}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "\" hx-target=\".response\" hx-on::after-request=\" if(event.detail.successful) this.reset()\"><input type=\"file\" name=\"file\"> <input type=\"submit\" value=\"Upload\"></form><div class=\"response\"><!-- could it be make a modal--><progress id=\"progress\" class=\"htmx-indicator\" value=\"0\" max=\"100\"></progress></div><!-- Modal for response --><div class=\"modal\" id=\"responseModal\"><p id=\"responseText\"></p></div></div><script>\n       htmx.on('#uploadForm', 'htmx:xhr:progress', function(evt) {\n          htmx.find('#progress').setAttribute('value', evt.detail.loaded/evt.detail.total * 100)\n        });\n\n\t\t// Handle response and show modal\n        document.body.addEventListener('htmx:response', function(event) {\n            const responseText = event.detail.xhr.responseText.trim();\n            const responseModal = document.getElementById('responseModal');\n            const responseTextElement = document.getElementById('responseText');\n\n\n            responseModal.style.display = 'block';\n\n            // Close modal after 3 seconds (adjust as needed)\n            setTimeout(function() {\n                responseModal.style.display = 'none';\n            }, 3000);\n        });\t\t\n    </script></body>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
