package main

import (
	"bytes"

	"github.com/a-h/templ"
	"github.com/kataras/iris/v12"
)

type todo struct {
	title     string
	completed bool
}

var todos []todo = make([]todo, 0)

func removeTodo(index int) []todo {
	return append(todos[:index], todos[index+1:]...)
}

func newTodo(title string) []todo {
	t := todo{title: title, completed: false}

	return append(todos, t)
}

type form struct {
	Title string `form:"title"`
}

func main() {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		Render(ctx, indexPage(todos))
	})

	app.Post("/todo", func(ctx iris.Context) {
		var form form
		err := ctx.ReadForm(&form)
		if err != nil {
			ctx.StopWithError(iris.StatusBadRequest, err)
			return
		}

		todos = newTodo(form.Title)

		ctx.Redirect("/")
	})

	app.Post("/todo/{id:int}/toggle", func(ctx iris.Context) {
		index := ctx.Params().GetIntDefault("id", 0)
		todo := todos[index]
		todo.completed = !todo.completed
		todos[index] = todo

		ctx.Redirect("/")
	})

	app.Post("/todo/{id:int}/delete", func(ctx iris.Context) {
		todos = removeTodo(ctx.Params().GetIntDefault("id", 0))

		ctx.Redirect("/")
	})

	app.Listen(":8080")
}

func Render(ctx iris.Context, c templ.Component) {
	buf := new(bytes.Buffer)

	c.Render(ctx, buf)

	ctx.HTML(buf.String())
}
