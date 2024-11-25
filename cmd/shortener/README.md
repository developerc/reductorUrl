# cmd/shortener

В данной директории будет содержаться код, который скомпилируется в бинарное приложение
git add .
git commit -m 'next commit13'
git push

go run cmd/shortener/main.go

curl -X GET -i http://localhost:8080/0

cd cmd/shortener/
go test -v
$ git push --set-upstream origin iter2