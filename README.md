## 🇺🇾 Uruguayan Truco (sync+async) game engine library (writen 100% in native Go)

Run the interactive example

  - using docker: `docker run -it --rm filevich/gotruco:latest -n 6 -timeout 60`

  - or simply compile and run: `go run cmd/example/*.go --timeout 120 -n 2`

![Example](https://i.imgur.com/e55VJwh.png "Example")

Other examples:

  - `cmd/bench` a simple benchmark to compare the speed and memory consumption 
  between other implementations of this same engine rewritten in other 
  languages. See `cmd/bench/README.md` for more info. 

  - `cmd/walker` example on how to traverse a **round-level** game tree