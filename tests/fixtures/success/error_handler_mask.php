<?php
function notices($level, $message) { echo "$level:$message|"; return true; }
set_error_handler('notices', E_USER_NOTICE | E_WARNING);
trigger_error('notice');
$missing;
echo $absent;
restore_error_handler();
echo "done\n";
