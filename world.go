package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lferth93/util/hashset"
)

//constantes para las diferentes tipos de celdas y colores
//los colores primarios son rojo, verde y azul
const (
	BLACK  = iota //0
	RED           //1
	GREEN         //2
	YELLOW        //...
	BLUE
	MAGENTA
	CYAN
	WHITE
	BLOCK
)

//funcion que determina la interaccion entre dos colores
func addColor(c1, c2 int8) int8 {
	//alguno de los dos colores es un bloque
	if c1 == BLOCK || c2 == BLOCK {
		return BLOCK
	}
	//comparten algun color primario y el resultado es el color dominante en la mezcla
	if c1&c2 != 0 {
		return c1 & c2
	}
	//no comparten ningun color primario
	return c1 | c2
}

//World parecido a una clase de java
type World struct {
	s       [][]int8         //matriz para guardar el estado del mundo
	a       [][]int8         //automata para saber como pintar usa Wireworld con vecindad de Von Neumann
	life    int              //generacioes que tiene de vida el mundo
	w, h    int              //dimensiones del mundo
	final   bool             //para saber si el mundo ya llego a un estado en el que no va a cambiar
	activos *hashset.Hashset //set que guarda las celdas que pueden pintar celdas vecinas que seran analizadas cuando se genere el siguiente estado
}

//NewWorld equivalente a un constructor de java
func NewWorld(w, h int) *World {
	var nw World
	nw.w, nw.h = w, h
	nw.final = false
	nw.s = make([][]int8, h)
	for i := range nw.s {
		nw.s[i] = make([]int8, w)
	}

	//rodea el mundo en blanco con bloques para delimitarlo
	for i := 0; i < w; i++ {
		nw.s[0][i] = BLOCK
	}
	for i := 0; i < h; i++ {
		nw.s[i][w-1] = BLOCK
	}
	for i := w - 1; i > 0; i-- {
		nw.s[h-1][i] = BLOCK
	}
	for i := h - 1; i > 0; i-- {
		nw.s[i][0] = BLOCK
	}

	return &nw
}

//Init metodo para inicializar alguna celda dentro del limite del mundo
//las celdas usables se indexan de (1,1) a (w,h)
func (w *World) Init(x, y int, c int8) {
	if 0 < x && x < w.w-1 && 0 < y && y < w.h-1 {
		w.s[y][x] = c
		if c != BLOCK && c != BLACK {
			if w.activos == nil {
				w.activos = hashset.New()
			}
			w.activos.Insert([2]int{x, y})
		}
	}
}

func (w *World) initAutomata() {
	w.a = make([][]int8, w.h)
	for i := range w.a {
		w.a[i] = make([]int8, w.w)
	}
	for i := range w.s {
		for j := range w.s[i] {
			switch {
			case w.s[i][j] == 8:
				w.a[i][j] = -1
			case w.s[i][j] > 0:
				w.a[i][j] = 1
			}
		}
	}
}

func (w *World) nextAutomata() {
	na := make([][]int8, w.h)
	for i := range w.a {
		na[i] = make([]int8, w.w)
		for j := range w.a[i] {
			switch w.a[i][j] {
			case -1:
				na[i][j] = -1
			case 0:
				c := 0
				if w.a[i-1][j] == 1 {
					c++
				}
				if w.a[i+1][j] == 1 {
					c++
				}
				if w.a[i][j-1] == 1 {
					c++
				}
				if w.a[i][j+1] == 1 {
					c++
				}
				if 0 < c && c < 4 {
					na[i][j] = 1
				}
			case 1:
				na[i][j] = 2
			case 2:
				na[i][j] = 0
			}
		}
	}
	w.a = na
}

//equivalente al metodo toString de java
func (w *World) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Life:%d\n", w.life))
	sb.WriteString(fmt.Sprintf("Actives:%d\n", w.activos.Size()))
	sb.WriteString(w.hashString())
	sb.WriteRune('\n')
	for i := range w.s {
		sb.WriteString(fmt.Sprintln(w.s[i]))
	}
	return sb.String()
}

//Next intenta avanzar el mundo al siguiente estado y devuelve true si es posible avanzar y false si el mundo llego a un estado final
func (w *World) Next() bool {
	if !w.final {
		nw := make([][]int8, w.h)
		for i := range nw {
			nw[i] = make([]int8, w.w)
		}
		for i := range w.a {
			for j, a := range w.a[i] {
				nw[i][j] = addColor(w.s[i][j], nw[i][j])
				if a == 1 {
					if w.s[i-1][j] != BLOCK && w.a[i-1][j] == 0 {
						nw[i-1][j] = addColor(nw[i-1][j], w.s[i][j])
					}
					if w.s[i+1][j] != BLOCK && w.a[i+1][j] == 0 {
						nw[i+1][j] = addColor(nw[i+1][j], w.s[i][j])
					}
					if w.s[i][j-1] != BLOCK && w.a[i][j-1] == 0 {
						nw[i][j-1] = addColor(nw[i][j-1], w.s[i][j])
					}
					if w.s[i][j+1] != BLOCK && w.a[i][j+1] == 0 {
						nw[i][j+1] = addColor(nw[i][j+1], w.s[i][j])
					}
				}
			}
		}
		w.s = nw
		w.nextAutomata()
	}
	return false
}

//revisa si la celda de cordenadas x,y puede pintar alguna de sus celdas vecinas
func (w *World) evolve(x, y int) bool {
	c := w.s[y][x]
	if c == BLOCK || c == BLACK {
		return false
	}
	return (w.s[y-1][x] != c && w.s[y-1][x] != BLOCK) ||
		(w.s[y][x-1] != c && w.s[y][x-1] != BLOCK) ||
		(w.s[y+1][x] != c && w.s[y+1][x] != BLOCK) ||
		(w.s[y][x+1] != c && w.s[y][x+1] != BLOCK)
}

func (w *World) hashString() string {
	var sb strings.Builder
	for _, r := range w.s {
		for _, b := range r {
			sb.WriteRune('a' + rune(b))
		}
	}
	return sb.String()
}

//LoadFile reads a world from a json file
func LoadFile(path string) *World {
	var w World
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var wd [][]int8

	err = json.Unmarshal(buf, &wd)
	w.s = wd
	w.h = len(wd)
	w.w = len(wd[0])
	w.activos = hashset.New()
	for y, r := range w.s {
		for x, c := range r {
			if c != BLOCK && c != BLACK {
				w.activos.Insert([2]int{x, y})
			}
		}
	}
	w.initAutomata()
	return &w
}

//Save guarda el estado de w en un archivo con formato json
func (w *World) Save(path string) error {
	buf, err := json.Marshal(w.s)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, buf, os.ModeDir)
	return err
}
