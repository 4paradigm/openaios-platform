package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	cache "github.com/4paradigm/openaios-platform/src/webterminal/gotty/cache/token"
	"github.com/4paradigm/openaios-platform/src/webterminal/gotty/pkg/randomstring"
	"github.com/4paradigm/openaios-platform/src/webterminal/gotty/webtty"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type GetTerminalApiResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type GetTerminalApiRequest struct {
	Namespace     string `json:"namespace"`
	PodName       string `json:"pod_name"`
	ContainerName string `json:"container_name"`
	BearerToken   string `json:"bearer_token"`
}

func (server *Server) generateHandleWS(ctx context.Context, cancel context.CancelFunc, counter *counter) http.HandlerFunc {
	once := new(int64)

	go func() {
		select {
		case <-counter.timer().C:
			cancel()
		case <-ctx.Done():
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		if server.options.Once {
			success := atomic.CompareAndSwapInt64(once, 0, 1)
			if !success {
				http.Error(w, "Server is shutting down", http.StatusServiceUnavailable)
				return
			}
		}

		num := counter.add(1)
		closeReason := "unknown reason"

		defer func() {
			num := counter.done()
			log.Printf(
				"Connection closed by %s: %s, connections: %d/%d",
				closeReason, r.RemoteAddr, num, server.options.MaxConnection,
			)

			if server.options.Once {
				cancel()
			}
		}()

		if int64(server.options.MaxConnection) != 0 {
			if num > server.options.MaxConnection {
				closeReason = "exceeding max number of connections"
				return
			}
		}

		log.Printf("New client connected: %s, connections: %d/%d", r.RemoteAddr, num, server.options.MaxConnection)

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		conn, err := server.upgrader.Upgrade(w, r, nil)
		if err != nil {
			closeReason = err.Error()
			return
		}
		defer conn.Close()

		err = server.processWSConn(ctx, conn)

		switch err {
		case ctx.Err():
			closeReason = "cancelation"
		case webtty.ErrSlaveClosed:
			closeReason = server.factory.Name()
		case webtty.ErrMasterClosed:
			closeReason = "client"
		case cache.InvalidToken:
			closeReason = fmt.Sprintf("an token error: %s", err)
			return
		case cache.NoTokenProvided:
			closeReason = fmt.Sprintf("an token error: %s", err)
			return
		default:
			closeReason = fmt.Sprintf("an error: %s", err)
		}
	}
}

func (server *Server) processWSConn(ctx context.Context, conn *websocket.Conn) error {
	typ, initLine, err := conn.ReadMessage()
	if err != nil {
		return errors.Wrapf(err, "failed to authenticate websocket connection")
	}
	if typ != websocket.TextMessage {
		return errors.New("failed to authenticate websocket connection: invalid message type")
	}

	var init InitMessage
	err = json.Unmarshal(initLine, &init)
	if err != nil {
		return errors.Wrapf(err, "failed to authenticate websocket connection")
	}
	if init.AuthToken != server.options.Credential {
		return errors.New("failed to authenticate websocket connection")
	}

	queryPath := "?"
	if server.options.PermitArguments && init.Arguments != "" {
		queryPath = init.Arguments
	}

	query, err := url.Parse(queryPath)
	if err != nil {
		return errors.Wrapf(err, "failed to parse arguments")
	}

	params := query.Query()
	params.Del("arg")
	windowTitle := ""
	arg := []string{}
	if len(params.Get("token")) > 0 {
		ttyParameter := server.cache.Get(params.Get("token"))
		cachedKey := params.Get("token")
		if ttyParameter != nil {
			windowTitle = ttyParameter.Title
			arg = ttyParameter.Arg
			server.cache.Delete(cachedKey)
		} else {
			return cache.InvalidToken
		}
	} else {
		return cache.NoTokenProvided
	}
	for i, _ := range arg {
		params.Add("arg", arg[i])
	}

	var slave Slave
	slave, err = server.factory.New(params)
	if err != nil {
		return errors.Wrapf(err, "failed to create backend")
	}
	defer slave.Close()

	titleVars := server.titleVariables(
		[]string{"server", "master", "slave"},
		map[string]map[string]interface{}{
			"server": server.options.TitleVariables,
			"master": map[string]interface{}{
				"remote_addr": conn.RemoteAddr(),
			},
			"slave": slave.WindowTitleVariables(),
		},
	)

	titleBuf := new(bytes.Buffer)
	err = server.titleTemplate.Execute(titleBuf, titleVars)
	if err != nil {
		return errors.Wrapf(err, "failed to fill window title template")
	}

	if len(windowTitle) > 0 {
		titleBuf.Reset()
		titleBuf.WriteString("Pineapple Web Terminal: " + windowTitle)
	}

	opts := []webtty.Option{
		webtty.WithWindowTitle(titleBuf.Bytes()),
	}
	if server.options.PermitWrite {
		opts = append(opts, webtty.WithPermitWrite())
	}
	if server.options.EnableReconnect {
		opts = append(opts, webtty.WithReconnect(server.options.ReconnectTime))
	}
	if server.options.Width > 0 {
		opts = append(opts, webtty.WithFixedColumns(server.options.Width))
	}
	if server.options.Height > 0 {
		opts = append(opts, webtty.WithFixedRows(server.options.Height))
	}
	if server.options.Preferences != nil {
		opts = append(opts, webtty.WithMasterPreferences(server.options.Preferences))
	}

	tty, err := webtty.New(&wsWrapper{conn}, slave, opts...)
	if err != nil {
		return errors.Wrapf(err, "failed to create webtty")
	}

	err = tty.Run(ctx)

	return err
}

func (server *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	titleVars := server.titleVariables(
		[]string{"server", "master"},
		map[string]map[string]interface{}{
			"server": server.options.TitleVariables,
			"master": map[string]interface{}{
				"remote_addr": r.RemoteAddr,
			},
		},
	)

	titleBuf := new(bytes.Buffer)
	err := server.titleTemplate.Execute(titleBuf, titleVars)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	indexVars := map[string]interface{}{
		"title": titleBuf.String(),
	}

	indexBuf := new(bytes.Buffer)
	err = server.indexTemplate.Execute(indexBuf, indexVars)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Write(indexBuf.Bytes())
}

func (server *Server) handleAuthToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	// @TODO hashing?
	w.Write([]byte("var gotty_auth_token = '" + server.options.Credential + "';"))
}

func (server *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write([]byte("var gotty_term = '" + server.options.Term + "';"))
}

// titleVariables merges maps in a specified order.
// varUnits are name-keyed maps, whose names will be iterated using order.
func (server *Server) titleVariables(order []string, varUnits map[string]map[string]interface{}) map[string]interface{} {
	titleVars := map[string]interface{}{}

	for _, name := range order {
		vars, ok := varUnits[name]
		if !ok {
			panic("title variable name error")
		}
		for key, val := range vars {
			titleVars[key] = val
		}
	}

	// safe net for conflicted keys
	for _, name := range order {
		titleVars[name] = varUnits[name]
	}

	return titleVars
}

func (server *Server) handleGetTerminalApi(w http.ResponseWriter, r *http.Request) {
	result := GetTerminalApiResponse{
		Success: false,
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if !strings.EqualFold(r.Method, "POST") {
		result.Message = "Method Not Allowed"
		json.NewEncoder(w).Encode(result)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		result.Message = "Invalid Request Body"
		json.NewEncoder(w).Encode(result)
		return
	}
	var request GetTerminalApiRequest
	if err = json.Unmarshal(body, &request); err != nil {
		fmt.Printf("Unmarshal err, %v\n", err)
		result.Message = "Invalid Request Body"
		json.NewEncoder(w).Encode(result)
		return
	}
	if len(request.BearerToken) < 10 {
		result.Message = "Invalid Bearer Token"
		json.NewEncoder(w).Encode(result)
		return
	}

	token := randomstring.Generate(20)
	arg := []string{
		"--token",
		request.BearerToken,
		"exec",
		"-it",
		"-n",
		request.Namespace,
		request.PodName,
		"-c",
		request.ContainerName,
		"--",
		"sh",
	}

	ttyParameter := cache.TtyParameter{
		Title: request.ContainerName + "@" + request.PodName,
		Arg:   arg,
	}
	if err := server.cache.Add(token, &ttyParameter, time.Duration(server.options.TokenExpiresDuration)*time.Second); err != nil {
		log.Printf("save token and ttyParam err:%s", err.Error())
		result.Success = false
		result.Message = err.Error()
		json.NewEncoder(w).Encode(result)
		return
	}
	result.Success = true
	result.Token = token
	json.NewEncoder(w).Encode(result)
}
