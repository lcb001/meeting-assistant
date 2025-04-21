package main

import (
	"context"
	"time"

	"meetingagent/handlers"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func main() {
	h := server.Default()
	h.Use(Logger())

	// Register API routes first
	h.POST("/meeting", handlers.CreateMeeting)
	h.GET("/meeting", handlers.ListMeetings)
	h.GET("/summary", handlers.GetMeetingSummary)
	h.GET("/chat", handlers.HandleChat)

	// Serve static files
	h.StaticFS("/", &app.FS{
		Root:               "./static",
		PathRewrite:        app.NewPathSlashesStripper(1),
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: true,
	})

	// Start server
	h.Spin()
}

func Logger() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		start := time.Now()
		path := string(ctx.Request.URI().Path())
		query := string(ctx.Request.URI().QueryString())
		if query != "" {
			path = path + "?" + query
		}

		// Process request
		ctx.Next(c)

		// Calculate latency
		latency := time.Since(start)

		// Get response status code
		statusCode := ctx.Response.StatusCode()

		// Log request details
		hlog.CtxInfof(c, "[HTTP] %s %s - %d - %v",
			ctx.Request.Method(),
			path,
			statusCode,
			latency,
		)
	}
}
