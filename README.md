# Go Sankey

This is a tool that can generate you a sankey chart from a single file

File format:

```
[[Config]]
ConnectionOpacity=50
Width=1920
Height=1080

[[Bars]]
[Title Here]
ID=Title By Default | STRING
Step=INT
Color=random | HEX

[[Connections]]
ID 1 -> ID 2: 55
ID 2 -> ID 3: 222
```