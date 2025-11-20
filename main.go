package main

import "fmt"

func main() {

	app.MecanicosTaller = []*Mecanico{
		{IDMecanico: 1, Nombre: "Laura", Especialidad: "mecánica", AniosExperiencia: 3, Activo: true},
		{IDMecanico: 2, Nombre: "Pedro", Especialidad: "eléctrica", AniosExperiencia: 5, Activo: true},
	}
	app.ClientesTaller = []*Cliente{}
	app.InicializarPlazas()

	var opcion int

	for {
		fmt.Println("\n===== MENU PRINCIPAL =====")
		fmt.Println("1. Gestionar clientes")
		fmt.Println("2. Gestionar vehículos")
		fmt.Println("3. Gestionar incidencias")
		fmt.Println("4. Gestionar mecánicos")
		fmt.Println("5. Asignar vehículo a plaza")
		fmt.Println("6. Consultar estado del taller")
		fmt.Println("7. Simular atención concurrente")
		fmt.Println("0. Salir")
		fmt.Print("Opción: ")

		fmt.Scanln(&opcion)

		switch opcion {
		case 1:
			menuClientes()
		case 2:
			menuVehiculos()
		case 3:
			menuIncidencias()
		case 4:
			menuMecanicos()
		case 5:
			asignarVehiculoAPlaza()
		case 6:
			consultarEstadoTaller()
		case 7:
			simularAtencionConcurrente()
		case 0:
			return
		default:
			fmt.Println("Opción inválida.")
		}
	}
}
