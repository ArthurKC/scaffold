package mold

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOFILE -package=$GOPACKAGE

type Service interface {
	ImportFrom(path string, dstName string) (*Mold, error)
}
