package mold

//go:generate mockgen -source=$GOFILE -destination=./mock_${GOFILE} -package=$GOPACKAGE

type MaterialService interface {
	GetMaterial(requirements []*MaterialRequirement) (Material, error)
}
