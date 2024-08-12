/*
Example:

```
$ eg -t tools_eg/http_get/http_get_template.go -w tools_eg/http_get/http_get.go
```
*/
package main

import (
	"io"
	"io/ioutil"
)

func before(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}

func after(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
