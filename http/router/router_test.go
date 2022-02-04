package router

import (
	"github.com/labstack/echo/v4"
	"testing"
)

func TestRouter(t *testing.T) {
	name := "MyRoute"
	otherName := "OtherRoute"
	r := &Router{}
	r.Routes = []Route{{Name: name}, {Name: otherName}}

	var tests = []struct {
		input    *Route
		expected string
	}{
		{r.FindRouteByName(name), name},
		{r.FindRouteByName(otherName), otherName},
		{r.FindRouteByName("does_not_exists"), ""},
	}

	for _, tc := range tests {
		if tc.input.Name != tc.expected {
			t.Errorf("got %s expected %s", tc.input, tc.expected)
		}
	}
}

func TestRegister(t *testing.T) {
	name := "MyRoute"
	otherName := "OtherRoute"
	r := &Router{}
	r.Routes = []Route{{Name: name}, {Name: otherName}}

	e := echo.New()
	Register(e, r)

	if len(r.Routes) != len(e.Routes()) {
		t.Errorf("got %d expected %d", len(r.Routes), len(e.Routes()))
	}
}
