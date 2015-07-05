### Passo a passo

- Faça download do SDK do `appengine` para `Go` em https://cloud.google.com/appengine/downloads
- Extraia o arquivo para algum lugar no seu file system
- Configure as variáveis de ambiente necessárias para disponibilizar o `goapp` e o `$GOPATH`
````bash
export APPENGINE=$HOME/path/to/go_appengine
export GOPATH=$HOME/path/to/go

export PATH=$PATH:$APPENGINE:$GOPATH
```

- Para baixar este repositório e suas dependências, digite no terminal:
```sh
$ goapp get github.com/rafbgarcia/go-api
```


### Faça requisições para a API
- Para criar um Post
```sh
$ curl -X POST http://go-api-rafa.appspot.com/posts -d '{"title": "My Awesome Post Title", "body": "An Amazing Post Body!"}' -H "Content-type: application/json"
```
- Para listar Posts
```sh
$ curl http://go-api-rafa.appspot.com/posts
```
