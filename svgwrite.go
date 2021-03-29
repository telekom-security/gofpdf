/*
 * Copyright (c) 2014 Kurt Jung (Gmail: kurt.w.jung)
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

 package gofpdf

 // SVGBasicWrite renders the paths encoded in the basic SVG image specified by
 // sb. The scale value is used to convert the coordinates in the path to the
 // unit of measure specified in New(). The current position (as set with a call
 // to SetXY()) is used as the origin of the image. The current line cap style
 // (as set with SetLineCapStyle()), line width (as set with SetLineWidth()),
 // and draw color (as set with SetDrawColor()) are used in drawing the image
 // paths.
 // styleStr can be "F" for filled, "D" for outlined only, or "DF" or "FD" for
 // outlined and filled. An empty string will be replaced with "D".
 // Path-painting operators as defined in the PDF specification are also
 // allowed: "S" (Stroke the path), "s" (Close and stroke the path),
 // "f" (fill the path, using the nonzero winding number), "f*"
 // (Fill the path, using the even-odd rule), "B" (Fill and then stroke
 // the path, using the nonzero winding number rule), "B*" (Fill and
 // then stroke the path, using the even-odd rule), "b" (Close, fill,
 // and then stroke the path, using the nonzero winding number rule) and
 // "b*" (Close, fill, and then stroke the path, using the even-odd
 // rule).
 func (f *Fpdf) SVGBasicWrite(sb *SVGBasicType, scale float64, styleStr string) {
	 originX, originY := f.GetXY()
	 var x, y, newX, newY float64
	 var cx0, cy0, cx1, cy1 float64
	 var path []SVGBasicSegmentType
	 var seg SVGBasicSegmentType
	 var startX, startY float64
	 sval := func(origin float64, arg int) float64 {
		 return origin + scale*seg.Arg[arg]
	 }
	 xval := func(arg int) float64 {
		 return sval(originX, arg)
	 }
	 yval := func(arg int) float64 {
		 return sval(originY, arg)
	 }
	 val := func(arg int) (float64, float64) {
		 return xval(arg), yval(arg + 1)
	 }
	 for j := 0; j < len(sb.Segments) && f.Ok(); j++ {
		 path = sb.Segments[j]
		 for k := 0; k < len(path) && f.Ok(); k++ {
			 seg = path[k]
			 switch seg.Cmd {
			 case 'M':
				 x, y = val(0)
				 startX, startY = x, y
				 f.MoveTo(x, y)
			 case 'L':
				 newX, newY = val(0)
				 f.LineTo(newX, newY)
			 case 'C':
				 cx0, cy0 = val(0)
				 cx1, cy1 = val(2)
				 newX, newY = val(4)
				 f.CurveBezierCubicTo(cx0, cy0, cx1, cy1, newX, newY)
				 x, y = newX, newY
			 case 'Q':
				 cx0, cy0 = val(0)
				 newX, newY = val(2)
				 f.CurveTo(cx0, cy0, newX, newY)
				 x, y = newX, newY
			 case 'H':
				 newX = xval(0)
				 f.LineTo(newX, f.GetY())
			 case 'V':
				 newY = yval(0)
				 f.LineTo(f.GetX(), newY)
			 case 'Z':
				 f.LineTo(startX, startY)
			 default:
				 f.SetErrorf("Unexpected path command '%c'", seg.Cmd)
			 }
		 }
	 }
	 f.DrawPath(styleStr)
 }
