### getcities
to run the program, call
```
go run main.go belarus ukrain
```
to run programm for double words countries, call
```
go run main.go "russian federation"
```
to compile on linux for another system:
GOOS=windows GOARCH=386 go build -o main386.exe main.go
