/** @import {ProjectCategory} from "$lib/types" */

/** @type {ProjectCategory} */
export let featuredProjects = {
  category: "Featured Projects",
  items: [
    {
      title: "Oolong",
      description:
        "Platform agnostic, next gen note taking application with automatic note linking. Uses a custom keyword extraction system to link notes and ideas, enabling their visualization in a 2D/3D force-directed graph.",
      image: "/images/oolong-graph-screenshot.png",
      link: "https://github.com/oolong-sh",
      time: "Fall 2024 - Present",
    },
    {
      title: "Pendulum-nvim",
      description:
        "Neovim plugin that monitors coding duration and compiles metrics such as git project name, file type, and additional parameters, providing insightful data for productivity analysis. Written in Go and Lua.",
      image:
        "https://github.com/ptdewey/pendulum-nvim/raw/main/assets/screenshot0.png",
      link: "https://github.com/ptdewey/pendulum-nvim",
      time: "Spring 2024 - Present",
    },

    {
      title: "PlantUML LSP",
      description:
        "An implementation of the language server protocol (LSP) for PlantUML, providing autocompletion, definitions, and diagnostics for PlantUML diagrams. Written in Go.",
      link: "https://github.com/ptdewey/plantuml-lsp",
      time: "Summer 2024 - Present",
    },
  ],
};

/** @type {ProjectCategory[]} */
export let projects = [
  {
    category: "Other Projects",
    items: [
      {
        title: "YankBank-nvim",
        description:
          "Versatile Neovim plugin that offers an enhanced clipboard history interface with a quick-access menu, featuring session persistence through SQLite. Written in Lua.",
        image:
          "https://github.com/ptdewey/yankbank-nvim/raw/main/assets/screenshot-2.png",
        link: "https://github.com/ptdewey/yankbank-nvim",
        time: "Spring 2024 - Present",
      },
      {
        title: "DarkEarth-nvim",
        description:
          "A dark and earthy color scheme for Neovim. Written in Lua.",
        image:
          "https://github.com/ptdewey/darkearth-nvim/raw/main/assets/color_bar.png",
        link: "https://github.com/ptdewey/darkearth-nvim",
        time: "Spring 2024 - Present",
      },
      {
        title: "FRC Scouting Database",
        description:
          "Cloud-deployable scouting system for FIRST Robotics Competition event stats and match predictions, achieving 80% match prediction accuracy. Written in Go.",
        link: "https://github.com/ptdewey/frc-scouting-database-v2",
        time: "Spring 2023 - Present",
      },
      {
        title: "Linux Dotfiles",
        description:
          "A collection of configuration files and scripts I use every day on my Linux desktop. Includes a quick setup script that allows me to get working on any system in minutes.",
        link: "https://github.com/ptdewey/dotfiles",
        time: "Summer 2016 - Present",
      },
      // {
      //   title: "Bluprinter",
      //   description:
      //     "An extensible template management tool with a beautiful terminal interface used for generating commonly used files. Written in Go using Bubble Tea.",
      //   // image:
      //   // "https://github.com/ptdewey/blueprinter/raw/main/assets/screenshot-1.png",
      //   link: "https://github.com/ptdewey/blueprinter",
      //   time: "Summer 2024 - Fall 2024",
      // },
      {
        title: "Visualizing *What* Neural Networks Learn",
        description:
          "Animated visualizations of neural network learning processes, built with Python and R.",
        link: "https://aanish-pradhan.github.io/CS-5764-Project/",
        time: "Spring 2024",
      },
      {
        title: "CUDA Neural Network",
        description:
          "Modular feed-forward neural network implementation in CUDA C++ with various activation and cost functions for classification and regression tasks.",
        link: "https://github.com/ptdewey/cuda-nn",
        time: "Fall 2023",
      },
    ],
  },
  {
    category: "Open Source Contributions",
    items: [
      {
        title: "XAMPPRocky/tokei",
        description:
          "CLI app counting lines of code in a project. Added support for PlantUML.",
        link: "https://github.com/XAMPPRocky/tokei",
      },
      {
        title: "nvim-lualine/lualine.nvim",
        description:
          "Customizable status bar plugin for Neovim. Added feature for filename display with parent directory in multi-buffer projects.",
        link: "https://github.com/nvim-lualine/lualine.nvim",
      },
      {
        title: "letieu/harpoon-lualine",
        description:
          "Harpoon extension for lualine integrating with Harpoon to show tracked files. Fixed a bug for empty Harpoon lists.",
        link: "https://github.com/letieu/harpoon-lualine",
      },
      {
        title: "rockerBOO/awesome-neovim",
        description:
          "Curated list of Neovim plugins. Contributed new and updated plugins.",
        link: "https://github.com/rockerBOO/awesome-neovim",
      },
      {
        title: "karthik/wesanderson",
        description:
          "Color palette library for R. Added Asteroid City-themed palettes.",
        link: "https://github.com/karthik/wesanderson",
      },
    ],
  },
];

export default { projects: projects, featuredProjects: featuredProjects };
