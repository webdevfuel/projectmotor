// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package layout

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Dashboard() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex\"><div class=\"flex-shrink-0 w-[288px] bg-slate-800 py-2.5 flex flex-col justify-between h-screen sticky top-0\"><div><p class=\"text-xl font-bold text-center dark:text-white mt-2.5\">ProjectMotor</p><ul class=\"mt-24 px-4\"><li><a href=\"/\" class=\"inline-flex items-center gap-2 w-full p-2 rounded-lg dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"lucide lucide-layout-dashboard\"><rect width=\"7\" height=\"9\" x=\"3\" y=\"3\" rx=\"1\"></rect> <rect width=\"7\" height=\"5\" x=\"14\" y=\"3\" rx=\"1\"></rect> <rect width=\"7\" height=\"9\" x=\"14\" y=\"12\" rx=\"1\"></rect> <rect width=\"7\" height=\"5\" x=\"3\" y=\"16\" rx=\"1\"></rect></svg> Dashboard</a></li><li><a href=\"/projects\" class=\"inline-flex items-center gap-2 w-full p-2 rounded-lg dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"lucide lucide-folder\"><path d=\"M20 20a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.9a2 2 0 0 1-1.69-.9L9.6 3.9A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13a2 2 0 0 0 2 2Z\"></path></svg> Projects</a></li><li><a href=\"/tasks\" class=\"inline-flex items-center gap-2 w-full p-2 rounded-lg dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"lucide lucide-list-todo\"><rect x=\"3\" y=\"5\" width=\"6\" height=\"6\" rx=\"1\"></rect> <path d=\"m3 17 2 2 4-4\"></path> <path d=\"M13 6h8\"></path> <path d=\"M13 12h8\"></path> <path d=\"M13 18h8\"></path></svg> Tasks</a></li></ul></div><div class=\"px-4 w-full\"><button hx-post=\"/logout\" class=\"py-3 px-4 w-full inline-flex justify-center items-center gap-x-2 text-sm rounded-lg border border-gray-200 bg-white text-gray-800 shadow-sm hover:bg-gray-50 disabled:opacity-50 disabled:pointer-events-none dark:bg-slate-900 dark:border-gray-700 dark:text-white dark:hover:bg-slate-800 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"16\" height=\"16\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"lucide lucide-log-out\"><path d=\"M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4\"></path> <polyline points=\"16 17 21 12 16 7\"></polyline> <line x1=\"21\" x2=\"9\" y1=\"12\" y2=\"12\"></line></svg> Log out</button></div></div><div class=\"px-12 py-6 flex-1\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = Base().Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
