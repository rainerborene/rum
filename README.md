# rookie.go

Decode Rails 4 encrypted cookies in Go.

## Example

```go
rookie := rookie.New(os.Getenv("SECRET_KEY_BASE"))
data, _ := rookie.Decode(cookie)
params := data.(map[string]interface{})
params["user_id"]
```

## References

- [A little dip into Ruby's Marshal format](http://jakegoulding.com/blog/2013/01/15/a-little-dip-into-rubys-marshal-format/)
- [Another dip into Ruby's Marshal format](http://jakegoulding.com/blog/2013/01/16/another-dip-into-rubys-marshal-format/)
- [A final dip into Ruby's Marshal format](http://jakegoulding.com/blog/2013/01/20/a-final-dip-into-rubys-marshal-format/)
- [MRI Marshal Documentation](http://rxr.whitequark.org/mri/source/doc/marshal.rdoc)

## License

MIT
