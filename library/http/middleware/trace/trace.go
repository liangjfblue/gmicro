package trace

import (
	"context"

	"github.com/liangjfblue/gmicro/library/config"

	"github.com/gin-gonic/gin"

	"github.com/liangjfblue/gmicro/library/pkg/tracer"
)

func OpenTracingMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span, err := tracer.TraceFromHeader(context.Background(), "api:"+c.Request.URL.Path, c.Request.Header)
		if err == nil {
			defer span.Finish()
			c.Set(config.Conf.TraceConf.TraceContext, ctx)
		} else {
			c.Set(config.Conf.TraceConf.TraceContext, context.Background())
		}

		c.Next()
	}
}
