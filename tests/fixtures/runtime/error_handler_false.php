<?php
function passThrough($level) { echo "handled:$level|"; return false; }
set_error_handler('passThrough');
echo $missing;
echo "after\n";
