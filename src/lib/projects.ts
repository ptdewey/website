export interface Project {
  name: string;
  href: string;
  desc: string;
}

export const featured: Project[] = [
  { name: 'arabica.social', href: 'https://alpha.arabica.social', desc: 'coffee journaling social app, built on the AT Protocol' },
  { name: 'shutter', href: 'https://github.com/ptdewey/shutter', desc: 'approval-based snapshot testing library for Go' },
  { name: 'pendulum-nvim', href: 'https://github.com/ptdewey/pendulum-nvim', desc: 'Neovim plugin for coding-time metrics: git project, file type, duration' },
];

export const other: Project[] = [
  { name: 'yankbank-nvim', href: 'https://github.com/ptdewey/yankbank-nvim', desc: 'Neovim plugin — quick access to clipboard history' },
  { name: 'darkearth-nvim', href: 'https://github.com/ptdewey/darkearth-nvim', desc: 'dark and earthy colorscheme for Neovim' },
  { name: 'plantuml-lsp', href: 'https://github.com/ptdewey/plantuml-lsp', desc: 'language server for PlantUML — completions, definitions, diagnostics' },
  { name: 'blueprinter', href: 'https://github.com/ptdewey/blueprinter', desc: 'CLI tool for generating files from templates' },
  { name: 'bluesky-comments (Svelte)', href: 'https://github.com/ptdewey/bluesky-comments-svelte', desc: 'Svelte library — Bluesky comments on any page' },
];

export const contributions: Project[] = [
  { name: 'gleam-lang/gleam', href: 'https://github.com/gleam-lang/gleam', desc: 'improved int/float exhaustiveness checks; fixed invalid JS codegen for large floats' },
  { name: 'aaaton/golem', href: 'https://github.com/aaaton/golem', desc: 'fixed BOM dictionary-encoding issue, added test cases' },
  { name: 'Myriad-Dreamin/tinymist', href: 'https://github.com/Myriad-Dreamin/tinymist', desc: 'added and fixed Neovim documentation' },
  { name: 'Cyboard/zmk-keyboards', href: 'https://github.com/Cyboard-DigitalTailor/zmk-keyboards', desc: 'additional keyboard layout — single-arc, three-column dactyl' },
  { name: 'XAMPPRocky/tokei', href: 'https://github.com/XAMPPRocky/tokei', desc: 'PlantUML support for the LOC counter' },
  { name: 'rockerBOO/awesome-neovim', href: 'https://github.com/rockerBOO/awesome-neovim', desc: 'additional Neovim plugins' },
  { name: 'shoenot/witchesbrew.nvim', href: 'https://github.com/shoenot/witchesbrew.nvim', desc: 'improved build tooling, additional highlight groups' },
  { name: 'fredrikaverpil/godoc.nvim', href: 'https://github.com/fredrikaverpil/godoc.nvim', desc: 'fzf-lua support; reviewed experimental tree-sitter support' },
  { name: 'karthik/wesanderson', href: 'https://github.com/karthik/wesanderson', desc: 'Asteroid City colorschemes' },
];
