# svgpoly is a package to support SVG to polygon conversion

## Overview

This package adds support for manipulating a
[polygon](http://zappem.net/pub/math/polygon/) object representation
of the content of SVGs. This
[package](http://zappem.net/pub/graphics/svgpoly/) is for converting
an SVG file input into to a collection of such polygons.

The primary use case is digesting SVG files from KiCad and automating
the outline generation of overlapping polygons.

## How to use

First confirm `gnuplot` is available on your system. Try:

```
$ gnuplot --version
```

If it is missing, install it (Fedora: `sudo dnf install gnuplot`,
Debian: `sudo apt install gnuplot`). Then:

```
$ git clone https://github.com/tinkerator/svgpoly.git
$ cd svgpoly
$ go run examples/outline.go -- --svg examples/test.svg --hatch 0.3 | gnuplot -p
```

Which should render this processed (union) image:

<img src="with-union-hatched.png" width="80%" alt="polygon outlines of shapes with hatch fill"/>

This is more faithful to the raw input SVG image in terms of the
overlapping polygons.

```
$ go run examples/outline.go --svg examples/test.svg --before --after=false | gnuplot -p
```

Which should render this image:

<img src="before-not-after.png" width="80%" alt="More raw SVG input"/>

We can inflate the polygons by the value specified with the `--inflate` option.

```
$ go run examples/outline.go --svg examples/test.svg --before --inflate 0.3 | gnuplot -p
```

Which should render this image:

<img src="before-and-inflate.png" width="80%" alt="Inflated union"/>

Finally, you can output `*polygons.Shapes` in the form of an SVG. The
SVG follows the conventions that `.Hole`s appear white and
non-`.Hole`s as blue:

```
$ go run examples/outline.go --svg examples/test.svg --osvg output.svg --hatch 0.3
```

Which renders as follows:

<img src="ref-output.svg" width="80%" alt="Inflated union with hatch"/>

Note how this output is flipped vertically from the `gnuplot`
output. This SVG output is faithful to the input format. The `gnuplot`
output is not. This is because the native input data is from an SVG
which has the opposite Y coordinate direction as `gnuplot`.

## License info

The `svgpoly` package is distributed with the same BSD 3-clause
[license](LICENSE) as that used by
[golang](https://golang.org/LICENSE) itself.

## Reporting bugs

This is a hobby project, so I can't guarantee a fix, but do use the
[github `svgpoly` bug
tracker](https://github.com/tinkerator/svgpoly/issues).
