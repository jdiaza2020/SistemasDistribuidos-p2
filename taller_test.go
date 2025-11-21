package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// RESULTADOS DE SIMULACIÓN

type ResultadoSimulacion struct {
	Nombre       string
	NumCoches    int
	NumMecMec    int
	NumMecElec   int
	NumMecCarr   int
	DuracionReal time.Duration
	TiempoMedio  float64
}

// WORKERS

func workerTipo(segundos int, unidad time.Duration, ch <-chan *Vehiculo, done chan<- *Vehiculo) {
	for v := range ch {
		time.Sleep(time.Duration(segundos) * unidad)
		v.TiempoAtencion += segundos
		done <- v
	}
}

// MOTOR DE SIMULACIÓN BASE

func ejecutarSimulacionEscenario(nombre string, vehiculos []*Vehiculo,
	numMecMec, numMecElec, numMecCarr int) ResultadoSimulacion {

	// 100 ms = 1 segundo simulado
	unidad := 100 * time.Millisecond

	chMec := make(chan *Vehiculo, len(vehiculos)*2)
	chElec := make(chan *Vehiculo, len(vehiculos)*2)
	chCarr := make(chan *Vehiculo, len(vehiculos)*2)
	done := make(chan *Vehiculo, len(vehiculos)*5)

	// Lanzar workers
	for i := 0; i < numMecMec; i++ {
		go workerTipo(5, unidad, chMec, done)
	}
	for i := 0; i < numMecElec; i++ {
		go workerTipo(7, unidad, chElec, done)
	}
	for i := 0; i < numMecCarr; i++ {
		go workerTipo(11, unidad, chCarr, done)
	}

	// Encolar coches según incidencia
	for _, v := range vehiculos {
		inc := v.GetIncidencia()
		if inc == nil {
			continue
		}

		tipo := strings.ToLower(strings.TrimSpace(inc.Tipo))
		switch tipo {
		case "mecánica", "mecanica":
			chMec <- v
		case "eléctrica", "electrica":
			chElec <- v
		case "carrocería", "carroceria":
			chCarr <- v
		default:
			fmt.Println("⚠ Incidencia desconocida:", inc.Tipo)
		}
	}

	total := len(vehiculos)
	atendidos := 0
	start := time.Now()

	// NO re-encolamos por prioridad: queremos que el test siempre termine
	for atendidos < total {
		v := <-done
		_ = v // para evitar warning si no lo usamos
		atendidos++
	}

	elapsed := time.Since(start)

	close(chMec)
	close(chElec)
	close(chCarr)
	close(done)

	// Tiempo medio simulado
	suma := 0
	for _, v := range vehiculos {
		suma += v.TiempoAtencion
	}

	media := float64(suma) / float64(total)

	return ResultadoSimulacion{
		Nombre:       nombre,
		NumCoches:    total,
		NumMecMec:    numMecMec,
		NumMecElec:   numMecElec,
		NumMecCarr:   numMecCarr,
		DuracionReal: elapsed,
		TiempoMedio:  media,
	}
}

// CREACIÓN DE ESCENARIOS

func crearVehiculosTipoUnico(numCoches int, tipo string) []*Vehiculo {
	var vs []*Vehiculo
	for i := 0; i < numCoches; i++ {
		v := &Vehiculo{
			Matricula: fmt.Sprintf("%s-%02d", tipo, i+1),
			Marca:     "Test",
			Modelo:    "Sim",
		}
		v.SetIncidencia(&Incidencia{
			IDIncidencia: i + 1,
			Tipo:         tipo,
			Prioridad:    "media",
			Descripcion:  "simulada",
			Estado:       "abierta",
		})
		vs = append(vs, v)
	}
	return vs
}

func crearVehiculosMixtos(numCoches int) []*Vehiculo {
	tipos := []string{"mecánica", "eléctrica", "carrocería"}
	var vs []*Vehiculo
	for i := 0; i < numCoches; i++ {
		tipo := tipos[i%3]
		v := &Vehiculo{
			Matricula: fmt.Sprintf("MIX-%02d", i+1),
			Marca:     "Test",
			Modelo:    "Mix",
		}
		v.SetIncidencia(&Incidencia{
			IDIncidencia: i + 1,
			Tipo:         tipo,
			Prioridad:    "media",
			Descripcion:  "simulada",
			Estado:       "abierta",
		})
		vs = append(vs, v)
	}
	return vs
}

// TEST PRINCIPAL

func TestComparativasTaller(t *testing.T) {
	fmt.Println("====================================")
	fmt.Println("   TESTS DE SIMULACIÓN DEL TALLER   ")
	fmt.Println("    (goroutines + channels + time)  ")
	fmt.Println("====================================")

	n := 10

	// 1) Duplicar nº de coches, misma plantilla
	res1 := ejecutarSimulacionEscenario("N coches (mecánica)", crearVehiculosTipoUnico(n, "mecánica"), 3, 0, 0)
	res2 := ejecutarSimulacionEscenario("2N coches (mecánica)", crearVehiculosTipoUnico(2*n, "mecánica"), 3, 0, 0)

	// 2) Duplicar plantilla, misma carga
	mix1 := crearVehiculosMixtos(30)
	mix2 := crearVehiculosMixtos(30)
	res3 := ejecutarSimulacionEscenario("Plantilla 1-1-1", mix1, 1, 1, 1)
	res4 := ejecutarSimulacionEscenario("Plantilla 2-2-2", mix2, 2, 2, 2)

	// 3) Distintas distribuciones de especialidad (mismo nº total de mecánicos)
	dist1 := crearVehiculosMixtos(30)
	dist2 := crearVehiculosMixtos(30)
	res5 := ejecutarSimulacionEscenario("3M/1E/1C", dist1, 3, 1, 1)
	res6 := ejecutarSimulacionEscenario("1M/3E/3C", dist2, 1, 3, 3)

	fmt.Println("\nRESULTADOS:")
	printRes(res1)
	printRes(res2)
	printRes(res3)
	printRes(res4)
	printRes(res5)
	printRes(res6)

	fmt.Println("\n(Estos datos pueden copiarse directamente a la memoria)")
}

func printRes(r ResultadoSimulacion) {
	fmt.Printf("- %-20s | Coches:%2d | MecMec:%d | MecElec:%d | MecCarr:%d | Duración:%4d ms | Tmedio: %.2f s\n",
		r.Nombre, r.NumCoches, r.NumMecMec, r.NumMecElec, r.NumMecCarr,
		r.DuracionReal.Milliseconds(), r.TiempoMedio)
}

// BENCHMARKS

func BenchmarkDuplicarCochesMecanica(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ejecutarSimulacionEscenario("N coches", crearVehiculosTipoUnico(10, "mecánica"), 3, 0, 0)
		ejecutarSimulacionEscenario("2N coches", crearVehiculosTipoUnico(20, "mecánica"), 3, 0, 0)
	}
}

func BenchmarkDuplicarPlantilla(b *testing.B) {
	for i := 0; i < b.N; i++ {
		coches1 := crearVehiculosMixtos(30)
		coches2 := crearVehiculosMixtos(30)
		ejecutarSimulacionEscenario("1-1-1", coches1, 1, 1, 1)
		ejecutarSimulacionEscenario("2-2-2", coches2, 2, 2, 2)
	}
}

func BenchmarkDistribucionEspecialidades(b *testing.B) {
	for i := 0; i < b.N; i++ {
		coches1 := crearVehiculosMixtos(30)
		coches2 := crearVehiculosMixtos(30)
		ejecutarSimulacionEscenario("3M1E1C", coches1, 3, 1, 1)
		ejecutarSimulacionEscenario("1M3E3C", coches2, 1, 3, 3)
	}
}
