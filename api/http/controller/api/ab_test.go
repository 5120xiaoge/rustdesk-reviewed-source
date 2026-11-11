package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func TestPeerUpdateRejectsNonStringID(t *testing.T) {
	oldLocalizer, oldService := global.Localizer, service.AllService
	t.Cleanup(func() {
		global.Localizer, service.AllService = oldLocalizer, oldService
	})
	bundle := i18n.NewBundle(language.English)
	bundle.AddMessages(language.English, &i18n.Message{ID: "ParamsError", Other: "Invalid parameters"})
	global.Localizer = func(string) *i18n.Localizer { return i18n.NewLocalizer(bundle, "en") }
	service.AllService = &service.Service{UserService: &service.UserService{}}

	for _, body := range []string{`{"id":123}`, `{"id":null}`} {
		t.Run(body, func(t *testing.T) {
			writer := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(writer)
			context.Request = httptest.NewRequest(http.MethodPut, "/ab/peer/update/1-1-0", strings.NewReader(body))
			context.Request.Header.Set("Content-Type", "application/json")
			context.Params = gin.Params{{Key: "guid", Value: "1-1-0"}}
			user := &model.User{GroupId: 1}
			user.Id = 1
			context.Set("curUser", user)
			new(Ab).PeerUpdate(context)
			if writer.Code != http.StatusBadRequest {
				t.Fatalf("expected bad request, got %d", writer.Code)
			}
		})
	}
}
