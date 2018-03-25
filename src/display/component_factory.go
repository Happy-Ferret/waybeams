package display

import (
	"errors"
)

type newComponent (func() Displayable)
type innerComponentFactory (func(b Builder, opts ...ComponentOption) (Displayable, error))

type ComponentFactory (func(c newComponent) innerComponentFactory)

// Returns a component factory that will properly accept options and register a
// component with the Builder.
//
// Also accepts zero or more ComponentOptions to set default values for components
// created from the returned function.
//
// Usage:
//   var MyComponent = NewComponentFactory(NewComponent, Padding(5))
//
// Callers can then:
//   box, err := MyComponent(NewBuilder(), FlexWidth(1), MaxWidth(100), MinWidth(10))
//
// Or:
//
//   root, err := MyComponent(NewBuilder(), Width(800), Height(600), Children(func(b Builder) {
//		Box(b, Id("one"), Height(80), FlexWidth(1))
//		Box(b, Id("two"), FlexHeight(1), FlexWidth(1))
//		Box(b, Id("three"), Height(60), FlexWidth(1))
//	 })
//
func NewComponentFactory(c newComponent, defaultOpts ...ComponentOption) innerComponentFactory {
	return func(b Builder, opts ...ComponentOption) (Displayable, error) {
		// Create a builder if we weren't provided with one. This makes tests much, much
		// more readable, but it not be expected
		if b == nil {
			return nil, errors.New("Compnent factory requires a Builder instance, try Component(NewBuilder()) or in the parent closure, add a (b Builder) argument and forward it to the child nodes.")
		}
		instance := c()
		opts = append(defaultOpts, opts...)
		// Instantiate the component from the provided factory function.
		// Apply all provided options to the component instance.
		for _, opt := range opts {
			err := opt(instance)
			if err != nil {
				// If an option error is found, bail with it.
				return nil, err
			}
		}

		// Send the instance to the provided builder for tree placement.
		b.Push(instance)
		// Everything worked great, return the instance.
		return instance, nil
	}
}

// Admittedly odd, poorly readable scheme for inheriting a set of default
// options from some other, known component factory.
// The idea here is that one factory may set up a set of defaults and another
// factory and add to that set. Here is a contrived example:
//
// Box := NewComponentFactory(NewComponent, MinWidth(100), MinHeight(200))
// LimitedBox := NewComponentFactoryFrom(Box, MaxWidth(200), MaxHeight(400))
//
// limited := LimitedBox(NewBuilder(), Width(300), Height(500))
//
// if limited.GetHeight() == 400 {
//     log.Printf("It worked!, the max height default was respected")
// }
//
func NewComponentFactoryFrom(baseFactory innerComponentFactory, defaultOpts ...ComponentOption) innerComponentFactory {
	return func(b Builder, opts ...ComponentOption) (Displayable, error) {
		opts = append(defaultOpts, opts...)
		instance, err := baseFactory(b, opts...)
		if err != nil {
			return nil, err
		}
		return instance, nil
	}
}

type ComponentOption (func(d Displayable) error)

func Id(value string) ComponentOption {
	return func(d Displayable) error {
		d.GetModel().Id = value
		return nil
	}
}

func Title(value string) ComponentOption {
	return func(d Displayable) error {
		d.GetModel().Title = value
		return nil
	}
}

func ExcludeFromLayout(value bool) ComponentOption {
	return func(d Displayable) error {
		d.ExcludeFromLayout(value)
		return nil
	}
}

func Width(value float64) ComponentOption {
	return func(d Displayable) error {
		d.Width(value)
		return nil
	}
}

func Height(value float64) ComponentOption {
	return func(d Displayable) error {
		d.Height(value)
		return nil
	}
}

func Size(width, height float64) ComponentOption {
	return func(d Displayable) error {
		d.Width(width)
		d.Height(height)
		return nil
	}
}

func MaxWidth(value float64) ComponentOption {
	return func(d Displayable) error {
		d.MaxWidth(value)
		return nil
	}
}

func MaxHeight(value float64) ComponentOption {
	return func(d Displayable) error {
		d.MaxHeight(value)
		return nil
	}
}

func MinWidth(value float64) ComponentOption {
	return func(d Displayable) error {
		d.MinWidth(value)
		return nil
	}
}

func MinHeight(value float64) ComponentOption {
	return func(d Displayable) error {
		d.MinHeight(value)
		return nil
	}
}

func PrefWidth(value float64) ComponentOption {
	return func(d Displayable) error {
		d.PrefWidth(value)
		return nil
	}
}

func PrefHeight(value float64) ComponentOption {
	return func(d Displayable) error {
		d.PrefHeight(value)
		return nil
	}
}

func FlexWidth(value float64) ComponentOption {
	return func(d Displayable) error {
		d.FlexWidth(value)
		return nil
	}
}

func FlexHeight(value float64) ComponentOption {
	return func(d Displayable) error {
		d.FlexHeight(value)
		return nil
	}
}

func HAlign(align Alignment) ComponentOption {
	return func(d Displayable) error {
		d.HAlign(align)
		return nil
	}
}

func VAlign(align Alignment) ComponentOption {
	return func(d Displayable) error {
		d.VAlign(align)
		return nil
	}
}

func X(pos float64) ComponentOption {
	return func(d Displayable) error {
		d.X(pos)
		return nil
	}
}

func Y(pos float64) ComponentOption {
	return func(d Displayable) error {
		d.Y(pos)
		return nil
	}
}

func Z(pos float64) ComponentOption {
	return func(d Displayable) error {
		d.Z(pos)
		return nil
	}
}

func LayoutType(layoutType LayoutTypeValue) ComponentOption {
	return func(d Displayable) error {
		d.LayoutType(layoutType)
		return nil
	}
}

func Padding(value float64) ComponentOption {
	return func(d Displayable) error {
		opts := d.GetModel()
		// Set the ComponentModel object directly
		if opts.PaddingBottom == 0 {
			opts.PaddingBottom = -1
		}
		if opts.PaddingLeft == 0 {
			opts.PaddingLeft = -1
		}
		if opts.PaddingRight == 0 {
			opts.PaddingRight = -1
		}
		if opts.PaddingTop == 0 {
			opts.PaddingTop = -1
		}
		opts.Padding = value
		return nil
	}
}

func PaddingBottom(value float64) ComponentOption {
	return func(d Displayable) error {
		d.PaddingBottom(value)
		return nil
	}
}

func PaddingLeft(value float64) ComponentOption {
	return func(d Displayable) error {
		d.PaddingLeft(value)
		return nil
	}
}

func PaddingRight(value float64) ComponentOption {
	return func(d Displayable) error {
		d.PaddingRight(value)
		return nil
	}
}

func PaddingTop(value float64) ComponentOption {
	return func(d Displayable) error {
		d.PaddingTop(value)
		return nil
	}
}

func AttrStyles(opts ...StyleOption) ComponentOption {
	styles := NewStyleDefinition()
	return func(d Displayable) error {
		for _, opt := range opts {
			opt(styles)
		}
		d.Styles(styles)
		return nil
	}
}

// Compose children onto the current component by providing a closure that
// either accepts zero arguments, or accepts a single argument which will
// be a function that, when called will invalidate the component instance
// for a future render.
func Children(composer interface{}) ComponentOption {
	return func(d Displayable) error {
		return d.Composer(composer)
	}
}