-- Bootstrap lazy.nvim plugin manager
local lazypath = vim.fn.stdpath 'data' .. '/lazy/lazy.nvim'
if not (vim.uv or vim.loop).fs_stat(lazypath) then
  local lazyrepo = 'https://github.com/folke/lazy.nvim.git'
  local out = vim.fn.system { 'git', 'clone', '--filter=blob:none', '--branch=stable', lazyrepo, lazypath }
  if vim.v.shell_error ~= 0 then
    error('Error cloning lazy.nvim:\n' .. out)
  end
end

---@type vim.Option
local rtp = vim.opt.rtp
rtp:prepend(lazypath)

-- Setup lazy.nvim
require('lazy').setup({
  -- Import all plugin modules
  { import = 'plugins.ui' },
  { import = 'plugins.editor' },
  { import = 'plugins.coding' },
  { import = 'plugins.git' },

  require 'hecate.plugins.autopairs',
  require 'hecate.plugins.debug',
  require 'hecate.plugins.gitsigns',
  require 'hecate.plugins.indent_line',
  require 'hecate.plugins.lint',
  require 'hecate.plugins.neo-tree',

  -- Import custom plugins (optional)
  -- { import = 'custom.plugins' },
}, {
  ui = {
    icons = vim.g.have_nerd_font and {} or {
      cmd = '⌘',
      config = '🛠',
      event = '📅',
      ft = '📂',
      init = '⚙',
      keys = '🗝',
      plugin = '🔌',
      runtime = '💻',
      require = '🌙',
      source = '📄',
      start = '🚀',
      task = '📌',
      lazy = '💤 ',
    },
  },
})
