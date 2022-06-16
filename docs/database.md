### database notes

you may see the ridiculous "overuse" of generics and nesting in this package. do not worry, I am aware of it myself. this generic overuse has arisen from a few properties that I have discovered on the database library. these properties make it useful for generics to exist.

1. **nested structs serialize properly** it's totally fine to put a nested struct into your struct, save it, and have it intact when you read it again
1. **generic nested structs work as expected** each type that is passed as a type parameter to a generic struct is treated as different type.
1. **generic nested structs work, even if the underlying type shape is the same** so if you passed a typedef for `int` to a generic struct and then passed another one for `int` just under a different name, they would be treated as separate types

if you wish to test these properties and ensure they do not break, please take a look at `tests/db`, although do make sure to test each case one by one.
