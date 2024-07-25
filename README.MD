# Jolang

_Jolang_ is Java transpiler and extender written on [_Go_](https://go.dev/)

## Try it (transpile Java to JS)

```shell
mkdir ~/Projects
cd ~/Projects
git clone https://github.com/jolangorg/jolang
git clone https://github.com/jbox2d/jbox2d
cd jolang

# Transpile all Java files
go run github.com/jolangorg/jolang/cmd/jo \
  -src="~/Projects/jbox2d/jbox2d-library/src/main/java" \
  -write-all

# You can also write AST files with -write-ast flag
go run github.com/jolangorg/jolang/cmd/jo \
  -src="~/Projects/jbox2d/jbox2d-library/src/main/java" \
  -write-all
  -write-ast

python3 -m http.server &

xdg-open http://localhost:8000/html/
```

It's totally incomplete ='(