--[[
=====================================================================
==================== HECATE NEOVIM CONFIGURATION ====================
=====================================================================

Welcome to the HecateShell Neovim configuration!

This is a modular, organized config built on kickstart.nvim principles.
All configuration is split into logical modules:

  • lua/config/options.lua   - Neovim options and settings
  • lua/config/keymaps.lua   - General keybindings
  • lua/config/autocmds.lua  - Autocommands
  • lua/config/lazy.lua      - Plugin manager setup

  • lua/plugins/ui.lua       - UI and visual plugins
  • lua/plugins/editor.lua   - Editor enhancement plugins
  • lua/plugins/coding.lua   - LSP, completion, formatting

  • lua/plugins/hecate-colors.lua - HecateShell theme (auto-generated)

For help, run :Tutor or :help
For plugin management, run :Lazy
For LSP/tool management, run :Mason

=====================================================================
--]]

-- Load core configuration
require 'config.options'
require 'config.keymaps'
require 'config.autocmds'
require 'config.lazy'

-- The line beneath this is called `modeline`. See `:help modeline`
-- vim: ts=2 sts=2 sw=2 et
