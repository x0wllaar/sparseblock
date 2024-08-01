vim.api.nvim_create_autocmd({"BufWritePost"}, {
  pattern = {"*.go"},
  callback = function(args) 
      local file_name = args.file
      vim.fn.system {'goimports', '-l', '-w', file_name}
      vim.fn.system {'gofumpt', '-l', '-w', file_name}
      vim.cmd("checktime")
  end,
})
