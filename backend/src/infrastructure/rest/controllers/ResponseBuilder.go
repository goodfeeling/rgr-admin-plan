package controllers

import "github.com/gbrayhan/microservices-go/src/domain"

func NewCommonResponseBuilder[T any]() *CommonResponseBuilder[T] {
	return &CommonResponseBuilder[T]{}
}

type CommonResponseBuilder[T any] struct {
	data    T
	message string
	status  int
}

func (b *CommonResponseBuilder[T]) Data(data T) *CommonResponseBuilder[T] {
	b.data = data
	return b
}

func (b *CommonResponseBuilder[T]) Message(msg string) *CommonResponseBuilder[T] {
	b.message = msg
	return b
}
func (b *CommonResponseBuilder[T]) Status(code int) *CommonResponseBuilder[T] {
	b.status = code
	return b
}

func (b *CommonResponseBuilder[T]) Build() domain.CommonResponse[T] {
	return domain.CommonResponse[T]{
		Data:    b.data,
		Message: b.message,
		Status:  b.status,
	}
}
