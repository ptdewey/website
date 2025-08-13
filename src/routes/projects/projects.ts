import type { ProjectCategory } from "$lib/types";

export let featuredProjects: ProjectCategory = {
  category: "Featured Projects",
  items: [
    {
      title: "PlantUML LSP",
      description:
        "An implementation of the language server protocol (LSP) for PlantUML, providing autocompletion, definitions, and diagnostics for PlantUML diagrams. Built accompanying Neovim and VSCode plugins for hassle-free editor integration.",
      link: "https://github.com/ptdewey/plantuml-lsp",
      time: "Summer 2024 - Present",
      languages: "Go",
    },
    {
      title: "Pendulum-nvim",
      description:
        "Neovim plugin that monitors coding duration and compiles metrics such as git project name, file type, and additional parameters, providing insightful data for productivity analysis.",
      image:
        "https://github.com/ptdewey/pendulum-nvim/raw/main/assets/screenshot0.png",
      link: "https://github.com/ptdewey/pendulum-nvim",
      time: "Spring 2024 - Present",
      languages: "Go, Lua",
    },
  ],
};

export let projects: ProjectCategory[] = [
  {
    category: "Other Projects",
    items: [
      {
        title: "YankBank-nvim",
        description:
          "Versatile Neovim plugin that offers an enhanced clipboard history interface with a quick-access menu, featuring session persistence through SQLite.",
        shortDescription:
          "Neovim plugin that provides quick access to clipboard history.",
        image:
          "https://github.com/ptdewey/yankbank-nvim/raw/main/assets/screenshot-2.png",
        link: "https://github.com/ptdewey/yankbank-nvim",
        time: "Spring 2024 - Present",
        languages: "Lua",
      },
      {
        title: "DarkEarth-nvim",
        description: "A dark and earthy color scheme for Neovim.",
        shortDescription: "A dark and earthy color scheme for Neovim.",
        image:
          "https://github.com/ptdewey/darkearth-nvim/raw/main/assets/color_bar.png",
        link: "https://github.com/ptdewey/darkearth-nvim",
        time: "Spring 2024 - Present",
        languages: "Lua",
      },
      {
        title: "MonaLisa-nvim",
        description:
          "A dark and colorful theme for Neovim based on the painting.",
        shortDescription:
          "A dark and colorful theme for Neovim based on the painting.",
        image:
          "https://github.com/ptdewey/monalisa-nvim/raw/main/assets/screenshot1.png",
        link: "https://github.com/ptdewey/monalisa-nvim",
        time: "Spring 2025 - Present",
        languages: "Lua",
      },
      {
        title: "Bluprinter",
        description:
          "An extensible template management tool with a beautiful terminal interface used for generating commonly used files. Written in Go using Bubble Tea.",
        shortDescription:
          "CLI tool for quickly generating files from templates.",
        image:
          "https://github.com/ptdewey/blueprinter/raw/main/assets/screenshot-1.png",
        link: "https://codeberg.org/ptdewey/blueprinter",
        time: "Summer 2024 - Fall 2024",
        languages: "Go",
      },
    ],
  },
  {
    category: "Open Source Contributions",
    items: [
      {
        title: "aaaton/golem",
        description:
          "Lemmatization library for Go. Fixed dictionary encoding issue where a zero-width byte order mark would be included in some outputs and added associated test cases.",
        link: "https://github.com/aaaton/golem",
      },
      {
        title: "Myriad-Dreamin/tinymist",
        description:
          "Language server for Typst. Added and fixed Neovim related documentation.",
        link: "https://github.com/Myriad-Dreamin/tinymist",
      },
      {
        title: "fredrikaverpil/godoc.nvim",
        description:
          "Fuzzy search Go packages/symbols and view docs from within Neovim. Improved support for fzf_lua.",
        link: "https://github.com/fredrikaverpil/godoc.nvim",
      },
      {
        title: "XAMPPRocky/tokei",
        description:
          "CLI app counting lines of code in a project. Added support for PlantUML.",
        link: "https://github.com/XAMPPRocky/tokei",
      },
      {
        title: "letieu/harpoon-lualine",
        description:
          "Harpoon extension for lualine integrating with Harpoon to show tracked files. Fixed a bug for empty Harpoon lists.",
        link: "https://github.com/letieu/harpoon-lualine",
      },
      {
        title: "karthik/wesanderson",
        description:
          "Color palette library for R. Added Asteroid City-themed palettes.",
        link: "https://github.com/karthik/wesanderson",
      },
      {
        title: "rockerBOO/awesome-neovim",
        description:
          "Curated list of Neovim plugins. Added and updated plugins.",
        link: "https://github.com/rockerBOO/awesome-neovim",
      },
    ],
  },
];

export default { projects: projects, featuredProjects: featuredProjects };
