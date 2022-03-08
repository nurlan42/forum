# Forum

> Every day is a learning day.. 
___

### Table of Contents

---

- [Description](#description)
- [Author Info](#author-info)
- [How-To-Use](#how-to-use)

---

## Description:
    This project consists in creating a web forum that allows:
        - communication between users.
        - associating categories to posts.
        - liking and disliking posts and comments.
        - filtering posts.

#### Technologies

    - Go Verison: 1.16.3
    - SQLite3
    - HTML, CSS
    - Docker
    - Makefile

---

## Author Info: 
    Nurlan
    Gmail - [nurlan.ikhsan@gmail.com]
    Discord - Nurlan #9886

## How-To-Use project: 
Run project: `make run` <br>
Build bin: `make build` <br>
Run from Docker: `make docker`<br>
Stop docker: `make stop` <br>

## Use sequre Database:
First build with auth extension: `go build --tags sqlite_userauth -o auth ./cmd` <br>
Run project: `./auth` <br>

other auth manipulations see go.doc: `https://pkg.go.dev/github.com/mattn/go-sqlite3` <br>
section: `Type SQLITEConn`

---

 
- [Back To The Top](#forum)

---

## Licence 

Alem School Licence 

Copyright &copy [2022] [Nurlan]

Permission is hereby granted, free of charge, to any person obtaining a copy of this 
software and associated documentation files(the "Software"), to deal in the Software 
without restriction, including the rights to use, copy, modify, merge, publish or sell. 

- [Back To The Top](#forum)

---