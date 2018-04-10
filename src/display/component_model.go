package display

type ComponentModel struct {
	// Application
	FramesPerSecond int

	// General
	ID          string
	Text        string
	Title       string
	Description string

	// Layout
	ActualHeight      float64
	ActualWidth       float64
	Builder           Builder
	Disabled          bool
	ExcludeFromLayout bool
	FlexHeight        float64
	FlexWidth         float64
	HAlign            Alignment
	Height            float64
	Hidden            bool
	LayoutType        LayoutTypeValue
	MaxHeight         float64
	MaxWidth          float64
	MinHeight         float64
	MinWidth          float64
	Padding           float64
	PaddingBottom     float64
	PaddingLeft       float64
	PaddingRight      float64
	PaddingTop        float64
	PrefHeight        float64
	PrefWidth         float64
	TraitNames        []string
	VAlign            Alignment
	Width             float64
	X                 float64
	Y                 float64
	Z                 float64

	// Style
	BgColor     int
	FontColor   int
	FontFace    string
	FontSize    int
	StrokeColor int
	StrokeSize  int
}
