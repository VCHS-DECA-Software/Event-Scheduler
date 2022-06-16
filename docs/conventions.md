### links

- when constructing new links, always use the simplest form of type parameters (so, the least # of wrappings possible) as the type parameter is not used for anything outside of differentiation
- link order is significant so always follow the orders below
    - `account (F) - grouping (T)`

### structs (database objects)

- when passing around structs, minimize use of pointers besides when reading / writing to the database.
    - there are reasons to use either method [see here](https://medium.com/a-journey-with-go/go-should-i-use-a-pointer-instead-of-a-copy-of-my-struct-44b43b104963), pointers are generally less performant because they use the heap and require GC. but if one requires extensive updates to a struct, it may be preferable to use a pointer and avoid redundant copies of data.
    - to actually determine which will be more effective however, will require more robust testing and profiling. in this project, values are generally considered immutable so pointers do not really add any value.

