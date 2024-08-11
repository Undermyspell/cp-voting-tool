// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"voting/bff/templates/components"
	voting_usecases "voting/voting/usecases"
)

func Hello(questions []voting_usecases.QuestionDto) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<ul>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, item := range questions {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(item.Text)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `bff/templates/pages/main.templ`, Line: 11, Col: 18}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" - ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(item.Creator)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `bff/templates/pages/main.templ`, Line: 11, Col: 37}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</ul>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func Main(title string, message string, questions []voting_usecases.QuestionDto) templ.Component {
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
		templ_7745c5c3_Var4 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var4 == nil {
			templ_7745c5c3_Var4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templ.JSONScript("questions", questions).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(title)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `bff/templates/pages/main.templ`, Line: 23, Col: 17}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</title><link href=\"/static/css/tailwind.css\" rel=\"stylesheet\"><link href=\"/static/test.css\" rel=\"stylesheet\"><style>\n\t\t\t\tbody {\n\t\t\t\t\tfont-family: Arial, sans-serif;\n\t\t\t\t\tmargin: 0;\n\t\t\t\t\tpadding: 0;\n\t\t\t\t\theight: 100vh;\n\t\t\t\t\tdisplay: grid;\n\t\t\t\t\tgrid-template-rows: 3rem 1fr;\n\t\t\t\t\tgrid-template-areas: \"header\" \"content\"\n\t\t\t\t}\n\t\t\t\t.content {\n\t\t\t\t\ttext-align: center;\n\t\t\t\t\tpadding: 20px;\n\t\t\t\t\tborder-radius: 8px;\n\t\t\t\t\tbox-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\n\t\t\t\t\tgrid-area: content\n\t\t\t\t}\n\t\t\t\t.header {\n\t\t\t\t\tgrid-area: header\n\t\t\t\t}\n\t\t\t</style><script defer src=\"https://cdn.jsdelivr.net/npm/@alpinejs/focus@3.x.x/dist/cdn.min.js\"></script><script defer src=\"https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js\"></script><script src=\"https://unpkg.com/centrifuge@5.0.1/dist/centrifuge.js\"></script><script src=\"https://unpkg.com/htmx.org@2.0.1\" integrity=\"sha384-QWGpdj554B4ETpJJC9z+ZHJcA/i59TyjxEPXiiUgN2WmTyV5OEZWCD6gQhgkdpB/\" crossorigin=\"anonymous\"></script><script>\n\t\t\t\tdocument.addEventListener('alpine:init', () => {\n\t\t\t\t\t\tAlpine.data(\"userData\", () =>({\n\t\t\t\t\t\t\tasync init() {\n\t\t\t\t\t\t\t\tthis.user = await (await fetch('/user')).json()\n\t\t\t\t\t\t\t\tinitCentrifugo(this.user)\n\t\t\t\t\t\t\t},\n\t\t\t\t\t\t\tquestions: JSON.parse(document.getElementById('questions').textContent),\n\t\t\t\t\t\t\tuser:  null,\n\t\t\t\t\t\t\taddQuestion() {\n\t\t\t\t\t\t\t\tthis.questions.push({\n\t\t\t\t\t\t\t\t\tId: \"asdasdaasasd1123\" + this.questions.length+1,\n\t\t\t\t\t\t\t\t\tText: \"This is an added question \" + this.questions.length+1,\n\t\t\t\t\t\t\t\t\tCreator: \"John Doe\"\n\t\t\t\t\t\t\t\t})\n\t\t\t\t\t\t\t}\n\t\t\t\t\t\t}))\n\t\t\t\t})\n\n\t\t\t\tconst initCentrifugo = async(user) => {\n\t\t\t\t\tconsole.log(\"init centrifugo\")\n\n\t\t\t\t\tconst centrifuge = new Centrifuge(\"ws://localhost:3333/api/v1/connection/websocket\", {\n\t\t\t\t\t\ttoken: user.Token\n\t\t\t\t\t});\n\t\t\t\t\tcentrifuge.on('connecting', function (ctx) {\n\t\t\t\t\t\tconsole.log(`connecting: ${ctx.code}, ${ctx.reason}`);\n\t\t\t\t\t}).on('connected', function (ctx) {\n\t\t\t\t\t\tconsole.log(`connected over ${ctx.transport}`);\n\t\t\t\t\t}).on('disconnected', function (ctx) {\n\t\t\t\t\t\tconsole.log(`disconnected: ${ctx.code}, ${ctx.reason}`);\n\t\t\t\t\t}).on('message', function (msg) {\n\t\t\t\t\t\tconsole.log(`message: ${JSON.stringify(msg)}`);\n\t\t\t\t\t}).connect();\n\t\t\t\t}\n\n\t\t\t\t\n\t\t\t\t</script></head><body><div class=\"header bg-[#CA8787] flex items-center px-10\"><h3>Voting Tool with templ, htmx, alpine.js and golang</h3></div><div class=\"content rounded-md bg-[#606676] overflow-hidden\"><button hx-get=\"/q/new\" hx-target=\"body\" hx-trigger=\"click\" hx-swap=\"beforeend\" class=\"bg-pink-600 p-4 border-4 border-blue-700\" type=\"button\">Add new Question Modal</button>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.QuestionListAlpine().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
