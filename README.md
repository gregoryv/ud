[![Build Status](https://travis-ci.org/gregoryv/ud.svg?branch=master)](https://travis-ci.org/gregoryv/ud)
[![codecov](https://codecov.io/gh/gregoryv/ud/branch/master/graph/badge.svg)](https://codecov.io/gh/gregoryv/ud)

[ud](https://godoc.org/github.com/gregoryv/ud) - Update document

Utility for replacing elements by id within html files.

## Quick start

    go get -u github.com/gregoryv/ud

Replace an element by id

    echo "<em>new thing</em>" | ud -w -i "someid" -html index.html

which is same as

    echo '<em id="someid">new thing</em>' | ud -w -html index.html

Note! when `-i` flag is not given `-c` has no effect, it will always
replace the identified element.

Replace content of element by id use the `-c` flag

    echo "<em>new thing</em>" | ud -w -c -i "someid" -html index.html


## Primary usecase

The primary reason for this tool was to simply generate and update
image maps within html when working with graphviz documents

    dot -Tcmapx somegraph.dot | ud -w -html index.html
