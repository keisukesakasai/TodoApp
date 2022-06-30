package controllers

import (
	"TodoApp/app/models"
	"TodoApp/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// otel collector
func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("TodoAPP"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "otel-collector-collector.tracing.svc.cluster.local:4318", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	/*
		traceExporter, err := stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithWriter(os.Stderr),
		)
	*/

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

func generateHTML(c *gin.Context, data interface{}, procname string, filenames ...string) {
	//tracer := otel.Tracer("generateHTML")
	_, span := tracer.Start(c.Request.Context(), "generateHTML: "+procname)
	defer span.End()

	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(c.Writer, "layout", data)
}

func session(c *gin.Context) (sess models.Session, err error) {
	_, span := tracer.Start(c.Request.Context(), "session")
	defer span.End()

	// cookie, err := c.Request.Cookie("_cookie")
	cookie, err := c.Cookie("_cookie")
	fmt.Println("===session===")
	//fmt.Println(cookie.Value)
	fmt.Println(cookie)
	fmt.Println("===session===")

	if err == nil {
		// sess = models.Session{UUID: cookie.Value}
		sess = models.Session{UUID: cookie}
		if ok, _ := sess.CheckSession(c); !ok {
			err = fmt.Errorf("invalid session")
		}
	}
	return sess, err
}

/*
func session(ctx context.Context, w http.ResponseWriter, r *http.Request) (sess models.Session, err error) {
	tracer := otel.Tracer("session")
	// ctx := r.Context()
	ctx, span := tracer.Start(ctx, "session")
	defer span.End()

	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess = models.Session{UUID: cookie.Value}
		if ok, _ := sess.CheckSession(ctx); !ok {
			err = fmt.Errorf("invalid session")
		}
	}
	return sess, err
}
*/

var validPath = regexp.MustCompile("^/todos/(edit|save|update|delete)/([0-9]+)$")

func parseURL(fn func(*gin.Context, int)) gin.HandlerFunc {
	return func(c *gin.Context) {

		_, span := tracer.Start(c.Request.Context(), "parseURL")
		defer span.End()

		fmt.Println("===parseURL")
		fmt.Println(c.Request.URL.Path)
		fmt.Println("===parseURL")

		q := validPath.FindStringSubmatch(c.Request.URL.Path)

		fmt.Println("===parseURL")
		fmt.Println(q)
		fmt.Println("===parseURL")

		if q == nil {
			http.NotFound(c.Writer, c.Request)
			return
		}

		id, _ := strconv.Atoi(q[2])
		fmt.Println(id)
		fn(c, id)
	}
}

// --otelcollecotr--
var tracer = otel.Tracer("controllers")

func StartMainServer() {
	fmt.Println("start server" + "port: " + config.Config.Port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// otel collector
	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	r := gin.New()
	r.Use(otelgin.Middleware("todoapp-server"))
	r.LoadHTMLGlob(config.Config.Static + "/templates/*")
	r.Static("/static/", config.Config.Static)

	//--- handler
	r.GET("/", top)
	r.GET("/signup", getSignup)
	r.POST("/signup", postSignup)
	r.GET("/login", login)
	r.GET("/logout", logout)
	r.POST("/authenticate", authenticate)

	r.GET("/todos", index)
	r.GET("/todos/new", todoNew)
	r.POST("/todos/save", todoSave)
	r.GET("/todos/edit/:id", parseURL(todoEdit))
	r.POST("/todos/update/:id", parseURL(todoUpdate))
	r.GET("/todos/delete/:id", parseURL(todoDelete))

	r.Run(":" + config.Config.Port)
}
