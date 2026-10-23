package domain

import "testing"

func TestRepartirComisiones(t *testing.T) {
	got := RepartirComisiones(100, TasasPorNivel)

	if len(got) != 3 {
		t.Fatalf("esperaba 3 niveles, obtuve %d", len(got))
	}
	// 10 % de 100 = 10 (nivel 1), 5 % = 5 (nivel 2), 2 % = 2 (nivel 3).
	esperado := []float64{10, 5, 2}
	for i, c := range got {
		if c.Nivel != i+1 {
			t.Errorf("nivel[%d] = %d, quiero %d", i, c.Nivel, i+1)
		}
		if c.Monto != esperado[i] {
			t.Errorf("monto[%d] = %.2f, quiero %.2f", i, c.Monto, esperado[i])
		}
	}
}

func TestRepartirComisionesVacio(t *testing.T) {
	if got := RepartirComisiones(50, nil); len(got) != 0 {
		t.Fatalf("sin tasas no debe haber comisiones, obtuve %d", len(got))
	}
}
