# Changelog

## 2.0.0 (2017-05-01)

- This was a big update and while it came too soon. It was as a result of some design tradeoff v1.0.0 made.
- The major change was to stop the use of `interface{}` and migrate to `[]byte`. This offers much more type safety.

## 1.0.0 (2017-04-29)

- Initial release
- Stores
 
  - Inmemory
  - Redis
  - Memcached
  - Filesystem

