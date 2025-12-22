return {
	{
		"RRethy/base16-nvim",
		priority = 1000,
		config = function()
			require('base16-colorscheme').setup({
				base00 = '#141311',
				base01 = '#211f1d',
				base02 = '#363432',
				base03 = '#4c463d',
				base04 = '#cec5ba',
				base05 = '#e7e1de',
				base06 = '#e7e1de',
				base07 = '#32302e',
				base08 = '#ffb4ab',
				base09 = '#c4c5d6',
				base0A = '#cec5ba',
				base0B = '#d5c4aa',
				base0C = '#cfc5b8',
				base0D = '#d5c4aa',
				base0E = '#cfc5b8',
				base0F = '#ffb4ab',
			})

			local function set_hl_multiple(groups, value)
				for _, v in pairs(groups) do vim.api.nvim_set_hl(0, v, value) end
			end

			vim.api.nvim_set_hl(0, 'Visual',
				{ bg = '#b3a48b', fg = '#221a09', bold = true })
			vim.api.nvim_set_hl(0, 'LineNr', { fg = '#4c463d' })
			vim.api.nvim_set_hl(0, 'CursorLineNr', { fg = '#d5c4aa', bold = true })

			-- Hot-reload watcher
			local current_file_path = vim.fn.stdpath("config") .. "/lua/plugins/hecate-colors.lua"

			if not _G._hecate_theme_watcher then
				local uv = vim.uv or vim.loop
				_G._hecate_theme_watcher = uv.new_fs_event()

				_G._hecate_theme_watcher:start(current_file_path, {}, vim.schedule_wrap(function()
					local new_spec = dofile(current_file_path)

					if new_spec and new_spec[1] and new_spec[1].config then
						new_spec[1].config()
						print("HecateShell: Colors reloaded!")
					end
				end))
			end
		end
	}
}
