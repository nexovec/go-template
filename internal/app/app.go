package app

import (
	"math/rand"
	"templates"
	"utils"

	"github.com/a-h/templ"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/gofiber/fiber/v3"
)

// generates random data for an echarts bar chart
func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}

func RegisterRoutes(root fiber.Router) error {
	// NOTE: this middleware is how I patched templ's .Render method assuming wrong MIME type in the http response.
	// This is also why there is a separate web handler.
	root.Use(utils.MwDefaultResponseMIME("text/html"))

	root.Get("/", func(c fiber.Ctx) error {
		imports := []templ.Component{
			templates.ScriptTag("/static/js/htmx.min.js", false, false, false),
			templates.ScriptTag("/static/js/htmx-ext-json-enc.js", false, false, false),
			templates.LinkTag("/static/css/bulma.min.css"),
		}
		return templates.Boilerplate(templates.Homepage(), templates.Concatenate(imports)).Render(c.Context(), c.Response().BodyWriter())
	})
	root.Get("/getData", func(c fiber.Ctx) error {
		return c.JSON(map[string]interface{}{
			"xAxis": map[string]interface{}{
				"type":        "category",
				"boundaryGap": false,
				"data":        []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
			},
			"yAxis": map[string]interface{}{
				"type": "value",
			},
			"series": []map[string]interface{}{
				{
					"data":      []int{820, 932, 901, 934, 1290, 1330, 1320},
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
			},
		})
	})
	return nil
}
