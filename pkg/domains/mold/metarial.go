package mold

//go:generate mockgen -source=$GOFILE -destination=./mock_${GOFILE} -package=$GOPACKAGE

type Material interface {
	Parameters() map[string]string
}
