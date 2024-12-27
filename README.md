# Licensmith

Effortlessly craft the perfect LICENSE for your Git repo in seconds with a single command!

## Usage
This command will generate ISC `LICENSE` file in your current directory, including current year, and your name read from Git configuration:
```bash
licensmith --license ISC
```

By default, Licensmith read your local repository looking for user details (name and e-mail), as a fallback it uses global configuration.

You can also specify different values using:
```bash
licensmith --license ISC --name "John Doe" --email "jdoe@example.com"
```

To list available templates run:
```bash
licensmith --list
```
