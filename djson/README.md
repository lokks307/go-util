# DJSON: Another JSON Library for who hates json.Unmarshal

## 1. HOW-TO

### 1.1. New JSON

```go
mJson := djson.NewDJSON()
```

### 1.2. Assign values to JSON

```go
// Array
mJson := djson.NewDJSON()
mJson.Put(djson.Array{
    1,2,3,4,5,6
})
```

```go
// Object
mJson := djson.NewDJSON()
mJson.Put(djson.Object{
    "name":  "Hery Victor",
    "idade": 32,
})
```

### 1.3. Assign table values to JSON
```go
mJson := djson.NewDJSON()
mJson.Put(
    djson.Array{
        djson.Object{
            "name":  "Ricardo Longa",
            "idade": 28,
            "skills": djson.Array{
                "Golang",
                "Android",
            },
        },
        djson.Object{
            "name":  "Hery Victor",
            "idade": 32,
            "skills": djson.Array{
                "Golang",
                "Java",
            },
        },
    },
)
```

### 1.4. Append values to existing JSON

```go
// Array
mJson := djson.NewDJSON()

mJson.Put(djson.Array{
    1,2,3,4,5,6,
})

mJson.Put(djson.Array{
    7,8,9,
})

fmt.Println(mJson.GetAsString()) // must be [1,2,3,4,5,6,7,8,9]
```

```go
// Object
mJson := djson.NewDJSON()

mJson.Put(djson.Object{
    "name":"Hery Victor",
})

mJson.Put(djson.Object{
    "idade": 32,
})

fmt.Println(mJson.GetAsString()) // must be like {"name":"Hery Victor","idade":55}
```

### 1.5. Parse existing JSON string

```go
jsonDoc := `[
    {
        "name":"Ricardo Longa",
        "idade":28,
        "skills":[
            "Golang","Android"
        ]
    },
    {
        "name":"Hery Victor",
        "idade":32,
        "skills":[
            "Golang",
            "Java"
        ]
    }
]`

mJson := NewDJSON().Parse(jsonDoc)
```

