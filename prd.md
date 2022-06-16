### datastructures

**note: all structures contain an "ID"**
- student: represents a student account
    - name
    - username
    - password (hashed)
- judge: represents a judge's account
    - name
    - username
    - password (hashed)
- admin: represents an admin account
    - name
    - username
    - password (hashed)

- team:
    - name
    - owner (a student ID)
- event
    - name
    - owner (an admin ID)
    - description
    - location
    - event type (enum)
    - start time
    - end times

- associations
    - student <-> team
    - team <-> event (constraints on # of associations are enforced)
    - judge <-> event (constraints on # of associations are enforced)

### methods

- require methods to CRUD all accounts
- require methods to CRUD all teams and their associations
- require methods to CRUD all events and their associations, as well as view all events collectively regardless of association
