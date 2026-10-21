<?php
$stored = 2;
$calls = 0;
function &sourceForParameter() {
    global $stored, $calls;
    $calls++;
    return $stored;
}
function incrementReference(&$value) {
    $value++;
}
incrementReference(sourceForParameter());
echo $stored, '|', $calls;
