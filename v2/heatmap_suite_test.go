package heatmap

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestV2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HeatMap Suite")
}

var _ = Describe("HeatMap tests", func() {

	Describe("validate", func() {
		var x, y int
		var data [][]float32
		var err error

		BeforeEach(func() {
			x = 3
			y = 2
			data = [][]float32{{0.1, 0.3}, {0, 0.2}, {0.01, 0.5}}
		})

		JustBeforeEach(func() {
			err = validate(x, y, data)
		})

		Context("when called with valid data", func() {
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when the x dimension is < 1", func() {
			BeforeEach(func() {
				x = 0
			})

			It("returns an error", func() {
				Expect(err).To(Not(BeNil()))
			})
		})

		Context("when the y dimension is < 1", func() {
			BeforeEach(func() {
				y = 0
			})

			It("returns an error", func() {
				Expect(err).To(Not(BeNil()))
			})
		})

		Context("when the data contains a value that is too high", func() {
			BeforeEach(func() {
				data = [][]float32{{1, 0}, {0, 1.01}, {0, 1}}
			})

			It("returns an error", func() {
				Expect(err).To(Not(BeNil()))
			})
		})

		Context("when the data contains a value that is too low", func() {
			BeforeEach(func() {
				data = [][]float32{{1, 0}, {0, -0.01}, {0, 1}}
			})

			It("returns an error", func() {
				Expect(err).To(Not(BeNil()))
			})
		})
		Context("when the x dimension is violated", func() {
			BeforeEach(func() {
				data = [][]float32{{0.1, -0.3}, {0, 0.2}, {0.01, 0.5}, {0.3, 0.4}}
			})

			It("returns an error", func() {
				Expect(err).To(Not(BeNil()))
			})
		})

		Context("when the y dimension is violated", func() {
			BeforeEach(func() {
				data = [][]float32{{0.1, -0.3}, {0, 0.2, 0.4}, {0.01, 0.5}}
			})

			It("returns an error", func() {
				Expect(err).To(Not(BeNil()))
			})
		})
	})

	Describe("Average", func() {
		var a, b *HeatMap
		var r *HeatMap
		var err error

		BeforeEach(func() {
			a, _ = New(3, 2, [][]float32{{0.1, 0.3}, {0, 0.2}, {1, 0.5}})
			b, _ = New(3, 2, [][]float32{{0.3, 0.3}, {0, 0}, {0, 0.1}})
		})

		Context("when called with no arguments", func() {
			JustBeforeEach(func() {
				r, err = Average()
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("when called with one argument", func() {
			JustBeforeEach(func() {
				r, err = Average(a)
			})

			It("returns the argument back", func() {
				Expect(r).To(Equal(a))
			})
		})

		Context("when called with multiple arguments", func() {
			var expected *HeatMap

			BeforeEach(func() {
				expected, _ = New(3, 2, [][]float32{{0.2, 0.3}, {0, 0.1}, {0.5, 0.3}})
			})

			JustBeforeEach(func() {
				r, err = Average(a, b)
			})

			It("returns an averaged HeatMap", func() {
				Expect(r).To(Equal(expected))
			})

			Context("whose dimensions do not match", func() {
				BeforeEach(func() {
					a.width = 5
				})

				It("returns an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})
		})
	})

	Describe("MaxY", func() {
		var a *HeatMap
		var r *HeatMap

		BeforeEach(func() {
			a, _ = New(3, 2, [][]float32{{0.1, 0.3}, {0, 0.2}, {0.01, 0.5}})
		})

		JustBeforeEach(func() {
			r = a.MaxY()
		})

		Context("when called", func() {
			var e *HeatMap

			BeforeEach(func() {
				e, _ = New(3, 2, [][]float32{{0, 1}, {0, 1}, {0, 1}})
			})

			It("returns a new HeatMap with only the max columns of the original set to 1, and 0 eveywhere else", func() {
				Expect(r).To(Equal(e))
			})
		})
	})
	Describe("MaxX", func() {
		var a *HeatMap
		var r *HeatMap

		BeforeEach(func() {
			a, _ = New(3, 2, [][]float32{{0.1, 0.3}, {0, 0.2}, {0.01, 0.5}})
		})

		JustBeforeEach(func() {
			r = a.MaxX()
		})

		Context("when called", func() {
			var e *HeatMap

			BeforeEach(func() {
				e, _ = New(3, 2, [][]float32{{1, 0}, {0, 0}, {0, 1}})
			})

			It("returns a new HeatMap with only the max rows of the original set to 1, and 0 eveywhere else", func() {
				Expect(r).To(Equal(e))
			})
		})
	})

	Describe("sliceToMax", func() {
		var a, r []float32

		JustBeforeEach(func() {
			r = sliceToMax(a)
		})

		Context("when called on an array with a single max", func() {
			var e []float32

			BeforeEach(func() {
				a = []float32{1, 2, -3, 0, 1.5}
				e = []float32{0, 1, 0, 0, 0}
			})

			It("return the expected result", func() {
				Expect(r).To(Equal(e))
			})
		})

		Context("when called on an array with multiple max", func() {
			var e []float32

			BeforeEach(func() {
				a = []float32{2, -3, 2, 1.5}
				e = []float32{1, 0, 1, 0}
			})

			It("return the expected result", func() {
				Expect(r).To(Equal(e))
			})
		})
	})

	Describe("transpose", func() {
		var a, r [][]float32

		JustBeforeEach(func() {
			r = transpose(a)
		})

		BeforeEach(func() {
			a = [][]float32{{1, 2}, {3, 4}, {5, 6}}
		})

		Context("when called", func() {
			var e [][]float32

			BeforeEach(func() {
				e = [][]float32{{1, 3, 5}, {2, 4, 6}}
			})

			It("returns the original data with the columns and rows swapped", func() {
				Expect(r).To(Equal(e))
			})
		})
	})
})
