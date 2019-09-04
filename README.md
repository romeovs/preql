# preql

A simple compile time query builder.


`preql` will search a package for sql queryies and generate simple helper
functions that perform that query at compile time so the runtime does not have
to parse queries anymore.

In addition `preql` can generate SQL scanners for the types in your code.
