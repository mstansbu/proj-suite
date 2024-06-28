// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Layout() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html><head><link rel=\"stylesheet\" href=\"main.css\"><script src=\"https://cdnjs.cloudflare.com/ajax/libs/jquery/3.7.1/jquery.slim.min.js\" integrity=\"sha512-sNylduh9fqpYUK5OYXWcBleGzbZInWj8yCJAU57r1dpSK9tP2ghf/SRYCMj+KsslFkCOt3TvJrX2AV/Gc3wOqA==\" crossorigin=\"anonymous\" referrerpolicy=\"no-referrer\"></script><script defer src=\"https://cdnjs.cloudflare.com/ajax/libs/alpinejs/3.14.1/cdn.min.js\" integrity=\"sha512-ytM6hP1K9BkRTjUQZpxZKFjJ2TvE4QXaK7phVymsm7NimaI5H09TWWW6f2JMbonLp4ftYU6xfwQGoe3C8jta9A==\" crossorigin=\"anonymous\" referrerpolicy=\"no-referrer\"></script><script src=\"https://cdnjs.cloudflare.com/ajax/libs/htmx/2.0.0/htmx.min.js\" integrity=\"sha512-Cpedvic0/Mgc3uRJ5apQ/ZYroPCZpatYEXGJayRaRNjKLaFualFxfxn97LJymznV+nC7y0/Hp/apHNwGpMimuw==\" crossorigin=\"anonymous\" referrerpolicy=\"no-referrer\"></script><script defer src=\"tic-tac-toe.js\"></script><title>Tic Tac Toe</title></head><body><div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Home().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
