
package pages

import (
  "fmt"
  "net/http"
  . "youdotodo/repo"
)

 func Index(w http.ResponseWriter, r *http.Request, logged_in bool, my_uid int64, todosByUser map[int64][]TodoDto) { 
 page := Pages[0] 
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
fmt.Fprint(w, "\n\n")
 if logged_in { 
fmt.Fprint(w, "\n\n<!-- <img src=\"static/웃는아기.png\"> -->\n<!-- 어서오세요 -->\n\n<style>\n  .cal {\n      gap: 0.1rem;\n      display: grid;\n      grid-template-columns: repeat(7, 1fr);\n      width: 250px;\n      background: var(--code-bg);\n  }\n  .cal-item:hover,\n  .cal-item.selected {\n      background: var(--fg);\n      color: var(--bg);\n  }\n  .done span {\n      text-decoration: line-through;\n      color: var(--small-fg);\n  }\n</style>\n\n<div>\n  <input id=\"year-month\" type=\"month\">\n  <div id=\"cal\" class=\"cal\"></div>\n</div>\n\n<hr>\n\n<form id=\"form\">\n  <input id=\"todo-input\">\n  <button type=\"submit\">추가</button>\n</form>\n<span id=\"msg\"></span>\n\n<br>\n\n")
 renderTodos := func(todos []TodoDto, delBtn bool) { 
fmt.Fprint(w, "\n<ol>\n  ")
 for _, t := range todos { 
fmt.Fprint(w, "\n    <li>\n      <div style=\"display:flex;\" ")
 if t.Done { 
fmt.Fprint(w, "class=\"done\"")
 } 
fmt.Fprint(w, ">\n        <input type=\"checkbox\" oninput=\"doneTodo(")
 echo(t.Id) 
fmt.Fprint(w, ", this.checked)\" ")
 if t.Done { 
fmt.Fprint(w, "checked")
 } 
fmt.Fprint(w, " ")
 if !delBtn { 
fmt.Fprint(w, "disabled")
 } 
fmt.Fprint(w, ">\n        <span>")
 echo(t.Content) 
fmt.Fprint(w, "</span>\n        ")
 if delBtn { 
fmt.Fprint(w, "\n          <button style=\"margin-left:auto;\" onclick=\"deleteTodo(")
 echo(t.Id) 
fmt.Fprint(w, ")\">X</button>\n        ")
 } 
fmt.Fprint(w, "\n      </div>\n    </li>\n  ")
 } 
fmt.Fprint(w, "\n</ol>\n")
 } 
fmt.Fprint(w, "\n\n<ol>\n  <fieldset>\n    <legend>You</legend>\n    ")
 renderTodos(todosByUser[my_uid], true) 
fmt.Fprint(w, "\n  </fieldset>\n  ")
 for uid, todos := range todosByUser { 
fmt.Fprint(w, "\n    ")
 if uid != my_uid { 
fmt.Fprint(w, "\n      <fieldset>\n        <legend>")
 echo(todos[0].UserName) 
fmt.Fprint(w, "</legend>\n        ")
 renderTodos(todos, false) 
fmt.Fprint(w, "\n      </fieldset>\n    ")
 } 
fmt.Fprint(w, "\n  ")
 } 
fmt.Fprint(w, "\n</ol>\n\n<script>\n")
fmt.Fprint(w, "const cal = document.querySelector(\"#cal\");\nconst week = [\"일\",\"월\",\"화\",\"수\",\"목\",\"금\",\"토\"];\n\nfunction get_max_day(year, month) {\n    return new Date(year, month, 0).getDate();\n}\nfunction render_cal(year, month, day) {\n    cal.innerHTML = \"\";\n    for (const w of week) {\n        const div = document.createElement(\"div\");\n        div.style.fontWeight = \"bold\";\n        div.textContent = w;\n        cal.append(div);\n    }\n    const max_day = get_max_day(year, month);\n    const date = new Date(year,month-1,1);\n    const day_number = date.getDay();\n    for (let i = 0; i < day_number; ++i) {\n        const div = document.createElement(\"div\");\n        div.textContent = \"\";\n        cal.append(div);\n    }\n    for (let i = 1; i <= max_day; ++i) {\n        const div = document.createElement(\"div\");\n        div.classList.add(\"cal-item\");\n        if (i==day) div.classList.add(\"selected\");\n        div.addEventListener(\"click\", () => {\n            location.href=`?year=${year}&month=${month}&day=${i}`;\n        });\n        div.textContent = i;\n        cal.append(div);\n    }\n}\n\nconst today = new Date();\nconst params = new URLSearchParams(window.location.search);\nconst year = parseInt(params.get(\"year\"), 10);\nconst month = parseInt(params.get(\"month\"), 10);\nconst day = parseInt(params.get(\"day\"), 10);\n\nconst year_month = document.querySelector(\"#year-month\");\nyear_month.addEventListener(\"input\", (e)=>{\n    const [year, month] = year_month.value.split(\"-\");\n    location.href = `?year=${year}&month=${month}&day=1`;\n});\n\nyear_month.value = `${year}-${String(month).padStart(2,'0')}`;\nrender_cal(year,month,day);\n\nconst form = document.querySelector(\"#form\");\nconst input = document.querySelector(\"#todo-input\");\ninput.focus();\nform.addEventListener(\"submit\", (e) => {\n    e.preventDefault();\n    addTodo(input.value, year, month, day);\n});\n\nconst span = document.querySelector(\"#msg\");\nconst nextIdx = ")
 echo(len(todosByUser[my_uid])) 
fmt.Fprint(w, ";\nfunction addTodo(content, year, month, day) {\n    fetch(\"")
 echo(api) 
fmt.Fprint(w, "/todos\", {\n        method: \"POST\",\n        headers: { \"Content-Type\": \"application/json\" },\n        body: JSON.stringify({content,year,month,day,idx:nextIdx}),\n    }).then(async res=>{\n        if (res.ok) {\n            location.reload();\n        } else {\n            span.textContent = await res.text();\n        }\n    });\n}\n\nfunction doneTodo(id,done) {\n    fetch(`")
 echo(api) 
fmt.Fprint(w, "/todos/${id}`,{\n        method:\"PATCH\",\n        headers: { \"Content-Type\": \"application/json\" },\n        body: JSON.stringify({done}),\n    }).then(async (res)=>{\n        if (res.ok) {\n            location.reload();\n        } else {\n            const msg = await res.text();\n            alert(msg);\n        }\n    });\n}\n\nfunction deleteTodo(id) {\n    fetch(`")
 echo(api) 
fmt.Fprint(w, "/todos/${id}`,{\n        method:\"DELETE\",\n    }).then(async (res)=>{\n        if (res.ok) {\n            location.reload();\n        } else {\n            const msg = await res.text();\n            alert(msg);\n        }\n    });\n}\n")
fmt.Fprint(w, "\n</script>\n\n")
 } else { 
fmt.Fprint(w, "\n\n<img src=\"static/우는아기.png\">\n로그인해주세요\n\n")
 } 
fmt.Fprint(w, "\n\n")
fmt.Fprint(w, "  </body>\n</html>\n")
fmt.Fprint(w, "\n")
 } 
