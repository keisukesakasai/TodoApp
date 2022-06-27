package controllers

import (
	"TodoApp/app/models"
	"log"
	"net/http"
)

func top(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "top")
	defer span.End()

	_, err := session(w, r)
	if err != nil {
		generateHTML(ctx, w, "hello", "top", "layout", "top", "public_navbar")
	} else {
		ctx = r.Context()
		ctx, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/todos", 302)
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "index")
	defer span.End()

	sess, err := session(w, r)
	if err != nil {
		ctx = r.Context()
		ctx, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/", 302)
	} else {
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		todos, _ := user.GetTodosByUser()
		user.Todos = todos
		generateHTML(ctx, w, user, "index", "layout", "private_navbar", "index")
	}
}

func todoNew(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "todoNew")
	defer span.End()

	_, err := session(w, r)
	if err != nil {
		ctx = r.Context()
		ctx, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/login", 302)
	} else {
		generateHTML(ctx, w, nil, "todoNew", "layout", "private_navbar", "todo_new")
	}
}

func todoSave(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "todoSave")
	defer span.End()

	sess, err := session(w, r)
	if err != nil {
		ctx = r.Context()
		ctx, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/login", 302)
	} else {
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		content := r.PostFormValue(("content"))
		if err := user.CreateTodo(content); err != nil {
			log.Println(err)
		}
		ctx = r.Context()
		ctx, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/todos", 302)
	}
}

func todoEdit(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "todoEdit")
	defer span.End()

	sess, err := session(w, r)
	if err != nil {
		ctx = r.Context()
		ctx, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/login", 302)
	} else {
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)

		}
		t, err := models.GetTodo(id)
		if err != nil {
			log.Fatalln(err)
		}
		generateHTML(ctx, w, t, "todoEdit", "layout", "private_navbar", "todo_edit")
	}
}

func todoUpdate(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "todoUpdate")
	defer span.End()

	sess, err := session(w, r)
	if err != nil {
		ctx = r.Context()
		ctx, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/login", 302)
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		content := r.PostFormValue("content")
		t := &models.Todo{ID: id, Content: content, UserID: user.ID}
		if err := t.UpdateTodo(); err != nil {
			log.Println(err)
		}
		ctx = r.Context()
		ctx, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/todos", 302)
	}
}

func todoDelete(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "todoDelete")
	defer span.End()

	sess, err := session(w, r)
	if err != nil {
		ctx = r.Context()
		ctx, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/login", 302)
	} else {
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		t, err := models.GetTodo(id)
		if err != nil {
			log.Println(err)
		}
		if err := t.DeleteTodo(); err != nil {
			log.Println(err)
		}
		ctx = r.Context()
		ctx, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/todos", 302)
	}
}
