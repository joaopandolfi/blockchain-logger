package log

import (
	"logger/services"
	"logger/web"
	"logger/web/controllers"
	"logger/web/server"
	"net/http"
	"strconv"

	"github.com/joaopandolfi/blackwhale/handlers"
	"github.com/joaopandolfi/blackwhale/remotes/jaeger"
	"github.com/joaopandolfi/blackwhale/utils"
	"github.com/opentracing/opentracing-go"
)

// --- Health ---

type controller struct {
	s                 *server.Server
	log               services.Logs
	blockchainService services.BlockChain
}

// New Health controller
func New() controllers.Controller {
	return &controller{
		s:                 nil,
		log:               services.NewLogs(),
		blockchainService: services.NewBlockChain(),
	}
}

func (c *controller) newLog(w http.ResponseWriter, r *http.Request) {
	ctx, span := jaeger.StartSpanFromRequest(opentracing.GlobalTracer(), r, "log")
	defer span.Finish()

	var p payload
	msg, err := handlers.UnmarshalSnakeCaseAndValidate(w, r, &p)
	if err != nil {
		utils.CriticalError("[New Log] parsing body", msg, err.Error())
		handlers.ResponseTypedError(w, web.ErrorCodeInvalidBody, web.ErrorMessageInvalidBody, err)
		span.SetTag("error", true)
		span.SetTag("err_msg", msg)
		return
	}

	newBlock, err := c.log.New(ctx, p.ToLog())
	if err != nil {
		utils.CriticalError("[New Log] saving log", msg, err.Error())
		handlers.ResponseTypedError(w, web.ErrorCodeSave, web.ErrorMessageSave, err)
		span.SetTag("error", true)
		span.SetTag("err_msg", err.Error())
		return
	}

	handlers.RESTResponse(w, newBlock)
}

func (c *controller) validate(w http.ResponseWriter, r *http.Request) {
	ctx, span := jaeger.StartSpanFromRequest(opentracing.GlobalTracer(), r, "log")
	defer span.Finish()

	err := c.blockchainService.Validate(ctx)
	if err != nil {
		utils.CriticalError("[Validate] validating chain", err.Error())
		handlers.ResponseTypedError(w, web.ErrorCodeInternal, web.ErrorMessageInternal, err)
		span.SetTag("error", true)
		return
	}

	handlers.Response(w, true, http.StatusOK)
}

func (c *controller) validateSegment(w http.ResponseWriter, r *http.Request) {
	ctx, span := jaeger.StartSpanFromRequest(opentracing.GlobalTracer(), r, "log")
	defer span.Finish()

	vars := handlers.GetVars(r)
	init, _ := strconv.Atoi(vars["init"])
	end, _ := strconv.Atoi(vars["end"])

	if end < init {
		handlers.ResponseTypedError(w, web.ErrorCodeInvalidBody, "end must be bigger than init", nil)
		return
	}

	err := c.blockchainService.ValidateSegment(ctx, init, end)
	if err != nil {
		utils.CriticalError("[Validate Segment] validating chain", err.Error())
		handlers.ResponseTypedError(w, web.ErrorCodeInternal, web.ErrorMessageInternal, err)
		span.SetTag("error", true)
		return
	}

	handlers.Response(w, true, http.StatusOK)
}
