<?php
if (false) {
    set_error_handler('missingFunction');
    trigger_error('unreachable');
}
echo "ok\n";
