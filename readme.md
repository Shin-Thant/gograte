# Gograte: Database Migration Tool

### Motivations

This project is highly inspired by [Goose](https://github.com/pressly/goose) which is also a database migration tool. I was susprised by the simplicity of this tool. So, I wanted to build by myself.

### Guides

Move to `cmd/gograte` directory and run `go install`.
Then you are good to go.
Move to your project and use those commands.

1. Init your migration.
    ```sh
    gograte init
    ```
2. Create your migration sql file.
    ```sh
    gograte create first
    ```
    This will create `sql` file in this format **[timestamp]-[name].sql**
3. Migrate with this commands.
    ```sh
    gograte migrate [db_driver] [db_url] [action]
    ```
    Action option can be the followings.
    1. `up`
    2. `up-one`
    3. `up-to` (example `up-to [version]`)
    4. `down`
    5. `down-one`
    6. `down-to` (example `down-to [version]`)
4. Help
    ```
    gograte --help
    ```
