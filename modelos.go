package main

// ESTRUCTURAS DE DATOS

type Taller struct {
	MaxPlazas       int
	ClientesTaller  []*Cliente
	MecanicosTaller []*Mecanico
	PlazasTaller    []*Plaza
}

type Plaza struct {
	IDPlaza  int
	ocupada  bool
	cliente  *Cliente
	mecanico *Mecanico
}

type Cliente struct {
	IDCliente int
	Nombre    string
	Telefono  string
	Email     string
	Vehiculos []*Vehiculo
}

type Vehiculo struct {
	Matricula      string
	Marca          string
	Modelo         string
	FechaEntrada   string
	FechaSalida    string
	incidencia     *Incidencia
	TiempoAtencion int
}

type Incidencia struct {
	IDIncidencia int
	mecanicos    []*Mecanico
	Tipo         string
	Prioridad    string
	Descripcion  string
	Estado       string
}

type Mecanico struct {
	IDMecanico       int
	Nombre           string
	Especialidad     string
	AniosExperiencia int
	Activo           bool
}

// MÉTODOS

func (t *Taller) InicializarPlazas() {
	t.MaxPlazas = 2 * len(t.MecanicosTaller)
	t.PlazasTaller = make([]*Plaza, t.MaxPlazas)
	for i := 0; i < t.MaxPlazas; i++ {
		t.PlazasTaller[i] = &Plaza{IDPlaza: i + 1}
	}
}

func (t *Taller) EstadoTaller() (ocupadas, libres int) {
	for _, p := range t.PlazasTaller {
		if p.ocupada {
			ocupadas++
		}
	}
	libres = len(t.PlazasTaller) - ocupadas
	return
}

func (t *Taller) BuscarVehiculo(matricula string) (*Cliente, *Vehiculo) {
	for _, c := range t.ClientesTaller {
		for _, v := range c.Vehiculos {
			if v.Matricula == matricula {
				return c, v
			}
		}
	}
	return nil, nil
}

func (t *Taller) ListarMecanicosDisponibles() []*Mecanico {
	var out []*Mecanico
	for _, m := range t.MecanicosTaller {
		if m.Activo {
			out = append(out, m)
		}
	}
	return out
}

// Plaza
func (p *Plaza) Ocupar(c *Cliente, m *Mecanico) {
	p.ocupada = true
	p.cliente = c
	p.mecanico = m
}
func (p *Plaza) Liberar() {
	p.ocupada = false
	p.cliente = nil
	p.mecanico = nil
}
func (p *Plaza) EstaLibre() bool        { return !p.ocupada }
func (p *Plaza) GetCliente() *Cliente   { return p.cliente }
func (p *Plaza) GetMecanico() *Mecanico { return p.mecanico }

// Vehiculo
func (v *Vehiculo) SetIncidencia(i *Incidencia) { v.incidencia = i }
func (v *Vehiculo) GetIncidencia() *Incidencia  { return v.incidencia }

// Incidencia
func (i *Incidencia) AsignarMecanico(m *Mecanico) {
	i.mecanicos = append(i.mecanicos, m)
}
func (i *Incidencia) GetMecanicos() []*Mecanico { return i.mecanicos }
func (i *Incidencia) SetEstado(estado string)   { i.Estado = estado }
func (i *Incidencia) GetEstado() string         { return i.Estado }
func (i *Incidencia) EsAltaPrioridad() bool     { return i.Prioridad == "alta" }

// Mecanico
func (m *Mecanico) CambiarEstado(activo bool) { m.Activo = activo }
func (m *Mecanico) Disponible() bool          { return m.Activo }

// VARIABLES GLOBALES

var app Taller
var nextIncID int = 1
