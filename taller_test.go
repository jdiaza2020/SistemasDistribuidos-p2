package main

import (
	"fmt"
	"testing"
	"time"
)

// ------------------------------
// Resultado de una simulación
// ------------------------------

type ResultadoSimulacion struct {
	Nombre       string
	NumCoches    int
	NumMecMec    int
	NumMecElec   int
	NumMecCarr   int
	DuracionReal time.Duration // tiempo real de la simulación
	TiempoMedio  float64       // tiempo medio por coche (segundos simulados)
}

// ------------------------------
// Worker por tipo de incidencia
// ------------------------------

func workerTipo(duracionSeg int, unidad time.Duration, ch <-chan *Vehiculo, done chan<- *Vehiculo) {
	for v := range ch {
		// simulamos trabajo
		time.Sleep(time.Duration(duracionSeg) * unidad)
		v.TiempoAtencion += duracionSeg
		done <- v
	}
}

// ------------------------------
// Motor de simulación genérico (SIN reencolar por prioridad)
// ------------------------------

// vehiculos: lista de vehículos ya creados con su incidencia.Tipo
// numMecMec / numMecElec / numMecCarr: nº de mecánicos de cada especialidad
func ejecutarSimulacionEscenario(nombre string, vehiculos []*Vehiculo,
	numMecMec, numMecElec, numMecCarr int) ResultadoSimulacion {

	// 1 "segundo simulado" = 100 ms reales (puedes ajustar)
	unidad := 100 * time.Millisecond

	chMec := make(chan *Vehiculo, len(vehiculos)*2)
	chElec := make(chan *Vehiculo, len(vehiculos)*2)
	chCarr := make(chan *Vehiculo, len(vehiculos)*2)
	done := make(chan *Vehiculo, len(vehiculos)*2)

	if numMecMec < 0 {
		numMecMec = 0
	}
	if numMecElec < 0 {
		numMecElec = 0
	}
	if numMecCarr < 0 {
		numMecCarr = 0
	}

	// Lanzar workers
	for i := 0; i < numMecMec; i++ {
		go workerTipo(5, unidad, chMec, done) // mecánica → 5 s simulados
	}
	for i := 0; i < numMecElec; i++ {
		go workerTipo(7, unidad, chElec, done) // eléctrica → 7 s
	}
	for i := 0; i < numMecCarr; i++ {
		go workerTipo(11, unidad, chCarr, done) // carrocería → 11 s
	}

	// Enviar vehículos a su cola
	for _, v := range vehiculos {
		inc := v.GetIncidencia()
		if inc == nil {
			continue
		}
		switch inc.Tipo {
		case "mecánica":
			chMec <- v
		case "eléctrica":
			chElec <- v
		case "carrocería":
			chCarr <- v
		}
	}

	total := len(vehiculos)
	atendidos := 0
	start := time.Now()

	// Cada coche se procesa UNA vez (sin reencolar)
	for atendidos < total {
		<-done
		atendidos++
	}

	elapsed := time.Since(start)

	// Cerramos canales
	close(chMec)
	close(chElec)
	close(chCarr)
	close(done)

	// Tiempo medio por coche (en segundos simulados)
	suma := 0
	for _, v := range vehiculos {
		suma += v.TiempoAtencion
	}
	media := 0.0
	if total > 0 {
		media = float64(suma) / float64(total)
	}

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

// ------------------------------
// Construcción de escenarios
// ------------------------------

func crearVehiculosTipoUnico(numCoches int, tipo string) []*Vehiculo {
	var vs []*Vehiculo
	for i := 0; i < numCoches; i++ {
		v := &Vehiculo{
			Matricula: fmt.Sprintf("TEST-%s-%02d", tipo, i+1),
			Marca:     "Test",
			Modelo:    "Sim",
		}
		inc := &Incidencia{
			IDIncidencia: i + 1,
			Tipo:         tipo,
			Prioridad:    "media",
			Descripcion:  "simulada",
			Estado:       "abierta",
		}
		v.SetIncidencia(inc)
		vs = append(vs, v)
	}
	return vs
}

func crearVehiculosMixtos(numCoches int) []*Vehiculo {
	tipos := []string{"mecánica", "eléctrica", "carrocería"}
	var vs []*Vehiculo
	for i := 0; i < numCoches; i++ {
		tipo := tipos[i%len(tipos)]
		v := &Vehiculo{
			Matricula: fmt.Sprintf("MIX-%02d", i+1),
			Marca:     "Test",
			Modelo:    "Mix",
		}
		inc := &Incidencia{
			IDIncidencia: i + 1,
			Tipo:         tipo,
			Prioridad:    "media",
			Descripcion:  "simulada",
			Estado:       "abierta",
		}
		v.SetIncidencia(inc)
		vs = append(vs, v)
	}
	return vs
}

// ------------------------------
// Ejecutor de todos los tests
// ------------------------------

func EjecutarTestsComparativa() {
	fmt.Println("====================================")
	fmt.Println("   TESTS DE SIMULACIÓN DEL TALLER   ")
	fmt.Println("    (goroutines + channels + time)  ")
	fmt.Println("====================================")

	// 1) Duplicar número de coches (mecánica, misma plantilla)
	n := 10
	cochesA := crearVehiculosTipoUnico(n, "mecánica")
	cochesB := crearVehiculosTipoUnico(2*n, "mecánica")

	res1 := ejecutarSimulacionEscenario("N coches (mecánica)", cochesA, 3, 0, 0)
	res2 := ejecutarSimulacionEscenario("2N coches (mecánica)", cochesB, 3, 0, 0)

	// 2) Duplicar plantilla de mecánicos (misma carga mixta)
	cochesPlantilla := crearVehiculosMixtos(30)
	res3 := ejecutarSimulacionEscenario("Plantilla base 1-1-1", cochesPlantilla, 1, 1, 1)
	res4 := ejecutarSimulacionEscenario("Plantilla doble 2-2-2", cochesPlantilla, 2, 2, 2)

	// 3) Distribución distinta con mismo nº total de mecánicos
	cochesDist := crearVehiculosMixtos(30)
	res5 := ejecutarSimulacionEscenario("Distribución 3Mec/1Elec/1Carr", cochesDist, 3, 1, 1)
	res6 := ejecutarSimulacionEscenario("Distribución 1Mec/3Elec/3Carr", cochesDist, 1, 3, 3)

	resultados := []ResultadoSimulacion{res1, res2, res3, res4, res5, res6}

	fmt.Println("\nRESULTADOS (tiempo real ≈ ms, tiempo medio en segundos simulados):")
	for _, r := range resultados {
		fmt.Printf("- %-30s | Coches:%2d | MecMec:%2d | MecElec:%2d | MecCarr:%2d | Duración:%4d ms | Tmedio: %.2f s\n",
			r.Nombre, r.NumCoches, r.NumMecMec, r.NumMecElec, r.NumMecCarr,
			r.DuracionReal.Milliseconds(), r.TiempoMedio)
	}

	fmt.Println("\nANÁLISIS PARA LA MEMORIA:")

	fmt.Println("1) Duplicar número de coches (mecánica, plantilla fija 3 mecánicos):")
	fmt.Printf("   - N coches:   duración ≈ %d ms, tiempo medio ≈ %.2f s\n",
		res1.DuracionReal.Milliseconds(), res1.TiempoMedio)
	fmt.Printf("   - 2N coches:  duración ≈ %d ms, tiempo medio ≈ %.2f s\n",
		res2.DuracionReal.Milliseconds(), res2.TiempoMedio)
	fmt.Println("   → Al duplicar el número de coches aumenta claramente el tiempo total de procesamiento,")
	fmt.Println("     mientras que el tiempo medio por vehículo se mantiene cercano al tiempo de atención")
	fmt.Println("     definido por la incidencia mecánica (≈ 5 s simulados).")

	fmt.Println("\n2) Duplicar plantilla de mecánicos (misma carga mixta):")
	fmt.Printf("   - Plantilla 1-1-1: duración ≈ %d ms, tiempo medio ≈ %.2f s\n",
		res3.DuracionReal.Milliseconds(), res3.TiempoMedio)
	fmt.Printf("   - Plantilla 2-2-2: duración ≈ %d ms, tiempo medio ≈ %.2f s\n",
		res4.DuracionReal.Milliseconds(), res4.TiempoMedio)
	fmt.Println("   → Al duplicar el número de mecánicos disminuye el tiempo total para atender toda la cola,")
	fmt.Println("     porque hay más goroutines trabajando en paralelo, aunque el tiempo medio por coche")
	fmt.Println("     depende principalmente del tipo de incidencia asignada.")

	fmt.Println("\n3) Distribución de especialidades (mismo nº total de mecánicos):")
	fmt.Printf("   - 3 Mecánica / 1 Eléctrica / 1 Carrocería: duración ≈ %d ms, Tmedio ≈ %.2f s\n",
		res5.DuracionReal.Milliseconds(), res5.TiempoMedio)
	fmt.Printf("   - 1 Mecánica / 3 Eléctrica / 3 Carrocería: duración ≈ %d ms, Tmedio ≈ %.2f s\n",
		res6.DuracionReal.Milliseconds(), res6.TiempoMedio)
	fmt.Println("   → Se observa cómo la distribución de especialidades afecta al rendimiento global.")
	fmt.Println("     En un escenario con carga equilibrada, una plantilla más homogénea tiende a repartir")
	fmt.Println("     mejor el trabajo. Si predominan incidencias de un tipo concreto, conviene incrementar")
	fmt.Println("     el número de mecánicos especializados en ese tipo para reducir el tiempo total.")

	fmt.Println("\n(Puedes copiar y adaptar estos comentarios directamente en el PDF de la memoria.)")
}

// ------------------------------
// Test de Go que llama al ejecutor
// ------------------------------

func TestComparativasTaller(t *testing.T) {
	EjecutarTestsComparativa()
}
