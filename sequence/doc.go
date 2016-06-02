package sequence

/* sequence holds the types and mechanism for generically storing, and
 * nucleotide or protine sequences
 *
 * I want this to lazily load the data; so we can run two passes for
 * quick checking without completly nuking the memory; (Purely an
 * idea right now, I don't know if its is a sane approach)

 * cite: https://github.com/gaurav/taxondna
 */
