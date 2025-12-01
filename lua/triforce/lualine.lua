local util = require('triforce.util')

---Lualine integration components for Triforce
---@class Triforce.Lualine
local Lualine = {}

-- STATE: Track the time coding value when Neovim first loaded.
-- This allows us to calculate "Session Time" as "Time since I opened Neovim"
-- even if we take breaks.
local initial_time_coding = nil

---Default configuration for lualine components
---@class Triforce.Lualine.Config
Lualine.config = {
  ---Level component config
  level = {
    prefix = 'Lv.',
    show_level = true,
    show_bar = true,
    show_percent = false,
    show_xp = false,
    bar_length = 8,
    bar_chars = { filled = '█', empty = '░' },
  },

  achievements = {
    icon = '',
    show_count = true,
  },

  streak = {
    icon = '',
    show_days = true,
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
  total_time = {
    icon = '󰔟',
    show_duration = true,
    format = 'human',
  },
}

---Setup lualine integration with custom config
---@param opts Triforce.Lualine.Config|nil User configuration
function Lualine.setup(opts)
  util.validate({ opts = { opts, { 'table', 'nil' }, true } })
  Lualine.config = vim.tbl_deep_extend('force', Lualine.config, opts or {})
end

---Get current stats safely and initialize baseline
---@return Stats|nil stats
local function get_stats()
  local ok, triforce = pcall(require, 'triforce')
  if not ok then return end

  local stats = triforce.get_stats()

  -- Initialize the baseline "Start of Session" time only once
  if stats and initial_time_coding == nil then
    initial_time_coding = stats.time_coding or 0
  end

  return stats
end

-- ... [create_progress_bar function remains the same] ...
local function create_progress_bar(current, max, length, chars)
  if max == 0 then return chars.empty:rep(length) end
  local filled = math.min(math.floor((current / max) * length), length)
  return chars.filled:rep(filled) .. chars.empty:rep(length - filled)
end

-- ... [format_time function remains the same] ...
local function format_time(seconds, format)
  -- Default to digital if invalid format provided
  if not vim.list_contains({ 'human', 'digital', 'clock' }, format) then
    format = 'digital'
  end

  local hours = math.floor(seconds / 3600)
  local minutes = math.floor((seconds % 3600) / 60)
  local secs = seconds % 60

  if format == 'human' then
    if hours > 0 then return ('%dh %dm'):format(hours, minutes)
    elseif minutes > 0 then return ('%dm'):format(minutes)
    else return ('%ds'):format(secs) end
  elseif format == 'clock' then
    return ('%02d:%02d'):format(hours, minutes)
  else
    if hours > 0 then return ('%02d:%02d:%02d'):format(hours, minutes, secs)
    else return ('%02d:%02d'):format(minutes, secs) end
  end
end

-- ... [Lualine.level, achievements, streak remain the same] ...
function Lualine.level(opts)
  -- (Same as previous code)
  local stats = get_stats()
  if not stats then return '' end
  local config = vim.tbl_deep_extend('force', Lualine.config.level, opts or {})
  -- XP Logic...
  local stats_module = require('triforce.stats')
  local xp_for_current = stats_module.xp_for_next_level(stats.level - 1)
  local xp_for_next = stats_module.xp_for_next_level(stats.level)
  local xp_needed = xp_for_next - xp_for_current
  local xp_progress = stats.xp - xp_for_current
  
  local parts = {}
  if config.show_level then table.insert(parts, not config.prefix and tostring(stats.level) or (config.prefix .. stats.level)) end
  if config.show_bar then table.insert(parts, create_progress_bar(xp_progress, xp_needed, config.bar_length, config.bar_chars)) end
  if config.show_percent then table.insert(parts, ('%d%%'):format(math.floor(xp_progress / xp_needed) * 100)) end
  if config.show_xp then table.insert(parts, ('%d/%d'):format(xp_progress, xp_needed)) end
  return table.concat(parts, ' ')
end

function Lualine.achievements(opts)
  local stats = get_stats()
  if not stats then return '' end
  local config = vim.tbl_deep_extend('force', Lualine.config.achievements, opts or {})
  local all_achievements = require('triforce.achievement').get_all_achievements(stats)
  local unlocked = 0
  for _, _ in ipairs(stats.achievements or {}) do unlocked = unlocked + 1 end
  local parts = {}
  if config.icon ~= '' then table.insert(parts, config.icon) end
  if config.show_count then table.insert(parts, ('%d/%d'):format(unlocked, #all_achievements)) end
  return table.concat(parts, ' ')
end

function Lualine.streak(opts)
  local stats = get_stats()
  if not stats then return '' end
  local streak = stats.current_streak or 0
  if streak == 0 then return '' end
  local config = vim.tbl_deep_extend('force', Lualine.config.streak, opts or {})
  local parts = {}
  if config.icon ~= '' then table.insert(parts, config.icon) end
  if config.show_days then table.insert(parts, tostring(streak)) end
  return table.concat(parts, ' ')
end

---Session time component - Shows time spent in THIS Neovim instance
---Pauses when idle, resumes when active.
---@param opts Triforce.Lualine.Config.SessionTime|nil Component-specific options
---@return string component
function Lualine.session_time(opts)
  util.validate({ opts = { opts, { 'table', 'nil' }, true } })

  local stats = get_stats()
  if not stats then
    return ''
  end

  local config = vim.tbl_deep_extend('force', Lualine.config.session_time, opts or {})

  -- 1. Calculate the real-time "Lifetime" coding duration
  local current_lifetime_coding = stats.time_coding or 0
  if stats.session_active and stats.last_session_start > 0 then
    -- Add the partial uncommitted time from the current burst
    current_lifetime_coding = current_lifetime_coding + (os.time() - stats.last_session_start)
  end

  -- 2. Subtract the baseline (time when we opened Neovim)
  -- This gives us "Time spent coding since I opened this window"
  local session_duration = math.max(0, current_lifetime_coding - (initial_time_coding or 0))

  -- 3. Determine Icon (Active vs Paused)
  local icon
  if stats.session_active then
    icon = config.status_icons.active -- Clock
  else
    icon = config.status_icons.paused -- Pause
  end

  -- Build component
  local parts = {}
  if icon ~= '' then
    table.insert(parts, icon)
  end

  if config.show_duration then
    table.insert(parts, format_time(session_duration, config.format))
  end

  return table.concat(parts, ' ')
end

---Total time component - Shows lifetime accumulated time
---@param opts Triforce.Lualine.Config.TotalTime|nil Component-specific options
---@return string component
function Lualine.total_time(opts)
  util.validate({ opts = { opts, { 'table', 'nil' }, true } })

  local stats = get_stats()
  if not stats then
    return ''
  end

  local config = vim.tbl_deep_extend('force', Lualine.config.total_time, opts or {})

  -- Calculate TOTAL duration (Historical + Current Session)
  local total_duration = stats.time_coding or 0
  if stats.session_active and stats.last_session_start > 0 then
    total_duration = total_duration + (os.time() - stats.last_session_start)
  end

  local parts = {}
  if config.icon ~= '' then
    table.insert(parts, config.icon)
  end
  if config.show_duration then
    table.insert(parts, format_time(total_duration, config.format))
  end

  return table.concat(parts, ' ')
end

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