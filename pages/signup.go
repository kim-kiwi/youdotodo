
package pages

import (
  "fmt"
  "net/http"
)

 func SignUp(w http.ResponseWriter, r *http.Request) { 
 page := Pages[1] 
 echo := Echo(w) 
fmt.Fprint(w, "\n")
fmt.Fprint(w, "<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n    <meta charset=\"UTF-8\" />\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=0.45\" />\n    <link rel=\"stylesheet\" href=\"static/style.css\">\n    <title>")
 fmt.Fprint(w, page.Title) 
fmt.Fprint(w, "</title>\n  </head>\n  <body>\n    <h1>")
 fmt.Fprintf(w, page.Title) 
fmt.Fprint(w, "</h1>\n")
fmt.Fprint(w, "\n")
fmt.Fprint(w, "<p>\n  ")
 for _, p := range Pages {
fmt.Fprint(w, "\n    ")
 if page.Link == p.Link { 
fmt.Fprint(w, "\n      <a href=")
 fmt.Fprint(w, p.Link) 
fmt.Fprint(w, " class=\"highlight\">")
 fmt.Fprint(w, p.Text) 
fmt.Fprint(w, "</a>\n    ")
 } else {
fmt.Fprint(w, "\n      <a href=")
 fmt.Fprint(w, p.Link) 
fmt.Fprint(w, ">")
 fmt.Fprint(w, p.Text) 
fmt.Fprint(w, "</a>\n    ")
 } 
fmt.Fprint(w, "\n  ")
 } 
fmt.Fprint(w, "\n  <a href=\"/logout\">logout</a>\n</p>\n")
fmt.Fprint(w, "\n\n<form>\n  <input>\n  <button type=\"submit\">signup</button>\n</form>\n\n<span></span>\n\n<script>\n  const form = document.querySelector(\"form\");\n  const input = document.querySelector(\"input\");\n  const span = document.querySelector(\"span\");\n  form.onsubmit = (e) => {\n      e.preventDefault();\n      fetch(\"")
 echo(api) 
fmt.Fprint(w, "/users\", {\n          method: \"POST\",\n          headers: { \"Content-Type\": \"application/json\" },\n          body: JSON.stringify({name: input.value}),\n      }).then(async (res)=>{\n          if (res.ok) {\n              span.textContent = \"successfully signed up. please log in.\";\n          } else {\n              span.textContent = await res.text();\n          }\n      }).catch(_=>{\n          span.textContent = \"something went wrong.\";\n      });\n  }\n</script>\n\n")
fmt.Fprint(w, "  </body>\n</html>\n")
fmt.Fprint(w, "\n")
 } 
