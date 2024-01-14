package ui

import (
	"html/template"
	"io"
)

// Define your template content
const uploadPage_template = `



<head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>File Upload Form</title>
</head>
<body>
	<style>
		.upload{
			text-align: center;
			padding: 10px;
            border: 1px solid #ccc;
            margin-top: 20px;
		}

		form {
			text-align: center;
			margin-top: 20px;
		}

		/* Modal styles */
        .modal {
            display: none;
            position: fixed;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            padding: 20px;
            background-color: #fff;
            border: 1px solid #ccc;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            z-index: 1000;

	</style>

	<div class="upload" >
		<form id='form' hx-encoding='multipart/form-data' hx-post={{.Endpoint}} hx-target=".response">
			<input type="file" name="file" />
			<input type="submit" value="Upload" />
			<progress id='progress' value='0' max='100'></progress>
		</form>

		<div class="response">
			<!-- could it be make a modal-->
		</div>

		<!-- Modal for response -->
        <div class="modal" id="responseModal">
            <p id="responseText"></p>
        </div>

	</div>



	<script>
        htmx.on('#form', 'htmx:xhr:progress', function(evt) {
          htmx.find('#progress').setAttribute('value', evt.detail.loaded/evt.detail.total * 100)
        });

		// Handle response and show modal
        document.body.addEventListener('htmx:response', function(event) {
            const responseText = event.detail.xhr.responseText.trim();
            const responseModal = document.getElementById('responseModal');
            const responseTextElement = document.getElementById('responseText');


            responseModal.style.display = 'block';

            // Close modal after 3 seconds (adjust as needed)
            setTimeout(function() {
                responseModal.style.display = 'none';
            }, 3000);
        });
		
    </script>
</body>

`

var uploadPage = template.Must(template.New("uploadForm").Parse(uploadPage_template))

type uploadPageData struct {
	Method   string
	Endpoint string
}

func RenderUploadPage(w io.Writer, uploadEndpoint string) error {
	upd := uploadPageData{
		Method:   "POST",
		Endpoint: uploadEndpoint,
	}
	return uploadPage.Execute(w, upd)
}
