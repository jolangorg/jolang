# Jolang2

_Jolang2_ is Java transpiler and extender written on [_Go_](https://go.dev/)

## Try it (transpile Java to JS)

```shell
mkdir ~/Projects
cd ~/Projects
git clone https://github.com/jolang2/jolang2
git clone https://github.com/jbox2d/jbox2d
cd jolang2

# Transpile all Java files
go run jolang2/cmd/jo \
  -src="~/Projects/jbox2d/jbox2d-library/src/main/java" \
  -write-all

# You can also write AST files with -write-ast flag
go run jolang2/cmd/jo \
  -src="~/Projects/jbox2d/jbox2d-library/src/main/java" \
  -write-all
  -write-ast

python3 -m http.server &

xdg-open http://localhost:8000/html/
```

It's totally incomplete ='(