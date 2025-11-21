# Jorge Díaz Alcojor

# SistemasDistribuidos-p2 – Taller concurrente en Go

Esta práctica amplía la **Práctica 1 del taller** añadiendo **concurrencia** con goroutines y channels, además de una **batería de pruebas y benchmarks** para estudiar el rendimiento del sistema.

El objetivo es simular cómo un taller atiende vehículos con distintas incidencias, repartiendo el trabajo entre varios mecánicos en paralelo.

---

## Estructura del proyecto

El proyecto está dividido en varios ficheros:

### `main.go`

Contiene:

* La función `main()`.
* El **menú principal** de la aplicación:

  * Gestión de clientes
  * Gestión de vehículos
  * Gestión de incidencias
  * Gestión de mecánicos
  * Asignación de vehículos a plazas
  * Consulta del estado del taller
  * (Práctica 1) Simulación manual si quieres lanzarla desde aquí

Es el punto de entrada del programa cuando se ejecuta con `go run .` o el binario.

---

### `modelos.go`

Define todas las **estructuras de datos** del taller:

* `Taller`
* `Plaza`
* `Cliente`
* `Vehiculo`
* `Incidencia`
* `Mecanico`

Incluye también métodos básicos sobre estas estructuras, por ejemplo:

* Getters y setters de incidencias.
* Comprobación de disponibilidad de mecánicos o plazas.
* Inicialización de plazas según número de mecánicos, etc.

Es la “capa de modelo” del sistema.

---

### `crud.go`

Incluye toda la parte de **lógica de negocio y menús**:

* Submenús de:

  * Clientes (crear, listar, modificar, eliminar)
  * Vehículos (crear, listar, modificar, eliminar, asociar incidencia, consultar incidencia)
  * Incidencias (crear, listar, modificar, eliminar, cambiar estado)
  * Mecánicos (crear, listar, modificar, eliminar, alta/baja)
* Gestión de plazas:

  * Asignar vehículo a plaza
  * Liberar plazas cuando se elimina cliente/mecánico
  * Consulta del estado actual del taller (plazas libres/ocupadas y porcentaje de ocupación)

Aquí se usan las estructuras definidas en `modelos.go` y se implementa el funcionamiento “normal” del programa (la parte de consola que el usuario ve).

---

### `simulacion.go`

Fichero dedicado a la **simulación concurrente** del taller usando goroutines y channels.

Incluye:

* Un **worker de mecánico** (goroutine) que atiende vehículos según el tipo de incidencia:

  * mecánica → 5 “segundos”
  * eléctrica → 7 “segundos”
  * carrocería → 11 “segundos”
* Una función de simulación que:

  * Encola los vehículos en canales según el tipo de incidencia.
  * Lanza una goroutine por mecánico activo.
  * Mide tiempos de atención y muestra por pantalla el proceso.

Este fichero sirve de base para entender cómo se reparte la carga entre mecánicos en paralelo.

---

### `taller_test.go`

Contiene los **tests y benchmarks** de la práctica 2:

* Un test principal `TestComparativasTaller` que:

  * Crea diferentes escenarios:

    * N coches vs 2N coches (misma plantilla).
    * Plantilla 1-1-1 vs 2-2-2 (misma carga).
    * Distintas distribuciones de especialidad (3M/1E/1C vs 1M/3E/3C).
  * Lanza la simulación concurrente para cada caso.
  * Muestra por pantalla:

    * Número de coches
    * Nº de mecánicos por especialidad
    * Duración real de la simulación (ms)
    * Tiempo medio simulado por vehículo (s)
* Benchmarks:

  * `BenchmarkDuplicarCochesMecanica`
  * `BenchmarkDuplicarPlantilla`
  * `BenchmarkDistribucionEspecialidades`

Estos benchmarks miden el tiempo que tarda Go en ejecutar cada escenario muchas veces y sirven para comparar el rendimiento de cada configuración.

---

### `go.mod`

Archivo de módulo de Go. Indica:

* El nombre del módulo (por ejemplo `module taller`)
* La versión mínima de Go

Permite compilar y ejecutar el proyecto con `go build` y `go run` desde la carpeta del proyecto sin problemas.

---

## Cómo compilar y ejecutar la práctica

Desde la carpeta raíz del proyecto (donde está `go.mod`):

### 1. Compilar el programa

```bash
go build .
```

Esto genera un ejecutable (por ejemplo `taller` o `taller.exe` según el sistema).

### 2. Ejecutar el programa

```bash
go run .
```

Se abrirá el **menú principal**, desde el que se puede:

* Dar de alta clientes, vehículos, incidencias y mecánicos.
* Asignar vehículos a plazas.
* Consultar el estado actual del taller.
* Lanzar la simulación concurrente de la práctica 2.

---

## Cómo ejecutar los tests de simulación

Para lanzar el test que compara escenarios:

```bash
go test -v
```

El `-v` hace que se vea por pantalla:

* El resumen de cada escenario.
* Las métricas de duración y tiempos medios.
* Comentarios que puedes reutilizar en la memoria.

---

## Cómo ejecutar los benchmarks

Para lanzar los benchmarks de rendimiento:

```bash
go test -bench=. -run=^$
```

* `-bench=.` ejecuta todas las funciones `Benchmark...`.
* `-run=^$` evita correr los tests normales, solo benchmarks.

La salida muestra el tiempo medio (en nanosegundos por ejecución) de cada escenario, lo que permite comparar:

* N vs 2N coches.
* Plantilla 1-1-1 vs 2-2-2.
* Distribuciones 3M/1E/1C vs 1M/3E/3C.
