<?php
$shared = 1;
$other = 5;
function updateShared() {
    global $shared;
    $alias =& $shared;
    $alias = 2;
}
function detachGlobal() {
    global $shared, $other;
    $shared =& $other;
    $shared = 7;
}
function updateStatic() {
    static $counter = 1;
    $alias =& $counter;
    $alias++;
    echo $counter, '|';
}
updateShared();
echo $shared, '|';
detachGlobal();
echo $shared, '|', $other, '|';
updateStatic();
updateStatic();
