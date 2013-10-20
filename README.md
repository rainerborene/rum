# rookie.go

Decode Rails 4 encrypted cookies in Go.

## Usage

```go
rookie := rookie.New(os.Getenv("SECRET_KEY_BASE"))
rookie.Iterations = 500
rookie.CookieSalt = []byte("cookie salt")
rookie.CookieSaltLength = 32
data, _ := rookie.Decode(cookie)
data.Get("user_id")
```

`rookie.New` initialize the above values with Rails defaults.

## License

MIT
