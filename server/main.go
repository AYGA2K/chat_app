package main

import (
	"log"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
)

const namespace = "default"

var serverEvents = websocket.Namespaces{
	namespace: websocket.Events{
		websocket.OnNamespaceConnected: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			// with `websocket.GetContext` you can retrieve the Iris' `Context`.
			ctx := websocket.GetContext(nsConn.Conn)

			log.Printf("[%s] connected to namespace [%s] with IP [%s]",
				nsConn, msg.Namespace,
				ctx.RemoteAddr())
			return nil
		},
		websocket.OnNamespaceDisconnect: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("[%s] disconnected from namespace [%s]", nsConn, msg.Namespace)
			return nil
		},
		"chat": func(nsConn *websocket.NSConn, msg websocket.Message) error {
			// room.String() returns -> NSConn.String() returns -> Conn.String() returns -> Conn.ID()
			log.Printf("[%s] sent: %s", nsConn, string(msg.Body))

			// Write message back to the client message owner with:
			// nsConn.Emit("chat", msg)
			// Write message to all except this client with:
			nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil
		},
	},
}

func main() {
	app := iris.New()
	websocketServer := websocket.New(
		websocket.DefaultGorillaUpgrader, /* DefaultGobwasUpgrader can be used too. */
		serverEvents)

	j := jwt.New(jwt.Config{
		// Extract by the "token" url,
		// so the client should dial with ws://localhost:8080/echo?token=$token
		Extractor: jwt.FromParameter("token"),

		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},

		// When set, the middleware verifies that tokens are signed
		// with the specific signing algorithm
		// If the signing method is not constant the
		// `Config.ValidationKeyGetter` callback field can be used
		// to implement additional checks
		// Important to avoid security issues described here:
		// https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	idGen := func(ctx iris.Context) string {
		if username := ctx.GetHeader("X-Username"); username != "" {
			return username
		}

		return websocket.DefaultIDGenerator(ctx)
	}

	websocketRoute := app.Get("/echo", websocket.Handler(websocketServer, idGen))
	websocketRoute.Use(j.Serve)
	app.Listen(":8080")
}
