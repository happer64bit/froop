# Froop

Froop is a project which aim is to share the file steamlessly and securely on the network event it's private or public.

Froop is 54.43x times faster than `python -m http.server`.

## Usage

### Start Server

```
froop serve -a <ADDRESS> -p <port>
```

#### With Auth

```
froop serve -a <ADDRESS> -p <port> --auth username:password
```

## Installation

### Via Golang Package Manager (Recommended)

```sh
go install github.com/happer64bit/froop
```

### Building from source

- golang (https://go.dev/)

```sh
git clone https://github.com/happer64bit/froop

cd froop

go build
```

### Adding to PATH (Linux)

```
sudo cp froop /usr/bin
```

### Adding to PATH (Windows)

For Windows, you may need to declare or add a system environment variable with the path of the executable

## Benchmarking

> [!INFO]
> These benchmarking result are tested by a tool called `oha`.

|  | Froop | http.server |
|---|---|---|
| Total | 22.7559s |1238.5834s |
| Fastest | 00009s | .0121s |
| Success Rate | 100.00% | 98.28% |

## Dontation

If you are looking ahead to support the development of the project, you can donate via crypto.

* TON - `UQBZXR35Z7KYeuwdqfZIDgAl6_wQxTbYuvRqT9I6CVonX0f_`
* SOLANA - `EK48PtUR2vXA7wTzUiNuUpmS6peXWWmmrqncNxFn3sL7`
* ERC20 - `0x6e47eDAdA0A25f38f2c2e3851256E455ed17A8A0`
