// Package domain contiene la lógica de negocio pura (sin I/O), fácil de testear.
package domain

// Comision es el resultado de repartir un monto entre los niveles del MLM.
type Comision struct {
	Nivel      int
	Porcentaje float64
	Monto      float64
}

// TasasPorNivel son los porcentajes de comisión por nivel (10 %, 5 %, 2 %).
// Se exponen como slice ordenado por nivel ascendente.
var TasasPorNivel = []float64{0.10, 0.05, 0.02}

// RepartirComisiones divide un monto total entre los niveles de la red.
// Es una función pura: mismo input → mismo output, sin tocar la base de datos.
// El nivel i (1-based) recibe tasas[i-1] del total.
func RepartirComisiones(total float64, tasas []float64) []Comision {
	comisiones := make([]Comision, 0, len(tasas))
	for i, tasa := range tasas {
		comisiones = append(comisiones, Comision{
			Nivel:      i + 1,
			Porcentaje: tasa,
			Monto:      total * tasa,
		})
	}
	return comisiones
}
