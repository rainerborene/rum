# rookie.go

Decode Rails 4 encrypted cookies in Go.

## Example

```go
var json map[string]interface{}
rookie := rookie.New(os.Getenv("SECRET_KEY_BASE"))
raw, _ := rookie.Decode(cookie)
json.Unmarshal(raw, &json)
```

You have to monkey patch Ruby on Rails' cookie based session store to use JSON
as its serializer instead of Marshal.

See [gist.github.com/jeffyip/4091166](https://gist.github.com/jeffyip/4091166)

## License

MIT
