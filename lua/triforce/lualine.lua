local util = require('triforce.util')

---Lualine integration components for Triforce
---@class Triforce.Lualine
local Lualine = {}

-- STATE: Track the time coding value when Neovim first loaded.
local initial_time_coding = nil

---Default configuration for lualine components
---@class Triforce.Lualine.Config
Lualine.config = {
  ---Level component config
  ---@class Triforce.Lualine.Config.Level
  level = {
    prefix = 'Lv.', ---@type string Text prefix before level number
    show_level = true, ---@type boolean Show level number
    show_bar = true, ---@type boolean Show progress bar
    show_percent = false, ---@type boolean Show percentage
    show_xp = false, ---@type boolean Show XP numbers (current/needed)
    bar_length = 8, ---@type integer Length of progress bar
    bar_chars = { filled = '█', empty = '░' }, ---@type { filled: string, empty: string } Bar characters
  },

  ---Achievements component config
  ---@class Triforce.Lualine.Config.Achievements
  achievements = {
    icon = '', ---@type string|'' Nerd Font trophy icon
    show_count = true, ---@type boolean Show unlocked/total count
  },

  ---Streak component config
  ---@class Triforce.Lualine.Config.Streak
  streak = {
    icon = '', ---@type string|'' Nerd Font flame icon
    show_days = true, ---@type boolean Show number of days
  },

  ---Session time component config (Current Session)
  ---@class Triforce.Lualine.Config.SessionTime
  session_time = {
    -- Icons for different states
    status_icons = {
      active = '', -- Clock (Running)
      paused = '⏸', -- Pause (Frozen)
    },
    show_duration = true,
    format = 'digital', ---@type 'human'|'digital'|'clock'
  },

  ---Total time component config (Lifetime Stats)
  ---@class Triforce.Lualine.Config.TotalTime
  total_time = {
    icon = '󰔟',
    show_duration = true,
    format = 'human', ---@type 'human'|'digital'|'clock'
  },
}

---Setup lualine integration with custom config
---@param opts Triforce.Lualine.Config|nil User configuration
function Lualine.setup(opts)
  util.validate({ opts = { opts, { 'table', 'nil' }, true } })
  Lualine.config = vim.tbl_deep_extend('force', Lualine.config, opts or {})
end

---Get current stats AND global config safely
---@return Stats|nil stats
---@return TriforceConfig|nil config
local function get_env()
  local ok, triforce = pcall(require, 'triforce')
  if not ok then
    return nil, nil
  end

  local stats = triforce.get_stats()

  -- Initialize the baseline "Start of Session" time only once
  if stats and initial_time_coding == nil then
    initial_time_coding = stats.time_coding or 0
  end

  return stats, triforce.config
end

---Generate progress bar
---@param current number Current value
---@param max number Maximum value
---@param length integer Bar length
---@param chars table<string, string> Characters for filled and empty
---@return string bar
local function create_progress_bar(current, max, length, chars)
  if max == 0 then
    return chars.empty:rep(length)
  end
  local filled = math.min(math.floor((current / max) * length), length)
  return chars.filled:rep(filled) .. chars.empty:rep(length - filled)
end

---Format time duration
local function format_time(seconds, format)
  -- Default to digital if invalid format provided
  if not vim.list_contains({ 'human', 'digital', 'clock' }, format) then
    format = 'digital'
  end

  local hours = math.floor(seconds / 3600)
  local minutes = math.floor((seconds % 3600) / 60)
  local secs = seconds % 60

  if format == 'human' then
    if hours > 0 then
      return ('%dh %dm'):format(hours, minutes)
    elseif minutes > 0 then
      return ('%dm'):format(minutes)
    else
      return ('%ds'):format(secs)
    end
  elseif format == 'clock' then
    return ('%02d:%02d'):format(hours, minutes)
  else
    if hours > 0 then
      return ('%02d:%02d:%02d'):format(hours, minutes, secs)
    else
      return ('%02d:%02d'):format(minutes, secs)
    end
  end
end

---Level component - Shows level and XP progress
---@param opts Triforce.Lualine.Config.Level|nil Component-specific options
---@return string component
function Lualine.level(opts)
  util.validate({ opts = { opts, { 'table', 'nil' }, true } })

  local stats = get_env()
  if not stats then
    return ''
  end

  local config = vim.tbl_deep_extend('force', Lualine.config.level, opts or {})

  local stats_module = require('triforce.stats')
  local xp_for_current = stats_module.xp_for_next_level(stats.level - 1)
  local xp_for_next = stats_module.xp_for_next_level(stats.level)
  local xp_needed = xp_for_next - xp_for_current
  local xp_progress = stats.xp - xp_for_current

  local parts = {}
  if config.show_level then
    table.insert(parts, not config.prefix and tostring(stats.level) or (config.prefix .. stats.level))
  end
  if config.show_bar then
    table.insert(parts, create_progress_bar(xp_progress, xp_needed, config.bar_length, config.bar_chars))
  end
  if config.show_percent then
    table.insert(parts, ('%d%%'):format(math.floor(xp_progress / xp_needed) * 100))
  end
  if config.show_xp then
    table.insert(parts, ('%d/%d'):format(xp_progress, xp_needed))
  end

  return table.concat(parts, ' ')
end

---Achievements component - Shows unlocked achievement count
---@param opts Triforce.Lualine.Config.Achievements|nil Component-specific options
---@return string component
function Lualine.achievements(opts)
  util.validate({ opts = { opts, { 'table', 'nil' }, true } })

  local stats = get_env()
  if not stats then
    return ''
  end

  local config = vim.tbl_deep_extend('force', Lualine.config.achievements, opts or {})
  local all_achievements = require('triforce.achievement').get_all_achievements(stats)
  local unlocked = 0
  for _, _ in ipairs(stats.achievements or {}) do
    unlocked = unlocked + 1
  end

  local parts = {}
  if config.icon ~= '' then
    table.insert(parts, config.icon)
  end
  if config.show_count then
    table.insert(parts, ('%d/%d'):format(unlocked, #all_achievements))
  end

  return table.concat(parts, ' ')
end

---Streak component - Shows current coding streak
---@param opts Triforce.Lualine.Config.Streak|nil Component-specific options
---@return string|'' component
function Lualine.streak(opts)
  util.validate({ opts = { opts, { 'table', 'nil' }, true } })

  local stats = get_env()
  if not stats then
    return ''
  end
  local streak = stats.current_streak or 0
  if streak == 0 then
    return ''
  end

  local config = vim.tbl_deep_extend('force', Lualine.config.streak, opts or {})
  local parts = {}
  if config.icon ~= '' then
    table.insert(parts, config.icon)
  end
  if config.show_days then
    table.insert(parts, tostring(streak))
  end

  return table.concat(parts, ' ')
end

---Session time component - Shows time spent in THIS Neovim instance
---@param opts Triforce.Lualine.Config.SessionTime|nil Component-specific options
---@return string component
function Lualine.session_time(opts)
  util.validate({ opts = { opts, { 'table', 'nil' }, true } })

  local stats, global_config = get_env()
  if not stats then
    return ''
  end

  local config = vim.tbl_deep_extend('force', Lualine.config.session_time, opts or {})

  local fmt = (opts and opts.format) or (global_config and global_config.time_format) or config.format

  local current_lifetime_coding = stats.time_coding or 0
  if stats.session_active and stats.last_session_start > 0 then
    current_lifetime_coding = current_lifetime_coding + (os.time() - stats.last_session_start)
  end

  local session_duration = math.max(0, current_lifetime_coding - (initial_time_coding or 0))

  local icon
  if stats.session_active then
    icon = config.status_icons.active
  else
    icon = config.status_icons.paused
  end

  local parts = {}
  if icon ~= '' then
    table.insert(parts, icon)
  end
  if config.show_duration then
    table.insert(parts, format_time(session_duration, fmt))
  end

  return table.concat(parts, ' ')
end

---Total time component - Shows lifetime accumulated time
---@param opts Triforce.Lualine.Config.TotalTime|nil Component-specific options
---@return string component
function Lualine.total_time(opts)
  util.validate({ opts = { opts, { 'table', 'nil' }, true } })

  local stats, global_config = get_env()
  if not stats then
    return ''
  end

  local config = vim.tbl_deep_extend('force', Lualine.config.total_time, opts or {})

  local fmt = (opts and opts.format) or (global_config and global_config.time_format) or config.format

  local total_duration = stats.time_coding or 0
  if stats.session_active and stats.last_session_start > 0 then
    total_duration = total_duration + (os.time() - stats.last_session_start)
  end

  local parts = {}
  if config.icon ~= '' then
    table.insert(parts, config.icon)
  end
  if config.show_duration then
    table.insert(parts, format_time(total_duration, fmt))
  end

  return table.concat(parts, ' ')
end

---Convenience function to get all components at once
---@param opts Triforce.Lualine.Config|nil Configuration for all components
---@return Triforce.Lualine.Config components Table with level, achievements, streak, session_time functions
function Lualine.components(opts)
  util.validate({ opts = { opts, { 'table', 'nil' }, true } })
  Lualine.setup(opts)
  return {
    level = Lualine.level,
    achievements = Lualine.achievements,
    streak = Lualine.streak,
    session_time = Lualine.session_time,
    total_time = Lualine.total_time,
  }
end

return Lualine
-- vim:ts=2:sts=2:sw=2:et:ai:si:sta:
