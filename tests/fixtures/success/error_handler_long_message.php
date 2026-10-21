<?php
function lengthOfWarning($level, $message) {
    echo strlen($message), "|";
    return true;
}
$hundred = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx";
$message = $hundred;
$message .= $message;
$message .= $message;
$message .= $message;
$message .= $hundred;
$message .= $hundred;
$message .= $hundred;
set_error_handler('lengthOfWarning');
trigger_error($message, E_USER_WARNING);
echo strlen($message), "\n";
