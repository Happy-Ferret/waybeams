package display

import "math"

type LayoutAxis int

const (
	LayoutHorizontal = iota
	LayoutVertical
)

// This pattern is probably not the way to go, but I'm having trouble finding a
// reasonable alternative. The problem here is that LayoutHandler types will not be
// user-extensible. Box definitions will only be able to refer to the
// Layouts that have been enumerated here. The benefit is that ComponentModel objects
// will remain serializable and simply be a bag of scalars. I'm definitely
// open to suggestions.
type LayoutTypeValue int

const (
	// GROSS! I'm sure I've done something wrong here, but the "zero value" for
	// an enum field (above) is 0. This means that not setting the enum will
	// automatically set it to the first value in this list. :barf:
	// DO NOT SORT THESE ALPHABETICALLY!
	StackLayoutType = iota
	// DO NOT SORT
	VerticalFlowLayoutType
	// DO NOT SORT
	HorizontalFlowLayoutType
	// DO NOT SORT
	RowLayoutType
)

// Constants to represent Alignment of Component children, text or any other
// alignable collections.
type Alignment int

const (
	BottomAlign = iota
	LeftAlign
	RightAlign
	TopAlign
)

// Concrete implementation of a given layout. These handlers are pure functions
// that accept a Displayable and manage the scale and position of the children
// for that element.
type LayoutHandler func(d Displayable)

// These entities are stateless bags of hooks that allow us to apply
// the exact same layout rules on both supported axes.
var hDelegate *horizontalDelegate
var vDelegate *verticalDelegate

// Instantiate each delegate once the declarations are ready
func init() {
	hDelegate = &horizontalDelegate{}
	vDelegate = &verticalDelegate{}
}

// Arrange children in a vertical flow and use displayStack for horizontal rules.
func StackLayout(d Displayable) {
	if d.GetChildCount() == 0 {
		return
	}

	if hDelegate.GetFixed(d) == 0 && hDelegate.GetFlex(d) == 0 {
		hDelegate.ActualSize(d, hDelegate.GetChildrenSize(d))
	}

	if vDelegate.GetFixed(d) == 0 && vDelegate.GetFlex(d) == 0 {
		vDelegate.ActualSize(d, vDelegate.GetChildrenSize(d))
	}

	stackScaleChildren(hDelegate, d)
	stackScaleChildren(vDelegate, d)

	stackPositionChildren(hDelegate, d)
	stackPositionChildren(vDelegate, d)
}

func HorizontalFlowLayout(d Displayable) {
	if d.GetChildCount() == 0 {
		return
	}

	flowScaleChildren(hDelegate, d)
	stackScaleChildren(vDelegate, d)

	flowPositionChildren(hDelegate, d)
	stackPositionChildren(vDelegate, d)
}

func VerticalFlowLayout(d Displayable) {
	if d.GetChildCount() == 0 {
		return
	}

	stackScaleChildren(hDelegate, d)
	flowScaleChildren(vDelegate, d)

	stackPositionChildren(hDelegate, d)
	flowPositionChildren(vDelegate, d)
}

func notExcludedFromLayout(d Displayable) bool {
	return !d.GetExcludeFromLayout()
}

func isFlexible(d Displayable) bool {
	return d.GetFlexWidth() > 0 || d.GetFlexHeight() > 0
}

// Collect the layoutable children of a Displayable
func getLayoutableChildren(d Displayable) []Displayable {
	return d.GetFilteredChildren(notExcludedFromLayout)
}

func getFlexibleChildren(delegate LayoutDelegate, d Displayable) []Displayable {
	return d.GetFilteredChildren(func(child Displayable) bool {
		return notExcludedFromLayout(child) && delegate.GetIsFlexible(child)
	})
}

func getNotExcludedFromLayoutChildren(delegate LayoutDelegate, d Displayable) []Displayable {
	return d.GetFilteredChildren(func(child Displayable) bool {
		return notExcludedFromLayout(child)
	})
}

func getStaticChildren(d Displayable) []Displayable {
	return d.GetFilteredChildren(func(child Displayable) bool {
		return notExcludedFromLayout(child) && !isFlexible(child)
	})
}

func getStaticSize(delegate LayoutDelegate, d Displayable) float64 {
	sum := 0.0
	staticChildren := getStaticChildren(d)
	for _, child := range staticChildren {
		sum += delegate.GetSize(child)
	}
	return sum
}

func flowScaleChildren(delegate LayoutDelegate, d Displayable) {
	flexibleChildren := getFlexibleChildren(delegate, d)
	unitSize := flowGetUnitSize(delegate, d, flexibleChildren)
	for _, child := range flexibleChildren {
		value := math.Floor(delegate.GetFlex(child) * unitSize)
		delegate.ActualSize(child, value)
	}
	// flowSpreadRemainder(delegate, flexibleChildren)
}

func flowPositionChildren(delegate LayoutDelegate, d Displayable) {
	children := getNotExcludedFromLayoutChildren(delegate, d)
	position := delegate.GetPaddingFirst(d)
	// gutter := delegate.GetGutter(d)
	for _, child := range children {
		delegate.Position(child, position)
		position = position + delegate.GetSize(child) // + gutter
	}
}

func flowSpreadRemainder(delegate LayoutDelegate, flexibleChildren []Displayable) {
	/*
			// TODO(lbayes): Introduce this when needed
		// Spread remainder pixels from right to left
		var difference:Number = (delegate.actual - delegate.padding) - aggregateActualChildrenSize(delegate);
		var index:int = kids.length - 1;
		if(index == -1) {
			return;
		}
		while(difference-- > 0) {
		if(index == 0) {
		// We've reached the first child,
		// go ahead and push the entire remainder
		kids[index].actual += difference;
		break;
		}
		kids[index].actual += 1;
		index--;
		}
	*/
}

func flowGetUnitSize(delegate LayoutDelegate, d Displayable, flexibleChildren []Displayable) float64 {
	availablePixels := getAvailablePixels(delegate, d)
	flexSum := flowGetFlexSum(delegate, flexibleChildren)
	return availablePixels / flexSum
}

func flowGetFlexSum(delegate LayoutDelegate, flexibleChildren []Displayable) float64 {
	sum := 0.0
	for _, child := range flexibleChildren {
		sum += delegate.GetFlex(child)
	}
	return sum
}

func stackScaleChildren(delegate LayoutDelegate, d Displayable) {
	flexChildren := getFlexibleChildren(delegate, d)

	if len(flexChildren) == 0 {
		return
	}

	availablePixels := getAvailablePixels(delegate, d)

	for _, child := range flexChildren {
		delegate.ActualSize(child, availablePixels)
	}
}

// Get the (Size - Padding) on delegated axis for STACK layouts.
// NOTE: Flow layouts will also take into account the non-flexible children.
func getAvailablePixels(delegate LayoutDelegate, d Displayable) float64 {
	return delegate.GetSize(d) - delegate.GetPadding(d)
}

func stackGetUnitSize(delegate LayoutDelegate, d Displayable, flexPixels float64) float64 {
	return delegate.GetFlex(d) * flexPixels
}

func stackPositionChildren(delegate LayoutDelegate, d Displayable) {
	// TODO(lbayes): Work with alignment (first, center, last == left, center, right or top, center, bottom)

	// Position all children in upper left of container
	pos := delegate.GetPaddingFirst(d)
	for _, child := range getLayoutableChildren(d) {
		delegate.Position(child, pos)
	}
}

// Delegate for all properties that are used for Horizontal layouts
type horizontalDelegate struct{}

func (h *horizontalDelegate) ActualSize(d Displayable, size float64) {
	d.ActualWidth(size)
}

func (h *horizontalDelegate) GetActualSize(d Displayable) float64 {
	return d.GetActualWidth()
}

func (h *horizontalDelegate) GetAlign(d Displayable) Alignment {
	return d.GetHAlign()
}

func (v *horizontalDelegate) GetAxis() LayoutAxis {
	return LayoutHorizontal
}

func (h *horizontalDelegate) GetChildrenSize(d Displayable) float64 {
	return 0.0
}

func (h *horizontalDelegate) GetFixed(d Displayable) float64 {
	return d.GetFixedWidth()
}

func (h *horizontalDelegate) GetFlex(d Displayable) float64 {
	return d.GetFlexWidth()
}

func (h *horizontalDelegate) GetIsFlexible(d Displayable) bool {
	return d.GetFlexWidth() > 0
}

func (h *horizontalDelegate) GetMinSize(d Displayable) float64 {
	return d.GetMinWidth()
}

func (h *horizontalDelegate) GetPadding(d Displayable) float64 {
	return d.GetHorizontalPadding()
}

func (h *horizontalDelegate) GetPaddingFirst(d Displayable) float64 {
	return d.GetPaddingLeft()
}

func (h *horizontalDelegate) GetPaddingLast(d Displayable) float64 {
	return d.GetPaddingRight()
}

func (h *horizontalDelegate) GetPosition(d Displayable) float64 {
	return d.GetX()
}

func (h *horizontalDelegate) GetPreferred(d Displayable) float64 {
	return d.GetPrefWidth()
}

func (h *horizontalDelegate) GetSize(d Displayable) float64 {
	return d.GetWidth()
}

func (h *horizontalDelegate) Position(d Displayable, pos float64) {
	d.X(pos)
}

// Delegate for all properties that are used for Vertical layouts
type verticalDelegate struct{}

func (v *verticalDelegate) ActualSize(d Displayable, size float64) {
	d.ActualHeight(size)
}

func (v *verticalDelegate) GetActualSize(d Displayable) float64 {
	return d.GetActualHeight()
}

func (v *verticalDelegate) GetAlign(d Displayable) Alignment {
	return d.GetVAlign()
}

func (v *verticalDelegate) GetAxis() LayoutAxis {
	return LayoutVertical
}

func (v *verticalDelegate) GetChildrenSize(d Displayable) float64 {
	return 0.0
}

func (v *verticalDelegate) GetFixed(d Displayable) float64 {
	return d.GetFixedHeight()
}

func (v *verticalDelegate) GetFlex(d Displayable) float64 {
	return d.GetFlexHeight()
}

func (v *verticalDelegate) GetIsFlexible(d Displayable) bool {
	return d.GetFlexHeight() > 0
}

func (v *verticalDelegate) GetMinSize(d Displayable) float64 {
	return d.GetMinHeight()
}

func (v *verticalDelegate) GetPadding(d Displayable) float64 {
	return d.GetVerticalPadding()
}

func (v *verticalDelegate) GetPaddingFirst(d Displayable) float64 {
	return d.GetPaddingTop()
}

func (v *verticalDelegate) GetPaddingLast(d Displayable) float64 {
	return d.GetPaddingBottom()
}

func (v *verticalDelegate) GetPosition(d Displayable) float64 {
	return d.GetY()
}

func (v *verticalDelegate) GetPreferred(d Displayable) float64 {
	return d.GetPrefHeight()
}

func (v *verticalDelegate) GetSize(d Displayable) float64 {
	return d.GetHeight()
}

func (v *verticalDelegate) GetStaticSize(d Displayable) float64 {
	return 0.0
}

func (v *verticalDelegate) Position(d Displayable, pos float64) {
	d.Y(pos)
}

type LayoutDelegate interface {
	ActualSize(d Displayable, size float64)
	GetActualSize(d Displayable) float64
	GetAlign(d Displayable) Alignment
	GetAxis() LayoutAxis
	GetChildrenSize(d Displayable) float64
	GetFixed(d Displayable) float64
	GetFlex(d Displayable) float64 // GetPercent?
	GetIsFlexible(d Displayable) bool
	GetMinSize(d Displayable) float64
	GetPadding(d Displayable) float64
	GetPaddingFirst(d Displayable) float64
	GetPaddingLast(d Displayable) float64
	GetPosition(d Displayable) float64
	GetPreferred(d Displayable) float64
	GetSize(d Displayable) float64
	Position(d Displayable, pos float64)
}