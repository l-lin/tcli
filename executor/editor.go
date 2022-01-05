//go:generate mockgen -source editor.go -package executor -destination editor_mock.go
package executor

const (
	markdownFileType = "md"
)

type Editor interface {
	Edit(in []byte, fileType string) (out []byte, err error)
}
