package main

import (
	"fmt"
	"math"
)

// HELPERS

func findClienteByID(id int) (*Cliente, int) {
	for idx, c := range app.ClientesTaller {
		if c.IDCliente == id {
			return c, idx
		}
	}
	return nil, -1
}

func findMecanicoByID(id int) (*Mecanico, int) {
	for idx, m := range app.MecanicosTaller {
		if m.IDMecanico == id {
			return m, idx
		}
	}
	return nil, -1
}

func liberarPlazasDeCliente(c *Cliente) {
	for _, p := range app.PlazasTaller {
		if p.ocupada && p.cliente == c {
			p.Liberar()
		}
	}
}

func liberarPlazasDeMecanico(m *Mecanico) {
	for _, p := range app.PlazasTaller {
		if p.ocupada && p.mecanico == m {
			p.Liberar()
		}
	}
}

// MENÚS CRUD

// CLIENTES
func menuClientes() {
	var op int
	for {
		fmt.Println("\n===== GESTIÓN DE CLIENTES =====")
		fmt.Println("1. Crear cliente")
		fmt.Println("2. Visualizar clientes")
		fmt.Println("3. Modificar cliente")
		fmt.Println("4. Eliminar cliente")
		fmt.Println("0. Volver")
		fmt.Print("Opción: ")
		fmt.Scanln(&op)

		switch op {
		case 1:
			crearCliente()
		case 2:
			listarClientes()
		case 3:
			modificarCliente()
		case 4:
			eliminarCliente()
		case 0:
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

// VEHÍCULOS
func menuVehiculos() {
	var op int
	for {
		fmt.Println("\n===== GESTIÓN DE VEHÍCULOS =====")
		fmt.Println("1. Crear vehículo")
		fmt.Println("2. Visualizar vehículos")
		fmt.Println("3. Modificar vehículo")
		fmt.Println("4. Eliminar vehículo")
		fmt.Println("5. Registrar incidencia a un vehículo")
		fmt.Println("6. Consultar incidencia de un vehículo")
		fmt.Println("0. Volver")
		fmt.Print("Opción: ")
		fmt.Scanln(&op)

		switch op {
		case 1:
			crearVehiculo()
		case 2:
			listarVehiculos()
		case 3:
			modificarVehiculo()
		case 4:
			eliminarVehiculo()
		case 5:
			registrarIncidenciaVehiculo()
		case 6:
			consultarIncidenciaVehiculo()
		case 0:
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

// INCIDENCIAS
func menuIncidencias() {
	var op int
	for {
		fmt.Println("\n===== GESTIÓN DE INCIDENCIAS =====")
		fmt.Println("1. Crear incidencia (vehículo)")
		fmt.Println("2. Visualizar incidencias")
		fmt.Println("3. Modificar incidencia")
		fmt.Println("4. Eliminar incidencia")
		fmt.Println("5. Cambiar estado de incidencia")
		fmt.Println("0. Volver")
		fmt.Print("Opción: ")
		fmt.Scanln(&op)

		switch op {
		case 1:
			registrarIncidenciaVehiculo()
		case 2:
			listarIncidencias()
		case 3:
			modificarIncidencia()
		case 4:
			eliminarIncidencia()
		case 5:
			cambiarEstadoIncidencia()
		case 0:
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

// MECANICOS
func menuMecanicos() {
	var op int
	for {
		fmt.Println("\n===== GESTIÓN DE MECÁNICOS =====")
		fmt.Println("1. Crear mecánico")
		fmt.Println("2. Visualizar mecánicos")
		fmt.Println("3. Modificar mecánico")
		fmt.Println("4. Eliminar mecánico")
		fmt.Println("5. Dar de alta/baja a un mecánico")
		fmt.Println("0. Volver")
		fmt.Print("Opción: ")
		fmt.Scanln(&op)

		switch op {
		case 1:
			crearMecanico()
		case 2:
			listarMecanicos()
		case 3:
			modificarMecanico()
		case 4:
			eliminarMecanico()
		case 5:
			cambiarEstadoMecanico()
		case 0:
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

// CRUD CLIENTES

func crearCliente() {
	var id int
	var nombre, telefono, email string
	fmt.Print("ID cliente: ")
	fmt.Scanln(&id)
	if _, idx := findClienteByID(id); idx != -1 {
		fmt.Println("Ya existe un cliente con ese ID.")
		return
	}
	fmt.Print("Nombre: ")
	fmt.Scanln(&nombre)
	fmt.Print("Teléfono: ")
	fmt.Scanln(&telefono)
	fmt.Print("Email: ")
	fmt.Scanln(&email)

	c := &Cliente{IDCliente: id, Nombre: nombre, Telefono: telefono, Email: email}
	app.ClientesTaller = append(app.ClientesTaller, c)
	fmt.Println("Cliente creado.")
}

func listarClientes() {
	if len(app.ClientesTaller) == 0 {
		fmt.Println("No hay clientes.")
		return
	}
	fmt.Println("Listado de clientes:")
	for _, c := range app.ClientesTaller {
		fmt.Printf("- ID:%d | %s | Tel:%s | Email:%s | Vehículos:%d\n",
			c.IDCliente, c.Nombre, c.Telefono, c.Email, len(c.Vehiculos))
	}
}

func modificarCliente() {
	var id int
	fmt.Print("ID cliente a modificar: ")
	fmt.Scanln(&id)
	c, _ := findClienteByID(id)
	if c == nil {
		fmt.Println("Cliente no encontrado.")
		return
	}
	var nombre, tel, email string
	fmt.Print("Nuevo nombre: ")
	fmt.Scanln(&nombre)
	fmt.Print("Nuevo teléfono: ")
	fmt.Scanln(&tel)
	fmt.Print("Nuevo email: ")
	fmt.Scanln(&email)
	c.Nombre, c.Telefono, c.Email = nombre, tel, email
	fmt.Println("Cliente modificado.")
}

func eliminarCliente() {
	var id int
	fmt.Print("ID cliente a eliminar: ")
	fmt.Scanln(&id)
	c, idx := findClienteByID(id)
	if c == nil {
		fmt.Println("Cliente no encontrado.")
		return
	}
	liberarPlazasDeCliente(c)
	app.ClientesTaller = append(app.ClientesTaller[:idx], app.ClientesTaller[idx+1:]...)
	fmt.Println("Cliente eliminado y plazas liberadas si correspondía.")
}

// CRUD VEHÍCULOS

func crearVehiculo() {
	var idCliente int
	fmt.Print("ID del cliente propietario: ")
	fmt.Scanln(&idCliente)
	c, _ := findClienteByID(idCliente)
	if c == nil {
		fmt.Println("Cliente no encontrado.")
		return
	}
	var mat, marca, modelo, fIn, fOut string
	fmt.Print("Matrícula: ")
	fmt.Scanln(&mat)
	if _, v := app.BuscarVehiculo(mat); v != nil {
		fmt.Println("Ya existe un vehículo con esa matrícula.")
		return
	}
	fmt.Print("Marca: ")
	fmt.Scanln(&marca)
	fmt.Print("Modelo: ")
	fmt.Scanln(&modelo)
	fmt.Print("Fecha de entrada: ")
	fmt.Scanln(&fIn)
	fmt.Print("Fecha de salida: ")
	fmt.Scanln(&fOut)

	v := &Vehiculo{Matricula: mat, Marca: marca, Modelo: modelo, FechaEntrada: fIn, FechaSalida: fOut}
	c.Vehiculos = append(c.Vehiculos, v)
	fmt.Println("Vehículo creado y asignado al cliente.")
}

func listarVehiculos() {
	encontrados := 0
	for _, c := range app.ClientesTaller {
		for _, v := range c.Vehiculos {
			encontrados++
			estadoInc := "sin incidencia"
			if v.GetIncidencia() != nil {
				estadoInc = "incidencia " + v.GetIncidencia().Estado
			}
			fmt.Printf("- [%s] %s %s | Cliente:%s | %s\n",
				v.Matricula, v.Marca, v.Modelo, c.Nombre, estadoInc)
		}
	}
	if encontrados == 0 {
		fmt.Println("No hay vehículos registrados.")
	}
}

func modificarVehiculo() {
	var mat string
	fmt.Print("Matrícula del vehículo a modificar: ")
	fmt.Scanln(&mat)
	c, v := app.BuscarVehiculo(mat)
	if v == nil {
		fmt.Println("Vehículo no encontrado.")
		return
	}
	var marca, modelo, fIn, fOut string
	fmt.Print("Nueva marca: ")
	fmt.Scanln(&marca)
	fmt.Print("Nuevo modelo: ")
	fmt.Scanln(&modelo)
	fmt.Print("Nueva fecha de entrada: ")
	fmt.Scanln(&fIn)
	fmt.Print("Nueva fecha de salida: ")
	fmt.Scanln(&fOut)
	v.Marca, v.Modelo, v.FechaEntrada, v.FechaSalida = marca, modelo, fIn, fOut
	fmt.Printf("Vehículo %s del cliente %s modificado.\n", v.Matricula, c.Nombre)
}

func eliminarVehiculo() {
	var mat string
	fmt.Print("Matrícula del vehículo a eliminar: ")
	fmt.Scanln(&mat)
	c, v := app.BuscarVehiculo(mat)
	if v == nil {
		fmt.Println("Vehículo no encontrado.")
		return
	}
	v.SetIncidencia(nil)

	pos := -1
	for i, vv := range c.Vehiculos {
		if vv == v {
			pos = i
			break
		}
	}
	if pos != -1 {
		c.Vehiculos = append(c.Vehiculos[:pos], c.Vehiculos[pos+1:]...)
	}
	fmt.Println("Vehículo eliminado.")
}

// CRUD INCIDENCIAS

func registrarIncidenciaVehiculo() {
	var mat string
	fmt.Print("Matrícula del vehículo: ")
	fmt.Scanln(&mat)
	c, v := app.BuscarVehiculo(mat)
	if v == nil {
		fmt.Println("Vehículo no encontrado.")
		return
	}
	if v.GetIncidencia() != nil {
		fmt.Println("Este vehículo ya tiene una incidencia.")
		return
	}

	var tipo, prio, desc string
	fmt.Print("Tipo (mecánica/eléctrica/carrocería): ")
	fmt.Scanln(&tipo)
	fmt.Print("Prioridad (baja/media/alta): ")
	fmt.Scanln(&prio)
	fmt.Print("Descripción: ")
	fmt.Scanln(&desc)

	inc := &Incidencia{
		IDIncidencia: nextIncID,
		Tipo:         tipo,
		Prioridad:    prio,
		Descripcion:  desc,
		Estado:       "abierta",
	}
	nextIncID++
	v.SetIncidencia(inc)

	fmt.Printf("Incidencia registrada al vehículo %s del cliente %s (ID=%d).\n",
		v.Matricula, c.Nombre, inc.IDIncidencia)
}

func consultarIncidenciaVehiculo() {
	var mat string
	fmt.Print("Matrícula del vehículo: ")
	fmt.Scanln(&mat)
	_, v := app.BuscarVehiculo(mat)
	if v == nil {
		fmt.Println("Vehículo no encontrado.")
		return
	}
	inc := v.GetIncidencia()
	if inc == nil {
		fmt.Println("El vehículo no tiene incidencia.")
		return
	}
	fmt.Printf("Incidencia ID:%d | Tipo:%s | Prio:%s | Estado:%s | Desc:%s | Mecánicos:%d\n",
		inc.IDIncidencia, inc.Tipo, inc.Prioridad, inc.Estado, inc.Descripcion, len(inc.GetMecanicos()))
}

func listarIncidencias() {
	total := 0
	for _, c := range app.ClientesTaller {
		for _, v := range c.Vehiculos {
			if inc := v.GetIncidencia(); inc != nil {
				total++
				fmt.Printf("- Vehículo [%s] de %s | IncID:%d | Tipo:%s | Prio:%s | Estado:%s\n",
					v.Matricula, c.Nombre, inc.IDIncidencia, inc.Tipo, inc.Prioridad, inc.Estado)
			}
		}
	}
	if total == 0 {
		fmt.Println("No hay incidencias registradas.")
	}
}

func modificarIncidencia() {
	var mat string
	fmt.Print("Matrícula del vehículo con incidencia: ")
	fmt.Scanln(&mat)
	_, v := app.BuscarVehiculo(mat)
	if v == nil || v.GetIncidencia() == nil {
		fmt.Println("Vehículo no encontrado o sin incidencia.")
		return
	}
	inc := v.GetIncidencia()
	var tipo, prio, desc string
	fmt.Print("Nuevo tipo: ")
	fmt.Scanln(&tipo)
	fmt.Print("Nueva prioridad: ")
	fmt.Scanln(&prio)
	fmt.Print("Nueva descripción: ")
	fmt.Scanln(&desc)
	inc.Tipo, inc.Prioridad, inc.Descripcion = tipo, prio, desc
	fmt.Println("Incidencia modificada.")
}

func eliminarIncidencia() {
	var mat string
	fmt.Print("Matrícula del vehículo con incidencia a eliminar: ")
	fmt.Scanln(&mat)
	_, v := app.BuscarVehiculo(mat)
	if v == nil || v.GetIncidencia() == nil {
		fmt.Println("Vehículo no encontrado o sin incidencia.")
		return
	}
	v.SetIncidencia(nil)
	fmt.Println("Incidencia eliminada.")
}

func cambiarEstadoIncidencia() {
	var mat, nuevo string
	fmt.Print("Matrícula del vehículo: ")
	fmt.Scanln(&mat)
	_, v := app.BuscarVehiculo(mat)
	if v == nil || v.GetIncidencia() == nil {
		fmt.Println("Vehículo no encontrado o sin incidencia.")
		return
	}
	fmt.Print("Nuevo estado: ")
	fmt.Scanln(&nuevo)
	v.GetIncidencia().SetEstado(nuevo)
	fmt.Println("Estado actualizado.")
}

// CRUD MECÁNICOS

func crearMecanico() {
	var id int
	var nombre, esp string
	var anios int
	fmt.Print("ID mecánico: ")
	fmt.Scanln(&id)
	if _, idx := findMecanicoByID(id); idx != -1 {
		fmt.Println("Ya existe un mecánico con ese ID.")
		return
	}
	fmt.Print("Nombre: ")
	fmt.Scanln(&nombre)
	fmt.Print("Especialidad: ")
	fmt.Scanln(&esp)
	fmt.Print("Años de experiencia: ")
	fmt.Scanln(&anios)

	m := &Mecanico{IDMecanico: id, Nombre: nombre, Especialidad: esp, AniosExperiencia: anios, Activo: true}
	app.MecanicosTaller = append(app.MecanicosTaller, m)
	app.InicializarPlazas()
	fmt.Println("Mecánico creado.")
}

func listarMecanicos() {
	if len(app.MecanicosTaller) == 0 {
		fmt.Println("No hay mecánicos.")
		return
	}
	for _, m := range app.MecanicosTaller {
		estado := "baja"
		if m.Activo {
			estado = "activo"
		}
		fmt.Printf("- ID:%d | %s | %s | %d años | %s\n",
			m.IDMecanico, m.Nombre, m.Especialidad, m.AniosExperiencia, estado)
	}
}

func modificarMecanico() {
	var id int
	fmt.Print("ID del mecánico a modificar: ")
	fmt.Scanln(&id)
	m, _ := findMecanicoByID(id)
	if m == nil {
		fmt.Println("No existe ese mecánico.")
		return
	}
	var nombre, esp string
	var anios int
	fmt.Print("Nuevo nombre: ")
	fmt.Scanln(&nombre)
	fmt.Print("Nueva especialidad: ")
	fmt.Scanln(&esp)
	fmt.Print("Años de experiencia: ")
	fmt.Scanln(&anios)
	m.Nombre, m.Especialidad, m.AniosExperiencia = nombre, esp, anios
	fmt.Println("Mecánico modificado.")
}

func eliminarMecanico() {
	var id int
	fmt.Print("ID del mecánico a eliminar: ")
	fmt.Scanln(&id)
	m, idx := findMecanicoByID(id)
	if m == nil {
		fmt.Println("No existe ese mecánico.")
		return
	}
	liberarPlazasDeMecanico(m)
	app.MecanicosTaller = append(app.MecanicosTaller[:idx], app.MecanicosTaller[idx+1:]...)
	app.InicializarPlazas()
	fmt.Println("Mecánico eliminado y plazas recalculadas.")
}

func cambiarEstadoMecanico() {
	var id int
	var op int
	fmt.Print("ID del mecánico: ")
	fmt.Scanln(&id)
	m, _ := findMecanicoByID(id)
	if m == nil {
		fmt.Println("No existe ese mecánico.")
		return
	}
	fmt.Print("1=Activar, 2=Baja: ")
	fmt.Scanln(&op)
	if op == 1 {
		m.CambiarEstado(true)
	} else if op == 2 {
		m.CambiarEstado(false)
	}
	app.InicializarPlazas()
	fmt.Println("Estado actualizado.")
}

// PLAZAS

func asignarVehiculoAPlaza() {
	ocupadas, libres := app.EstadoTaller()
	if libres == 0 {
		fmt.Println("No hay plazas libres: taller lleno.")
		return
	}
	var mat string
	fmt.Print("Matrícula del vehículo a asignar: ")
	fmt.Scanln(&mat)
	cli, veh := app.BuscarVehiculo(mat)
	if veh == nil {
		fmt.Println("Vehículo no encontrado.")
		return
	}
	var idm int
	fmt.Print("ID del mecánico para asignar: ")
	fmt.Scanln(&idm)
	mec, _ := findMecanicoByID(idm)
	if mec == nil || !mec.Activo {
		fmt.Println("Mecánico inexistente o no activo.")
		return
	}
	for _, p := range app.PlazasTaller {
		if p.EstaLibre() {
			p.Ocupar(cli, mec)
			fmt.Printf("Vehículo %s asignado a plaza #%d con mecánico %s. (Ocupadas:%d→%d)\n",
				veh.Matricula, p.IDPlaza, mec.Nombre, ocupadas, ocupadas+1)
			return
		}
	}
	fmt.Println("No se encontró plaza libre (estado desactualizado).")
}

func consultarEstadoTaller() {
	ocupadas, libres := app.EstadoTaller()
	total := len(app.PlazasTaller)
	pct := 0.0
	if total > 0 {
		pct = math.Round((float64(ocupadas)/float64(total))*100.0 + 0.00001)
	}

	fmt.Printf("Ocupadas: %d | Libres: %d | Total: %d | %.0f%% ocupado.\n",
		ocupadas, libres, total, pct)

	for _, p := range app.PlazasTaller {
		if p.ocupada {
			fmt.Printf(" - Plaza %d: OCUPADA (Cliente:%s, Mecánico:%s)\n",
				p.IDPlaza, p.GetCliente().Nombre, p.GetMecanico().Nombre)
		} else {
			fmt.Printf(" - Plaza %d: libre\n", p.IDPlaza)
		}
	}
}
