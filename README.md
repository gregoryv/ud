ud - Update document

Small utility for replacing elements by id within html documents.

## Quick start

    go get -u github.com/gregoryv/ud

Replace an element by id

    echo "<em>new thing</em>" | ud -w -i "someid" -html index.html

Replace content of element by id use the `-c` flag

    echo "<em>new thing</em>" | ud -w -c -i "someid" -html index.html
