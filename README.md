# Hunk parser (and general hunk info)
Aim is to add support for Amiga Hunk format to Clang/LLVM.  
we should be able to cross compile for Amiga on any modern computer

The parser is simply an exercise to understand the format but could be used to validate generated hunks

## unanswered questions
 * how do reloc hunks actually work? look at the Amigados manual page 371
 * how is the entry point defined?  is there a convention?
 * what are external symbols?
 * how are the code/data/bss sections represented by LLVM/Clang?
 * what is the minimum viable support needed?
 * how would we control things like hunk names, chip/fast mem requirements?  config file/pragmas?

## plan of attack
 * create simple hello world application (no external lib deps)
 * compile and link with sas/c
 * examine hunk output
 * disassemble / decompile code hunk
 * repeat for increasingly complex examples
   * strings / constants
   * variables (inintialized / uninitialized)
   * structs

## references
 * SAS/C 650 vol1 chapter 12 for how the hunk sections are generated / controlled
 * http://amiga-dev.wikidot.com/file-format:hunk
 * http://fileformats.archiveteam.org/wiki/Amiga_Hunk
 * https://en.wikipedia.org/wiki/Amiga_Hunk
 * The AmigaDos Manual chapter 10 for the hunk format definition
 * The AmigaDos Manual chapter 7 for details on what the linker does to create load files