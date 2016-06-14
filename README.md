
TODO

- [x] Read a Fasta file, output a fasta file
- [x] Species Name and Gene Name schemas
- [x] Read a Fasta File, output a TNT file
      ccode and cgroup can be ignored
- [x] Support blocks and cnames in TNT
- [ ] Support single 'block' in tnt. This needs conditional using of xgroups
      when the number of blocks  == 1, and blocks when greater (verify this)
- [ ] Support Outgroup definition in TNT.
- [ ] In depth handling of '-h' from the interface; the simple one line usages
      are not enough.
- [ ] Structure configuration in such a way that reproducable pipelines can be
      easily set up, and the pipeline can be saved as a byproduct of a manual
      run.
- [ ] Coherent Errors
- [ ] Refactor out the sequence specific stuf from tnt into sequence
- [ ] Guess the Species from the name
  - [x] Just use the Name
  - [ ] Regexp rule
- [ ] Read a Fasta File, output a Nexus File
- [ ] Identify potentially missnamed species ( species names off by
      white space, special characters, or a couple characters
      by some language disntance metric
- [ ] Support Interleaving of Fasta
- [ ] Support Interleaving of TNT
