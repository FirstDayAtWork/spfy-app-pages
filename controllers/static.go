package controllers

import (
	"net/http"

	"github.com/FirstDayAtWork/mustracker/views"
)

type Static struct {
	Template *views.Template
	Data     *views.TemplateData
}

func (st Static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	st.Template.Execute(w, st.Data)
}
