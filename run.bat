IF NOT EXIST "pages" (
    mkdir pages
)

tt pages.tt/index.html.tt config.ini > pages/index.go
tt pages.tt/signup.html.tt config.ini > pages/signup.go
tt pages.tt/login.html.tt config.ini > pages/login.go

go run .
