<?php
$value = 10;
function change() {
    global $value;
    unset($value);
    $value = 3;
    echo $value, '|';
}
change();
echo $value, "\n";
