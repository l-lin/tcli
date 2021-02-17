//go:generate mockgen -source editor.go -package executor -destination editor_mock.go
package executor

type Editor interface {
	Edit(in []byte) (out []byte, err error)
}
