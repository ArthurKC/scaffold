package mold

//go:generate mockgen -source=$GOFILE -destination=./mock_${GOFILE} -package=$GOPACKAGE

type Repository interface {
	Save(m *Mold) error
	FindByName(name string) (*Mold, error)
}
