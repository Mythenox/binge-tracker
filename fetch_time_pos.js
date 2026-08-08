function fetch_time(){
    time_pos = mp.get_property("time-pos");
    if (time_pos != null) {
        mp.commandv("script-message", time_pos);
    }
}

mp.add_hook("on_unload", 50, fetch_time);