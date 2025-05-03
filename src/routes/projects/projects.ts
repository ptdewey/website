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
      title: "Oolong",
      description:
        "Platform agnostic, next gen note taking application with automatic note linking. Uses an NLP-based keyword extraction system to link notes and ideas, enabling their visualization in a 2D/3D force-directed graph.",
      image: "/images/oolong-graph-screenshot.png",
      link: "https://github.com/oolong-sh",
      time: "Fall 2024 - Present",
      languages: "Go, TypeScript",
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
        title: "Bluprinter",
        description:
          "An extensible template management tool with a beautiful terminal interface used for generating commonly used files. Written in Go using Bubble Tea.",
        shortDescription:
          "CLI tool for quickly generating files from templates.",
        image:
          "https://github.com/ptdewey/blueprinter/raw/main/assets/screenshot-1.png",
        link: "https://github.com/ptdewey/blueprinter",
        time: "Summer 2024 - Fall 2024",
        languages: "Go",
      },
      {
        title: "Rooibos",
        description: "Programmatic resume generator.",
        link: "https://github.com/ptdewey/rooibos",
        time: "Winter 2025",
        languages: "Go, Lua, Typst",
      },
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
        title: "FRC Scouting Database",
        description:
          "Cloud-deployable scouting system for FIRST Robotics Competition event stats and match predictions, achieving 80% match prediction accuracy.",
        shortDescription:
          "Scouting and match prediction application for FIRST Robotics.",
        link: "https://github.com/ptdewey/frc-scouting-database-v2",
        time: "Spring 2023 - Present",
        languages: "Go",
      },
      {
        title: "Linux Dotfiles",
        description:
          "A collection of configuration files and scripts I use every day on my Linux desktop. Includes a quick setup script that allows me to get working on any system in minutes.",
        shortDescription: "Configuration files and scripts I use every day.",
        link: "https://github.com/ptdewey/dotfiles",
        time: "Summer 2016 - Present",
        languages: "Lua, Bash, Nix",
      },
      {
        title: "Visualizing *What* Neural Networks Learn",
        description:
          "Animated visualizations of neural network learning processes.",
        shortDescription: "A visual study of what neural networks learn.",
        link: "https://aanish-pradhan.github.io/CS-5764-Project/",
        time: "Spring 2024",
        languages: "Python, R",
      },
      {
        title: "CUDA Neural Network",
        description:
          "Modular feed-forward neural network implementation in CUDA C++ with various activation and cost functions for classification and regression tasks.",
        shortDescription:
          "Modular feed-forward neural network implementation in CUDA.",
        link: "https://github.com/ptdewey/cuda-nn",
        time: "Fall 2023",
        languages: "CUDA, C++",
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
          "Language server for Typst. Various Neovim related documentation additions and fixes.",
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
          "Curated list of Neovim plugins. Contributed new and updated plugins.",
        link: "https://github.com/rockerBOO/awesome-neovim",
      },
    ],
  },
];

export default { projects: projects, featuredProjects: featuredProjects };
