# Jury

[WIP] A project designed to re-write Gavel using modern technologies aimed at optimizing the user experience of judges and admin users. This is NOT a rewrite of the underlying formulas from the [original Gavel by anishathalye](https://github.com/anishathalye/gavel). Refer to [this excellent article](https://www.anishathalye.com/2015/03/07/designing-a-better-judging-system/) for more details on the underlying math workings of Gavel!

# External System Requirements

> Note: This might be mitigated with the use of a dockerized development environment (WIP)

* [GNU Scientific Library (GSL)](https://www.gnu.org/software/gsl/)

# Developing

Requirements:

* [yarn](https://yarnpkg.com/)
* [cargo](https://doc.rust-lang.org/cargo/)

**THIS PROCESS WILL BE DOCKERIZED _VERY_ SOON**

Client dev server:

```
cd client
yarn install
yarn start
```

Backend dev server:

```
cargo run
```
