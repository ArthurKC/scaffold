package mold

//go:generate mockgen -source=$GOFILE -destination=./mock_${GOFILE} -package=$GOPACKAGE

type ProductService interface {
	SaveProduct(destDir string, product *Product) error
}
