package main

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/geometry"
)

func NewSwathe(points [][]int) *geometry.Geometry{
	s := geometry.NewGeometry()
	n := len(points)
	//m := len(points[0])
	normap := make(map[math32.Vector3]math32.Vector3)

	positions := math32.NewArrayF32(0, 0)
	normals := math32.NewArrayF32(0, 0)
	uvs := math32.NewArrayF32(0, 0)
	indices := math32.NewArrayU32(0, 0)

	// go through the map, make a vertex for each point
	// VERTICES
	for i:=0; i<len(points); i++ {
		for j:=0; j<len(points[i]); j++ {
			var vertex math32.Vector3
			vertex.X = float32(i)
			vertex.Y = float32(points[i][j])
			vertex.Z = float32(j)
			positions.AppendVector3(&vertex)
			uvs.Append(float32(i), float32(j))
		}
	}

	// INDICES
	for i:=0; i<(len(positions)/3)-n-1; i+=n {
		for j:=i; j<i+n-1; j++ {
			a := j
			b := j+1
			c := j+n
			d := j+n+1
			indices.Append(uint32(a), uint32(b), uint32(c))
			indices.Append(uint32(c), uint32(b), uint32(d))

			/*
			f1 := 3*a
			f2 := 3*b
			f3 := 3*c
			p1 := math32.NewVector3(positions[f1], positions[f1+1],
				positions[f1+2])
			p2 := math32.NewVector3(positions[f2], positions[f2+1],
				positions[f2+2])
			p3 := math32.NewVector3(positions[f3], positions[f3+1],
				positions[f3+2])
			A := p2.Sub(p1)
			B := p3.Sub(p1)
			N := A.Cross(B).Normalize()
			normals.AppendVector3(N)
			if j == (i+n-2) {
				normals.AppendVector3(N)
			}
*/
		}
	}

	// NORMALS
	for i:=0; i<len(positions)-(3*(n+2)); i+=3*(n+1) {
		for j:=i; j<i+(3*n); j+=3 {
			a := math32.NewVector3(positions[j], positions[j+1],
				positions[j+2])
			b := math32.NewVector3(positions[j+3], positions[j+4],
				positions[j+5])
			c := math32.NewVector3(positions[j+3*(n+1)],
				positions[j+(3*(n+1))+1],
				positions[j+(3*(n+1))+2])
			d := math32.NewVector3(positions[j+3*(n+2)],
				positions[j+(3*(n+2))+1],
				positions[j+(3*(n+2))+2])

			A := b.Sub(a)
			B := c.Sub(a)
			N := A.Cross(B).Normalize()

			if val, ok := normap[*a]; ok {
				normap[*a] = *(AvgVec3(&val, N).Normalize())
			} else {
				normap[*a] = *N
			}
			if val, ok := normap[*b]; ok {
				normap[*b] = *(AvgVec3(&val, N).Normalize())
			} else {
				normap[*b] = *N
			}
			if val, ok := normap[*c]; ok {
				normap[*c] = *(AvgVec3(&val, N).Normalize())
			} else {
				normap[*c] = *N
			}

			A = c.Sub(d)
			B = b.Sub(d)
			N = A.Cross(B).Normalize()

			if val, ok := normap[*d]; ok {
				normap[*d] = *(AvgVec3(&val, N).Normalize())
			} else {
				normap[*d] = *N
			}
			if val, ok := normap[*c]; ok {
				normap[*c] = *(AvgVec3(&val, N).Normalize())
			} else {
				normap[*c] = *N
			}
			if val, ok := normap[*b]; ok {
				normap[*b] = *(AvgVec3(&val, N).Normalize())
			} else {
				normap[*b] = *N
			}
		}
	}

	for i:=0; i<len(positions); i+=3 {
		vertex := math32.NewVector3(positions[i], positions[i+1],
			positions[i+2])
		norm := normap[*vertex]
		normals.AppendVector3(&norm)
	}
	
	/*
	// DO NORMALS for last row, quick n dirty
	for i := (len(positions)/3)-(2*n)+1; i < (len(positions)/3)-n+1; i++ {
		normals.Append(normals[i])
	}*/

	s.SetIndices(indices)
	s.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	s.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	s.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))
	
	return s
}

func AvgVec3(a, b *math32.Vector3) *math32.Vector3{
	return a.Add(b).DivideScalar(2)
}
