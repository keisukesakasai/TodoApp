package controllers

import (
	"TodoApp/app/models"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
)

func top(w http.ResponseWriter, r *http.Request) {
	//tracer := otel.Tracer("top")
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "top")
	defer span.End()

	_, err := session(ctx, w, r)
	if err != nil {
		generateHTML(ctx, w, "hello", "top", "layout", "top", "public_navbar")
	} else {
		ctx = r.Context()
		_, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/todos", http.StatusFound)
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("index")
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "index")
	defer span.End()

	sess, err := session(ctx, w, r)
	if err != nil {
		ctx = r.Context()
		_, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		user, err := sess.GetUserBySession(ctx)
		if err != nil {
			log.Println(err)
		}
		todos, _ := user.GetTodosByUser(ctx)
		user.Todos = todos
		generateHTML(ctx, w, user, "index", "layout", "private_navbar", "index")
	}
}

func todoNew(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("todoNew")
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "todoNew")
	defer span.End()

	_, err := session(ctx, w, r)
	if err != nil {
		ctx = r.Context()
		_, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		generateHTML(ctx, w, nil, "todoNew", "layout", "private_navbar", "todo_new")
	}
}

func todoSave(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("todoSave")
	ctx := r.Context()
	_, span := tracer.Start(ctx, "todoSave")
	defer span.End()

	sess, err := session(ctx, w, r)
	if err != nil {
		ctx = r.Context()
		_, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user, err := sess.GetUserBySession(ctx)
		if err != nil {
			log.Println(err)
		}
		content := r.PostFormValue(("content"))
		if err := user.CreateTodo(ctx, content); err != nil {
			log.Println(err)
		}
		ctx = r.Context()
		_, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
}

func todoEdit(w http.ResponseWriter, r *http.Request, id int) {
	tracer := otel.Tracer("todoEdit")
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "todoEdit")
	defer span.End()

	sess, err := session(ctx, w, r)
	if err != nil {
		ctx = r.Context()
		_, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		_, err := sess.GetUserBySession(ctx)
		if err != nil {
			log.Println(err)

		}
		t, err := models.GetTodo(ctx, id)
		if err != nil {
			log.Fatalln(err)
		}
		generateHTML(ctx, w, t, "todoEdit", "layout", "private_navbar", "todo_edit")
	}
}

func todoUpdate(w http.ResponseWriter, r *http.Request, id int) {
	tracer := otel.Tracer("todoUpdate")
	ctx := r.Context()
	_, span := tracer.Start(ctx, "todoUpdate")
	defer span.End()

	sess, err := session(ctx, w, r)
	if err != nil {
		ctx = r.Context()
		_, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user, err := sess.GetUserBySession(ctx)
		if err != nil {
			log.Println(err)
		}
		content := r.PostFormValue("content")
		t := &models.Todo{ID: id, Content: content, UserID: user.ID}
		if err := t.UpdateTodo(ctx); err != nil {
			log.Println(err)
		}
		ctx = r.Context()
		_, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
}

func todoDelete(w http.ResponseWriter, r *http.Request, id int) {
	tracer := otel.Tracer("todoDelete")
	ctx := r.Context()
	_, span := tracer.Start(ctx, "todoDelete")
	defer span.End()

	sess, err := session(ctx, w, r)
	if err != nil {
		ctx = r.Context()
		_, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		_, err := sess.GetUserBySession(ctx)
		if err != nil {
			log.Println(err)
		}
		t, err := models.GetTodo(ctx, id)
		if err != nil {
			log.Println(err)
		}
		if err := t.DeleteTodo(ctx); err != nil {
			log.Println(err)
		}
		ctx = r.Context()
		_, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
}
