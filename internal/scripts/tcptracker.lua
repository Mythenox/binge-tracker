local host = "127.0.0.1"
local port = 9099

local last_time = 0
local was_playing = false

function send_tcp_data(time_val)
    -- convert from microseconds to seconds
    local time_sec = time_val / 1000000
    local message = '{"event": "playback_stopped", "time_seconds": ' .. tostring(time_sec) .. '}\n'

    -- connect to tcp socket as client
    local fd = vlc.net.connect_tcp(host, port)
    
    -- fd will be non-negative if the connection succeeds
    if fd and fd >= 0 then
        vlc.net.send(fd, message)
        vlc.net.close(fd)
        vlc.msg.info("Successfully sent time via TCP: " .. tostring(time_sec) .. "s")
    else
        vlc.msg.err("Failed to connect to TCP server at " .. host .. ":" .. tostring(port))
    end
end

vlc.msg.info("TCP Close Tracker Interface started.")

-- Polling loop
while true do
    local input = vlc.object.input()
    
    if input then
        was_playing = true
        local current_time = vlc.var.get(input, "time")
        
        -- Cache the time as long as it's valid
        if current_time and current_time > 0 then
            last_time = current_time
        end
    else
        -- If the input disappears, the video was stopped or finished naturally
        if was_playing then
            send_tcp_data(last_time)
            was_playing = false
            last_time = 0
        end
    end
    
    local ok, err = pcall(function()
        vlc.misc.mwait(vlc.misc.mdate() + 500000)
    end)
    
    if not ok then
        break
    end
end

-- If VLC itself is closed while a video is actively playing
if was_playing then
    send_tcp_data(last_time)
end

vlc.msg.info("TCP Close Tracker Interface shut down.")