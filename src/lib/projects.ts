export interface Project {
  name: string;
  href: string;
  desc: string | string[];
}

export const featured: Project[] = [
  {
    name: "arabica.social",
    href: "https://alpha.arabica.social",
    desc: "coffee journaling social app, built on the AT Protocol",
  },
  {
    name: "shutter",
    href: "https://github.com/ptdewey/shutter",
    desc: "approval-based snapshot testing library for Go",
  },
];

export const other: Project[] = [
  {
    name: "pendulum-nvim",
    href: "https://github.com/ptdewey/pendulum-nvim",
    desc: "Neovim plugin for coding-time metrics: git project, file type, duration",
  },
  {
    name: "yankbank-nvim",
    href: "https://github.com/ptdewey/yankbank-nvim",
    desc: "Neovim plugin: quick access to clipboard history",
  },
  {
    name: "darkearth-nvim",
    href: "https://github.com/ptdewey/darkearth-nvim",
    desc: "dark and earthy colorscheme for Neovim",
  },
  {
    name: "plantuml-lsp",
    href: "https://github.com/ptdewey/plantuml-lsp",
    desc: "language server for PlantUML",
  },
];

export const contributions: Project[] = [
  {
    name: "tangled.org/core",
    href: "https://tangled.org/tangled.org/core",
    desc: [
      'added repo "starrers" page',
      "show issue/pull request counts on search",
    ],
  },
  {
    name: "gleam-lang/gleam",
    href: "https://github.com/gleam-lang/gleam",
    desc: [
      "expanded int/float exhaustiveness checks",
      "fixed invalid JS codegen for large floats",
    ],
  },
  {
    name: "teal-fm/piper",
    href: "https://github.com/teal-fm/piper",
    desc: ["added nix package and nixos module"],
  },
  {
    name: "aaaton/golem",
    href: "https://github.com/aaaton/golem",
    desc: "fixed BOM dictionary-encoding issue, added supporting test cases",
  },
  {
    name: "karthik/wesanderson",
    href: "https://github.com/karthik/wesanderson",
    desc: "Asteroid City themed colorschemes",
  },
];
