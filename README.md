# Stats Keeper

`stats-keeper` will be the backend of an application that users can store, update and view any number of personal statistics that they want.

## Structure

The system shall support a certain set of statistic types such as `counter`, `date`, `multi-value` etc. It shall be built in a way that it will allow extending the existing set of types without altering them.

It will be a Rest API, implemented in `Go`. As DBMS, `MongoDB` shall be used for its document-based structure and that the system won't need the strict relational structure of SQL-based DBMS's.

Data types shall be defined in `Protocol Buffers` files and will be compiled into `.go` files for use.
