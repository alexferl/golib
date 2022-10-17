package router

import (
	"testing"

	"github.com/labstack/echo/v4"
)

func TestRouter(t *testing.T) {
	name := "MyRoute"
	otherName := "OtherRoute"
	r := &Router{}
	r.Routes = []*Route{{Name: name}, {Name: otherName}}

	tests := []struct {
		input    *Route
		expected string
	}{
		{r.FindRouteByName(name), name},
		{r.FindRouteByName(otherName), otherName},
		{r.FindRouteByName("does_not_exists"), ""},
	}

	for _, tc := range tests {
		if tc.input != nil && tc.input.Name != tc.expected {
			t.Errorf("got %v expected %s", tc.input, tc.expected)
		}
	}
}

func TestRegister(t *testing.T) {
	name := "MyRoute"
	otherName := "OtherRoute"
	r := &Router{}
	h := func(c echo.Context) error { return nil }
	r.Routes = []*Route{
		{Name: name, Pattern: "/1", HandlerFunc: h},
		{Name: otherName, Pattern: "/2", HandlerFunc: h},
	}

	e := echo.New()
	Register(e, r)

	if len(r.Routes) != len(e.Routes()) {
		t.Errorf("got %d expected %d", len(r.Routes), len(e.Routes()))
	}
}
