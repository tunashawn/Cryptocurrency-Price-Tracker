package gin_test_setup

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

type GinTestContext struct {
	*gin.Context
	w *httptest.ResponseRecorder
}

// Usage
/*
	ctx := NewGinTestContext("GET", "/price/history")

	p := &Controller{}
	p.GetLatestPrice(ctx.Context)

	responseBody := ctx.GetResponseBody()

	var res response.ResponseData
	json.Unmarshal([]byte(responseBody), &res)
*/
func NewGinTestContext(method, path string) GinTestContext {
	// Set gin mode to test
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	// Create test gin context for running tests, which will have httptest.NewRecorder as an argument
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(method, path, nil)

	ctx.Request.Method = method

	return GinTestContext{
		Context: ctx,
		w:       w,
	}
}

/*
	 Set path params
	     params := []gin.Param{
			{
				Key:   "id",
				Value: "1",
			},
		}
*/
//func (c *GinTestContext) CustomQueryParams(params gin.Params) {
//	c.Params = params
//}
//
//func (c *GinTestContext) CustomBody(content any) {
//	jsonBytes, err := json.Marshal(content)
//	if err != nil {
//		panic(err)
//	}
//
//	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
//}
//
//func (c *GinTestContext) GetResponseCode() int {
//	return c.w.Code
//}

func (c *GinTestContext) GetResponseBody() string {
	return c.w.Body.String()
}
