package route

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"regexp"
)

type contextKey int

const (
	VarsKey contextKey = iota
)

type regexData struct {
	pattern       *regexp.Regexp
	pathParamKeys []string
}

type route struct {
	regexData regexData
	handler   http.Handler
	methods   []string
}

type RegexpHandler struct {
	routes []*route
}

func Vars(r *http.Request) map[string]string {
	if rv := r.Context().Value(VarsKey); rv != nil {
		return rv.(map[string]string)
	}
	return nil
}

func (route *route) matchRoute(r *http.Request) bool {

	if route.regexData.pattern.MatchString(r.URL.Path) {
		for _, method := range route.methods {
			if r.Method == method {
				return true
			}
		}
	}

	return false
}

func (route *route) Methods(methods ...string) *route {
	for _, m := range methods {
		route.methods = append(route.methods, m)
	}
	return route
}

func (route *route) withPathParams(r *http.Request) *http.Request {

	pathParamMap := map[string]string{}
	allString := route.regexData.pattern.FindStringSubmatch(r.URL.Path)
	for i, v := range allString {
		if i > 0 {
			pathParamMap[route.regexData.pathParamKeys[i-1]] = v
		}
	}

	r = r.WithContext(context.WithValue(r.Context(), VarsKey, pathParamMap))

	return r
}

func (h *RegexpHandler) Handler(urlPattern string, handler http.Handler) *route {
	regexPattern, _ := UrlRegexConverter(urlPattern)
	r := route{regexPattern, handler, make([]string, 0)}
	h.routes = append(h.routes, &r)
	return &r
}

func (h *RegexpHandler) HandleFunc(urlPattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	regexPattern, _ := UrlRegexConverter(urlPattern)
	r := route{regexPattern, http.HandlerFunc(handler), make([]string, 0)}
	h.routes = append(h.routes, &r)
	return &r
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()

	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

	//expected to get a child span, as the client is sending the trace id, didn't work
	//span, _ := tracer.StartSpan(r.URL.Path), ext.RPCServerOption(spanCtx)

	//creating either a new span or a child span
	var span opentracing.Span
	if spanCtx != nil {
		span = tracer.StartSpan(
			r.URL.Path,
			opentracing.ChildOf(spanCtx),
		)
	} else {
		span = tracer.StartSpan(r.URL.Path)
	}

	var matchingRoute *route
	for _, route := range h.routes {
		if route.matchRoute(r) {
			matchingRoute = route
		}
	}

	if matchingRoute != nil {
		r = matchingRoute.withPathParams(r)
		matchingRoute.handler.ServeHTTP(w, r)
	} else {
		// no pattern matched; send 404 response
		http.NotFound(w, r)
	}
	span.Finish()
}

func UrlRegexConverter(urlPattern string) (regexData, error) {
	regexData := regexData{}
	pathParamExpr, _ := regexp.Compile("({[A-Za-z]*})")
	pathParamKey, _ := regexp.Compile("{([A-Za-z]*)}")

	// first prepare the url for matching an url
	literalString := pathParamExpr.ReplaceAllLiteralString(urlPattern, "(\\w+)")
	r, _ := regexp.Compile("^" + literalString + "$")
	regexData.pattern = r

	//second find all path param key for mapping to value of a matching url
	allString := pathParamKey.FindStringSubmatch(urlPattern)
	for i, v := range allString {
		if i > 0 {
			regexData.pathParamKeys = append(regexData.pathParamKeys, v)
		}
	}

	return regexData, nil
}
