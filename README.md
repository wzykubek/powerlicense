# Licensmith

Crafting the ideal license for your Git repository in seconds!

## Getting Started

Licensmith, a streamlined tool, allows you to create an `LICENSE` file for your Git repository with ease, using just one command. This tool is designed to save you time and effort.

### Usage

To generate an ISC `LICENSE` file with the current year and your name, run the following command:

```bash
licensmith add ISC
```

By default, Licensmith searches for user details in your local repository (name and email) as a fallback option. It uses global configuration if no local information is found.

You can customize this process by providing specific values using the following command:

```bash
licensmith add ISC --name "John Doe" --email "jdoe@example.com"
```

To view available templates, run the following command:

```bash
licensmith list
```

To display a license summary, use:

```bash
licensmith show ISC
```

### Installation

Licensmith can be installed using various methods:

1. **Prebuilt Binaries:**
    - For stable versions, visit the [Releases](https://github.com/wzykubek/licensmith/releases) page.
    - To access development versions, check out the [Actions](https://github.com/wzykubek/licensmith/actions).

2. **Installation from Source:** Refer to the [compilation section](#compilation) for step-by-step instructions.

### Compilation

To build Licensmith from source, follow these steps:

```bash
git clone https://github.com/wzykubek/licensmith
cd licensmith
go build -v ./...
```
