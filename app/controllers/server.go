package controllers

import (
	"TodoApp/app/SessionInfo"
	"TodoApp/config"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var LoginInfo SessionInfo.Session

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

	var tracerProvider *sdktrace.TracerProvider
	if config.Config.Deploy == "local" {

		traceExporter, err := stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			// stdouttrace.WithWriter(os.Stderr),
			stdouttrace.WithWriter(io.Discard),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create trace exporter: %w", err)
		}
		bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
		tracerProvider := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(res),
			sdktrace.WithSpanProcessor(bsp),
		)
		otel.SetTracerProvider(tracerProvider)
		otel.SetTextMapPropagator(propagation.TraceContext{})
	}

	if config.Config.Deploy == "prod" {
		conn, err := grpc.DialContext(ctx, "otel-collector-collector.tracing.svc.cluster.local:4318", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
		}

		// Set up a trace exporter
		traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, fmt.Errorf("failed to create trace exporter: %w", err)
		}

		// idg := xray.NewIDGenerator()

		bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
		tracerProvider := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(res),
			sdktrace.WithSpanProcessor(bsp),
			// sdktrace.WithIDGenerator(idg),
		)

		otel.SetTracerProvider(tracerProvider)
		otel.SetTextMapPropagator(propagation.TraceContext{})
	}

	return tracerProvider.Shutdown, nil
}

func generateHTML(c *gin.Context, data interface{}, procname string, filenames ...string) {
	_, span := tracer.Start(c.Request.Context(), "generateHTML : "+procname)
	defer span.End()

	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(c.Writer, "layout", data)
}

var validPath = regexp.MustCompile("^/menu/todos/(edit|save|update|delete)/([0-9]+)$")

func parseURL(fn func(*gin.Context, int)) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, span := tracer.Start(c.Request.Context(), "parseURL")
		defer span.End()

		fmt.Println(c.Request.URL.Path)
		q := validPath.FindStringSubmatch(c.Request.URL.Path)
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
	fmt.Println("info: Start Server" + "port: " + config.Config.Port)

	// ????????????????????????
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Otel Collecotor ??????????????????
	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	// router ??????
	r := gin.New()

	// Custom Middleware ??????
	r.Use(otelgin.Middleware("todoapp-server"))

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// template ??????
	r.LoadHTMLGlob(config.Config.Static + "/templates/*")
	r.Static("/static/", config.Config.Static)

	//--- handler ??????
	r.GET("/", top)
	r.GET("/login", getLogin)
	r.POST("/login", postLogin)

	r.GET("/signup", getSignup)
	r.POST("/signup", postSignup)

	rTodos := r.Group("/menu")
	rTodos.Use(checkSession())
	{
		rTodos.GET("/todos", index)
		rTodos.GET("/todos/new", todoNew)
		rTodos.POST("/todos/save", todoSave)
		rTodos.GET("/todos/edit/:id", parseURL(todoEdit))
		rTodos.POST("/todos/update/:id", parseURL(todoUpdate))
		rTodos.GET("/todos/delete/:id", parseURL(todoDelete))
	}
	r.GET("/logout", getLogout)

	r.Run(":" + config.Config.Port)
}

func checkSession() gin.HandlerFunc {
	return func(c *gin.Context) {

		_, span := tracer.Start(c.Request.Context(), "?????????????????????????????????")
		defer span.End()

		log.Println("?????????????????????????????????")

		session := sessions.Default(c)
		LoginInfo.UserID = session.Get("UserId")

		if LoginInfo.UserID == nil {
			log.Println("??????????????????????????????")

			c.Redirect(http.StatusMovedPermanently, "/login")
			c.Abort()
		} else {
			c.Set("UserId", LoginInfo.UserID) // ?????????ID????????????
			c.Next()
		}

		_, span = tracer.Start(c.Request.Context(), "?????????????????????????????????")
		defer span.End()

		log.Println("?????????????????????????????????")
	}
}
