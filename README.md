# Gavel+

[WIP] A project designed to re-write Gavel using modern technologies aimed at optimizing the user experience of judges and admin users. This is NOT a rewrite of the underlying algorithm, and the formulas from the [original Gavel by anishathalye](https://github.com/anishathalye/gavel) are still used. Refer to [this excellent article](https://www.anishathalye.com/2015/03/07/designing-a-better-judging-system/) for more details on the underlying workings of Gavel!

What IS different is the architecture that I'm using. I will be developing this project in primarily Rust as opposed to Python. Additionally, I will be using MongoDB instead of Postgresql. Finally, the main deployment method will be through a Docker cluster instead of Heroku.

# Requirements

* [GNU Scientific Library (GSL)](https://www.gnu.org/software/gsl/)

