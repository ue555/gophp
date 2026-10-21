<?php
function namedHandler($level, $message) { echo "$level:$message|"; return true; }
set_error_handler(error_levels: E_USER_WARNING, callback: 'namedHandler');
trigger_error(error_level: E_USER_WARNING, message: 'named');
restore_error_handler();
echo "done\n";
