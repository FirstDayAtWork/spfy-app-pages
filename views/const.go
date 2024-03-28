package views

// CSS Paths
const (
	TemplateCSS = "../static/styles/template.css"
	BaseCSS     = "../static/styles/style.css"
	RegisterCSS = "../static/styles/register-style.css"
)

// JS Paths
const (
	TemplateJS = "../static/scripts/template-script.js"
	LoginJS    = "../static/scripts/app-login.js"
	RegisterJS = "../static/scripts/app-register.js"
	AboutJS    = "../static/scripts/app-about.js"
	DonateJS   = "../static/scripts/app-donate.js"
)

// Templates
const (
	BaseTemplate        = "base.html"
	RegisterTemplate    = "register.html"
	LoginTemplate       = "login.html"
	IndexTemplate       = "index.html"
	AboutTemplate       = "about.html"
	DonateTemplate      = "donate.html"
	GuestNavbarTemplate = "navbar-guest.html"
	UserNavbarTemplate  = "navbar-user.html"
)

// Titles
const (
	RegisterTitle = "Register"
	IndexTitle    = "Index"
	LoginTitle    = "Login"
	AboutTitle    = "About"
	DonateTitle   = "Donate"
)

// Error messages
const (
	TemplateRenderError = "Error rendering page template"
	NoFilePathsError    = "no filepaths provided"
)
