package main

type RequestError struct {
  Message string
}

func (re *RequestError) Error() string {
  return re.Message
}
