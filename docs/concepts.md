### object (represented by the `object` package)

this is a wrapper for any resource that should have CRUD methods (Create, Read, Update, Delete) these resources are differentiated by the type passed into the generic parameter

### link (represented by the `links` package)

this is an implementation of an association table found in your typical database. the difference is that each association is represented by just one struct `links.Link`

### account (represented by the `users` package)

this is a generic account, it allows one to store any metadata within it but shares the logic of creation and authentication. this abstracts the `student`, `judge` and `admin` structs from legacy code.

### grouping (represented by the `groupings` package)

this is a collection of objects, it does not contain the associations within itself (that is what a link does) but stores the metadata for the collection. this abstracts the `team` and `event` structs from legacy code.
