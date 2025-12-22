-- Keymaps configuration
-- See `:help vim.keymap.set()`

-- Clear search highlights on <Esc> in normal mode
vim.keymap.set('n', '<Esc>', '<cmd>nohlsearch<CR>')

-- Diagnostic quickfix list
vim.keymap.set('n', '<leader>q', vim.diagnostic.setloclist, { desc = 'Open diagnostic [Q]uickfix list' })

-- Exit terminal mode easier
vim.keymap.set('t', '<Esc><Esc>', '<C-\\><C-n>', { desc = 'Exit terminal mode' })

-- Split navigation with CTRL+hjkl
vim.keymap.set('n', '<C-h>', '<C-w><C-h>', { desc = 'Move focus to the left window' })
vim.keymap.set('n', '<C-l>', '<C-w><C-l>', { desc = 'Move focus to the right window' })
vim.keymap.set('n', '<C-j>', '<C-w><C-j>', { desc = 'Move focus to the lower window' })
vim.keymap.set('n', '<C-k>', '<C-w><C-k>', { desc = 'Move focus to the upper window' })

-- ====================
-- Familiar Keybinds
-- ====================

-- Save with Ctrl+S
vim.keymap.set({ 'n', 'i', 'v' }, '<C-s>', '<cmd>w<CR><Esc>', { desc = 'Save file' })

-- Quit with Ctrl+Q
vim.keymap.set('n', '<C-q>', '<cmd>q<CR>', { desc = 'Quit' })

-- Undo with Ctrl+Z
vim.keymap.set({ 'n', 'i' }, '<C-z>', '<cmd>undo<CR>', { desc = 'Undo' })

-- Redo with Ctrl+Y (common in many editors)
vim.keymap.set({ 'n', 'i' }, '<C-y>', '<cmd>redo<CR>', { desc = 'Redo' })

-- Select all with Ctrl+A
vim.keymap.set('n', '<C-a>', 'ggVG', { desc = 'Select all' })

-- Copy (already handled by clipboard sync, but explicit binds)
vim.keymap.set('v', '<C-c>', '"+y', { desc = 'Copy to clipboard' })

-- Cut
vim.keymap.set('v', '<C-x>', '"+d', { desc = 'Cut to clipboard' })

-- Paste with Ctrl+V (normal mode)
vim.keymap.set('n', '<C-v>', '"+p', { desc = 'Paste from clipboard' })
-- Paste with Ctrl+V (insert mode)
vim.keymap.set('i', '<C-v>', '<C-r>+', { desc = 'Paste from clipboard' })
-- Paste with Ctrl+V (visual mode - replace selection)
vim.keymap.set('v', '<C-v>', '"+p', { desc = 'Paste from clipboard' })

-- ====================
-- Jump Navigation
-- ====================

-- Jump up/down 10 lines with Ctrl+Up/Down
vim.keymap.set({ 'n', 'v' }, '<C-Up>', '10k', { desc = 'Jump up 10 lines' })
vim.keymap.set({ 'n', 'v' }, '<C-Down>', '10j', { desc = 'Jump down 10 lines' })
vim.keymap.set('i', '<C-Up>', '<C-o>10k', { desc = 'Jump up 10 lines' })
vim.keymap.set('i', '<C-Down>', '<C-o>10j', { desc = 'Jump down 10 lines' })

-- Jump to start/end of line with Shift+Left/Right
vim.keymap.set({ 'n', 'v' }, '<S-Left>', '^', { desc = 'Jump to start of line' })
vim.keymap.set({ 'n', 'v' }, '<S-Right>', '$', { desc = 'Jump to end of line' })
vim.keymap.set('i', '<S-Left>', '<C-o>^', { desc = 'Jump to start of line' })
vim.keymap.set('i', '<S-Right>', '<C-o>$', { desc = 'Jump to end of line' })

-- Jump to top/bottom of file with Shift+Up/Down
vim.keymap.set({ 'n', 'v' }, '<S-Up>', 'gg', { desc = 'Jump to top of file' })
vim.keymap.set({ 'n', 'v' }, '<S-Down>', 'G', { desc = 'Jump to bottom of file' })
vim.keymap.set('i', '<S-Up>', '<C-o>gg', { desc = 'Jump to top of file' })
vim.keymap.set('i', '<S-Down>', '<C-o>G', { desc = 'Jump to bottom of file' })

-- ====================
-- Selection Navigation (Ctrl+Shift)
-- ====================

-- Select while jumping up/down 10 lines with Ctrl+Shift+Up/Down
vim.keymap.set('n', '<C-S-Up>', 'v10k', { desc = 'Select up 10 lines' })
vim.keymap.set('n', '<C-S-Down>', 'v10j', { desc = 'Select down 10 lines' })
vim.keymap.set('v', '<C-S-Up>', '10k', { desc = 'Select up 10 lines' })
vim.keymap.set('v', '<C-S-Down>', '10j', { desc = 'Select down 10 lines' })
vim.keymap.set('i', '<C-S-Up>', '<C-o>v10k', { desc = 'Select up 10 lines' })
vim.keymap.set('i', '<C-S-Down>', '<C-o>v10j', { desc = 'Select down 10 lines' })

-- Select to start/end of line with Ctrl+Shift+Left/Right
vim.keymap.set('n', '<C-S-Left>', 'v^', { desc = 'Select to start of line' })
vim.keymap.set('n', '<C-S-Right>', 'v$', { desc = 'Select to end of line' })
vim.keymap.set('v', '<C-S-Left>', '^', { desc = 'Select to start of line' })
vim.keymap.set('v', '<C-S-Right>', '$', { desc = 'Select to end of line' })
vim.keymap.set('i', '<C-S-Left>', '<C-o>v^', { desc = 'Select to start of line' })
vim.keymap.set('i', '<C-S-Right>', '<C-o>v$', { desc = 'Select to end of line' })

-- Select to top/bottom of file with Ctrl+Shift+Home/End (alternative)
vim.keymap.set('n', '<C-S-Home>', 'vgg', { desc = 'Select to top of file' })
vim.keymap.set('n', '<C-S-End>', 'vG', { desc = 'Select to bottom of file' })
vim.keymap.set('v', '<C-S-Home>', 'gg', { desc = 'Select to top of file' })
vim.keymap.set('v', '<C-S-End>', 'G', { desc = 'Select to bottom of file' })
vim.keymap.set('i', '<C-S-Home>', '<C-o>vgg', { desc = 'Select to top of file' })
vim.keymap.set('i', '<C-S-End>', '<C-o>vG', { desc = 'Select to bottom of file' })

-- Ctrl+Backspace to delete word backwards (like in most editors)
vim.keymap.set('i', '<C-BS>', '<C-w>', { desc = 'Delete word backwards' })
vim.keymap.set('i', '<C-H>', '<C-w>', { desc = 'Delete word backwards (terminal fallback)' })

-- NOTE: Uncomment these to disable arrow keys in normal mode
-- vim.keymap.set('n', '<left>', '<cmd>echo "Use h to move!!"<CR>')
-- vim.keymap.set('n', '<right>', '<cmd>echo "Use l to move!!"<CR>')
-- vim.keymap.set('n', '<up>', '<cmd>echo "Use k to move!!"<CR>')
-- vim.keymap.set('n', '<down>', '<cmd>echo "Use j to move!!"<CR>')
