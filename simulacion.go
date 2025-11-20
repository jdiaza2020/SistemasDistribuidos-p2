package main

import (
	"fmt"
	"time"
)

// WORKER DE MECÁNICO (GOROUTINE)

func workerMecanico(m *Mecanico, cChan <-chan *Vehiculo, doneChan chan<- *Vehiculo) {
	for coche := range cChan {

		inc := coche.GetIncidencia()
		if inc == nil {
			continue
		}

		switch inc.Tipo {
		case "mecánica":
			time.Sleep(5 * time.Second)
			coche.TiempoAtencion += 5
		case "eléctrica":
			time.Sleep(7 * time.Second)
			coche.TiempoAtencion += 7
		case "carrocería":
			time.Sleep(11 * time.Second)
			coche.TiempoAtencion += 11
		default:
			// tipo desconocido, no hacemos nada
		}

		doneChan <- coche
	}
}

// SIMULACIÓN PRINCIPAL

func simularAtencionConcurrente() {
	// Si no hay datos, generamos un escenario de prueba
	if len(app.ClientesTaller) == 0 || len(app.MecanicosTaller) == 0 {
		fmt.Println("No hay datos en el taller. Generando datos de prueba para la simulación...")
		cargarDatosPruebaSimulacion()
	}

	// Recoger todos los vehículos con incidencia
	var coches []*Vehiculo
	for _, c := range app.ClientesTaller {
		for _, v := range c.Vehiculos {
			if v.GetIncidencia() != nil {
				v.TiempoAtencion = 0 // reset por si se lanza varias veces
				coches = append(coches, v)
			}
		}
	}

	if len(coches) == 0 {
		fmt.Println("No hay vehículos con incidencia para simular.")
		return
	}

	// Crear canales
	cChan := make(chan *Vehiculo, len(coches))
	doneChan := make(chan *Vehiculo, len(coches))

	// Lanzar goroutines de mecánicos activos
	for _, m := range app.MecanicosTaller {
		if m.Activo {
			go workerMecanico(m, cChan, doneChan)
		}
	}

	// Enviar todos los coches iniciales a la cola
	for _, v := range coches {
		cChan <- v
	}

	// Procesar resultados y aplicar prioridad
	atendidos := 0
	totalInicial := len(coches)

	for atendidos < totalInicial {
		coche := <-doneChan

		fmt.Printf("\n Vehículo %s atendido (%d segundos acumulados)\n",
			coche.Matricula, coche.TiempoAtencion)

		// PRIORIDAD: si supera 15 segundos de atención total
		if coche.TiempoAtencion > 15 {
			fmt.Printf(" Vehículo %s supera 15s: prioridad activada.\n", coche.Matricula)

			// Buscar mecánico de la especialidad de la incidencia
			inc := coche.GetIncidencia()
			mec := buscarMecanicoEspecialidad(inc.Tipo)

			if mec == nil {
				// No hay mecánico disponible de esa especialidad: contratamos uno nuevo
				idNuevo := len(app.MecanicosTaller) + 1
				nuevo := &Mecanico{
					IDMecanico:       idNuevo,
					Nombre:           fmt.Sprintf("Mec%d", idNuevo),
					Especialidad:     inc.Tipo,
					AniosExperiencia: 1,
					Activo:           true,
				}
				app.MecanicosTaller = append(app.MecanicosTaller, nuevo)
				go workerMecanico(nuevo, cChan, doneChan)
				fmt.Printf(" Contratado mecánico %s de especialidad %s\n",
					nuevo.Nombre, nuevo.Especialidad)
			}

			// Reencolamos el coche con prioridad
			cChan <- coche
			continue
		}

		// Si ya no necesita prioridad, lo contamos como finalizado
		atendidos++
	}

	// Cerrar canales
	close(cChan)
	close(doneChan)

	fmt.Println("\n Simulación finalizada")
}

// AYUDAS

func buscarMecanicoEspecialidad(tipo string) *Mecanico {
	for _, m := range app.MecanicosTaller {
		if m.Activo && m.Especialidad == tipo {
			return m
		}
	}
	return nil
}

// Genera un conjunto de datos de prueba para la simulación:
// - 3 mecánicos (uno de cada especialidad)
// - 2 clientes
// - 3 vehículos con incidencias (mecánica, eléctrica, carrocería)
func cargarDatosPruebaSimulacion() {
	// Mecánicos
	app.MecanicosTaller = []*Mecanico{
		{IDMecanico: 1, Nombre: "Laura", Especialidad: "mecánica", AniosExperiencia: 3, Activo: true},
		{IDMecanico: 2, Nombre: "Pedro", Especialidad: "eléctrica", AniosExperiencia: 5, Activo: true},
		{IDMecanico: 3, Nombre: "Ana", Especialidad: "carrocería", AniosExperiencia: 4, Activo: true},
	}
	app.InicializarPlazas()

	// Clientes
	c1 := &Cliente{IDCliente: 1, Nombre: "Carlos", Telefono: "111111111", Email: "carlos@ejemplo.com"}
	c2 := &Cliente{IDCliente: 2, Nombre: "Lucía", Telefono: "222222222", Email: "lucia@ejemplo.com"}

	// Vehículos de prueba con incidencias
	v1 := &Vehiculo{Matricula: "AAA111", Marca: "Seat", Modelo: "Ibiza", FechaEntrada: "01/01/2025"}
	inc1 := &Incidencia{
		IDIncidencia: 1,
		Tipo:         "mecánica",
		Prioridad:    "media",
		Descripcion:  "motor",
		Estado:       "abierta",
	}
	v1.SetIncidencia(inc1)

	v2 := &Vehiculo{Matricula: "BBB222", Marca: "Renault", Modelo: "Clio", FechaEntrada: "01/01/2025"}
	inc2 := &Incidencia{
		IDIncidencia: 2,
		Tipo:         "eléctrica",
		Prioridad:    "alta",
		Descripcion:  "batería",
		Estado:       "abierta",
	}
	v2.SetIncidencia(inc2)

	v3 := &Vehiculo{Matricula: "CCC333", Marca: "Ford", Modelo: "Focus", FechaEntrada: "01/01/2025"}
	inc3 := &Incidencia{
		IDIncidencia: 3,
		Tipo:         "carrocería",
		Prioridad:    "baja",
		Descripcion:  "golpe lateral",
		Estado:       "abierta",
	}
	v3.SetIncidencia(inc3)

	// Asignar vehículos a clientes
	c1.Vehiculos = []*Vehiculo{v1, v2}
	c2.Vehiculos = []*Vehiculo{v3}

	app.ClientesTaller = []*Cliente{c1, c2}

	// Actualizamos el siguiente ID de incidencia
	nextIncID = 4

	fmt.Println("Datos de prueba cargados: 3 mecánicos, 2 clientes, 3 vehículos con incidencia.")
}
