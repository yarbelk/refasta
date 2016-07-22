# Refasta (temporary name)

This is a program to convert various biology formats from one into another.

*Warning* this project is Very Alpha, and its interface will change frequently.

This is born out of the complexity that arrises from the abuse and missuse of
biology file formats, such as [fasta](https://en.wikipedia.org/wiki/FASTA),
or the complexity of the formats, such as [TNT](http://www.lillo.org.ar/phylogeny/tnt/).


## Installation

if you have `go` installed on your system, you can `go get github.com/yarbelk/refasta`
Otherwise, look at the releases page.

TODO

- [x] Read a Fasta file, output a fasta file
- [x] Species Name and Gene Name schemas
- [x] Read a Fasta File, output a TNT file
      ccode and cgroup can be ignored
- [x] Support blocks and cnames in TNT
- [ ] Support single 'block' in tnt. This needs conditional using of xgroups
      when the number of blocks  == 1, and blocks when greater (verify this)
      - Question: What is the difference between xgroup and block?
- [ ] Support Outgroup definition in TNT (using outgroup command)
- [x] In depth handling of '-h' from the interface; the simple one line usages
      are not enough.
- [ ] Structure configuration in such a way that reproducable pipelines can be
      easily set up, and the pipeline can be saved as a byproduct of a manual
      run.
  - [x] switch to using [cli](https://github.com/urfave/cli) for the cli: this
        supports loading all arguments from the ENV or yaml files.
  - [ ] Implement loading and saving of pipelines using cli.
  - [ ] Document said usage
- [ ] Coherent Errors: All failure modes must have human readable errors, that
      the bioinformation can use to identify where the bad data is.
- [ ] Refactor out the sequence specific stuf from tnt into sequence
- [ ] Guess the Species from the name. This is also very specific to one
      kind of usage of the FASTA format.  Specifically using it as an interchange
      between something and TNT.  This should probably be a flag.
  - [x] Just use the Name
  - [ ] Regexp rule
- [ ] Read a Fasta File, output a Nexus File
- [ ] Identify potentially missnamed species ( species names off by
      white space, special characters, or a couple characters
      by some language disntance metric
- [ ] Support Interleaving of Fasta
- [ ] Support Interleaving of TNT
