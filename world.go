package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

//Color type uint8
type Color uint8

//Determina el color resultante de la mezcla
func addColor(cs [5]Color) (nc Color) {
	cnt := []int{0, 0, 0}
	cl := []Color{RED, GREEN, BLUE}
	max := 0
	for _, c := range cs {
		for i := range cl {
			if c&cl[i] > 0 {
				cnt[i]++
			}
		}
	}
	for i := range cnt {
		if cnt[i] > max {
			max = cnt[i]
		}
	}
	for i := range cnt {
		if cnt[i] == max {
			nc |= cl[i]
		}
	}
	return
}

//World parecido a una clase de java
type World struct {
	s     [][]Color //matriz para guardar el estado del mundo
	a     [][]int8  //automata para saber como pintar usa Wireworld con vecindad de Von Neumann
	life  int       //generacioes que tiene de vida el mundo
	w, h  int       //dimensiones del mundo
	final bool      //para saber si el mundo ya llego a un estado en el que no va a cambiar
}

//NewWorld equivalente a un constructor de java
func NewWorld(w, h int) *World {
	var nw World
	nw.w, nw.h = w, h
	nw.final = false
	nw.s = make([][]Color, h)
	for i := range nw.s {
		nw.s[i] = make([]Color, w)
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
func (w *World) Init(x, y int, c Color) {
	if 0 < x && x < w.w-1 && 0 < y && y < w.h-1 {
		w.s[y][x] = c
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
		nw := make([][]Color, w.h)
		for i := range nw {
			nw[i] = make([]Color, w.w)
		}
		for i := range w.a {
			for j, s := range w.a[i] {
				if s == 0 {
					count := 0
					var cs [5]Color
					cs[0] = w.s[i][j]
					if w.a[i][j-1] == 1 {
						cs[1] = w.s[i][j-1]
						count++
					}
					if w.a[i][j+1] == 1 {
						cs[2] = w.s[i][j+1]
						count++
					}
					if w.a[i-1][j] == 1 {
						cs[3] = w.s[i-1][j]
						count++
					}
					if w.a[i+1][j] == 1 {
						cs[4] = w.s[i+1][j]
						count++
					}
					if count > 0 {
						nw[i][j] = addColor(cs)
					} else {
						nw[i][j] = w.s[i][j]
					}
				} else {
					nw[i][j] = w.s[i][j]
				}
			}
		}
		w.s = nw
		w.nextAutomata()
	}
	return false
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
	var wd [][]Color

	err = json.Unmarshal(buf, &wd)
	w.s = wd
	w.h = len(wd)
	w.w = len(wd[0])
	w.initAutomata()
	return &w
}

//Save guarda el estado de w en un archivo con formato json
func (w *World) Save(path string) error {
	buf, err := json.Marshal(w.s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, buf, os.ModeDir)
}
